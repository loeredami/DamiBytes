package main

import (
	"strconv"
	"strings"
)

var tokenTypes struct {
	Identifier uint64
	Label      uint64
	String     uint64
	Number     uint64
} = struct {
	Identifier uint64
	Label      uint64
	String     uint64
	Number     uint64
}{
	0,
	1,
	2,
	4,
}

type Token interface {
	getType() uint64
}

type TokenWithString struct {
	val string
	tokenType uint64
}

type TokenWithNumber struct {
	val uint64
}

func (stringToken *TokenWithString) getType() uint64 {
	return stringToken.tokenType
}

func (stringToken *TokenWithString) getValue() string {
	return stringToken.val
}

func (numberToken *TokenWithNumber) getType() uint64 {
	return tokenTypes.Number
}

func (numberToken *TokenWithNumber) getValue() uint64 {
	return numberToken.val
}

type Lexer struct {
	text    string
	curChar byte
	idx     uint64
	done    bool
	tokens  []Token
}

func makeLexer(text string) *Lexer {
	lexer := &Lexer{text+" ", ' ', 0, false, []Token{}}

	lexer.updateChar()

	return lexer
}

func (lexer *Lexer) updateChar() {
	if lexer.idx >= uint64(len(lexer.text)) {
		lexer.done = true
		return
	}
	lexer.curChar = lexer.text[lexer.idx]
}

func (lexer *Lexer) advance() {
	lexer.idx += 1

	lexer.updateChar()
}

func (lexer *Lexer) collectIdentifierOrLabel() {
	identifier := ""
	tokenType := tokenTypes.Identifier

	for {
		if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_", rune(lexer.curChar)) {
			break
		}

		identifier += string(lexer.curChar)
		lexer.advance()
	}

	if lexer.curChar == ':' {
		lexer.advance()
		tokenType = tokenTypes.Label
	}

	lexer.tokens = append(lexer.tokens, &TokenWithString{identifier, tokenType})
}

func (lexer *Lexer) getSpecialCharacter() byte {
	lexer.advance()

	switch lexer.curChar {
	case 'n': return '\n'
	case 'r': return '\r'
	case 't': return '\t'
	case 'a': return '\a'
	case 'b': return '\b'
	}

	return lexer.curChar
}

func (lexer *Lexer) collectString() {
	lexer.advance()

	stringVal := ""

	for {
		if lexer.curChar == '"' {
			break
		}

		if lexer.curChar == '\\' {
			stringVal += string(lexer.getSpecialCharacter())
			lexer.advance()
			continue
		}

		stringVal += string(lexer.curChar)
		lexer.advance()
	}

	lexer.advance()

	lexer.tokens = append(lexer.tokens, &TokenWithString{stringVal, tokenTypes.String})
}

func (lexer *Lexer) collectNumber() {
	stringVal := ""

	BASE := 10

	for {
		if !strings.ContainsRune("0123456789xABCDEFabcdef", rune(lexer.curChar)) {
			break
		}

		if lexer.curChar == 'b' {
			BASE = 2
			stringVal = ""
			lexer.advance()
			continue
		}

		if lexer.curChar == 'x' {
			BASE = 16
			stringVal = ""
			lexer.advance()
			continue
		}

		stringVal += string(lexer.curChar)
		lexer.advance()
	}

	num, err := strconv.ParseUint(stringVal, BASE, 64)

	if err != nil {
		panic(err)
	}

	lexer.tokens = append(lexer.tokens, &TokenWithNumber{num})
}

func (lexer *Lexer) skipLine() {
	for {
		if (lexer.curChar == '\n') {
			break
		}
		lexer.advance()
	}

	lexer.advance()
}

func (lexer *Lexer) lex() []Token {
	for {
		if lexer.done {
			break
		}

		if strings.ContainsRune(" \n\r\t", rune(lexer.curChar)) {
			lexer.advance()
			continue
		}

		if strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_", rune(lexer.curChar)) {
			lexer.collectIdentifierOrLabel()
			continue
		}

		if lexer.curChar == '"' {
			lexer.collectString()
			continue
		}

		if strings.ContainsRune("0123456789", rune(lexer.curChar)) {
			lexer.collectNumber()
			continue
		}

		if lexer.curChar == '#' {
			lexer.skipLine()
		}
	}

	return lexer.tokens
}