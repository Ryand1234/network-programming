package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"unicode"
)

type JSONValue interface{}
type JSONObject map[JSONValue]JSONValue
type JSONArray []JSONValue

func main() {
	folderPath := "./test"

	// Open the folder
	dirEntries, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Error opening folder: %v", err)
	}
	pass, fail := 0, 0
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
		status := checkValidity(string(fileContent))
		if !status {
			if strings.Contains(entry.Name(), "fail") {
				pass++
			} else {
				fmt.Printf("Contents of %s:\n%s\n", entry.Name(), string(fileContent))
				fail++
			}
		} else {
			if strings.Contains(entry.Name(), "pass") {
				pass++
			} else {
				fmt.Printf("Contents of %s:\n%s\n", entry.Name(), string(fileContent))
				fail++
			}
		}
	}
	fmt.Println("Success: ", pass, " Fail: ", fail)
}

type TokenType int

const (
	TokenEOF = iota
	TokenNumber
	TokenLeftCurlyBracket
	TokenRightCurlyBracket
	TokenLeftRoundBracket
	TokenRightRoundBracket
	TokenLeftSquareBracket
	TokenRightSquareBracket
	TokenIdentifier
	TokenSingleQuote
	TokenDoubleQuote
	TokenColon
	TokenComma
	TokenUnknown
	TokenTrue
	TokenFalse
	TokenNull
)

type Token struct {
	Type  TokenType
	Value JSONValue
}

type Lexer struct {
	pos     int
	input   string
	curChar rune
}

type Parser struct {
	lexer *Lexer
	cur   Token
	peek  Token
	error []string
}

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{lexer: lexer, error: []string{}}
	parser.nextToken()
	return parser
}

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.lexer.getNextToken()
}

func (p *Parser) parseArray() (JSONArray, error) {
	list := make(JSONArray, 0)
	isLastComma := false
	for {
		val, err := p.parseValue()

		if err != nil {
			return nil, err
		}
		if p.cur.Type == TokenRightSquareBracket && !isLastComma {
			return list, nil
		}
		if p.cur.Type == TokenRightSquareBracket && isLastComma {
			return nil, fmt.Errorf("Expected element after comma")
		}
		list = append(list, val)
		p.nextToken()

		if p.cur.Type == TokenRightSquareBracket {
			return list, nil
		}
		if p.cur.Type != TokenComma {
			return nil, fmt.Errorf("Expected comma after element but got %v", p.cur.Value)
		} else {
			isLastComma = true
		}
	}
}

func (p *Parser) parseObject() (JSONObject, error) {
	obj := make(JSONObject)
	isLastTokenComma := false
	for {
		p.nextToken()
		if p.cur.Type == TokenRightCurlyBracket && !isLastTokenComma {
			return obj, nil
		}
		if p.cur.Type == TokenRightCurlyBracket && isLastTokenComma {
			return nil, fmt.Errorf("Expected key after comma")
		}
		// if p.cur.Type != TokenDoubleQuote {
		// 	return nil, fmt.Errorf("Expected double quotes bug got %v", p.cur.Value)
		// }
		// p.nextToken()
		if p.cur.Type != TokenIdentifier {
			return nil, fmt.Errorf("Expected identifier but got %v", p.cur.Value)
		}
		key := p.cur.Value
		p.nextToken()
		if p.cur.Type != TokenColon {
			return nil, fmt.Errorf("Expected colon bug got %v", p.cur.Value)
		}
		// p.nextToken()
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = value
		p.nextToken()
		if p.cur.Type == TokenRightCurlyBracket {
			return obj, nil
		}
		if p.cur.Type != TokenComma {
			return nil, fmt.Errorf("Expected Comma but got %v", p.cur.Value)
		} else {
			isLastTokenComma = true
		}
	}
}

func (p *Parser) parseValue() (JSONValue, error) {
	// stack := Stack{}

	p.nextToken()
	switch p.cur.Type {
	case TokenIdentifier:
		return p.cur.Value, nil
	case TokenTrue:
		return true, nil
	case TokenFalse:
		return false, nil
	case TokenNull:
		return p.cur.Value, nil
	case TokenNumber:
		return p.cur.Value, nil
	case TokenLeftCurlyBracket:
		return p.parseObject()
	case TokenLeftSquareBracket:
		return p.parseArray()
	case TokenRightSquareBracket:
		return p.cur, nil
	case TokenEOF:
		return TokenEOF, nil
	default:
		return nil, fmt.Errorf("Unexpected token %v", p.cur.Value)
	}
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

func (l *Lexer) readString() (TokenType, string) {
	start := l.pos
	l.readChar()
	for l.curChar != '"' {
		if l.curChar == 0 {
			break
		}
		l.readChar()
	}
	if l.curChar == TokenEOF {
		return TokenUnknown, ""
	}
	value := l.input[start : l.pos-1]
	l.readChar()
	return TokenIdentifier, value
}

func (l *Lexer) getNextToken() Token {
	l.skipWhiteSpace()
	// fmt.Println("CUR: ", l.curChar, l.pos)
	switch {
	case l.curChar == '{':
		token := Token{Type: TokenLeftCurlyBracket, Value: l.curChar}
		l.readChar()
		return token
	case l.curChar == '}':
		token := Token{Type: TokenRightCurlyBracket, Value: l.curChar}
		l.readChar()
		return token
	case l.curChar == '"':
		tokenType, value := l.readString()
		token := Token{Type: tokenType, Value: value}
		return token
	case l.curChar == '(':
		token := Token{Type: TokenLeftRoundBracket, Value: l.curChar}
		l.readChar()
		return token
	case l.curChar == ')':
		token := Token{Type: TokenRightRoundBracket, Value: l.curChar}
		l.readChar()
		return token
	case l.curChar == '[':
		token := Token{Type: TokenLeftSquareBracket, Value: l.curChar}
		l.readChar()
		return token
	case l.curChar == ']':
		token := Token{Type: TokenRightSquareBracket, Value: l.curChar}
		l.readChar()
		return token
	case l.curChar == ',':
		token := Token{Type: TokenComma, Value: l.curChar}
		l.readChar()
		return token
	case l.curChar == 't':
		token := Token{Type: TokenTrue, Value: true}
		for i := 0; i < 4; i++ {
			l.readChar()
		}
		return token
	case l.curChar == 'f':
		token := Token{Type: TokenFalse, Value: false}
		for i := 0; i < 5; i++ {
			l.readChar()
		}
		return token
	case l.curChar == 'n':
		token := Token{Type: TokenNull, Value: TokenNull}
		for i := 0; i < 4; i++ {
			l.readChar()
		}
		return token
	case unicode.IsDigit(l.curChar):
		val := 0
		for unicode.IsDigit(l.curChar) {
			temp := int(l.curChar - '0')
			val = val*10 + temp
			l.readChar()
		}
		token := Token{Type: TokenNumber, Value: val}
		return token

	case l.curChar == ':':
		token := Token{Type: TokenColon, Value: l.curChar}
		l.readChar()
		return token
	case l.curChar == 0:
		return Token{Type: TokenEOF, Value: ""}
		// default:
		// 	l.readChar()
	default:
		// l.readChar()
		return Token{Type: TokenUnknown, Value: l.curChar}
	}
}

func checkValidity(fileContent string) bool {
	lexer := NewLexer(fileContent)
	// for {
	// 	token := lexer.getNextToken()
	// 	fmt.Printf("Token: Type=%v, Value='%v'\n", token.Type, token.Value)
	// 	if token.Type == TokenEOF {
	// 		break
	// 	}
	// }
	parser := NewParser(lexer)
	parsedValue, err := parser.parseValue()
	if err != nil {
		fmt.Println("Invalid JSON")
		fmt.Println(err)
		return false
	}
	if parsedValue == TokenEOF {
		fmt.Println("Invalid JSON")
		return false
	}
	return true
	// fmt.Println(parsedValue)
	// valid := parser.parseJson()
	// if valid {
	// 	fmt.Println("Valid JSON file")
	// } else {
	// 	fmt.Println("Invalid JSON file")
	// }
}
