import {CONFIG_EXAMPLES, PAYLOAD_EXAMPLES} from "./examples";
import Split from 'split.js'
import {basicSetup, EditorView} from "codemirror"
import {json, jsonParseLinter} from "@codemirror/lang-json"
import {yaml} from "@codemirror/lang-yaml"
import {indentWithTab} from "@codemirror/commands"
import {keymap} from "@codemirror/view"
import {linter, lintGutter} from '@codemirror/lint'
import * as jsondiffpatch from 'jsondiffpatch';
import * as htmlFormatter from 'jsondiffpatch/formatters/html';
import * as annotatedFormatter from 'jsondiffpatch/formatters/annotated';

const runButton = document.querySelector('#btn-run');
const resultPanel = document.querySelector('#result-panel');
const showUnchangedInput = document.querySelector('#show-unchanged-input');
const dataTypeInput = document.querySelector('#data-type-input');
const diffDeltaSelect = document.querySelector('#diff-delta-select');
const loadExampleButton = document.querySelector('#load-example-button');
const exampleInput = document.querySelector('#example-input');
const evaluatorSelect = document.querySelector("#evaluator")

window.addEventListener('load', function () {
    document.querySelector('#loading').remove();
    document.querySelector('#playground').style.visibility = "";
})

Split(['#input-config-panel', '#input-sample-panel'], {
    direction: 'vertical',
})

Split(['#left-panel', '#right-panel'])

const go = new Go();
WebAssembly.instantiateStreaming(
    fetch("ottlplayground.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
        document.querySelector('#run-btn-ready').style.display = 'block';
        document.querySelector('#run-btn-loading').remove();
        runButton.disabled = false;
    }
);

const configInput = new EditorView({
    extensions: [
        basicSetup,
        keymap.of([indentWithTab]),
        EditorView.lineWrapping,
        yaml(),
    ],
    parent: document.querySelector('#config-input'),
})

const otlpDataInput = new EditorView({
    extensions: [
        basicSetup,
        keymap.of([indentWithTab]),
        linter(jsonParseLinter()),
        lintGutter(),
        EditorView.lineWrapping,
        json(),
    ],
    parent: document.querySelector('#otlp-data-input'),
})

const setEditorValue = (editor, value) => {
    editor.dispatch({
        changes: {from: 0, to: editor.state.doc.length, insert: value}
    });
}

const setResultText = (text, pre = false) => {
    let escaped = text
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");

    let value  = pre ? `<pre>${escaped}</pre>` : escaped
    resultPanel.innerHTML = `<div class="text">${value}</div>`
}

const loadOTLPExample = () => {
    let type = dataTypeInput.value;
    if (type === "") {
        return;
    }
    setEditorValue(otlpDataInput, JSON.stringify(JSON.parse(PAYLOAD_EXAMPLES[type]), null, 2));
}

const setResultError = (message) => {
    delete window.lastRunData;
    diffDeltaSelect.value = "visual";
    setResultText(`Error executing statements: ${message}`);
}

const runStatements = (event) => {
    let config = configInput.state.doc.toString()
    let type = dataTypeInput.value || "logs"
    let input = otlpDataInput.state.doc.toString()
    let evaluator = evaluatorSelect.value || "transform_processor"

    try {
        JSON.parse(input)
    } catch (e) {
        setResultError(`Invalid OTLP JSON payload: ${e}`);
        return
    }

    let result = executeStatements(config, type, input, evaluator)
    if (result?.hasOwnProperty("error")) {
        setResultError(result.error);
        console.error("Go error: ", result)
        return;
    }

    window.lastRunData = {
        input: input,
        result: result,
    }

    setResult(input, result);
}

const setLogResult = (result) => {
    resultPanel.innerHTML = "";
    let editor = new EditorView({
        extensions: [
            basicSetup,
            EditorView.editable.of(false),
            json(),
        ],
        parent: resultPanel,
    });

    setEditorValue(editor, result.logs);
}

const setJsonResult = (input, result) => {
    if (!result.value) {
        setResultError("Empty result value");
        return;
    }

    let left = JSON.parse(input);
    let right = JSON.parse(result.value);
    let selectedDelta = diffDeltaSelect.value;

    // Plain JSON result
    if (selectedDelta === "json") {
        resultPanel.innerHTML = "";
        let editor = new EditorView({
            extensions: [
                basicSetup,
                EditorView.editable.of(false),
                EditorView.lineWrapping,
                json(),
            ],
            parent: resultPanel,
        });
        setEditorValue(editor, JSON.stringify(right, null, 2));
        return
    }

    // Comparable JSON results
    const delta = jsondiffpatch.diff(left, right);
    if (!delta) {
        setResultText('No changes');
        return;
    }

    let formatter = selectedDelta === "annotated_json"
        ? annotatedFormatter
        : htmlFormatter;

    if (formatter.showUnchanged && formatter.hideUnchanged) {
        if (showUnchangedInput.checked) {
            formatter.showUnchanged();
        } else {
            formatter.hideUnchanged();
        }
    }

    resultPanel.innerHTML = formatter.format(
        delta,
        left,
    );
}

const setResult = (input, result) => {
    if (diffDeltaSelect.value === "logs") {
        setLogResult(result);
        return;
    }

    setJsonResult(input, result);
}

const loadConfigExamples = () => {
    let examples = CONFIG_EXAMPLES[evaluatorSelect.value]
    while (exampleInput.children.length > 1) {
        exampleInput.children[exampleInput.children.length - 1].remove();
    }

    for (let i = 0; i < examples.length; i++) {
        let example = examples[i];
        let opt = document.createElement('option');
        opt.value = i.toString();
        opt.innerHTML = example.name;
        exampleInput.appendChild(opt);
    }

    exampleInput.value = "";
}

const showUnchangedInputChanged = function (event) {
    if (event.target.checked) {
        htmlFormatter.showUnchanged();
    } else {
        htmlFormatter.hideUnchanged();
    }
}

const selectedDeltaChanged = (event) => {
    let delta = event.target.value;
    showUnchangedInput.disabled = delta !== "visual";
    showUnchangedInput.checked = false;
    let data = window.lastRunData;
    if (data) {
        setResult(data.input, data.result);
    }
}

const selectedExampleChanged = (event) => {
    let example = CONFIG_EXAMPLES[evaluatorSelect.value][parseInt(event.target.value)]
    if (example) {
        setEditorValue(configInput, example.config);
        dataTypeInput.value = example.otlp_type;
        loadOTLPExample()
    }
}

// Listeners
runButton.addEventListener('click', runStatements)
showUnchangedInput.addEventListener('change', showUnchangedInputChanged);
diffDeltaSelect.addEventListener('change', selectedDeltaChanged);
exampleInput.addEventListener('change', selectedExampleChanged);
loadExampleButton.addEventListener('click', () => loadOTLPExample());
evaluatorSelect.addEventListener('change', () => loadConfigExamples())

setEditorValue(otlpDataInput, "{}")

loadConfigExamples();