// SPDX-License-Identifier: Apache-2.0

export const getJsonPayloadType = (payload) => {
  let json = JSON.parse(payload);
  if (json['resourceLogs']) {
    return 'logs';
  } else if (json['resourceSpans']) {
    return 'traces';
  } else if (json['resourceMetrics']) {
    return 'metrics';
  } else {
    throw new Error(
      'document must include an OTLP ["resourceLogs", "resourceSpans", "resourceMetrics"] root element'
    );
  }
};
