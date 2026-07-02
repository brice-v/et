package lexer

import (
	"encoding/hex"
	"et/consts"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input     []rune
	pos       int  // current pos. in input (points to current char)
	readPos   int  // current reading pos. in input (after current char)
	ch        rune // current char under examination
	prevCh    rune // previous char read
	posInLine int

	TabCount int

	identsMap map[string]consts.HlStyleType
	operators string

	commentToken []rune
	stringTokens [][]rune
}

type TokenType int

const (
	TTEof TokenType = iota
	TTIllegal
	TTIdent
	TTNumber
	TTString
	TTComment
)

type Token struct {
	Type     TokenType
	Literal  string
	Position int

	HlStyleType consts.HlStyleType
}

// New returns a pointer to the lexer struct
func New(input []rune, idents map[string]consts.HlStyleType, operators string,
	comentToken string, stringTokens []string) *Lexer {
	st := make([][]rune, len(stringTokens))
	for i, s := range stringTokens {
		st[i] = []rune(s)
	}
	l := &Lexer{
		input:        input,
		identsMap:    idents,
		operators:    operators,
		commentToken: []rune(comentToken),
		stringTokens: st,
	}
	l.readChar()
	l.posInLine = 0
	return l
}

func (l *Lexer) LookupIdent(s string) (consts.HlStyleType, bool) {
	hlStyle, ok := l.identsMap[s]
	return hlStyle, ok
}

// readChar gives us the next character and advances out position
// in the input string
func (l *Lexer) readChar() {
	if l.ch == '\n' {
		l.posInLine = 0
	} else {
		l.posInLine++
	}
	l.prevCh = l.ch
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos += utf8.RuneCountInString(string(l.ch))
}

// peekChar will return the rune that is in the readPos without consuming any input
func (l *Lexer) peekChar() rune {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// peekSecondChar will return the rune right after the readPos+1 without consuming any input
func (l *Lexer) peekSecondChar() rune {
	if l.readPos >= len(l.input) || l.readPos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.readPos+1]
}

// readHexNumber will read the hex number as a helper
func (l *Lexer) readHexNumber(ogPos int) (TokenType, string) {
	// consume the 0 and x and continue to the number
	l.readChar()
	l.readChar()
	for isHexChar(l.ch) || (l.ch == '_' && isHexChar(l.peekChar())) {
		l.readChar()
	}
	return TTNumber, string(l.input[ogPos:l.pos])
}

// readOctalNumber will read the hex number as a helper
func (l *Lexer) readOctalNumber(ogPos int) (TokenType, string) {
	// consume the 0 and the o and continue to the number
	l.readChar()
	l.readChar()
	for isOctalChar(l.ch) || (l.ch == '_' && isOctalChar(l.peekChar())) {
		l.readChar()
	}
	return TTNumber, string(l.input[ogPos:l.pos])
}

func (l *Lexer) readBinaryNumber(ogPos int) (TokenType, string) {
	// consume the 0 and the b and continue to the number
	l.readChar()
	l.readChar()
	for isBinaryChar(l.ch) || (l.ch == '_' && isBinaryChar(l.peekChar())) {
		l.readChar()
	}
	return TTNumber, string(l.input[ogPos:l.pos])
}

// readNumber will keep consuming valid digits of the input according to `isDigit`
// and return the string
func (l *Lexer) readNumber() (TokenType, string) {
	position := l.pos
	if l.ch == '0' {
		if l.peekChar() == 'x' && isHexChar(l.peekSecondChar()) {
			return l.readHexNumber(position)
		} else if l.peekChar() == 'o' && isOctalChar(l.peekSecondChar()) {
			return l.readOctalNumber(position)
		} else if l.peekChar() == 'b' && isBinaryChar(l.peekSecondChar()) {
			return l.readBinaryNumber(position)
		} else if l.peekChar() == 'u' && isDigit(l.peekSecondChar()) {
			// Skip over 0 and u
			l.readChar()
			l.readChar()
			_, uis := l.readNumber()
			return TTNumber, uis
		}
	}
	dotFlag := false
	eFlag := false
	for isDigit(l.ch) || (l.ch == '_' && isDigit(l.peekChar())) {
		if l.peekChar() == '.' && !dotFlag && l.peekSecondChar() != '.' {
			dotFlag = true
			l.readChar()
			l.readChar()
		}
		if (l.peekChar() == 'e' || l.peekChar() == 'E') && (l.peekSecondChar() == '+' || l.peekSecondChar() == '-' || isDigit(l.peekSecondChar())) {
			eFlag = true
			// skip e
			l.readChar()
			if l.peekChar() == '+' || l.peekChar() == '-' {
				// skip + or -
				l.readChar()
			}
		}
		l.readChar()
	}
	isBigVal := l.ch == 'n' && (l.isWs(l.peekChar()) || !isLetter(l.peekChar()))
	tok := TTNumber

	if dotFlag || eFlag {
		if isBigVal {
			l.readChar()
		}
	} else {
		if isBigVal {
			l.readChar()
		}
	}
	return tok, string(l.input[position:l.pos])
}

// readIdentifier will keep consuming valid letters out of the input according to `isLetter`
// and return the string
func (l *Lexer) readIdentifier() (string, bool) {
	position := l.pos
	// Note: We can only do this because we check if the first char is a 'letter'
	// That includes underscores which is why 1 of the lexer tests changes to accomodate that
	for isLetter(l.ch) || unicode.IsNumber(l.ch) || l.ch == '?' || l.ch == '!' {
		l.readChar()
	}
	return string(l.input[position:l.pos]), l.ch == '('
}

// readSingleLineComment will continue to consume input until the EOL is reached
func (l *Lexer) readSingleLineComment() {
	for l.ch != 0 && l.ch != '\n' {
		l.readChar()
	}
}

// readString will consume tokens until the string is fully read
func (l *Lexer) readString() (string, error) {
	b := &strings.Builder{}

	stringStart := string(l.ch)
	strStart := l.ch
	for {
		l.readChar()

		// Support some basic escapes like \"
		if l.ch == '\\' {
			switch l.peekChar() {
			case strStart:
				b.WriteRune(strStart)
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case '\\':
				b.WriteByte('\\')
			case 'x':
				// Skip over the the '\\', 'x' and the next two bytes (hex)
				l.readChar()
				l.readChar()
				l.readChar()
				src := string([]rune{l.prevCh, l.ch})
				dst, err := hex.DecodeString(src)
				if err != nil {
					return "", err
				}
				b.Write(dst)
				continue
			}

			// Skip over the '\\' and the matched single escape char
			l.readChar()
			continue
		} else {
			if string(l.ch) == stringStart || l.ch == 0 {
				break
			}
		}

		b.WriteRune(l.ch)
	}

	if l.ch == 0 {
		return "", fmt.Errorf("string is not ended")
	}

	return b.String(), nil
}

func (l *Lexer) isWs(ch rune) bool {
	if ch == '\t' {
		l.TabCount++
	}
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// hasPrefix returns true if the input starting at the current position matches prefix
func (l *Lexer) hasPrefix(prefix []rune) bool {
	if len(prefix) == 0 {
		return false
	}
	if l.pos+len(prefix) > len(l.input) {
		return false
	}
	for i, r := range prefix {
		if l.input[l.pos+i] != r {
			return false
		}
	}
	return true
}

func (l *Lexer) isOperator(ch rune) bool {
	for _, c := range l.operators {
		if c == ch {
			return true
		}
	}
	return false
}

// skipWhitespace will continue to advance if the current byte is considered
// a whitespace character such as ' ', '\t', '\n', '\r'
func (l *Lexer) skipWhitespace() {
	for l.isWs(l.ch) {
		l.readChar()
	}
}

// NextToken matches against a byte and if it succeeds it will
// read the next char and return a token struct
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if len(l.commentToken) > 0 && l.hasPrefix(l.commentToken) {
		pos := l.pos
		tok := Token{Position: l.posInLine, HlStyleType: consts.HlCom}
		for range l.commentToken {
			l.readChar()
		}
		l.readSingleLineComment()
		endPos := min(l.readPos, len(l.input))
		tok.Literal = string(l.input[pos:endPos])
		tok.Type = TTComment
		return tok
	}

	tok := Token{Position: l.posInLine}

	var stringMatch bool
	for _, st := range l.stringTokens {
		if len(st) == 1 && l.ch == st[0] {
			stringMatch = true
			str, err := l.readString()
			if err != nil {
				tok = Token{Type: TTIllegal, Literal: "", Position: l.posInLine}
			} else {
				tok.Type = TTString
				tok.Literal = string(st) + str + string(st)
				tok.HlStyleType = consts.HlStr
			}
			break
		}
	}

	if !stringMatch {
		if l.ch == 0 {
			tok.Literal = ""
			tok.Type = TTEof
		} else if isLetter(l.ch) {
			lit, isCall := l.readIdentifier()
			tok.Literal = lit
			hlStyle, ok := l.LookupIdent(tok.Literal)
			tok.Type = TTIdent
			if ok {
				tok.HlStyleType = hlStyle
			} else if isCall {
				tok.HlStyleType = consts.Hl3
			} else {
				tok.Type = TTIllegal
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Type, tok.Literal = l.readNumber()
			tok.HlStyleType = consts.HlSpc
			return tok
		} else if l.isOperator(l.ch) {
			tok.Type, tok.Literal = TTIdent, string(l.ch)
			tok.HlStyleType = consts.Hl1
			l.readChar()
			return tok
		} else {
			tok = Token{Type: TTIllegal, Literal: "", Position: l.posInLine}
		}
	}

	l.readChar()
	return tok
}

// isLetter will return true if the rune given matches the pattern below
// 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
func isLetter(ch rune) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

// isDigit will return true if the rune give is 0-9 or any unicode Digit
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// isHexChar will return true if the rune given is a hex character
func isHexChar(ch rune) bool {
	return 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F' || '0' <= ch && ch <= '9'
}

// isOctalChar will return true if the rune given is an octal character
func isOctalChar(ch rune) bool {
	return '0' <= ch && ch <= '7'
}

// isBinaryChar will return true if the rune given is a binary character
func isBinaryChar(ch rune) bool {
	return ch == '0' || ch == '1'
}
