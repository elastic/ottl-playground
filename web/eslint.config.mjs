import globals from "globals";
import pluginJs from "@eslint/js";


export default [
    {
        languageOptions: {globals: globals.browser},
    },
    {
        ignores: [
            "node_modules/*",
            "rollup.config.js",
            "**/wasm_exec.js",
            "**/bundle.js",
        ]
    },
    pluginJs.configs.recommended,
];