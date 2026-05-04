import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
  },
  optimizeDeps: {
    // The explorer dynamically imports the wasm-pack-generated JS at runtime
    // from /wasm/, which is served from public/. Excluding it keeps Vite's
    // dep optimizer from trying to prebundle the package's wasm subpath.
    exclude: ["@fractals/explorer"],
  },
});
