/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. See the NOTICE file distributed with
 * this work for additional information regarding copyright
 * ownership. Elasticsearch B.V. licenses this file to you under
 * the Apache License, Version 2.0 (the "License"); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

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
