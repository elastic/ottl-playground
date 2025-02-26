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

export function utf8ToBase64(str) {
  const encoder = new TextEncoder();
  const utf8Array = encoder.encode(str);
  let binaryStr = '';
  utf8Array.forEach((byte) => {
    binaryStr += String.fromCharCode(byte);
  });
  return btoa(binaryStr);
}

export function base64ToUtf8(str) {
  const binaryStr = atob(str);
  const binaryLen = binaryStr.length;
  const utf8Array = new Uint8Array(binaryLen);
  for (let i = 0; i < binaryLen; i++) {
    utf8Array[i] = binaryStr.charCodeAt(i);
  }
  const decoder = new TextDecoder();
  return decoder.decode(utf8Array);
}
