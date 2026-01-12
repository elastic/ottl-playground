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

package internal

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// CompletionContext describes the cursor context for autocomplete.
type CompletionContext struct {
	// InFunctionArgs is true if cursor is inside function arguments (unbalanced parens)
	InFunctionArgs bool `json:"inFunctionArgs"`
	// FunctionName is the name of the function we're inside (if InFunctionArgs)
	FunctionName string `json:"functionName"`
	// ArgIndex is the 0-based argument index (based on comma count)
	ArgIndex int `json:"argIndex"`
	// AfterDot is true if the last non-whitespace token was a dot
	AfterDot bool `json:"afterDot"`
	// AfterWhere is true if we're after a 'where' keyword with balanced parens
	AfterWhere bool `json:"afterWhere"`
	// AtStatementStart is true if we're at the beginning of a statement
	AtStatementStart bool `json:"atStatementStart"`
	// LastToken is the type of the last non-whitespace token
	LastToken string `json:"lastToken"`
	// ParenDepth is the current parenthesis nesting depth
	ParenDepth int `json:"parenDepth"`
}

// ottlLexer is the OTTL lexer definition, matching the collector's grammar.go
var ottlLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: `Bytes`, Pattern: `0x[a-fA-F0-9]+`},
	{Name: `Float`, Pattern: `[-+]?\d*\.\d+([eE][-+]?\d+)?`},
	{Name: `Int`, Pattern: `[-+]?\d+`},
	{Name: `String`, Pattern: `"(\\.|[^\\"])*"`},
	{Name: `OpNot`, Pattern: `\b(not)\b`},
	{Name: `OpOr`, Pattern: `\b(or)\b`},
	{Name: `OpAnd`, Pattern: `\b(and)\b`},
	{Name: `Where`, Pattern: `\b(where)\b`},
	{Name: `OpComparison`, Pattern: `==|!=|>=|<=|>|<`},
	{Name: `OpAddSub`, Pattern: `\+|\-`},
	{Name: `OpMultDiv`, Pattern: `\/|\*`},
	{Name: `Boolean`, Pattern: `\b(true|false)\b`},
	{Name: `Equal`, Pattern: `=`},
	{Name: `LParen`, Pattern: `\(`},
	{Name: `RParen`, Pattern: `\)`},
	{Name: `LBrace`, Pattern: `\{`},
	{Name: `RBrace`, Pattern: `\}`},
	{Name: `Colon`, Pattern: `\:`},
	{Name: `Comma`, Pattern: `,`},
	{Name: `Dot`, Pattern: `\.`},
	{Name: `LBracket`, Pattern: `\[`},
	{Name: `RBracket`, Pattern: `\]`},
	{Name: `Uppercase`, Pattern: `[A-Z][a-zA-Z0-9_]*`},
	{Name: `Lowercase`, Pattern: `[a-z][a-z0-9_]*`},
	{Name: `Whitespace`, Pattern: `\s+`},
})

// GetCompletionContext analyzes a partial OTTL statement and returns context
// information for autocomplete. Uses the same lexer as the OTTL parser to
// properly handle strings, comments, and other tokens.
func GetCompletionContext(statement string) map[string]any {
	ctx := analyzeCompletionContext(statement)
	return map[string]any{
		"inFunctionArgs":   ctx.InFunctionArgs,
		"functionName":     ctx.FunctionName,
		"argIndex":         ctx.ArgIndex,
		"afterDot":         ctx.AfterDot,
		"afterWhere":       ctx.AfterWhere,
		"atStatementStart": ctx.AtStatementStart,
		"lastToken":        ctx.LastToken,
		"parenDepth":       ctx.ParenDepth,
	}
}

func analyzeCompletionContext(statement string) CompletionContext {
	ctx := CompletionContext{
		AtStatementStart: true,
	}

	if statement == "" {
		return ctx
	}

	// Tokenize the statement
	lex, err := ottlLexer.LexString("", statement)
	if err != nil {
		// If lexer fails, return default context
		return ctx
	}

	// Track state while iterating through tokens
	var tokens []lexer.Token
	for {
		tok, err := lex.Next()
		if err != nil || tok.EOF() {
			break
		}
		tokens = append(tokens, tok)
	}

	if len(tokens) == 0 {
		return ctx
	}

	// Stack to track function calls for nested parens
	type funcCall struct {
		name     string
		argIndex int
	}
	var funcStack []funcCall
	var lastNonWS lexer.Token
	sawWhere := false
	parenDepthAtWhere := 0

	symbolNames := ottlLexer.Symbols()
	lparen := symbolNames["LParen"]
	rparen := symbolNames["RParen"]
	comma := symbolNames["Comma"]
	dot := symbolNames["Dot"]
	where := symbolNames["Where"]
	whitespace := symbolNames["Whitespace"]
	uppercase := symbolNames["Uppercase"]
	lowercase := symbolNames["Lowercase"]

	for _, tok := range tokens {
		// Skip whitespace for context tracking
		if tok.Type == whitespace {
			continue
		}

		ctx.AtStatementStart = false

		switch tok.Type {
		case lparen:
			// Opening paren - push to function stack
			funcName := ""
			if lastNonWS.Type == uppercase || lastNonWS.Type == lowercase {
				funcName = lastNonWS.Value
			}
			funcStack = append(funcStack, funcCall{name: funcName, argIndex: 0})
			ctx.ParenDepth++

		case rparen:
			// Closing paren - pop from function stack
			if len(funcStack) > 0 {
				funcStack = funcStack[:len(funcStack)-1]
			}
			if ctx.ParenDepth > 0 {
				ctx.ParenDepth--
			}

		case comma:
			// Comma increments arg index in current function
			if len(funcStack) > 0 {
				funcStack[len(funcStack)-1].argIndex++
			}

		case where:
			sawWhere = true
			parenDepthAtWhere = ctx.ParenDepth
		}

		lastNonWS = tok
	}

	// Determine final context state
	ctx.InFunctionArgs = len(funcStack) > 0
	if ctx.InFunctionArgs {
		top := funcStack[len(funcStack)-1]
		ctx.FunctionName = top.name
		ctx.ArgIndex = top.argIndex
	}

	// Check if last token was a dot
	if lastNonWS.Type == dot {
		ctx.AfterDot = true
	}

	// Check if we're after 'where' with balanced parens
	if sawWhere && ctx.ParenDepth == parenDepthAtWhere && !ctx.InFunctionArgs {
		ctx.AfterWhere = true
	}

	// Set last token name (reverse lookup from Symbols map)
	for name, rune := range symbolNames {
		if rune == lastNonWS.Type {
			ctx.LastToken = name
			break
		}
	}

	return ctx
}
