// Copy the wasm-pack output into dist/wasm so consumers receive it alongside
// the JS bundle. Also writes dist/wasm/fractals_wasm.js with a tiny patch: the
// default export takes an optional URL override so host apps can point at a
// static-asset path of their choosing.
import { copyFileSync, cpSync, existsSync, mkdirSync, readdirSync, statSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const root = join(__dirname, "..");
const src = join(root, "..", "..", "fractals-wasm", "pkg");
const dst = join(root, "dist", "wasm");

if (!existsSync(src)) {
  console.error(`wasm-pack output missing at ${src}. Run build.sh first.`);
  process.exit(1);
}

mkdirSync(dst, { recursive: true });
for (const f of readdirSync(src)) {
  if (f === "package.json" || f === "README.md") continue;
  const s = join(src, f);
  const d = join(dst, f);
  if (statSync(s).isDirectory()) {
    cpSync(s, d, { recursive: true });
    console.log("copied", f + "/");
  } else {
    copyFileSync(s, d);
    console.log("copied", f);
  }
}
