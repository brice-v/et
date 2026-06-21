package lexer

import (
	"encoding/hex"
	"et/consts"
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer is the struct that contains members for
// lexing needs
type Lexer struct {
	input        []rune
	position     int  // current pos. in input (points to current char)
	readPosition int  // current reading pos. in input (after current char)
	ch           rune // current char under examination
	prevCh       rune // previous char read
	posInLine    int

	TabCount int

	identsMap map[string]consts.HlStyleType
	operators string
}

type TokenType int

const (
	TTEof TokenType = iota
	TTIllegal
	TTIdent
	TTNumber
	TTString
)

type Token struct {
	Type     TokenType
	Literal  string
	Position int

	HlStyleType consts.HlStyleType
}

// New returns a pointer to the lexer struct
func New(input []rune, idents map[string]consts.HlStyleType, operators string) *Lexer {
	l := &Lexer{input: input, identsMap: idents, operators: operators}
	l.readChar()
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
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += utf8.RuneCountInString(string(l.ch))
}

// peekChar will return the rune that is in the readPosition without consuming any input
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// peekSecondChar will return the rune right after the readPosition+1 without consuming any input
func (l *Lexer) peekSecondChar() rune {
	if l.readPosition >= len(l.input) || l.readPosition+1 >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+1]
}

// readHexNumber will read the hex number as a helper
func (l *Lexer) readHexNumber(ogPos int) (TokenType, string) {
	// consume the 0 and x and continue to the number
	l.readChar()
	l.readChar()
	for isHexChar(l.ch) || (l.ch == '_' && isHexChar(l.peekChar())) {
		l.readChar()
	}
	return TTNumber, string(l.input[ogPos:l.position])
}

// readOctalNumber will read the hex number as a helper
func (l *Lexer) readOctalNumber(ogPos int) (TokenType, string) {
	// consume the 0 and the o and continue to the number
	l.readChar()
	l.readChar()
	for isOctalChar(l.ch) || (l.ch == '_' && isOctalChar(l.peekChar())) {
		l.readChar()
	}
	return TTNumber, string(l.input[ogPos:l.position])
}

func (l *Lexer) readBinaryNumber(ogPos int) (TokenType, string) {
	// consume the 0 and the b and continue to the number
	l.readChar()
	l.readChar()
	for isBinaryChar(l.ch) || (l.ch == '_' && isBinaryChar(l.peekChar())) {
		l.readChar()
	}
	return TTNumber, string(l.input[ogPos:l.position])
}

// readNumber will keep consuming valid digits of the input according to `isDigit`
// and return the string
func (l *Lexer) readNumber() (TokenType, string) {
	position := l.position
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
			tok, uis := l.readNumber()
			tok = TTNumber
			return tok, uis
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
	return tok, string(l.input[position:l.position])
}

// readIdentifier will keep consuming valid letters out of the input according to `isLetter`
// and return the string
func (l *Lexer) readIdentifier() string {
	position := l.position
	// Note: We can only do this because we check if the first char is a 'letter'
	// That includes underscores which is why 1 of the lexer tests changes to accomodate that
	for isLetter(l.ch) || unicode.IsNumber(l.ch) || l.ch == '?' || l.ch == '!' {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

// readMultiLineComment will continue to consume input until the end multiline token is reached
func (l *Lexer) readMultiLineComment() {
	for l.ch != 0 {
		if l.ch == 0 {
			break // break on EOF
		}
		if l.ch == '#' && l.peekChar() == '#' && l.peekSecondChar() == '#' {
			l.readChar()
			l.readChar()
			l.readChar()
			// fmt.Println(l.ch)
			break
		}
		l.readChar()
	}
}

// readSingleLineComment will continue to consume input until the EOL is reached
func (l *Lexer) readSingleLineComment() {
	for l.ch != 0 {
		if l.ch == 0 {
			break
		}
		if l.ch == '#' {
			l.readChar()
			for l.ch != '\n' {
				l.readChar()
				if l.ch == 0 {
					break
				}
			}
			break
		}
		l.readChar()
	}
}

func (l *Lexer) readExecString() string {
	b := strings.Builder{}
	for {
		l.readChar()
		if l.ch == '`' || l.ch == 0 {
			l.readChar()
			break
		}
		b.WriteRune(l.ch)
	}
	return b.String()
}

func (l *Lexer) readRawString() string {
	b := &strings.Builder{}
	// Skip the first 2 " chars
	l.readChar()
	l.readChar()
	for {
		l.readChar()
		if (l.ch == '"' && l.peekChar() == '"' && l.peekSecondChar() == '"') || l.ch == 0 {
			l.readChar()
			l.readChar()
			break
		}
		b.WriteRune(l.ch)
	}
	return b.String()
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

	if l.ch != '"' && l.ch != '\'' {
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

func (l *Lexer) isOperator(ch rune) bool {
	log.Printf("operators = %s", l.operators)
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
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '`':
		tok.Position = l.posInLine
		tok.Type = TTString
		tok.Literal = l.readExecString()
		return tok
	case 0:
		tok.Position = l.posInLine
		tok.Literal = ""
		tok.Type = TTEof
	case '"':
		if l.peekChar() == '"' && l.peekSecondChar() == '"' {
			tok.Position = l.posInLine
			str := l.readRawString()
			tok.Type = TTString
			tok.Literal = str
		} else {
			str, err := l.readString()
			if err != nil {
				tok = Token{
					Type:     TTIllegal,
					Literal:  "",
					Position: l.posInLine,
				}
			} else {
				tok.Position = l.posInLine
				tok.Type = TTString
				tok.Literal = str
			}
		}
	case '\'':
		str, err := l.readString()
		if err != nil {
			tok = Token{
				Type:     TTIllegal,
				Literal:  "",
				Position: l.posInLine,
			}
		} else {
			tok.Type = TTString
			tok.Literal = str
			tok.Position = l.posInLine
		}
	default:
		if isLetter(l.ch) {
			tok.Position = l.posInLine
			tok.Literal = l.readIdentifier()
			hlStyle, ok := l.LookupIdent(tok.Literal)
			if ok {
				tok.Type = TTIdent
				tok.HlStyleType = hlStyle
			} else {
				tok.Type = TTIllegal
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Position = l.posInLine
			tok.Type, tok.Literal = l.readNumber()
			// TODO: Will panic because no style for number at the moment
			return tok
		} else if l.isOperator(l.ch) {
			tok.Position = l.posInLine
			tok.Type, tok.Literal = TTIdent, string(l.ch)
			tok.HlStyleType = consts.Hl1
			l.readChar()
			return tok
		}
		tok = Token{
			Type:     TTIllegal,
			Literal:  "",
			Position: l.posInLine,
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

// isImportChar will return true if the rune given is allowed as part of an import path
//
// Note: numbers are not allowed in the filename because they are not allowed in identifiers
// this is a design decision and prevents issues.  The reason why '.' is allowed is because
// that will signifiy the path separation in the import path.
//
// We could just use a basic string which would solve most of these problems but i like
// the look of python's import syntax :)
func isImportChar(ch rune) bool {
	return isLetter(ch) || ch == '.'
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
	return '0' == ch || '1' == ch
}
