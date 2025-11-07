import fs from "node:fs/promises";
import net from "node:net";
import { createRequire } from "node:module";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const projectRoot = path.resolve(__dirname, "..");
const publicAppRoot = path.join(projectRoot, "apps", "public");
const previewMetadataDir = path.join(projectRoot, ".lighthouseci");
const previewPidFile = path.join(previewMetadataDir, "preview.pid");

const requireFromPublic = createRequire(path.join(publicAppRoot, "package.json"));
const previewMode = process.env.LHCI_PREVIEW_MODE ?? "production";

async function fileExists(filePath) {
  try {
    await fs.access(filePath);
    return true;
  } catch {
    return false;
  }
}

async function wait(ms) {
  await new Promise(resolve => setTimeout(resolve, ms));
}

async function isPortFree(host, port) {
  return new Promise(resolve => {
    const server = net.createServer();
    server.once("error", () => {
      resolve(false);
    });
    server.once("listening", () => {
      server.close(() => resolve(true));
    });
    server.listen(port, host);
  });
}

async function waitForPortRelease(host, port, attempts = 15) {
  for (let attempt = 0; attempt < attempts; attempt += 1) {
    if (await isPortFree(host, port)) {
      return true;
    }
    await wait(200);
  }
  return false;
}

async function cleanupPidFile() {
  await fs.unlink(previewPidFile).catch(() => {});
}

async function ensurePreviewPortAvailable(host, port) {
  const hasPidFile = await fileExists(previewPidFile);

  if (hasPidFile) {
    const rawPid = await fs.readFile(previewPidFile, "utf8").catch(() => "");
    const pid = Number.parseInt(rawPid, 10);

    if (Number.isInteger(pid)) {
      try {
        process.kill(pid, "SIGTERM");
      } catch (error) {
        if (error?.code !== "ESRCH") {
          console.warn(`[preview] Failed to stop previous preview process (${pid}): ${error.message}`);
        }
      }
    }

    const released = await waitForPortRelease(host, port);
    if (!released) {
      throw new Error(
        `[preview] Port ${port} on ${host} is still in use after attempting to stop the previous preview server.`,
      );
    }

    await cleanupPidFile();
    return;
  }

  if (!(await isPortFree(host, port))) {
    throw new Error(
      `[preview] Port ${port} on ${host} is already in use. Stop the existing process or set LHCI_PREVIEW_PORT.`,
    );
  }
}

async function startPreviewServer() {
  const { preview } = requireFromPublic("vite");

  const port = Number(process.env.LHCI_PREVIEW_PORT ?? 4173);
  const host = process.env.LHCI_PREVIEW_HOST ?? "127.0.0.1";

  await ensurePreviewPortAvailable(host, port);

  const previewServer = await preview({
    root: publicAppRoot,
    configFile: path.join(publicAppRoot, "vite.config.ts"),
    mode: previewMode,
    preview: {
      host,
      port,
      strictPort: true,
      open: false
    }
  });

  await fs.mkdir(previewMetadataDir, { recursive: true });
  await fs.writeFile(previewPidFile, String(process.pid));

  const localUrl =
    previewServer.resolvedUrls?.local?.[0] ?? `http://${host}:${previewServer.config.preview.port}/`;
  console.log(`Local: ${localUrl}`);

  const shutdown = async () => {
    await new Promise((resolve, reject) =>
      previewServer.httpServer.close(error => (error ? reject(error) : resolve()))
    );
    await cleanupPidFile();
    process.exit(0);
  };

  process.on("SIGTERM", () => {
    shutdown().catch(error => {
      console.error("[preview] Failed to stop server", error);
      process.exit(1);
    });
  });
  process.on("SIGINT", () => {
    shutdown().catch(error => {
      console.error("[preview] Failed to stop server", error);
      process.exit(1);
    });
  });

  process.on("exit", () => {
    cleanupPidFile().catch(error => {
      console.warn("[preview] Failed to clean up preview PID file", error);
    });
  });

  // Keep the process alive until LHCI stops it.
  await new Promise(() => {});
}

startPreviewServer().catch(error => {
  console.error("[preview] Failed to start server", error);
  process.exit(1);
});
