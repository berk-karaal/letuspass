import "dotenv/config";
import { defineConfig } from "orval";

export default defineConfig({
  letuspass: {
    input: process.env.VITE_BACKEND_BASE_URL + "/swagger/doc.json",
    output: {
      mode: "split",
      target: "./src/api/letuspass.ts",
      override: {
        mutator: {
          path: "./src/api/mutator/custom-instance.ts",
          name: "customInstance",
        },
      },
    },
    hooks: {
      afterAllFilesWrite: "prettier --write",
    },
  },
});
