import fs from "node:fs";
import os from "node:os";
import path from "node:path";

function isExecutable(filePath) {
  try {
    fs.accessSync(filePath, fs.constants.X_OK);
    return true;
  } catch {
    return false;
  }
}

export function findChromiumBinary() {
  const overridden = (process.env.CHROME_PATH || "").trim();
  if (overridden.length > 0) {
    return overridden;
  }

  const home = os.homedir();
  const playwrightCache = path.join(home, ".cache", "ms-playwright");

  let entries = [];
  try {
    entries = fs.readdirSync(playwrightCache, { withFileTypes: true });
  } catch {
    return "";
  }

  const chromiumDirs = entries
    .filter((entry) => entry.isDirectory() && entry.name.startsWith("chromium"))
    .map((entry) => ({
      fullPath: path.join(playwrightCache, entry.name),
      name: entry.name
    }))
    .sort((a, b) => {
      const aIsHeadless = a.name.includes("headless");
      const bIsHeadless = b.name.includes("headless");
      if (aIsHeadless !== bIsHeadless) {
        return aIsHeadless - bIsHeadless;
      }
      return a.name > b.name ? -1 : 1;
    });

  for (const entry of chromiumDirs) {
    const chromeBinary = path.join(entry.fullPath, "chrome-linux", "chrome");
    if (isExecutable(chromeBinary)) {
      return chromeBinary;
    }

    const headlessShell = path.join(entry.fullPath, "chrome-linux", "headless_shell");
    if (isExecutable(headlessShell)) {
      return headlessShell;
    }
  }

  return "";
}

if (import.meta.url === `file://${process.argv[1]}`) {
  const candidate = findChromiumBinary();
  if (candidate) {
    process.stdout.write(candidate);
  }
}
