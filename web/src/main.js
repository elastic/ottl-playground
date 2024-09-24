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

const setResultText = (text) => {
    resultPanel.innerHTML = `<small>${text}</small>`
}

const loadOTLPExample = () => {
    let type = dataTypeInput.value;
    if (type === "") {
        return;
    }
    setEditorValue(otlpDataInput, JSON.stringify(JSON.parse(PAYLOAD_EXAMPLES[type]), null, 2));
}

const runStatements = (event) => {
    let config = configInput.state.doc.toString()
    let type = dataTypeInput.value || "logs"
    let input = otlpDataInput.state.doc.toString()

    try {
        JSON.parse(input)
    } catch (e) {
        setResultText(`Invalid OTLP JSON payload: ${e}`)
        return
    }

    let result = executeStatements(config, type, input,)
    if (result?.hasOwnProperty("error")) {
        console.error("Go error: ", result)
        setResultText(`Error executing statements: ${result.error}`);
        return;
    }

    window.lastRunData = {
        input: input,
        result: result,
    }

    performResultDiff(input, result);
}

const performResultDiff = (leftContent, rightContent) => {
    let left = JSON.parse(leftContent);
    let right = JSON.parse(rightContent);

    let selectedDelta = diffDeltaSelect.value;
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

    const delta = jsondiffpatch.diff(left, right);
    // No changes
    if (!delta) {
        setResultText('<i>&nbsp;No changes</i>');
        return;
    }

    if (selectedDelta === "annotated_json") {
        resultPanel.innerHTML = annotatedFormatter.format(
            delta,
            left,
        );
        return;
    }

    // Visual
    if (showUnchangedInput.checked) {
        htmlFormatter.showUnchanged();
    } else {
        htmlFormatter.hideUnchanged();
    }

    resultPanel.innerHTML = htmlFormatter.format(
        delta,
        left,
    );
}

const loadConfigExamples = () => {
    for (let i = 0; i < CONFIG_EXAMPLES.length; i++) {
        let example = CONFIG_EXAMPLES[i];
        let opt = document.createElement('option');
        opt.value = i.toString();
        opt.innerHTML = example.name;
        exampleInput.appendChild(opt);
    }
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
        performResultDiff(data.input, data.result);
    }
}

const selectedExampleChanged = (event) => {
    let example = CONFIG_EXAMPLES[parseInt(event.target.value)]
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
loadExampleButton.addEventListener('click', (e) => loadOTLPExample());

setEditorValue(otlpDataInput, "{}")
loadConfigExamples();