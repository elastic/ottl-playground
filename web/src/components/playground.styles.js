// SPDX-License-Identifier: Apache-2.0

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
