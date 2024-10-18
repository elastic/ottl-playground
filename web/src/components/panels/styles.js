import {css} from 'lit-element';
import {globalStyles} from '../../styles';

const codePanelsStyle = css`
  .code-panel-parent {
    display: flex;
    flex-flow: column;
    height: 100%;
  }

  .code-panel-controls {
    flex-wrap: nowrap;
    border-left: #f5a800 4px solid;
    border-top: #eee 1px solid;
    border-bottom: #eee 1px solid;
    width: 99%;
    overflow-x: clip;
  }

  .code-panel-controls div {
    float: left;
    text-align: center;
    margin: 10px 5px 10px 5px;
    text-decoration: none;
    font-size: 17px;
  }

  .code-panel-controls div.right {
    float: right;
    overflow: hidden;
    align-items: center;

    div:not(:last-child) {
      margin-right: 4px;
    }
  }

  .code-panel-controls-header {
    width: 50%;
    overflow: hidden;
  }

  .code-editor-container {
    display: flex;
    height: 100%;
    overflow: auto;
  }

  .code-editor-container .wrapper {
    width: 100%;
  }

  .cm-editor {
    height: calc(100%);
  }

  .cm-scroller {
    overflow: auto;
  }
`;

export const codePanelsStyles = [codePanelsStyle, globalStyles];
