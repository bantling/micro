package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/go/iter"
	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

func TestLexString_(t *testing.T) {
	assert.Equal(t, util.Of2Error(token{tString, ``}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`""`))))
	assert.Equal(t, util.Of2Error(token{tString, `a`}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"a"`))))
	assert.Equal(t, util.Of2Error(token{tString, `a`}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"a"b`))))
	assert.Equal(t, util.Of2Error(token{tString, `abc`}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"abc"`))))
	assert.Equal(t, util.Of2Error(token{tString, `abc`}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"abc"b`))))
	assert.Equal(t, util.Of2Error(token{tString, `ab c`}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"ab c"b`))))

	assert.Equal(t, util.Of2Error(token{tString, "a\"\\/\b\f\n\r\tb"}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"a\"\\\/\b\f\n\r\tb"`))))

	assert.Equal(t, util.Of2Error(token{tString, `A`}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"\u0041"`))))
	assert.Equal(t, util.Of2Error(token{tString, `A`}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"\u0041"b`))))
	assert.Equal(t, util.Of2Error(token{tString, `abc`}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"a\u0062c"`))))
	assert.Equal(t, util.Of2Error(token{tString, "\U0001D11E"}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"\uD834\udd1e"`))))
	assert.Equal(t, util.Of2Error(token{tString, "a\U0001D11Eb"}, nil), util.Of2Error(lexString(iter.OfStringAsRunes(`"a\uD834\udd1eb"`))))

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("The ascii control character 0x05 is not valid in a string")), util.Of2Error(lexString(iter.OfStringAsRunes("\"\x05"))))

	for _, strs := range [][]string{
		{`"\u0`, `\u0`},
		{`"\u00`, `\u00`},
		{`"\u000`, `\u000`},
	} {
		assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("Incomplete string escape in %s", strs[1])), util.Of2Error(lexString(iter.OfStringAsRunes(strs[0]))))
	}

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("Illegal string escape \\uz")), util.Of2Error(lexString(iter.OfStringAsRunes(`"\uz`))))
	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("Incomplete string escape in \\")), util.Of2Error(lexString(iter.OfStringAsRunes(`"\`))))

	for _, strs := range [][]string{
		{`"\uD834`, `\uD834`},
		{`"\uD834z`, `\uD834`},
		{`"\uD834\`, `\uD834`},
		{`"\uD834\z`, `\uD834`},
	} {
		assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("The surrogate string escape %s must be followed by another surrogate escape to form valid UTF-16", strs[1])), util.Of2Error(lexString(iter.OfStringAsRunes(strs[0]))))
	}

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("The surrogate string escape \\uD834 cannot be followed by the non-surrogate escape \\u0061")), util.Of2Error(lexString(iter.OfStringAsRunes(`"\uD834\u0061`))))
	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("The surrogate string escape pair \\udd1e\\uD834 is not a valid UTF-16 surrogate pair")), util.Of2Error(lexString(iter.OfStringAsRunes(`"\udd1e\uD834"`))))
	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("Illegal string escape \\d")), util.Of2Error(lexString(iter.OfStringAsRunes(`"\d"`))))
	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("Incomplete string \": a string must be terminated by a \"")), util.Of2Error(lexString(iter.OfStringAsRunes(`"`))))

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
		assert.Equal(t, util.Of2Error(token{}, anErr), util.Of2Error(lexString(iter.SetError(iter.OfStringAsRunes(str), anErr))))
	}
}

func TestLexNumber_(t *testing.T) {
	assert.Equal(t, util.Of2Error(token{tNumber, "1"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "1"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1a"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-1"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-1"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-1"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-1a"))))

	assert.Equal(t, util.Of2Error(token{tNumber, "1.2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1.2"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "1.2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1.2a"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-1.2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-1.2"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-1.2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-1.2a"))))

	assert.Equal(t, util.Of2Error(token{tNumber, "1e2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1e2"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "1e2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1e2a"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-1e2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-1e2"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-1e2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-1e2a"))))

	assert.Equal(t, util.Of2Error(token{tNumber, "1e+2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1e+2"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "1e-2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1e-2a"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-1e+2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-1e+2"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-1e-2"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-1e-2a"))))

	assert.Equal(t, util.Of2Error(token{tNumber, "1.2e3"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1.2e3"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "1.2e3"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1.2e3a"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "1.2e+3"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1.2e+3"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "1.2e-3"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("1.2e-3a"))))

	assert.Equal(t, util.Of2Error(token{tNumber, "123"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("123"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-123"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-123a"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "123.456"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("123.456"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-123.456"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-123.456a"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "123.456e789"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("123.456e789"))))
	assert.Equal(t, util.Of2Error(token{tNumber, "-123.456e+789"}, nil), util.Of2Error(lexNumber(iter.OfStringAsRunes("-123.456e+789a"))))

	for _, strs := range [][]string{
		{"-", "-"},
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
		assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("Invalid number %s: a number must satisfy the regex -?[0-9]+([.][0-9]+)?([eE][0-9]+)?", strs[1])), util.Of2Error(lexNumber(iter.OfStringAsRunes(strs[0]))))
	}

	// Problem errors
	anErr := fmt.Errorf("An err")
	for _, str := range []string{
		"1e",
		"1",
		"1.",
	} {
		assert.Equal(t, util.Of2Error(token{}, anErr), util.Of2Error(lexNumber(iter.SetError(iter.OfStringAsRunes(str), anErr))))
	}
}

func TestLexBooleanNull_(t *testing.T) {
	assert.Equal(t, util.Of2Error(token{tBoolean, "true"}, nil), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("true"))))
	assert.Equal(t, util.Of2Error(token{tBoolean, "false"}, nil), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("false"))))
	assert.Equal(t, util.Of2Error(token{tNull, "null"}, nil), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("null"))))
	assert.Equal(t, util.Of2Error(token{tNull, "null"}, nil), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("null1"))))

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("Invalid sequence zippy: an array, object, string, number, boolean, or null was expected")), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("zippy"))))

	// Problem errors
	anErr := fmt.Errorf("An err")
	assert.Equal(t, util.Of2Error(token{}, anErr), util.Of2Error(lexBooleanNull(iter.SetError(iter.Of([]rune("t")...), anErr))))
}

func TestLex_(t *testing.T) {
	it := iter.OfStringAsRunes(`[]{},:"a"-1,1.25,1e2,1.25e2true,false,null`)

	assert.Equal(t, util.Of2Error(token{tOBracket, "["}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tCBracket, "]"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tOBrace, "{"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tCBrace, "}"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tComma, ","}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tColon, ":"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tString, "a"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tNumber, "-1"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tComma, ","}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tNumber, "1.25"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tComma, ","}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tNumber, "1e2"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tComma, ","}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tNumber, "1.25e2"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tBoolean, "true"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tComma, ","}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tBoolean, "false"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tComma, ","}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{tNull, "null"}, nil), util.Of2Error(lex(it)))
	assert.Equal(t, util.Of2Error(token{}, iter.EOI), util.Of2Error(lex(it)))

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf("Invalid character +: an array, object, string, number, boolean, or null was expected")), util.Of2Error(lex(iter.OfStringAsRunes("+"))))

	// Problem errors
	anErr := fmt.Errorf("An err")
	assert.Equal(t, util.Of2Error(token{}, anErr), util.Of2Error(lex(iter.SetError(iter.Of([]rune(" ")...), anErr))))
}

func TestLexer_(t *testing.T) {
	it := lexer(iter.OfStringAsRunes(`[`))
	assert.Equal(t, util.Of2Error(tokOBracket, nil), util.Of2Error(it.Next()))
	assert.Equal(t, util.Of2Error(token{}, iter.EOI), util.Of2Error(it.Next()))
}
