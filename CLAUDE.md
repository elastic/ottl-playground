# OTTL Playground

Interactive browser-based playground for experimenting with [OTTL (OpenTelemetry Transformation Language)](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/ottl). Users write OTTL configs and run them against OTLP JSON payloads entirely in-browser — no server round-trips for evaluation.

## Architecture

The project has three distinct layers:

**Go core (`internal/`)** — implements the `Executor` interface for each OTel Collector processor. Currently two executors: `TransformProcessorExecutor` and `FilterProcessorExecutor`. Each supports all four signal types: logs, traces, metrics, profiles.

**WebAssembly bridge (`wasm/`)** — compiles the Go core to WASM (`GOARCH=wasm GOOS=js`). Exposes two JS globals on startup:
- `execute(config, payloadType, payload, executorId, debug)` → result JSON. The `debug` boolean activates breakpoint mode; when true, `result.debug === true` and `result.value` contains per-line results rather than a final payload.
- `getExecutors()` → list of available executors, each with `id`, `name`, `type`, `version`, `docsURL`, and `examples` fields

**Frontend (`web/src/`)** — vanilla JS with [Lit](https://lit.dev/) Web Components, currently bundled by Rollup into `web/public/bundle.js`. Main component is `playground-stage` (`playground.js`), which fetches `wasm/versions.json`, loads the correct `.wasm` file at runtime, and wires up the editor panels. **Migration in progress** to Vue 3 + Web Awesome + Vite (see Frontend Libraries and UI Improvement Plan sections).

All build outputs land in `web/public/` — this directory is the static site root.

## Requirements

- Go 1.24+
- Node.js (for the frontend)

## Common Commands

```bash
# Full build (WASM + frontend)
make build

# Build only the WASM
make build-wasm

# Build only the frontend
make build-web

# Run a local dev server (after building)
go run main.go

# Frontend dev with live reload
cd web && npm run watch   # rebuild on change
cd web && npm run start   # serve public/

# Go tests
go test ./internal/...

# Frontend linting / formatting
cd web && npm run lint:eslint
cd web && npm run format
```

## Version Management

The collector contrib version (e.g. `v0.138.0`) is the single source of truth, derived from `go.mod`. The `ci-tools/main.go` CLI manages version lifecycle:

- `generate-constants` — regenerates `internal/versions.go` (do not edit by hand)
- `register-wasm` — appends a version entry to `web/public/wasm/versions.json`
- `get-unregistered-versions` — queries the GitHub releases API to find new versions to build
- `validate-registered-versions` — checks that every entry in `versions.json` has a corresponding `.wasm` file

To bump the collector version:
```bash
export PROCESSORS_VERSION=vX.Y.Z
make update-processor-version   # updates go.mod/go.sum
make build-wasm                 # compiles new .wasm
make register-version           # updates versions.json
```

## Key Files

| Path | Purpose |
|------|---------|
| `internal/executor.go` | `Executor` interface + registry |
| `internal/transformprocessorexecutor.go` | Transform processor executor |
| `internal/filterprocessorexecutor.go` | Filter processor executor |
| `internal/versions.go` | Generated — current collector version constant |
| `wasm/main.go` | WASM entry point; registers `execute()` and `getExecutors()` JS globals |
| `web/src/components/playground.js` | Root `playground-stage` LitElement |
| `web/public/wasm/versions.json` | WASM version manifest read at runtime |
| `ci-tools/main.go` | Build/version automation CLI |
| `Makefile` | All build targets |

## Frontend Libraries

### Lit (current — v3.x)

[Lit](https://lit.dev/) is the Web Components base library used for all existing components. Every component extends `LitElement` and is registered as a native custom element.

Key patterns used in this codebase:

```js
import {css, html, LitElement} from 'lit';
import {nothing} from 'lit';
import {repeat} from 'lit/directives/repeat.js';

export class MyComponent extends LitElement {
  // Declare reactive properties — any change triggers re-render
  static properties = {
    value: {type: String},
    loading: {type: Boolean},
    items: {type: Object},           // arrays/objects passed as JSON strings from HTML
    hidden: {type: Boolean, attribute: 'hide-thing'},  // kebab-case HTML attr
    _internal: {state: true},        // private state, no attribute reflection
  };

  static get styles() {
    return css`...`;  // Shadow DOM scoped — no leakage in or out
  }

  render() {
    return html`<div>${this.loading ? html`<span>...</span>` : nothing}</div>`;
  }
}
customElements.define('my-component', MyComponent);
```

**Event discipline**: custom events crossing shadow DOM boundaries need `bubbles: true, composed: true`. Events that stay within a component's own shadow tree don't need `composed`.

```js
this.dispatchEvent(new CustomEvent('value-changed', {
  detail: {value: newVal},
  bubbles: true,
  composed: true,
}));
```

**Slot pattern**: the playground uses named slots extensively for customisation points (e.g. `<slot name="app-title-text">`). Slotted content lives in the light DOM and is styled by the consumer, not the component.

---

### Web Awesome (migration target)

[Web Awesome](https://webawesome.com/) is the successor to Shoelace — a framework-agnostic Web Components UI kit. It provides production-ready components (buttons, selects, tabs, dialogs, split panels, icons, etc.) as native custom elements with `<wa-*>` tags.

Since Web Awesome components are standard custom elements, they drop in alongside Lit components and Vue components with no special bridge.

**Install:**
```bash
npm install @web-awesome/components
```

**Usage in templates (Lit or plain HTML):**
```html
<wa-button variant="primary">Run</wa-button>
<wa-select label="Evaluator">
  <wa-option value="transform_processor">Transform processor</wa-option>
</wa-select>
<wa-split-panel>
  <div slot="start">Left</div>
  <div slot="end">Right</div>
</wa-split-panel>
<wa-tab-group>
  <wa-tab slot="nav" panel="logs">Logs</wa-tab>
  <wa-tab-panel name="logs">...</wa-tab-panel>
</wa-tab-group>
```

**Theming**: Web Awesome uses CSS custom properties (`--wa-color-primary`, `--wa-border-radius-medium`, etc.) defined at `:root`. Override them globally or scope them to a container.

**Relevant components for this project**: `wa-button`, `wa-select`/`wa-option`, `wa-tab-group`/`wa-tab`/`wa-tab-panel`, `wa-split-panel` (can replace `split.js`), `wa-spinner`, `wa-tooltip`, `wa-icon`.

---

### Vue 3 (confirmed framework — migration in progress)

[Vue 3](https://vuejs.org/) is the application framework replacing Lit. Use **Composition API with `<script setup>`** — it is the recommended style for new components and integrates cleanly with TypeScript if added later.

**SFC structure:**
```vue
<script setup>
import {ref, computed} from 'vue'
import '@web-awesome/components/wa-button.js'  // side-effect import registers the custom element

const loading = ref(false)
const label = computed(() => loading.value ? 'Running...' : 'Run')

function run() { /* ... */ }
</script>

<template>
  <wa-button :disabled="loading" @click="run">{{ label }}</wa-button>
</template>

<style scoped>
/* scoped styles don't pierce into Web Awesome shadow DOM — use CSS vars for theming */
</style>
```

**Build tool**: Vite replaces Rollup. Vitest (Vite's test runner) is used for frontend tests.

**Web Components in Vue**: Vue warns about unknown custom elements unless you tell it which tags to skip component resolution for. Configure this in `vite.config.js`:

```js
// vite.config.js
import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => tag.startsWith('wa-') || tag.startsWith('playground-'),
        },
      },
    }),
  ],
})
```

**v-model with Web Components**: Vue's `v-model` doesn't bind to custom element properties by default. Use explicit prop + event binding:
```vue
<wa-select :value="evaluator" @sl-change="evaluator = $event.target.value" />
```

**Emitting events from Vue components used inside Lit**: if a Vue component needs to be wrapped in a `LitElement` for the existing shell, emit native DOM events (`new CustomEvent(...)`) from `onMounted` / lifecycle hooks rather than relying on Vue's `emit`.

---

## UI Improvement Plan

Four phases, informed by structural analysis of `device-builder-frontend` and `ui` reference repos.

### Phase 1 — Directory restructure + WASM bridge extraction

**Part A — Restructure `web/src/`**

Move from the current flat file list to a folder-per-component model. Also establish the parallel `web/test/` directory that mirrors `src/` exactly — this is the pattern used by all reference repos.

```
web/
  src/
    components/
      playground/          ← was playground.js + playground.styles.js
      navbar/              ← already exists
      controls/            ← already exists
      panels/
        config-panel/      ← was config-panel.js
        payload-panel/     ← was payload-panel.js
        result-panel/      ← was result-panel.js + result-panel.styles.js
      user-consent-banner/ ← already exists
    wasm/                  ← NEW — extracted from playground.js
      bridge.js
      executor.js
      versions.js
    utils/                 ← already exists (base64, json-payload, escape-html)
    styles/                ← replaces styles.js — global CSS tokens/themes
  test/                    ← NEW — mirrors src/ exactly
    components/
      playground/
      controls/
      panels/
        config-panel/
        payload-panel/
        result-panel/
    wasm/
    utils/
```

**Part B — Extract WASM bridge**

Move all WASM lifecycle logic out of `playground.js` into `src/wasm/`:
- `bridge.js` — loads `wasm_exec.js`, fetches and instantiates the `.wasm` binary, fires `playground-wasm-ready`
- `executor.js` — wraps `execute(config, payloadType, payload, executorId, debug)` and `getExecutors()` globals
- `versions.js` — fetches and parses `wasm/versions.json`

No visible change to the user. Creates the clean seams needed for Vue and Vitest.

---

### Phase 2 — Vite + Vitest baseline

Switch build tool from Rollup to Vite (`vite.config.js` replaces `rollup.config.js`).

**Vitest config** — `node` environment globally, `happy-dom` opted in per-file:

```js
// vitest.config.js
import {defineConfig} from 'vitest/config'

export default defineConfig({
  test: {
    environment: 'node',       // pure logic tests run lean by default
    include: ['test/**/*.test.js'],
    globals: false,
  },
})
```

For test files that need a simulated browser (DOM, custom elements), add one line at the top of that specific file:

```js
// @vitest-environment happy-dom
```

**First tests** — no DOM needed, start here:
- `test/utils/base64.test.js`
- `test/utils/json-payload.test.js`
- `test/utils/escape-html.test.js`
- `test/wasm/versions.test.js` — mock `fetch`, verify JSON parsing
- `test/wasm/executor.test.js` — mock `execute()` and `getExecutors()` globals, verify error handling and debug flag behaviour

---

### Phase 3 — TypeScript + Vue 3 + Web Awesome

**Add TypeScript first**

All three reference repos use TypeScript. The Vue migration is the natural moment to add it — Vue SFCs work natively with `<script setup lang="ts">`. Add `tsconfig.json` and convert files incrementally as they are rewritten; do not try to convert everything at once.

**Verify Web Awesome package name before installing**

`@web-awesome/components` (the package recorded in this file) does not appear in either reference repo. The reference repos use:
- `@home-assistant/webawesome` — Home Assistant fork
- `@awesome.me/webawesome` — HOT (Humanitarian OpenStreetMap) fork

Confirm the correct official package on npm before running `npm install`. The install command below is a placeholder until verified:

```bash
npm install @web-awesome/components   # ⚠️ verify package name first
npm install vue @vitejs/plugin-vue
```

**Component migration** — each Lit component becomes a Vue SFC in its own folder:

```
src/components/controls/
  index.vue          ← the component (Vue SFC)
```

Test file lives in the mirrored test folder (established in Phase 1):

```
test/components/controls/
  index.test.js
```

**Web Awesome replacements** (confirmed against reference repos):

| Remove | Replace with | Confirmed in reference repos |
|--------|-------------|------------------------------|
| `split.js` + manual `Split()` | `<wa-split-panel>` | Yes |
| Native `<select>` (version, evaluator, examples) | `<wa-select>` / `<wa-option>` | Yes |
| Custom `.run-button` | `<wa-button>` | Yes |
| Custom SVG spinner | `<wa-spinner>` | Yes |
| Custom `.tooltip` CSS | `<wa-tooltip>` | Yes |
| User consent banner | `<wa-dialog>` | Yes |

Move all theming to CSS custom properties (`--wa-color-primary`, etc.) in a single `:root` block in `src/styles/`.

---

### Phase 4 — Component tests

Two distinct patterns apply depending on what is being tested.

**Web Components** (any remaining Lit elements, Web Awesome `wa-*` wrappers) — `document.createElement` pattern from `device-builder-frontend`:

```js
// @vitest-environment happy-dom
import {describe, it, expect, afterEach} from 'vitest'

afterEach(() => { document.body.innerHTML = '' })

it('renders the run button', async () => {
  const el = document.createElement('playground-controls')
  document.body.appendChild(el)
  await el.updateComplete          // wait for Lit to finish rendering
  expect(el.shadowRoot.querySelector('#btn-run')).not.toBeNull()
})
```

**Vue SFCs** — `@vue/test-utils` `mount()` pattern:

```js
// @vitest-environment happy-dom
import {mount} from '@vue/test-utils'
import Controls from '../../src/components/controls/index.vue'

it('emits evaluator-changed on select', async () => {
  const wrapper = mount(Controls, {props: {evaluator: 'transform_processor'}})
  await wrapper.find('wa-select').trigger('wa-change')
  expect(wrapper.emitted('evaluator-changed')).toBeTruthy()
})
```

**Mocking Web Awesome in tests** — some Web Awesome components (particularly `wa-dialog`) use browser features that `happy-dom` does not support. Mock them at the module level:

```js
vi.mock('@web-awesome/components/dist/components/dialog/dialog.js', () => ({}))
```

---

## Adding a New Executor

1. Create `internal/<name>executor.go` implementing the `Executor` interface.
2. Register it in the `Executors()` slice in `internal/executor.go`.
3. Add config/payload examples in `web/src/components/examples.js` under the executor's ID.
4. The WASM rebuild will automatically expose it via `statementsExecutors()`.

## Ongoing Review Responsibilities

At the start of each session, check the current state of the following and flag anything relevant to the migration or the project's health:

- **This repo** — `git log --oneline -10`, open PRs, and any untracked changes. Note if the migration branch (`emiliaFer/webawesomeAndVueMigration`) has diverged significantly from `main`.
- **Lit** — check the [Lit changelog](https://github.com/lit/lit/blob/main/packages/lit/CHANGELOG.md) for new releases since the currently installed version (`lit@3.x` in `web/package.json`). Flag breaking changes or deprecations that affect this codebase.
- **Web Awesome** — check [webawesome.com](https://webawesome.com/) or the [Web Awesome GitHub](https://github.com/shoelace-style/web-awesome) for new releases, component additions, or breaking changes. Note that the npm package name is still being established (previously `@shoelace-style/shoelace`).
- **Vue 3** — check the [Vue changelog](https://github.com/vuejs/core/blob/main/CHANGELOG.md) for new releases since the target version. Flag anything that affects the Composition API, Web Components integration (`isCustomElement`), or `@vitejs/plugin-vue` compatibility.
- **Vite / Vitest** — check for new releases of `vite` and `vitest`. Flag breaking changes to `vite.config.js` or the Vitest API that would affect the test setup.

- **Frontend design** — run the `nr-ui:ds-reviewer` skill against any modified or newly created frontend files (`web/src/**`) to audit design system compliance, catch component anti-patterns, and apply automated refactoring where applicable. This applies both during the regular session review and whenever frontend files are changed.

Surface findings as a short bullet list at the top of the session before proceeding with any task.

## Important Notes

- `internal/versions.go` is code-generated — never edit it manually.
- The WASM binary is large; always serve `web/public/` with gzip or brotli compression in production. The included Docker image uses Nginx with static brotli.
- Versions ≤ v0.110.0 of collector-contrib fail to compile due to breaking changes — they are intentionally excluded from version registration.
- The frontend uses the URL hash (base64-encoded JSON) for shareable state; this also enables embedding the playground in iframes.