import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import {terser} from 'rollup-plugin-terser';
import * as path from 'path';

// `npm run build` -> `production` is true
// `npm run dev` -> `production` is false
const production = !process.env.ROLLUP_WATCH;
const webOutputDir = path.join(
  process.env.WEB_OUTPUT_DIR || 'public',
  'bundle.js'
);

export default {
  input: 'src/main.js',
  output: {
    file: webOutputDir,
    format: 'iife', // immediately-invoked function expression — suitable for <script> tags
    sourcemap: !production,
  },
  plugins: [
    resolve(), // tells Rollup how to find node_modules
    commonjs(), // converts to ES modules
    production && terser(), // minify, but only in production
  ],
};
