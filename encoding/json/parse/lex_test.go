package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/iter"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestLexString_(t *testing.T) {
	assert.Equal(t, union.OfResult(token{tString, ``}), union.OfResultError(lexString(iter.OfStringAsRunes(`""`))))
	assert.Equal(t, union.OfResult(token{tString, `a`}), union.OfResultError(lexString(iter.OfStringAsRunes(`"a"`))))
	assert.Equal(t, union.OfResult(token{tString, `a`}), union.OfResultError(lexString(iter.OfStringAsRunes(`"a"b`))))
	assert.Equal(t, union.OfResult(token{tString, `abc`}), union.OfResultError(lexString(iter.OfStringAsRunes(`"abc"`))))
	assert.Equal(t, union.OfResult(token{tString, `abc`}), union.OfResultError(lexString(iter.OfStringAsRunes(`"abc"b`))))
	assert.Equal(t, union.OfResult(token{tString, `ab c`}), union.OfResultError(lexString(iter.OfStringAsRunes(`"ab c"b`))))

	assert.Equal(t, union.OfResult(token{tString, "a\"\\/\b\f\n\r\tb"}), union.OfResultError(lexString(iter.OfStringAsRunes(`"a\"\\\/\b\f\n\r\tb"`))))

	assert.Equal(t, union.OfResult(token{tString, `A`}), union.OfResultError(lexString(iter.OfStringAsRunes(`"\u0041"`))))
	assert.Equal(t, union.OfResult(token{tString, `A`}), union.OfResultError(lexString(iter.OfStringAsRunes(`"\u0041"b`))))
	assert.Equal(t, union.OfResult(token{tString, `abc`}), union.OfResultError(lexString(iter.OfStringAsRunes(`"a\u0062c"`))))
	assert.Equal(t, union.OfResult(token{tString, "\U0001D11E"}), union.OfResultError(lexString(iter.OfStringAsRunes(`"\uD834\udd1e"`))))
	assert.Equal(t, union.OfResult(token{tString, "\U0001D11E"}), union.OfResultError(lexString(iter.OfStringAsRunes(`"\udd1e\uD834"`))))
	assert.Equal(t, union.OfResult(token{tString, "a\U0001D11Eb"}), union.OfResultError(lexString(iter.OfStringAsRunes(`"a\uD834\udd1eb"`))))
	assert.Equal(t, union.OfResult(token{tString, "a\U0001D11Eb"}), union.OfResultError(lexString(iter.OfStringAsRunes(`"a\udd1e\uD834b"`))))

	assert.Equal(t, union.OfError[token](fmt.Errorf("The ascii control character 0x05 is not valid in a string")), union.OfResultError(lexString(iter.OfStringAsRunes("\"\x05"))))

	for _, strs := range [][]string{
		{`"\u0`, `\u0`},
		{`"\u00`, `\u00`},
		{`"\u000`, `\u000`},
	} {
		assert.Equal(t, union.OfError[token](fmt.Errorf("Incomplete string escape in %s", strs[1])), union.OfResultError(lexString(iter.OfStringAsRunes(strs[0]))))
	}

	assert.Equal(t, union.OfError[token](fmt.Errorf("Illegal string escape \\uz")), union.OfResultError(lexString(iter.OfStringAsRunes(`"\uz`))))
	assert.Equal(t, union.OfError[token](fmt.Errorf("Incomplete string escape in \\")), union.OfResultError(lexString(iter.OfStringAsRunes(`"\`))))

	for _, strs := range [][]string{
		{`"\uD834`, `\uD834`},
		{`"\uD834z`, `\uD834`},
		{`"\uD834\`, `\uD834`},
		{`"\uD834\z`, `\uD834`},
	} {
		assert.Equal(t, union.OfError[token](fmt.Errorf("The surrogate string escape %s must be followed by another surrogate escape to form valid UTF-16", strs[1])), union.OfResultError(lexString(iter.OfStringAsRunes(strs[0]))))
	}

	assert.Equal(t, union.OfError[token](fmt.Errorf("The surrogate string escape \\uD834 cannot be followed by the non-surrogate escape \\u0061")), union.OfResultError(lexString(iter.OfStringAsRunes(`"\uD834\u0061`))))
	assert.Equal(t, union.OfError[token](fmt.Errorf("Illegal string escape \\d")), union.OfResultError(lexString(iter.OfStringAsRunes(`"\d"`))))
	assert.Equal(t, union.OfError[token](fmt.Errorf("Incomplete string \": a string must be terminated by a \"")), union.OfResultError(lexString(iter.OfStringAsRunes(`"`))))

	// Problem errors
	anErr := fmt.Errorf("An err")
	for _, str := range []string{
		`"\u`,
		`"`,
		`"\`,
		`"\uD834`,
		`"\uD834\`,
		`"\uD834\u`,
	} {
		assert.Equal(t, union.OfError[token](anErr), union.OfResultError(lexString(iter.SetError(iter.OfStringAsRunes(str), anErr))))
	}
}

func TestLexNumber_(t *testing.T) {
	assert.Equal(t, union.OfResult(token{tNumber, "0"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("0"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-0"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-0"))))
	assert.Equal(t, union.OfResult(token{tNumber, "1"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1"))))
	assert.Equal(t, union.OfResult(token{tNumber, "1"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1a"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-1"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-1"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-1"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-1a"))))

	assert.Equal(t, union.OfResult(token{tNumber, "1.2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1.2"))))
	assert.Equal(t, union.OfResult(token{tNumber, "1.2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1.2a"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-1.2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-1.2"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-1.2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-1.2a"))))

	assert.Equal(t, union.OfResult(token{tNumber, "1e2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1e2"))))
	assert.Equal(t, union.OfResult(token{tNumber, "1e2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1e2a"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-1e2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-1e2"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-1e2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-1e2a"))))

	assert.Equal(t, union.OfResult(token{tNumber, "1e+2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1e+2"))))
	assert.Equal(t, union.OfResult(token{tNumber, "1e-2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1e-2a"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-1e+2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-1e+2"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-1e-2"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-1e-2a"))))

	assert.Equal(t, union.OfResult(token{tNumber, "1.2e3"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1.2e3"))))
	assert.Equal(t, union.OfResult(token{tNumber, "1.2e3"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1.2e3a"))))
	assert.Equal(t, union.OfResult(token{tNumber, "1.2e+3"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1.2e+3"))))
	assert.Equal(t, union.OfResult(token{tNumber, "1.2e-3"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("1.2e-3a"))))

	assert.Equal(t, union.OfResult(token{tNumber, "123"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("123"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-123"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-123a"))))
	assert.Equal(t, union.OfResult(token{tNumber, "123.456"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("123.456"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-123.456"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-123.456a"))))
	assert.Equal(t, union.OfResult(token{tNumber, "123.456e789"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("123.456e789"))))
	assert.Equal(t, union.OfResult(token{tNumber, "-123.456e+789"}), union.OfResultError(lexNumber(iter.OfStringAsRunes("-123.456e+789a"))))

	for _, strs := range [][]string{
		{"-", "-"},
		{"01", "01"},
		{"-.", "-."},
		{"-e", "-e"},
		{"-a", "-"},
		{"1.", "1."},
		{"1.e", "1.e"},
		{"1.a", "1."},
		{"1e", "1e"},
		{"1ea", "1e"},
		{"1e+", "1e+"},
		{"1e-a", "1e-"},
		{"1.2e", "1.2e"},
		{"1.2ea", "1.2e"},
		{"1.2e+", "1.2e+"},
		{"1.2e-a", "1.2e-"},
	} {
		assert.Equal(t, union.OfError[token](fmt.Errorf("Invalid number %s: a number must satisfy the regex -?(0|[1-9][0-9]+)([.][0-9]+)?([eE][0-9]+)?", strs[1])), union.OfResultError(lexNumber(iter.OfStringAsRunes(strs[0]))))
	}

	// Problem errors
	anErr := fmt.Errorf("An err")
	for _, str := range []string{
		"1e",
		"1",
		"1.",
	} {
		assert.Equal(t, union.OfError[token](anErr), union.OfResultError(lexNumber(iter.SetError(iter.OfStringAsRunes(str), anErr))))
	}
}

func TestLexBooleanNull_(t *testing.T) {
	assert.Equal(t, union.OfResult(token{tBoolean, "true"}), union.OfResultError(lexBooleanNull(iter.OfStringAsRunes("true"))))
	assert.Equal(t, union.OfResult(token{tBoolean, "false"}), union.OfResultError(lexBooleanNull(iter.OfStringAsRunes("false"))))
	assert.Equal(t, union.OfResult(token{tNull, "null"}), union.OfResultError(lexBooleanNull(iter.OfStringAsRunes("null"))))
	assert.Equal(t, union.OfResult(token{tNull, "null"}), union.OfResultError(lexBooleanNull(iter.OfStringAsRunes("null1"))))

	assert.Equal(t, union.OfError[token](fmt.Errorf("Invalid sequence zippy: an array, object, string, number, boolean, or null was expected")), union.OfResultError(lexBooleanNull(iter.OfStringAsRunes("zippy"))))

	// Problem errors
	anErr := fmt.Errorf("An err")
	assert.Equal(t, union.OfError[token](anErr), union.OfResultError(lexBooleanNull(iter.SetError(iter.Of([]rune("t")...), anErr))))
}

func TestLex_(t *testing.T) {
	it := iter.OfStringAsRunes(`[]{},:"a"-1,1.25,1e2,1.25e2true,false,null`)

	assert.Equal(t, union.OfResult(token{tOBracket, "["}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tCBracket, "]"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tOBrace, "{"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tCBrace, "}"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tComma, ","}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tColon, ":"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tString, "a"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tNumber, "-1"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tComma, ","}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tNumber, "1.25"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tComma, ","}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tNumber, "1e2"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tComma, ","}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tNumber, "1.25e2"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tBoolean, "true"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tComma, ","}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tBoolean, "false"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tComma, ","}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfResult(token{tNull, "null"}), union.OfResultError(lex(it)))
	assert.Equal(t, union.OfError[token](iter.EOI), union.OfResultError(lex(it)))

	assert.Equal(t, union.OfError[token](fmt.Errorf("Invalid character +: an array, object, string, number, boolean, or null was expected")), union.OfResultError(lex(iter.OfStringAsRunes("+"))))

	// Problem errors
	anErr := fmt.Errorf("An err")
	assert.Equal(t, union.OfError[token](anErr), union.OfResultError(lex(iter.SetError(iter.Of([]rune(" ")...), anErr))))
}

func TestLexer_(t *testing.T) {
	it := lexer(iter.OfStringAsRunes(`[`))
	assert.Equal(t, union.OfResult(tokOBracket), union.OfResultError(it.Next()))
	assert.Equal(t, union.OfError[token](iter.EOI), union.OfResultError(it.Next()))
}
