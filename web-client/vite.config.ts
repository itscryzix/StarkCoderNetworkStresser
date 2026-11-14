import react from "@vitejs/plugin-react";
import path from "path";
import { fileURLToPath } from "url";
import { defineConfig } from "vite";

// __dirname is not available under ESM ("type": "module"), so derive it.
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  optimizeDeps: {
    exclude: [],
  },
  build: {
    outDir: path.resolve(__dirname, "dist/public"),
  },
  server: {
    strictPort: true,
    // Proxy API + Socket.IO to the Go backend during `npm run dev`.
    proxy: {
      "/socket.io": {
        target: "http://localhost:3000",
        ws: true,
      },
      "/attacks": "http://localhost:3000",
      "/configuration": "http://localhost:3000",
    },
  },
});
