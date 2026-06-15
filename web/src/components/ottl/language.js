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

import {parseMixed} from '@lezer/common';
import {yamlLanguage} from '@codemirror/lang-yaml';

import {
  LRLanguage,
  LanguageSupport,
  syntaxHighlighting,
  HighlightStyle,
} from '@codemirror/language';
import {tags as t} from '@lezer/highlight';
import {parser} from './ottl.grammar';
import {
  highlight,
  ottlStringTag,
  ottlFuncTag,
  ottlKeywordTag,
  ottlContextNameTag,
  ottlPathSegmentTag,
} from './highlight';

const ottlLanguage = LRLanguage.define({
  parser: parser.configure({props: [highlight]}),
  languageData: {},
});

const ottlHighlightStyle = HighlightStyle.define([
  {tag: t.meta, color: '#404740'},
  {tag: t.strong, fontWeight: 'bold'},
  {tag: t.keyword, color: '#708'},
  {
    tag: [t.atom, t.bool, t.url, t.contentSeparator, t.labelName],
    color: '#A52BCB',
  },
  {tag: t.definition(t.propertyName), color: '#00c'},
  {tag: [t.string], color: '#067D17FF'},
  {tag: [t.typeName], color: '#085'},
  {tag: [t.labelName], color: '#221199'},
  {tag: t.number, color: '#0550ae'},
  {tag: t.bool, color: '#A52BCB'},
  {tag: t.comment, color: '#8C8C8CFF', fontStyle: 'italic'},
  {tag: t.compareOperator, color: '#0550ae'},
  {tag: t.derefOperator, color: '#808080'},
  {tag: t.squareBracket, color: '#808080'},

  {tag: ottlKeywordTag, color: '#0550ae'},
  {tag: ottlStringTag, color: '#067D17FF'},
  {tag: ottlFuncTag, color: '#eb4a0f'},
  {tag: ottlContextNameTag, color: '#000'},
  {tag: ottlPathSegmentTag, color: '#59595b'},
]);

function buildSegmentsUpFrom(n, input) {
  const segments = [];
  let cur = n;
  while (cur) {
    if (cur.name === 'Pair') {
      const keyNode = cur.getChild('Key');
      if (keyNode) {
        const keyText = input.read(keyNode.from, keyNode.to).trim();
        segments.unshift(keyText);
      }
    } else if (
      cur.name === 'Sequence' ||
      cur.name === 'BlockSequence' ||
      cur.name === 'FlowSequence'
    ) {
      segments.unshift('*');
    }

    cur = cur.parent;
  }
  return segments;
}

function buildYamlPathFromPair(node, input) {
  if (!node) {
    return null;
  }

  const pushPairSegment = (pairNode) => {
    const keyNode = pairNode.getChild('Key');
    if (keyNode) {
      const keyText = input.read(keyNode.from, keyNode.to).trim();
      const rest = buildSegmentsUpFrom(pairNode.parent, input);
      return [...rest, keyText].filter(Boolean).join('.');
    }
  };

  // Literal or QuotedLiteral means the node is at the value of a Pair
  if (node.name === 'Literal' || node.name === 'QuotedLiteral') {
    if (node.parent?.name === 'Pair') {
      return pushPairSegment(node.parent);
    }
    return null;
  }

  // Pair: treat as single key path from root to this key, e.g. "foo.bar"
  if (node.name === 'Pair') {
    return pushPairSegment(node);
  }

  if (node.name === 'Item') {
    return buildSegmentsUpFrom(node.parent, input).filter(Boolean).join('.');
  }

  return null;
}

function buildPathsPatterns(configuredPatterns) {
  let result = [];
  for (let config of configuredPatterns) {
    let compiledPatterns = [];
    for (let pattern of config.patterns) {
      compiledPatterns.push(new RegExp(pattern));
    }
    result.push({
      grammar: config.grammar,
      patterns: compiledPatterns,
    });
  }
  return result;
}

function isSupportedYamlNode(node) {
  const isValidNodeType =
    node.name === 'Item' ||
    node.name === 'Literal' ||
    node.name === 'QuotedLiteral';

  if (!isValidNodeType) {
    return false;
  }

  const isValidNodeParent =
    node.node.parent.name === 'Pair' ||
    node.node.parent.name === 'BlockSequence' ||
    node.node.parent.name === 'Sequence' ||
    node.node.parent.name === 'FlowSequence';

  if (!isValidNodeParent) {
    return false;
  }

  const isLiteralValue =
    node.name === 'Literal' || node.name === 'QuotedLiteral';

  const isArrayLiteralValue =
    node.name === 'Item' && !node.node.getChild('BlockMapping');

  return isLiteralValue || isArrayLiteralValue;
}

function createMixedLanguage(configuredPatterns) {
  if (!configuredPatterns?.length) {
    return yamlLanguage;
  }

  const patterns = buildPathsPatterns(configuredPatterns);
  let seen = {};
  return yamlLanguage.configure({
    wrap: parseMixed((node, input) => {
      if (!isSupportedYamlNode(node)) {
        return null;
      }
      const nodeYamlPath = buildYamlPathFromPair(node.node, input);
      if (!nodeYamlPath) {
        return null;
      }
      let matchedPattern = seen[nodeYamlPath];
      if (!matchedPattern) {
        matchedPattern = patterns.find((p) =>
          p.patterns.some((v) => v.test(nodeYamlPath))
        );
        seen[nodeYamlPath] = matchedPattern;
      }

      return matchedPattern
        ? {parser: ottlLanguage.parser.configure({top: matchedPattern.grammar})}
        : null;
    }),
  });
}

export function yamlWithOTTL(configuredPatterns) {
  return [
    new LanguageSupport(
      createMixedLanguage(configuredPatterns),
      syntaxHighlighting(ottlHighlightStyle),
      yamlLanguage.support
    ),
  ];
}
