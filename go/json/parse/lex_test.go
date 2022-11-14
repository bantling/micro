package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/go/iter"
	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

func TestLexString(t *testing.T) {
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

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errControlCharInStringMsg, 0x05)), util.Of2Error(lexString(iter.OfStringAsRunes("\"\x05"))))

	for _, strs := range [][]string{
		{`"\u0`, `\u0`},
		{`"\u00`, `\u00`},
		{`"\u000`, `\u000`},
	} {
		assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errIncompleteStringEscapeMsg, strs[1])), util.Of2Error(lexString(iter.OfStringAsRunes(strs[0]))))
	}

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errIllegalStringEscapeMsg, `\uz`)), util.Of2Error(lexString(iter.OfStringAsRunes(`"\uz`))))
	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errIncompleteStringEscapeMsg, `\`)), util.Of2Error(lexString(iter.OfStringAsRunes(`"\`))))

	for _, strs := range [][]string{
		{`"\uD834`, `\uD834`},
		{`"\uD834z`, `\uD834`},
		{`"\uD834\`, `\uD834`},
		{`"\uD834\z`, `\uD834`},
	} {
		assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errOneSurrogateEscapeMsg, strs[1])), util.Of2Error(lexString(iter.OfStringAsRunes(strs[0]))))
	}

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errSurrogateNonSurrogateEscapeMsg, `\uD834`, `\u0061`)), util.Of2Error(lexString(iter.OfStringAsRunes(`"\uD834\u0061`))))
	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errSurrogateDecodeEscapeMsg, `\udd1e\uD834`)), util.Of2Error(lexString(iter.OfStringAsRunes(`"\udd1e\uD834"`))))
	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errIllegalStringEscapeMsg, `\d`)), util.Of2Error(lexString(iter.OfStringAsRunes(`"\d"`))))
	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errIncompleteStringMsg, "")), util.Of2Error(lexString(iter.OfStringAsRunes(`"`))))

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

func TestLexNumber(t *testing.T) {
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
		assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errInvalidNumberMsg, strs[1])), util.Of2Error(lexNumber(iter.OfStringAsRunes(strs[0]))))
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

func TestLexBooleanNull(t *testing.T) {
	assert.Equal(t, util.Of2Error(token{tBoolean, "true"}, nil), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("true"))))
	assert.Equal(t, util.Of2Error(token{tBoolean, "false"}, nil), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("false"))))
	assert.Equal(t, util.Of2Error(token{tNull, "null"}, nil), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("null"))))
	assert.Equal(t, util.Of2Error(token{tNull, "null"}, nil), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("null1"))))

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errInvalidBooleanNullMsg, "zippy")), util.Of2Error(lexBooleanNull(iter.OfStringAsRunes("zippy"))))

	// Problem errors
	anErr := fmt.Errorf("An err")
	assert.Equal(t, util.Of2Error(token{}, anErr), util.Of2Error(lexBooleanNull(iter.SetError(iter.Of([]rune("t")...), anErr))))
}

func TestLex(t *testing.T) {
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

	assert.Equal(t, util.Of2Error(token{}, fmt.Errorf(errInvalidCharMsg, '+')), util.Of2Error(lex(iter.OfStringAsRunes("+"))))

	// Problem errors
	anErr := fmt.Errorf("An err")
	assert.Equal(t, util.Of2Error(token{}, anErr), util.Of2Error(lex(iter.SetError(iter.Of([]rune(" ")...), anErr))))
}

func TestLexer(t *testing.T) {
	it := lexer(iter.OfStringAsRunes(`[`))
	assert.Equal(t, util.Of2Error(tokOBracket, nil), util.Of2Error(it.Next()))
	assert.Equal(t, util.Of2Error(token{}, iter.EOI), util.Of2Error(it.Next()))
}
