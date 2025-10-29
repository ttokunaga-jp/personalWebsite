import { spawn } from "node:child_process";
import { createRequire } from "node:module";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { findChromiumBinary } from "./find-chrome.mjs";

const require = createRequire(import.meta.url);
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

  if (!env.LHCI_LOG_LEVEL) {
    env.LHCI_LOG_LEVEL = "verbose";
  }

  const lhciCliPath = require.resolve("@lhci/cli/src/cli.js");

  const child = spawn(
    process.execPath,
    ["--experimental-json-modules", lhciCliPath, "autorun"],
    {
      cwd: projectRoot,
      env,
      stdio: "inherit"
    }
  );

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
