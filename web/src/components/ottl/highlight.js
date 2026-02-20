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

import {styleTags, tags as t} from '@lezer/highlight';

export const ottlStringTag = t.special(t.string);
export const ottlFuncTag = t.special(t.variableName);
export const ottlContextNameTag = t.special(t.propertyName);
export const ottlPathSegmentTag = t.special(t.attributeName);
export const ottlKeywordTag = t.special(t.keyword);

export const highlight = styleTags({
  'Where OpNot OpOr OpAnd': ottlKeywordTag,
  Nil: t.null,
  Boolean: t.bool,

  Bytes: t.number,
  Float: t.number,
  Int: t.number,
  String: ottlStringTag,

  EnumSymbol: t.atom,
  'EnumSymbol/Uppercase': t.atom,
  'Editor/Lowercase': ottlFuncTag,
  'Converter/Uppercase': ottlFuncTag,
  'Path/Lowercase': ottlContextNameTag,
  'PathTail/Lowercase': ottlPathSegmentTag,

  OpComparison: t.compareOperator,
  OpAddSub: t.arithmeticOperator,
  OpMultDiv: t.arithmeticOperator,
  Equal: t.definitionOperator,
  '.': t.derefOperator,
  ',': t.separator,
  ':': t.separator,
  '( )': t.paren,
  '[ ]': t.squareBracket,
  '{ }': t.brace,

  List: t.list,
  MapValue: t.brace,
  Statement: t.content,
  BooleanExpression: t.logicOperator,
  Comparison: t.compareOperator,
  MathExpression: t.arithmeticOperator,
});
