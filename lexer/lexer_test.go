package lexer

import (
	"et/consts"
	"testing"
)

func newTestLexer(input string) *Lexer {
	return New([]rune(input), testIdents, testOperators, testCommentToken, testStringTokens)
}

var testIdents = map[string]consts.HlStyleType{
	"fn":    consts.Hl1,
	"let":   consts.Hl1,
	"int":   consts.Hl2,
	"rune":  consts.Hl2,
	"true":  consts.Hl3,
	"false": consts.Hl3,
}

const testOperators = "+-*/!|^&%=~{}[]:()"
const testCommentToken = "//"

var testStringTokens = []string{`"`, `'`}

type tokenTest struct {
	typ      TokenType
	literal  string
	hlStyle  consts.HlStyleType
	position int
}

func tokensEqual(t *testing.T, expected tokenTest, actual Token) bool {
	t.Helper()
	if expected.typ != actual.Type {
		t.Errorf("Type: expected %v, got %v", expected.typ, actual.Type)
		return false
	}
	if expected.literal != actual.Literal {
		t.Errorf("Literal: expected %q, got %q", expected.literal, actual.Literal)
		return false
	}
	if expected.hlStyle != actual.HlStyleType {
		t.Errorf("HlStyleType: expected %v, got %v", expected.hlStyle, actual.HlStyleType)
		return false
	}
	if actual.Type != TTEof && expected.position != 0 && expected.position != actual.Position {
		t.Errorf("Position: expected %d, got %d", expected.position, actual.Position)
		return false
	}
	return true
}

func TestNextTokenBasic(t *testing.T) {
	input := `fn main() {
	let x int = 42
	x + y * 5
}`
	tests := []tokenTest{
		{TTIdent, "fn", consts.Hl1, 0},
		{TTIdent, "main", consts.Hl3, 0},
		{TTIdent, "(", consts.Hl1, 0},
		{TTIdent, ")", consts.Hl1, 0},
		{TTIdent, "{", consts.Hl1, 0},
		{TTIdent, "let", consts.Hl1, 0},
		{TTIllegal, "x", consts.HlBase, 0},
		{TTIdent, "int", consts.Hl2, 0},
		{TTIdent, "=", consts.Hl1, 0},
		{TTNumber, "42", consts.HlSpc, 0},
		{TTIllegal, "x", consts.HlBase, 0},
		{TTIdent, "+", consts.Hl1, 0},
		{TTIllegal, "y", consts.HlBase, 0},
		{TTIdent, "*", consts.Hl1, 0},
		{TTNumber, "5", consts.HlSpc, 0},
		{TTIdent, "}", consts.Hl1, 0},
		{TTEof, "", 0, 0},
	}

	l := newTestLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if !tokensEqual(t, tt, tok) {
			t.Errorf("test[%d] failed", i)
		}
	}
}

func TestNextTokenStrings(t *testing.T) {
	tests := []struct {
		input   string
		typ     TokenType
		literal string
	}{
		{`"hello"`, TTString, `"hello"`},
		{`'world'`, TTString, `'world'`},
		{`"hello world"`, TTString, `"hello world"`},
		{`'hello world'`, TTString, `'hello world'`},
		{`"quote\"inside"`, TTString, `"quote"inside"`},
		{`"new\nline"`, TTString, "\"new\nline\""},
		{`"carriage\rreturn"`, TTString, "\"carriage\rreturn\""},
		{`"tab\tchar"`, TTString, "\"tab\tchar\""},
		{`"back\\slash"`, TTString, `"back\slash"`},
		{`"hex\x41"`, TTString, `"hexA"`},
		{`"unclosed`, TTIllegal, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := newTestLexer(tt.input)
			tok := l.NextToken()
			if tok.Type != tt.typ {
				t.Errorf("expected type %v, got %v (literal=%q)", tt.typ, tok.Type, tok.Literal)
			}
			if tok.Literal != tt.literal {
				t.Errorf("expected literal %q, got %q", tt.literal, tok.Literal)
			}
			if tok.Type == TTString && tok.HlStyleType != consts.HlStr {
				t.Errorf("expected HlStr, got %v", tok.HlStyleType)
			}
		})
	}
}

func TestNextTokenUnclosedString(t *testing.T) {
	l := newTestLexer(`"hello`)
	tok := l.NextToken()
	if tok.Type != TTIllegal {
		t.Errorf("expected TTIllegal for unclosed string, got %v", tok.Type)
	}
}

func TestNextTokenComments(t *testing.T) {
	tests := []struct {
		input  string
		tokens []tokenTest
	}{
		{"// hello", []tokenTest{
			{TTComment, "// hello", consts.HlCom, 0},
			{TTEof, "", 0, 0},
		}},
		{"// hello\nfn", []tokenTest{
			{TTComment, "// hello\n", consts.HlCom, 0},
			{TTIdent, "fn", consts.Hl1, 0},
		}},
		{"//", []tokenTest{
			{TTComment, "//", consts.HlCom, 0},
			{TTEof, "", 0, 0},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := newTestLexer(tt.input)
			for i, expected := range tt.tokens {
				tok := l.NextToken()
				if !tokensEqual(t, expected, tok) {
					t.Errorf("test[%d] failed", i)
				}
			}
		})
	}
}

func TestNextTokenEmptyInput(t *testing.T) {
	l := newTestLexer("")
	tok := l.NextToken()
	if tok.Type != TTEof {
		t.Errorf("expected TTEof, got %v", tok.Type)
	}
}

func TestNextTokenWhitespace(t *testing.T) {
	l := newTestLexer("   \t  \n  ")
	tok := l.NextToken()
	if tok.Type != TTEof {
		t.Errorf("expected TTEof, got %v", tok.Type)
	}
}

func TestNextTokenNumbers(t *testing.T) {
	tests := []struct {
		input   string
		literal string
	}{
		{"42", "42"},
		{"0", "0"},
		{"3.14", "3.14"},
		{"1e10", "1e10"},
		{"1.5e-3", "1.5e-3"},
		{"1_000", "1_000"},
		{"1_000_000", "1_000_000"},
		{"0xff", "0xff"},
		{"0xFF", "0xFF"},
		{"0xDEAD_BEEF", "0xDEAD_BEEF"},
		{"0o77", "0o77"},
		{"0o777", "0o777"},
		{"0b1010", "0b1010"},
		{"0b1010_0101", "0b1010_0101"},
		{"42n", "42n"},
		{"1_000n", "1_000n"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := newTestLexer(tt.input)
			tok := l.NextToken()
			if tok.Type != TTNumber {
				t.Errorf("expected TTNumber, got %v (literal=%q)", tok.Type, tok.Literal)
			}
			if tok.Literal != tt.literal {
				t.Errorf("expected literal %q, got %q", tt.literal, tok.Literal)
			}
			if tok.HlStyleType != consts.HlSpc {
				t.Errorf("expected HlSpc, got %v", tok.HlStyleType)
			}
		})
	}
}

func TestNextTokenUnicodeNumber(t *testing.T) {
	tests := []struct {
		input   string
		literal string
	}{
		{"0u1", "1"},
		{"0u12", "12"},
		{"0u0", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := newTestLexer(tt.input)
			tok := l.NextToken()
			if tok.Type != TTNumber {
				t.Errorf("expected TTNumber, got %v (literal=%q)", tok.Type, tok.Literal)
			}
			if tok.Literal != tt.literal {
				t.Errorf("expected literal %q, got %q", tt.literal, tok.Literal)
			}
		})
	}
}

func TestNextTokenBigValFloat(t *testing.T) {
	l := newTestLexer("3.14n")
	tok := l.NextToken()
	if tok.Type != TTNumber {
		t.Errorf("expected TTNumber, got %v", tok.Type)
	}
	if tok.Literal != "3.14n" {
		t.Errorf("expected literal %q, got %q", "3.14n", tok.Literal)
	}
}

func TestNextTokenIdentifiers(t *testing.T) {
	tests := []struct {
		input   string
		literal string
		hlStyle consts.HlStyleType
	}{
		{"foo", "foo", consts.HlBase},
		{"_foo", "_foo", consts.HlBase},
		{"foo_bar", "foo_bar", consts.HlBase},
		{"fn", "fn", consts.Hl1},
		{"int", "int", consts.Hl2},
		{"true", "true", consts.Hl3},
		{"x123", "x123", consts.HlBase},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := newTestLexer(tt.input)
			tok := l.NextToken()
			if tok.Literal != tt.literal {
				t.Errorf("expected literal %q, got %q", tt.literal, tok.Literal)
			}
			if tok.HlStyleType != tt.hlStyle {
				t.Errorf("expected hlStyle %v, got %v", tt.hlStyle, tok.HlStyleType)
			}
		})
	}
}

func TestNextTokenIdentWithCall(t *testing.T) {
	l := newTestLexer("foo(")
	tok := l.NextToken()
	if tok.Type != TTIdent {
		t.Errorf("expected TTIdent, got %v", tok.Type)
	}
	if tok.Literal != "foo" {
		t.Errorf("expected literal %q, got %q", "foo", tok.Literal)
	}
	if tok.HlStyleType != consts.Hl3 {
		t.Errorf("expected Hl3 for unknown call, got %v", tok.HlStyleType)
	}

	l2 := newTestLexer("int(")
	tok = l2.NextToken()
	if tok.Type != TTIdent {
		t.Errorf("expected TTIdent, got %v", tok.Type)
	}
	if tok.HlStyleType != consts.Hl2 {
		t.Errorf("expected Hl2 for known keyword with call, got %v", tok.HlStyleType)
	}
}

func TestNextTokenUnicodeIdent(t *testing.T) {
	l := newTestLexer("π")
	tok := l.NextToken()
	if tok.Literal != "π" {
		t.Errorf("expected literal %q, got %q", "π", tok.Literal)
	}
}

func TestNextTokenQuestion(t *testing.T) {
	l := newTestLexer("foo? x!")
	tok1 := l.NextToken()
	if tok1.Literal != "foo?" {
		t.Errorf("expected literal %q, got %q", "foo?", tok1.Literal)
	}
	tok2 := l.NextToken()
	if tok2.Literal != "x!" {
		t.Errorf("expected literal %q, got %q", "x!", tok2.Literal)
	}
}

func TestNextTokenOperators(t *testing.T) {
	ops := "+-*/!|^&%=~{}[]:()"
	for _, op := range ops {
		t.Run(string(op), func(t *testing.T) {
			l := newTestLexer(string(op))
			tok := l.NextToken()
			if tok.Type != TTIdent {
				t.Errorf("expected TTIdent for operator %q, got %v", op, tok.Type)
			}
			if tok.Literal != string(op) {
				t.Errorf("expected literal %q, got %q", string(op), tok.Literal)
			}
			if tok.HlStyleType != consts.Hl1 {
				t.Errorf("expected Hl1 for operator, got %v", tok.HlStyleType)
			}
		})
	}
}

func TestNextTokenIllegalChar(t *testing.T) {
	l := newTestLexer("@")
	tok := l.NextToken()
	if tok.Type != TTIllegal {
		t.Errorf("expected TTIllegal, got %v", tok.Type)
	}
}

func TestNextTokenMultipleTokens(t *testing.T) {
	l := newTestLexer(`"a" 'b' 123 // comment`)
	tests := []tokenTest{
		{TTString, `"a"`, consts.HlStr, 0},
		{TTString, `'b'`, consts.HlStr, 0},
		{TTNumber, "123", consts.HlSpc, 0},
		{TTComment, "// comment", consts.HlCom, 0},
	}

	for i, tt := range tests {
		tok := l.NextToken()
		if !tokensEqual(t, tt, tok) {
			t.Errorf("test[%d] failed", i)
		}
	}
}

func TestNextTokenMultipleLines(t *testing.T) {
	input := "fn foo()\n{\n\tx\n}\n"
	l := newTestLexer(input)

	expected := []struct {
		typ     TokenType
		literal string
	}{
		{TTIdent, "fn"},
		{TTIdent, "foo"},
		{TTIdent, "("},
		{TTIdent, ")"},
		{TTIdent, "{"},
		{TTIllegal, "x"},
		{TTIdent, "}"},
	}

	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp.typ {
			t.Errorf("test[%d]: expected type %v, got %v (literal=%q)", i, exp.typ, tok.Type, tok.Literal)
		}
		if tok.Literal != exp.literal {
			t.Errorf("test[%d]: expected literal %q, got %q", i, exp.literal, tok.Literal)
		}
	}
}

func TestNextTokenTabCount(t *testing.T) {
	l := newTestLexer("\tx")
	l.NextToken()
	if l.TabCount != 1 {
		t.Errorf("expected TabCount=1, got %d", l.TabCount)
	}
}

func TestNextTokenNumberWithLeadingDot(t *testing.T) {
	l := newTestLexer("0.5")
	tok := l.NextToken()
	if tok.Type != TTNumber {
		t.Errorf("expected TTNumber, got %v", tok.Type)
	}
	if tok.Literal != "0.5" {
		t.Errorf("expected literal %q, got %q", "0.5", tok.Literal)
	}
}

func TestNextTokenBadEscape(t *testing.T) {
	l := newTestLexer(`"\xZZ"`)
	tok := l.NextToken()
	if tok.Type != TTIllegal {
		t.Errorf("expected TTIllegal for bad hex escape, got %v", tok.Type)
	}
}

func TestNextTokenEmptyStringTokens(t *testing.T) {
	l := New([]rune("hello"), testIdents, testOperators, testCommentToken, nil)
	tok := l.NextToken()
	if tok.Type != TTIllegal {
		t.Errorf("expected TTIllegal, got %v", tok.Type)
	}
}

func TestNextTokenEmptyCommentToken(t *testing.T) {
	l := New([]rune("hello"), testIdents, testOperators, "", testStringTokens)
	tok := l.NextToken()
	if tok.Type != TTIllegal {
		t.Errorf("expected TTIllegal, got %v", tok.Type)
	}
}

func TestLookupIdent(t *testing.T) {
	l := newTestLexer("fn")
	hl, ok := l.LookupIdent("fn")
	if !ok {
		t.Errorf("expected ok=true for 'fn'")
	}
	if hl != consts.Hl1 {
		t.Errorf("expected Hl1 for 'fn', got %v", hl)
	}

	_, ok = l.LookupIdent("nonexistent")
	if ok {
		t.Errorf("expected ok=false for 'nonexistent'")
	}
}

func TestNewWithNilIdents(t *testing.T) {
	l := New([]rune("x"), nil, testOperators, testCommentToken, testStringTokens)
	tok := l.NextToken()
	if tok.Type != TTIllegal {
		t.Errorf("expected TTIllegal, got %v", tok.Type)
	}
}

func TestNextTokenReadStringError(t *testing.T) {
	l := newTestLexer(`"\x"`)
	tok := l.NextToken()
	if tok.Type != TTIllegal {
		t.Errorf("expected TTIllegal for malformed hex escape, got %v", tok.Type)
	}
}

func TestNextTokenEdgeCases(t *testing.T) {
	t.Run("newline after comment", func(t *testing.T) {
		l := newTestLexer("// line one\n// line two")
		tok1 := l.NextToken()
		if tok1.Type != TTComment || tok1.Literal != "// line one\n" {
			t.Errorf("expected comment %q, got %q", "// line one\\n", tok1.Literal)
		}
		tok2 := l.NextToken()
		if tok2.Type != TTComment || tok2.Literal != "// line two" {
			t.Errorf("expected comment %q, got %q", "// line two", tok2.Literal)
		}
	})

	t.Run("string at EOF", func(t *testing.T) {
		l := newTestLexer(`"hello`)
		tok := l.NextToken()
		if tok.Type != TTIllegal {
			t.Errorf("expected TTIllegal for unclosed string, got %v", tok.Type)
		}
	})

	t.Run("identifier with number", func(t *testing.T) {
		l := newTestLexer("foo123bar")
		tok := l.NextToken()
		if tok.Literal != "foo123bar" {
			t.Errorf("expected literal 'foo123bar', got %q", tok.Literal)
		}
	})

	t.Run("multiple operators", func(t *testing.T) {
		l := newTestLexer("+-*/")
		tok1 := l.NextToken()
		if tok1.Type != TTIdent || tok1.Literal != "+" {
			t.Errorf("expected '+', got %q", tok1.Literal)
		}
		tok2 := l.NextToken()
		if tok2.Type != TTIdent || tok2.Literal != "-" {
			t.Errorf("expected '-', got %q", tok2.Literal)
		}
	})
}

func TestNextTokenMultipleStringTokens(t *testing.T) {
	l := New([]rune(`"a" 'b' "c"`), testIdents, testOperators, testCommentToken, []string{`"`, `'`})
	tests := []struct {
		typ     TokenType
		literal string
	}{
		{TTString, `"a"`},
		{TTString, `'b'`},
		{TTString, `"c"`},
	}
	for _, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.typ || tok.Literal != tt.literal {
			t.Errorf("expected %v %q, got %v %q", tt.typ, tt.literal, tok.Type, tok.Literal)
		}
	}
}

func TestNextTokenPositionTracking(t *testing.T) {
	l := newTestLexer("a bb ccc")
	tests := []struct {
		literal  string
		position int
	}{
		{"a", 1},
		{"bb", 3},
		{"ccc", 6},
	}

	for _, tt := range tests {
		tok := l.NextToken()
		if tok.Literal != tt.literal {
			t.Errorf("expected literal %q, got %q", tt.literal, tok.Literal)
		}
		if tok.Position != tt.position {
			t.Errorf("expected position %d for %q, got %d", tt.position, tt.literal, tok.Position)
		}
	}
}
