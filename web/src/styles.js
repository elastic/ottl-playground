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

export const globalStyles = css`
  input[type='text'],
  select,
  textarea {
    width: 100%;
    padding: 8px;
    border: 1px solid #ccc;
    border-radius: 4px;
    resize: vertical;
  }

  button:disabled,
  button[disabled] {
    opacity: 0.6;
    cursor: not-allowed;
  }

  label {
    padding: 10px 10px 10px 0;
    display: inline-block;
  }

  .hidden-overflow {
    overflow: hidden;
  }

  .h-full {
    height: 100%;
  }

  .w-full {
    width: 100%;
  }

  .full-size {
    height: 100%;
    width: 100%;
  }

  .tooltip {
    position: relative;
    display: inline-block;
  }

  .tooltip .tooltip-text {
    visibility: hidden;
    width: 100%;
    background-color: rgb(84, 84, 84);
    color: #fff;
    text-align: center;
    padding: 5px 0;
    border-radius: 6px;
    position: absolute;
    z-index: 1;
    font-size: small;
  }

  .tooltip:hover .tooltip-text {
    visibility: visible;
  }

  .tooltip-text-position-right {
    top: -5px;
    left: 105%;
  }

  .tooltip-text-position-left {
    top: -5px;
    right: 105%;
  }

  .tooltip-text-position-bottom {
    width: 120px;
    top: 70%;
    left: 50%;
    margin-left: -60px;
  }
`;
