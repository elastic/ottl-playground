import {css} from 'lit-element';
import {globalStyles} from '../styles';

const playgroundStyle = css`
  .playground {
    height: 100%;
    padding: 0 10px 10px 10px;
  }

  .split-horizontal {
    display: flex;
    flex-direction: row;
    height: calc(100% - 90px);
    border: #eeeeee 1px solid;
  }

  .split-vertical {
    height: 100%;
  }

  .gutter {
    background-color: #eee;
    background-repeat: no-repeat;
    background-position: 50%;
  }

  .gutter.gutter-horizontal {
    background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAeCAYAAADkftS9AAAAIklEQVQoU2M4c+bMfxAGAgYYmwGrIIiDjrELjpo5aiZeMwF+yNnOs5KSvgAAAABJRU5ErkJggg==');
    cursor: col-resize;
  }

  .gutter.gutter-vertical {
    background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAB4AAAAFAQMAAABo7865AAAABlBMVEVHcEzMzMzyAv2sAAAAAXRSTlMAQObYZgAAABBJREFUeF5jOAMEEAIEEFwAn3kMwcB6I2AAAAAASUVORK5CYII=');
    cursor: row-resize;
  }

  .delta-select {
    padding: 5px;
  }

  .result-panel-controls {
    overflow: hidden;
    flex-wrap: nowrap;
    border-left: darkgreen 4px solid;
    border-top: #eee 1px solid;
    border-bottom: #eee 1px solid;
    border-right: #eee 1px solid;
  }

  .result-panel-controls div {
    float: left;
    text-align: center;
    margin: 5px 5px 5px 5px;
    text-decoration: none;
    font-size: 17px;
  }

  .result-panel-controls div.right {
    float: right;
    overflow: hidden;

    div:not(:last-child) {
      margin-right: 4px;
    }
  }

  .result-panel-content {
    overflow: auto;
    height: calc(100% - 77px);
    border: #eee 1px solid;
    padding-top: 2px;
  }

  .result-panel-content .text {
    font-size: smaller;
    margin-left: 4px;
  }

  .result-panel-delta {
    display: grid;
    grid-template-columns: repeat(auto-fill, 135px);
    align-items: center;
    justify-content: start;
    border-left: gray 4px solid;
    padding: 5px;
    gap: 10px;
    font-size: 12px;
    overflow: auto;
  }

  #loading {
    width: 100%;
    height: 100%;
    padding-top: 25px;
    text-align: center;
  }
`;

const jsondiffpatchStyle = css`
  /*
    The MIT License
    
    Copyright (c) 2018 Benjamin Eidelman, https://twitter.com/beneidel
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
    THE SOFTWARE.
    
    https://github.com/benjamine/jsondiffpatch/tree/master/packages/jsondiffpatch/src/formatters/styles
    */

  .jsondiffpatch-delta {
    font-family: 'Bitstream Vera Sans Mono', 'DejaVu Sans Mono', Monaco, Courier,
      monospace;
    font-size: 12px;
    margin: 0;
    padding: 0 0 0 12px;
    display: inline-block;
  }

  .jsondiffpatch-delta pre {
    font-family: 'Bitstream Vera Sans Mono', 'DejaVu Sans Mono', Monaco, Courier,
      monospace;
    font-size: 12px;
    margin: 0;
    padding: 0;
    display: inline-block;
  }

  ul.jsondiffpatch-delta {
    list-style-type: none;
    padding: 0 0 0 20px;
    margin: 0;
  }

  .jsondiffpatch-delta ul {
    list-style-type: none;
    padding: 0 0 0 20px;
    margin: 0;
  }

  .jsondiffpatch-added .jsondiffpatch-property-name,
  .jsondiffpatch-added .jsondiffpatch-value pre,
  .jsondiffpatch-modified .jsondiffpatch-right-value pre,
  .jsondiffpatch-textdiff-added {
    background: #bbffbb;
  }

  .jsondiffpatch-deleted .jsondiffpatch-property-name,
  .jsondiffpatch-deleted pre,
  .jsondiffpatch-modified .jsondiffpatch-left-value pre,
  .jsondiffpatch-textdiff-deleted {
    background: #ffbbbb;
    text-decoration: line-through;
  }

  .jsondiffpatch-unchanged,
  .jsondiffpatch-movedestination {
    color: gray;
  }

  .jsondiffpatch-unchanged,
  .jsondiffpatch-movedestination > .jsondiffpatch-value {
    transition: all 0.5s;
    -webkit-transition: all 0.5s;
    overflow-y: hidden;
  }

  .jsondiffpatch-unchanged-showing .jsondiffpatch-unchanged,
  .jsondiffpatch-unchanged-showing
    .jsondiffpatch-movedestination
    > .jsondiffpatch-value {
    max-height: 100px;
  }

  .jsondiffpatch-unchanged-hidden .jsondiffpatch-unchanged,
  .jsondiffpatch-unchanged-hidden
    .jsondiffpatch-movedestination
    > .jsondiffpatch-value {
    max-height: 0;
  }

  .jsondiffpatch-unchanged-hiding
    .jsondiffpatch-movedestination
    > .jsondiffpatch-value,
  .jsondiffpatch-unchanged-hidden
    .jsondiffpatch-movedestination
    > .jsondiffpatch-value {
    display: block;
  }

  .jsondiffpatch-unchanged-visible .jsondiffpatch-unchanged,
  .jsondiffpatch-unchanged-visible
    .jsondiffpatch-movedestination
    > .jsondiffpatch-value {
    max-height: 100px;
  }

  .jsondiffpatch-unchanged-hiding .jsondiffpatch-unchanged,
  .jsondiffpatch-unchanged-hiding
    .jsondiffpatch-movedestination
    > .jsondiffpatch-value {
    max-height: 0;
  }

  .jsondiffpatch-unchanged-showing .jsondiffpatch-arrow,
  .jsondiffpatch-unchanged-hiding .jsondiffpatch-arrow {
    display: none;
  }

  .jsondiffpatch-value {
    display: inline-block;
  }

  .jsondiffpatch-property-name {
    display: inline-block;
    padding-right: 5px;
    vertical-align: top;
  }

  .jsondiffpatch-property-name:after {
    content: ': ';
  }

  .jsondiffpatch-child-node-type-array > .jsondiffpatch-property-name:after {
    content: ': [';
  }

  .jsondiffpatch-child-node-type-array:after {
    content: '],';
  }

  div.jsondiffpatch-child-node-type-array:before {
    content: '[';
  }

  div.jsondiffpatch-child-node-type-array:after {
    content: ']';
  }

  .jsondiffpatch-child-node-type-object > .jsondiffpatch-property-name:after {
    content: ': {';
  }

  .jsondiffpatch-child-node-type-object:after {
    content: '},';
  }

  div.jsondiffpatch-child-node-type-object:before {
    content: '{';
  }

  div.jsondiffpatch-child-node-type-object:after {
    content: '}';
  }

  .jsondiffpatch-value pre:after {
    content: ',';
  }

  li:last-child > .jsondiffpatch-value pre:after,
  .jsondiffpatch-modified > .jsondiffpatch-left-value pre:after {
    content: '';
  }

  .jsondiffpatch-modified .jsondiffpatch-value {
    display: inline-block;
  }

  .jsondiffpatch-modified .jsondiffpatch-right-value {
    margin-left: 5px;
  }

  .jsondiffpatch-moved .jsondiffpatch-value {
    display: none;
  }

  .jsondiffpatch-moved .jsondiffpatch-moved-destination {
    display: inline-block;
    background: #ffffbb;
    color: #888;
  }

  .jsondiffpatch-moved .jsondiffpatch-moved-destination:before {
    content: ' => ';
  }

  ul.jsondiffpatch-textdiff {
    padding: 0;
  }

  .jsondiffpatch-textdiff-location {
    color: #bbb;
    display: inline-block;
    min-width: 60px;
  }

  .jsondiffpatch-textdiff-line {
    display: inline-block;
  }

  .jsondiffpatch-textdiff-line-number:after {
    content: ',';
  }

  .jsondiffpatch-error {
    background: red;
    color: white;
    font-weight: bold;
  }

  /* */

  .jsondiffpatch-annotated-delta {
    font-family: 'Bitstream Vera Sans Mono', 'DejaVu Sans Mono', Monaco, Courier,
      monospace;
    font-size: 12px;
    margin: 0;
    padding: 0 0 0 12px;
    display: inline-block;
  }

  .jsondiffpatch-annotated-delta pre {
    font-family: 'Bitstream Vera Sans Mono', 'DejaVu Sans Mono', Monaco, Courier,
      monospace;
    font-size: 12px;
    margin: 0;
    padding: 0;
    display: inline-block;
  }

  .jsondiffpatch-annotated-delta td {
    margin: 0;
    padding: 0;
  }

  .jsondiffpatch-annotated-delta td pre:hover {
    font-weight: bold;
  }

  td.jsondiffpatch-delta-note {
    font-style: italic;
    padding-left: 10px;
  }

  .jsondiffpatch-delta-note > div {
    margin: 0;
    padding: 0;
  }

  .jsondiffpatch-delta-note pre {
    font-style: normal;
  }

  .jsondiffpatch-annotated-delta .jsondiffpatch-delta-note {
    color: #777;
  }

  .jsondiffpatch-annotated-delta tr:hover {
    background: #ffc;
  }

  .jsondiffpatch-annotated-delta tr:hover > td.jsondiffpatch-delta-note {
    color: black;
  }

  .jsondiffpatch-error {
    background: red;
    color: white;
    font-weight: bold;
  }
`;

export const playgroundStyles = [
  globalStyles,
  playgroundStyle,
  jsondiffpatchStyle,
];
