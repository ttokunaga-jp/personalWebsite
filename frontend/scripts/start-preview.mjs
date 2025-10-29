import { createRequire } from "node:module";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const projectRoot = path.resolve(__dirname, "..");
const publicAppRoot = path.join(projectRoot, "apps", "public");

const requireFromPublic = createRequire(path.join(publicAppRoot, "package.json"));

async function startPreviewServer() {
  const { preview } = requireFromPublic("vite");

  const port = Number(process.env.LHCI_PREVIEW_PORT ?? 4173);
  const host = process.env.LHCI_PREVIEW_HOST ?? "127.0.0.1";

  const previewServer = await preview({
    root: publicAppRoot,
    configFile: path.join(publicAppRoot, "vite.config.ts"),
    preview: {
      host,
      port,
      strictPort: true,
      open: false
    }
  });

  const localUrl =
    previewServer.resolvedUrls?.local?.[0] ?? `http://${host}:${previewServer.config.preview.port}/`;
  console.log(`Local: ${localUrl}`);

  const shutdown = async () => {
    await new Promise((resolve, reject) =>
      previewServer.httpServer.close(error => (error ? reject(error) : resolve()))
    );
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

  // Keep the process alive until LHCI stops it.
  await new Promise(() => {});
}

startPreviewServer().catch(error => {
  console.error("[preview] Failed to start server", error);
  process.exit(1);
});
