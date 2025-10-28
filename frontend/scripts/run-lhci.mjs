import { spawn } from "node:child_process";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { findChromiumBinary } from "./find-chrome.mjs";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const projectRoot = path.resolve(__dirname, "..");

async function run() {
  const env = { ...process.env };

  const chromeBinary = findChromiumBinary();
  if (chromeBinary) {
    env.CHROME_PATH = chromeBinary;
    env.LIGHTHOUSE_CHROMIUM_PATH = chromeBinary;
    console.log(`[lhci] Using Chromium binary at ${chromeBinary}`);
  }

  env.NODE_OPTIONS = "--experimental-json-modules";
  if (!env.LHCI_LOG_LEVEL) {
    env.LHCI_LOG_LEVEL = "verbose";
  }

  const child = spawn("pnpm", ["exec", "lhci", "autorun"], {
    cwd: projectRoot,
    env,
    stdio: "inherit"
  });

  child.on("close", (code, signal) => {
    if (signal) {
      process.kill(process.pid, signal);
    } else {
      process.exit(code ?? 1);
    }
  });
}

run().catch((error) => {
  console.error(error);
  process.exit(1);
});
