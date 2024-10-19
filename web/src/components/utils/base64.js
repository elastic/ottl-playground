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
