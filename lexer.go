package laxpars

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Lexer struct {
	input *bufio.Reader
	current rune
}

func NewLexer(input io.Reader) *Lexer {
	return &Lexer{
		input: bufio.NewReader(input),
	}
}

const (
	TOKEN_TYPE_DECLARATION = "DECLARATION"
	TOKEN_TYPE_COMMENT     = "COMMENT"
	TOKEN_TYPE_START_TAG_OPEN   = "TAG_START_OPEN"
	TOKEN_TYPE_END_TAG_OPEN   = "TAG_END_OPEN"
	TOKEN_TYPE_TAG_CLOSE    = "TAG_CLOSE"
	TOKEN_TYPE_EQUAL = "EQUAL"
	TOKEN_TYPE_STRING = "STRING"
	TOKEN_TYPE_NAME    = "NAME"
	TOKEN_TYPE_TEXT    = "TEXT"
	TOKEN_TYPE_UNKNOWN = "UNKNOWN"
	TOKEN_TYPE_EOF = "EOF"
)

type Token struct {
	Type  string
	Value string
}

func (this *Token) String() string {
    return fmt.Sprintf("Token(%s, %s)", this.Type, this.Value)
}

func endsWith(source string, contain string) bool {
	if len(contain) > len(source) {
		return false
	}

	return source[len(source) - len(contain):] == contain
}

func (this *Lexer) readNextRune() {
	var err error

	this.current, _, err = this.input.ReadRune()
	if err != nil {
		panic(err)
	}
}

func (this *Lexer) peekNextRune() rune {
	original := this.current

	this.readNextRune()

	nextRune := this.current

	err := this.input.UnreadRune()
	if err != nil {
		panic(err)
	}
	this.current = original

	return nextRune
}

func (this *Lexer) isFinished() bool {
	_, _, err := this.input.ReadRune()
	if err != nil {
		return err == io.EOF
	}

	err = this.input.UnreadRune()
	if err != nil {
		panic(err)
	}

	return false
}

func (this *Lexer) consume(expected rune) {
	if this.current != expected {
		panic(fmt.Errorf("unexpected %c", expected))
	}

	this.readNextRune()
}

func (this *Lexer) tokenizeDeclaration() *Token {
	this.consume('!')

	declaration := ""
	for ; this.current != '>'; this.readNextRune() {
		declaration += string(this.current)
	}

	return &Token {
		Type: TOKEN_TYPE_DECLARATION,
		Value: declaration,
	}
}

func (this *Lexer) tokenizeComment() *Token {
	this.consume('!')
	this.consume('-')
	this.consume('-')

	comment := ""
	for ; ; this.readNextRune() {
		comment += string(this.current)

		if endsWith(comment, "-->") {
			break
		}
	}

	return &Token{
		Type: TOKEN_TYPE_COMMENT,
		Value: strings.TrimSuffix(comment, "-->"),
	}
}

func (this *Lexer) tokenizeString() *Token {
	start := this.current

	content := string(start)
	this.readNextRune()

	for ; this.current != start; this.readNextRune() {
		content += string(this.current)
	}

	content += string(start)

	return &Token {
		Type: TOKEN_TYPE_STRING,
		Value: content,
	}
}

func isValidNameRune(r rune) bool {
    switch r {
    case '"', '\'', '`', '=', '<', '>':
        return false
    }

    return !unicode.IsSpace(r)
}

func (this *Lexer) tokenizeName() *Token {
	name := ""

	for ; ; this.readNextRune() {
		name += string(this.current)


		if !isValidNameRune(this.peekNextRune()) {
			return &Token{
				Type: TOKEN_TYPE_NAME,
				Value: name,
			}
		}
	}
}

func (this *Lexer) Next() *Token {
	if this.isFinished() {
		return &Token {
			Type: TOKEN_TYPE_EOF,
		}
	}

	this.readNextRune()

	for ; unicode.IsSpace(this.current); this.readNextRune() {
		if this.isFinished() {
			return &Token {
				Type: TOKEN_TYPE_EOF,
			}
		}

		continue
	}

	if this.current == '<' {
		nextRune := this.peekNextRune()

		if nextRune == '/' {
			this.readNextRune()

			return &Token{
				Type: TOKEN_TYPE_END_TAG_OPEN,
			}
		}

		if nextRune == '!' {
			this.readNextRune()

			if this.peekNextRune() == '-' {
				return this.tokenizeComment()
			}

			return this.tokenizeDeclaration()
		}

		return &Token{
			Type: TOKEN_TYPE_START_TAG_OPEN,
		}
	}

	if this.current == '/' && this.peekNextRune() == '>' {
		this.Next()

		return &Token{
			Type: TOKEN_TYPE_TAG_CLOSE,
		}
	}

	if this.current == '>' {
		return &Token{
			Type: TOKEN_TYPE_TAG_CLOSE,
		}
	}

	if this.current == '=' {
		return &Token{
			Type: TOKEN_TYPE_EQUAL,
		}
	}

	if this.current == '\'' || this.current == '"' {
		return this.tokenizeString()
	}

	if isValidNameRune(this.current) {
		return this.tokenizeName()
	}

	return &Token {
		Type: TOKEN_TYPE_UNKNOWN,
	}
}

func (this *Lexer) GetText() string {
	text := ""

	for {
		if this.peekNextRune() == '<' {
			return text
		}

		this.readNextRune()
		text += string(this.current)
	}
}

func (this *Lexer) GetRawText(endtag string) string {
	this.readNextRune()

	endtag = "</" + endtag + ">"

	text := ""
	for ; ; this.readNextRune() {
		text += string(this.current)

		if endsWith(text, endtag) {
			return strings.TrimSuffix(text, endtag)
		}
	}
}