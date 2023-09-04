package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"unicode"
)

func main() {
	folderPath := "./tests/step1"

	// Open the folder
	dirEntries, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Error opening folder: %v", err)
	}

	// Iterate over the files in the folder
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue // Skip directories
		}

		// Create the full file path
		filePath := filepath.Join(folderPath, entry.Name())

		// Read the file
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading %s: %v", filePath, err)
			continue
		}

		// Process the file content as needed
		fmt.Printf("Contents of %s:\n%s\n", entry.Name(), string(fileContent))
		checkValidity(string(fileContent))
	}
}

type TokenType int

const (
	TokenEOF = iota
	TokenNumber
	TokenLeftCurlyBracket
	TokenRightCurlyBracket
	TokenLeftSquareBracket
	TokenRightSquareBracket
	TokenString
	TokenSingleQuote
	TokenDoubleQuote
	TokenColon
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	pos     int
	input   string
	curChar rune
}

func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.curChar = 0
	} else {
		l.curChar = rune(l.input[l.pos])
	}
	l.pos++
}

func (l *Lexer) skipWhiteSpace() {
	for unicode.IsSpace(l.curChar) {
		l.readChar()
	}
}

func (l *Lexer) getNextToken() Token {
	l.skipWhiteSpace()
	switch l.curChar {
	case '{':
		token := Token{Type: TokenLeftCurlyBracket, Value: string(l.curChar)}
		l.readChar()
		return token
	case '}':
		token := Token{Type: TokenRightCurlyBracket, Value: string(l.curChar)}
		l.readChar()
		return token
	case 0:
		return Token{Type: TokenEOF, Value: ""}
	}
	return Token{Type: TokenEOF, Value: ""}
}

func checkValidity(fileContent string) {
	lexer := NewLexer(fileContent)

	for {
		token := lexer.getNextToken()
		fmt.Printf("Token: Type=%v, Value='%v'\n", token.Type, token.Value)
		if token.Type == TokenEOF {
			break
		}
	}
}
