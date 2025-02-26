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
import {globalStyles} from '../styles';

const playgroundStyle = css`
  .playground {
    height: 100%;
    padding: 0 10px 0 10px;
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

  #loading {
    top: 0;
    padding-top: 25px;
    position: fixed;
    display: block;
    width: 100%;
    height: 100%;
    text-align: center;
    background-color: #fff;
    z-index: 99;
  }

  .beta-box {
    font-size: 10px !important;
    font-weight: 300 !important;
    color: gray;
    border: gray solid 1px;
    padding-left: 2px;
    padding-right: 2px;
  }
`;

export const playgroundStyles = [globalStyles, playgroundStyle];
