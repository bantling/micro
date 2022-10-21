package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/iter"
	"github.com/stretchr/testify/assert"
)

func TestLexString(t *testing.T) {
	assert.Equal(t, token{tString, ``}, lexString(iter.OfStringAsRunes(`""`)))
	assert.Equal(t, token{tString, `a`}, lexString(iter.OfStringAsRunes(`"a"`)))
	assert.Equal(t, token{tString, `a`}, lexString(iter.OfStringAsRunes(`"a"b`)))
	assert.Equal(t, token{tString, `abc`}, lexString(iter.OfStringAsRunes(`"abc"`)))
	assert.Equal(t, token{tString, `abc`}, lexString(iter.OfStringAsRunes(`"abc"b`)))
	assert.Equal(t, token{tString, `ab c`}, lexString(iter.OfStringAsRunes(`"ab c"b`)))

	assert.Equal(t, token{tString, "a\"\\/\b\f\n\r\tb"}, lexString(iter.OfStringAsRunes(`"a\"\\\/\b\f\n\r\tb"`)))

	assert.Equal(t, token{tString, `A`}, lexString(iter.OfStringAsRunes(`"\u0041"`)))
	assert.Equal(t, token{tString, `A`}, lexString(iter.OfStringAsRunes(`"\u0041"b`)))
	assert.Equal(t, token{tString, `abc`}, lexString(iter.OfStringAsRunes(`"a\u0062c"`)))
	assert.Equal(t, token{tString, "\U0001D11E"}, lexString(iter.OfStringAsRunes(`"\uD834\udd1e"`)))
	assert.Equal(t, token{tString, "a\U0001D11Eb"}, lexString(iter.OfStringAsRunes(`"a\uD834\udd1eb"`)))

	funcs.TryTo(
		func() {
			lexString(iter.OfStringAsRunes("\"\x05"))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errControlCharInStringMsg, 0x05), e) },
	)

	for _, strs := range [][]string{
		{`"\uz"`, `\uz`},
		{`"\u0`, `\u0`},
		{`"\u00`, `\u00`},
		{`"\u000`, `\u000`},
	} {
		funcs.TryTo(
			func() {
				lexString(iter.OfStringAsRunes(strs[0]))
				assert.Fail(t, "Must die")
			},
			func(e any) { assert.Equal(t, fmt.Errorf(errIllegalStringEscapeMsg, strs[1]), e) },
		)
	}

	funcs.TryTo(
		func() {
			lexString(iter.OfStringAsRunes(`"\`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errIncompleteStringEscapeMsg, `\`), e) },
	)

	for _, strs := range [][]string{
		{`"\uD834`, `\uD834`},
		{`"\uD834z`, `\uD834`},
		{`"\uD834\`, `\uD834`},
		{`"\uD834\z`, `\uD834`},
	} {
		funcs.TryTo(
			func() {
				lexString(iter.OfStringAsRunes(strs[0]))
				assert.Fail(t, "Must die")
			},
			func(e any) { assert.Equal(t, fmt.Errorf(errOneSurrogateEscapeMsg, strs[1]), e) },
		)
	}

	funcs.TryTo(
		func() {
			lexString(iter.OfStringAsRunes(`"\uD834\u0061`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errSurrogateNonSurrogateEscapeMsg, `\uD834`, `\u0061`), e) },
	)

	funcs.TryTo(
		func() {
			lexString(iter.OfStringAsRunes(`"\udd1e\uD834"`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errSurrogateDecodeEscapeMsg, `\udd1e\uD834`), e) },
	)

	funcs.TryTo(
		func() {
			lexString(iter.OfStringAsRunes(`"\d"`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errIllegalStringEscapeMsg, `\d`), e) },
	)

	funcs.TryTo(
		func() {
			lexString(iter.OfStringAsRunes(`"`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errIncompleteStringMsg, `"`), e) },
	)
}

func TestLexNumber(t *testing.T) {
	assert.Equal(t, token{tNumber, "1"}, lexNumber(iter.OfStringAsRunes("1")))
	assert.Equal(t, token{tNumber, "1"}, lexNumber(iter.OfStringAsRunes("1a")))
	assert.Equal(t, token{tNumber, "-1"}, lexNumber(iter.OfStringAsRunes("-1")))
	assert.Equal(t, token{tNumber, "-1"}, lexNumber(iter.OfStringAsRunes("-1a")))

	assert.Equal(t, token{tNumber, "1.2"}, lexNumber(iter.OfStringAsRunes("1.2")))
	assert.Equal(t, token{tNumber, "1.2"}, lexNumber(iter.OfStringAsRunes("1.2a")))
	assert.Equal(t, token{tNumber, "-1.2"}, lexNumber(iter.OfStringAsRunes("-1.2")))
	assert.Equal(t, token{tNumber, "-1.2"}, lexNumber(iter.OfStringAsRunes("-1.2a")))

	assert.Equal(t, token{tNumber, "1e2"}, lexNumber(iter.OfStringAsRunes("1e2")))
	assert.Equal(t, token{tNumber, "1e2"}, lexNumber(iter.OfStringAsRunes("1e2a")))
	assert.Equal(t, token{tNumber, "-1e2"}, lexNumber(iter.OfStringAsRunes("-1e2")))
	assert.Equal(t, token{tNumber, "-1e2"}, lexNumber(iter.OfStringAsRunes("-1e2a")))

	assert.Equal(t, token{tNumber, "1e+2"}, lexNumber(iter.OfStringAsRunes("1e+2")))
	assert.Equal(t, token{tNumber, "1e-2"}, lexNumber(iter.OfStringAsRunes("1e-2a")))
	assert.Equal(t, token{tNumber, "-1e+2"}, lexNumber(iter.OfStringAsRunes("-1e+2")))
	assert.Equal(t, token{tNumber, "-1e-2"}, lexNumber(iter.OfStringAsRunes("-1e-2a")))

	assert.Equal(t, token{tNumber, "1.2e3"}, lexNumber(iter.OfStringAsRunes("1.2e3")))
	assert.Equal(t, token{tNumber, "1.2e3"}, lexNumber(iter.OfStringAsRunes("1.2e3a")))
	assert.Equal(t, token{tNumber, "1.2e+3"}, lexNumber(iter.OfStringAsRunes("1.2e+3")))
	assert.Equal(t, token{tNumber, "1.2e-3"}, lexNumber(iter.OfStringAsRunes("1.2e-3a")))

	assert.Equal(t, token{tNumber, "123"}, lexNumber(iter.OfStringAsRunes("123")))
	assert.Equal(t, token{tNumber, "-123"}, lexNumber(iter.OfStringAsRunes("-123a")))
	assert.Equal(t, token{tNumber, "123.456"}, lexNumber(iter.OfStringAsRunes("123.456")))
	assert.Equal(t, token{tNumber, "-123.456"}, lexNumber(iter.OfStringAsRunes("-123.456a")))
	assert.Equal(t, token{tNumber, "123.456e789"}, lexNumber(iter.OfStringAsRunes("123.456e789")))
	assert.Equal(t, token{tNumber, "-123.456e+789"}, lexNumber(iter.OfStringAsRunes("-123.456e+789a")))

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
		funcs.TryTo(
			func() {
				lexNumber(iter.OfStringAsRunes(strs[0]))
				assert.Fail(t, "Must die")
			},
			func(e any) { assert.Equal(t, fmt.Errorf(errInvalidNumberMsg, strs[1]), e) },
		)
	}
}

func TestLexBooleanNull(t *testing.T) {
	assert.Equal(t, token{tBoolean, "true"}, lexBooleanNull(iter.OfStringAsRunes("true")))
	assert.Equal(t, token{tBoolean, "false"}, lexBooleanNull(iter.OfStringAsRunes("false")))
	assert.Equal(t, token{tNull, "null"}, lexBooleanNull(iter.OfStringAsRunes("null")))

	funcs.TryTo(
		func() {
			lexBooleanNull(iter.OfStringAsRunes("zippy"))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errInvalidBooleanNullMsg, "zippy"), e) },
	)
}

func TestLex(t *testing.T) {
	it := iter.OfStringAsRunes(`[]{},:"a"-1,1.25,1e2,1.25e2true,false,null`)
	assert.Equal(t, token{tOBracket, "["}, lex(it))
	assert.Equal(t, token{tCBracket, "]"}, lex(it))
	assert.Equal(t, token{tOBrace, "{"}, lex(it))
	assert.Equal(t, token{tCBrace, "}"}, lex(it))
	assert.Equal(t, token{tComma, ","}, lex(it))
	assert.Equal(t, token{tColon, ":"}, lex(it))
	assert.Equal(t, token{tString, "a"}, lex(it))
	assert.Equal(t, token{tNumber, "-1"}, lex(it))
	assert.Equal(t, token{tComma, ","}, lex(it))
	assert.Equal(t, token{tNumber, "1.25"}, lex(it))
	assert.Equal(t, token{tComma, ","}, lex(it))
	assert.Equal(t, token{tNumber, "1e2"}, lex(it))
	assert.Equal(t, token{tComma, ","}, lex(it))
	assert.Equal(t, token{tNumber, "1.25e2"}, lex(it))
	assert.Equal(t, token{tBoolean, "true"}, lex(it))
	assert.Equal(t, token{tComma, ","}, lex(it))
	assert.Equal(t, token{tBoolean, "false"}, lex(it))
	assert.Equal(t, token{tComma, ","}, lex(it))
	assert.Equal(t, token{tNull, "null"}, lex(it))
	assert.Equal(t, tokEOF, lex(it))

	funcs.TryTo(
		func() {
			lex(iter.OfStringAsRunes("+"))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errInvalidCharMsg, '+'), e) },
	)
}

func TestLexer(t *testing.T) {
	it := lexer(iter.OfStringAsRunes(`[`))

	tok, haveIt := it.NextValue()
	assert.Equal(t, tokOBracket, tok)
	assert.True(t, haveIt)

	tok, haveIt = it.NextValue()
	assert.False(t, haveIt)
}
