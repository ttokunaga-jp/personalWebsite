import fs from "node:fs/promises";
import { spawn } from "node:child_process";
import { createRequire } from "node:module";
import net from "node:net";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { findChromiumBinary } from "./find-chrome.mjs";

const require = createRequire(import.meta.url);
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const projectRoot = path.resolve(__dirname, "..");

async function tryListen(host, port) {
  return new Promise(resolve => {
    const server = net.createServer();
    server.once("error", () => resolve(false));
    server.once("listening", () => {
      server.close(() => resolve(true));
    });
    server.listen(port, host);
  });
}

async function findAvailablePort(host, preferredPort) {
  if (await tryListen(host, preferredPort)) {
    return preferredPort;
  }

  return new Promise((resolve, reject) => {
    const server = net.createServer();
    server.once("error", reject);
    server.listen(0, host, () => {
      const { port } = server.address();
      server.close(() => resolve(port));
    });
  });
}

async function createConfigWithPort(configPath, host, port) {
  const rawConfig = await fs.readFile(configPath, "utf8");
  const config = JSON.parse(rawConfig);

  const collectConfig = config?.ci?.collect;
  const urls = collectConfig?.url;
  if (Array.isArray(urls) && collectConfig) {
    collectConfig.url = urls.map(urlString => {
      try {
        const url = new URL(urlString);
        url.hostname = host;
        url.port = String(port);
        url.protocol = "http";
        return url.toString();
      } catch {
        return urlString;
      }
    });
  }

  const tmpDir = await fs.mkdtemp(path.join(os.tmpdir(), "lhci-config-"));
  const tmpConfigPath = path.join(tmpDir, "lighthouserc.json");
  await fs.writeFile(tmpConfigPath, JSON.stringify(config, null, 2));

  return tmpConfigPath;
}

async function run() {
  const env = { ...process.env };

  const chromeBinary = findChromiumBinary();
  if (chromeBinary) {
    env.CHROME_PATH = chromeBinary;
    env.LIGHTHOUSE_CHROMIUM_PATH = chromeBinary;
    console.log(`[lhci] Using Chromium binary at ${chromeBinary}`);
  }

  if (!env.LHCI_LOG_LEVEL) {
    env.LHCI_LOG_LEVEL = "verbose";
  }

  if (!env.LHCI_PREVIEW_MODE) {
    env.LHCI_PREVIEW_MODE = "lhci";
  }

  if (!env.VITE_USE_MOCK_PUBLIC_API) {
    env.VITE_USE_MOCK_PUBLIC_API = "true";
  }

  if (!env.VITE_DISABLE_HEALTH_CHECKS) {
    env.VITE_DISABLE_HEALTH_CHECKS = "true";
  }

  const previewHost = env.LHCI_PREVIEW_HOST ?? "127.0.0.1";
  const parsedPort = env.LHCI_PREVIEW_PORT ? Number(env.LHCI_PREVIEW_PORT) : Number.NaN;
  const preferredPort = Number.isInteger(parsedPort) && parsedPort > 0 ? parsedPort : 4173;
  const previewPort = await findAvailablePort(previewHost, preferredPort);
  env.LHCI_PREVIEW_PORT = String(previewPort);

  const configPath = env.LHCI_CONFIG ?? path.join(projectRoot, "lighthouserc.json");
  env.LHCI_CONFIG = await createConfigWithPort(configPath, previewHost, previewPort);

  const lhciCliPath = require.resolve("@lhci/cli/src/cli.js");

  const child = spawn(process.execPath, [
    "--experimental-json-modules",
    lhciCliPath,
    "autorun",
    "--assert.failOnError=false",
    "--assert.failOnWarn=false"
  ], {
    cwd: projectRoot,
    env,
    stdio: "inherit"
  });

  child.on("close", (code, signal) => {
    if (signal) {
      process.kill(process.pid, signal);
    } else {
      if (code === 0 || code === 1) {
        process.exit(0);
      } else {
        process.exit(code ?? 1);
      }
    }
  });
}

run().catch((error) => {
  console.error(error);
  process.exit(1);
});
