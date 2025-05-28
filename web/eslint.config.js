import js from "@eslint/js";
import tailwind from "@hyoban/eslint-plugin-tailwindcss";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import fs from "fs";
import globals from "globals";
import path from "path";
import tseslint from "typescript-eslint";

function findTailwindImportCss(dir) {
  const entries = fs.readdirSync(dir, { withFileTypes: true });

  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);

    if (entry.isDirectory()) {
      const found = findTailwindImportCss(fullPath);
      if (found) return found;
    } else if (entry.isFile() && entry.name.endsWith(".css")) {
      // read & scan lines
      const lines = fs.readFileSync(fullPath, "utf8").split(/\r?\n/);
      for (let line of lines) {
        if (line.trim().startsWith('@import "tailwindcss')) {
          return fullPath;
        }
      }
    }
  }

  return null;
}

export default tseslint.config(
  { ignores: ["dist"] },
  {
    extends: [js.configs.recommended, ...tseslint.configs.recommended],
    files: ["**/*.{ts,tsx}"],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
    plugins: {
      "react-hooks": reactHooks,
      "react-refresh": reactRefresh,
    },
    rules: {
      ...reactHooks.configs.recommended.rules,
      "react-refresh/only-export-components": [
        "warn",
        { allowConstantExport: true },
      ],
    },
  },
  // Tailwind plugin config (split into two objects)
  ...tailwind.configs["flat/recommended"],
  {
    settings: {
      tailwindcss: {
        config: findTailwindImportCss(process.cwd()), // <- uses process.cwd() to start looking down from where eslint is run
      },
    },
  }
);
