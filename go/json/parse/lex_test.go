package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/iter"
	"github.com/stretchr/testify/assert"
)

func mkIter(str string) *iter.Iter[rune] {
	return iter.OfReaderAsRunes(strings.NewReader(str))
}

func TestLexString(t *testing.T) {
	assert.Equal(t, token{tString, ``}, lexString(mkIter(`""`)))
	assert.Equal(t, token{tString, `a`}, lexString(mkIter(`"a"`)))
	assert.Equal(t, token{tString, `a`}, lexString(mkIter(`"a"b`)))
	assert.Equal(t, token{tString, `abc`}, lexString(mkIter(`"abc"`)))
	assert.Equal(t, token{tString, `abc`}, lexString(mkIter(`"abc"b`)))
	assert.Equal(t, token{tString, `ab c`}, lexString(mkIter(`"ab c"b`)))

	assert.Equal(t, token{tString, "a\"\\/\b\f\n\r\tb"}, lexString(mkIter(`"a\"\\\/\b\f\n\r\tb"`)))

	assert.Equal(t, token{tString, `A`}, lexString(mkIter(`"\u0041"`)))
	assert.Equal(t, token{tString, `A`}, lexString(mkIter(`"\u0041"b`)))
	assert.Equal(t, token{tString, `abc`}, lexString(mkIter(`"a\u0062c"`)))
	assert.Equal(t, token{tString, "\U0001D11E"}, lexString(mkIter(`"\uD834\udd1e"`)))
	assert.Equal(t, token{tString, "a\U0001D11Eb"}, lexString(mkIter(`"a\uD834\udd1eb"`)))

	funcs.TryTo(
		func() {
			lexString(mkIter("\"\x05"))
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
				lexString(mkIter(strs[0]))
				assert.Fail(t, "Must die")
			},
			func(e any) { assert.Equal(t, fmt.Errorf(errIllegalStringEscapeMsg, strs[1]), e) },
		)
	}

	funcs.TryTo(
		func() {
			lexString(mkIter(`"\`))
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
				lexString(mkIter(strs[0]))
				assert.Fail(t, "Must die")
			},
			func(e any) { assert.Equal(t, fmt.Errorf(errOneSurrogateEscapeMsg, strs[1]), e) },
		)
	}

	funcs.TryTo(
		func() {
			lexString(mkIter(`"\uD834\u0061`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errSurrogateNonSurrogateEscapeMsg, `\uD834`, `\u0061`), e) },
	)

	funcs.TryTo(
		func() {
			lexString(mkIter(`"\udd1e\uD834"`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errSurrogateDecodeEscapeMsg, `\udd1e\uD834`), e) },
	)

	funcs.TryTo(
		func() {
			lexString(mkIter(`"\d"`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errIllegalStringEscapeMsg, `\d`), e) },
	)

	funcs.TryTo(
		func() {
			lexString(mkIter(`"`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errIncompleteStringMsg, `"`), e) },
	)
}

func TestLexNumber(t *testing.T) {
	assert.Equal(t, token{tNumber, "1"}, lexNumber(mkIter("1")))
	assert.Equal(t, token{tNumber, "1"}, lexNumber(mkIter("1a")))
	assert.Equal(t, token{tNumber, "-1"}, lexNumber(mkIter("-1")))
	assert.Equal(t, token{tNumber, "-1"}, lexNumber(mkIter("-1a")))

	assert.Equal(t, token{tNumber, "1.2"}, lexNumber(mkIter("1.2")))
	assert.Equal(t, token{tNumber, "1.2"}, lexNumber(mkIter("1.2a")))
	assert.Equal(t, token{tNumber, "-1.2"}, lexNumber(mkIter("-1.2")))
	assert.Equal(t, token{tNumber, "-1.2"}, lexNumber(mkIter("-1.2a")))

	assert.Equal(t, token{tNumber, "1e2"}, lexNumber(mkIter("1e2")))
	assert.Equal(t, token{tNumber, "1e2"}, lexNumber(mkIter("1e2a")))
	assert.Equal(t, token{tNumber, "-1e2"}, lexNumber(mkIter("-1e2")))
	assert.Equal(t, token{tNumber, "-1e2"}, lexNumber(mkIter("-1e2a")))

	assert.Equal(t, token{tNumber, "1e+2"}, lexNumber(mkIter("1e+2")))
	assert.Equal(t, token{tNumber, "1e-2"}, lexNumber(mkIter("1e-2a")))
	assert.Equal(t, token{tNumber, "-1e+2"}, lexNumber(mkIter("-1e+2")))
	assert.Equal(t, token{tNumber, "-1e-2"}, lexNumber(mkIter("-1e-2a")))

	assert.Equal(t, token{tNumber, "1.2e3"}, lexNumber(mkIter("1.2e3")))
	assert.Equal(t, token{tNumber, "1.2e3"}, lexNumber(mkIter("1.2e3a")))
	assert.Equal(t, token{tNumber, "1.2e+3"}, lexNumber(mkIter("1.2e+3")))
	assert.Equal(t, token{tNumber, "1.2e-3"}, lexNumber(mkIter("1.2e-3a")))

	assert.Equal(t, token{tNumber, "123"}, lexNumber(mkIter("123")))
	assert.Equal(t, token{tNumber, "-123"}, lexNumber(mkIter("-123a")))
	assert.Equal(t, token{tNumber, "123.456"}, lexNumber(mkIter("123.456")))
	assert.Equal(t, token{tNumber, "-123.456"}, lexNumber(mkIter("-123.456a")))
	assert.Equal(t, token{tNumber, "123.456e789"}, lexNumber(mkIter("123.456e789")))
	assert.Equal(t, token{tNumber, "-123.456e+789"}, lexNumber(mkIter("-123.456e+789a")))

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
				lexNumber(mkIter(strs[0]))
				assert.Fail(t, "Must die")
			},
			func(e any) { assert.Equal(t, fmt.Errorf(errInvalidNumberMsg, strs[1]), e) },
		)
	}
}

func TestLexBooleanNull(t *testing.T) {
	assert.Equal(t, token{tBoolean, "true"}, lexBooleanNull(mkIter("true")))
	assert.Equal(t, token{tBoolean, "false"}, lexBooleanNull(mkIter("false")))
	assert.Equal(t, token{tNull, "null"}, lexBooleanNull(mkIter("null")))

	funcs.TryTo(
		func() {
			lexBooleanNull(mkIter("zippy"))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errInvalidBooleanNullMsg, "zippy"), e) },
	)
}

func TestLex(t *testing.T) {
	it := mkIter(`[]{},:"a"-1,1.25,1e2,1.25e2true,false,null`)
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
			lex(mkIter("+"))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errInvalidCharMsg, '+'), e) },
	)
}
