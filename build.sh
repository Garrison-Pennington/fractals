#!/usr/bin/env bash
#
# Full pipeline: Rust → WASM → React component package → demo assets.
#
# Steps:
#   1. cargo test (workspace) — validates the core.
#   2. wasm-pack build --target web fractals-wasm — emits fractals-wasm/pkg/.
#   3. Install explorer package deps and run tsc + copy-wasm.mjs.
#   4. Copy wasm + JS into demo/public/wasm so the demo can serve them.
#
# The wasm-pack build enables the threading-friendly target features so a
# future Web Worker pool can use SharedArrayBuffer without rebuilding.
#
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

echo "[1/5] cargo test (workspace core)"
cargo test -p fractals-core --quiet

echo "[2/5] wasm-pack build --target web fractals-wasm"
# The threading-capable target features below are compatible with per-worker
# WASM instances sharing output buffers via SharedArrayBuffer. We opt in only
# when a nightly toolchain with -Z build-std is available (required for
# thread-safe std); otherwise we fall back to the default wasm32 target which
# still supports SharedArrayBuffer across separate WASM instances per worker.
if rustup toolchain list | grep -q '^nightly' && rustup +nightly component list --installed 2>/dev/null | grep -q '^rust-src'; then
    echo "  (threaded build: using nightly + build-std)"
    RUSTUP_TOOLCHAIN=nightly \
    RUSTFLAGS="-C target-feature=+atomics,+bulk-memory,+mutable-globals" \
        wasm-pack build --release --target web fractals-wasm \
        -- -Z build-std=panic_abort,std
else
    echo "  (single-instance build: nightly/rust-src not installed; worker pool will run independent WASM instances)"
    wasm-pack build --release --target web fractals-wasm
fi

echo "[3/5] npm install + tsc in packages/explorer"
pushd packages/explorer > /dev/null
if [ ! -d node_modules ]; then
    npm install --no-audit --no-fund
fi
rm -rf dist
npx tsc -p tsconfig.json
node scripts/copy-wasm.mjs
popd > /dev/null

echo "[4/5] copy wasm artefacts to demo/public/wasm"
mkdir -p demo/public/wasm
cp -rf packages/explorer/dist/wasm/* demo/public/wasm/

echo "[5/5] done."
echo
echo "Next: cd demo && npm install && npm run dev"
