import {html, LitElement} from 'lit-element';
import {codePanelsStyles} from './styles';
import {basicSetup, EditorView} from 'codemirror';
import {keymap} from '@codemirror/view';
import {indentWithTab} from '@codemirror/commands';
import {yaml} from '@codemirror/lang-yaml';
import {CONFIG_EXAMPLES} from '../examples';
import {nothing} from "lit";

export class PlaygroundConfigPanel extends LitElement {
    static properties = {
        // attributes
        evaluator: {type: String, attribute: true},
        hideExamples: {type: Boolean, attribute: 'hide-examples'},
        // states
        _editor: {type: Object, state: true},
    };

    constructor() {
        super();
        this.hideExamples = false;
    }

    static get styles() {
        return codePanelsStyles;
    }

    get text() {
        if (!this._editor) return '';
        return this._editor.state.doc.toString();
    }

    set text(value) {
        this._editor.dispatch({
            changes: {from: 0, to: this._editor.state.doc.length, insert: value},
        });
    }

    updated(changedProperties) {
        if (changedProperties.has('evaluator')) {
            this._loadConfigExamples();
        }
        super.updated(changedProperties);
    }

    firstUpdated(_changedProperties) {
        this._initCodeEditor();
        this._loadConfigExamples();
    }

    render() {
        return html`
            <div class="code-panel-parent" id="input-config-panel">
                <div class="code-panel-controls">
                    <div class="header">
            <span
            ><strong>Configuration</strong>
              <sup><small>YAML</small></sup></span>
                    </div>
                    <div class="right" style="display: flex">
                        ${this.hideExamples ? nothing : html`
                        <select
                                id="example-input"
                                @change="${this._selectedExampleChanged}"
                        >
                            <option disabled selected value="">Example</option>
                        </select>
                        `}
                        <slot name="custom-components">
                        </slot>
                    </div>
                </div>
                <div class="code-editor-container">
                    <div class="wrapper" id="config-input"></div>
                </div>
            </div>
        `;
    }

    _loadConfigExamples() {
        if (this.hideExamples) return;

        let exampleInput = this.shadowRoot.querySelector('#example-input');
        while (exampleInput.children.length > 1) {
            exampleInput.children[exampleInput.children.length - 1].remove();
        }

        let examples = CONFIG_EXAMPLES[this.evaluator];
        for (let i = 0; i < examples.length; i++) {
            let example = examples[i];
            let opt = document.createElement('option');
            opt.value = i.toString();
            opt.innerHTML = example.name;
            exampleInput.appendChild(opt);
        }

        exampleInput.value = '';
    }

    _selectedExampleChanged(event) {
        let idx = parseInt(event.target.value);
        let example = CONFIG_EXAMPLES[this.evaluator][idx];
        if (!example) return;

        this.text = example.config;
        this.dispatchEvent(
            new CustomEvent('playground-load-example', {
                detail: {
                    evaluator: this.evaluator,
                    example: example,
                },
                bubbles: true,
                composed: true,
                cancelable: true,
            })
        );
    }

    _initCodeEditor() {
        this._editor = new EditorView({
            extensions: [
                basicSetup,
                keymap.of([indentWithTab]),
                EditorView.lineWrapping,
                yaml(),
            ],
            parent: this.shadowRoot.querySelector('#config-input'),
        });
    }
}

customElements.define('playground-config-panel', PlaygroundConfigPanel);
