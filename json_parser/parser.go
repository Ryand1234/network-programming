package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"unicode"
)

func main() {
	folderPath := "./tests/step3"

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
	TokenBoolean
	TokenNull
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

type Parser struct {
	lexer *Lexer
	cur   Token
	peek  Token
	error []string
}

type Stack struct {
	data []Token
}

func (s *Stack) Push(t Token) {
	s.data = append(s.data, t)
}

func (s *Stack) Pop() Token {
	if len(s.data) == 0 {
		return Token{Type: TokenEOF, Value: ""}
	}

	top := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return top
}

func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack) TestRow() bool {

	return true
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

type JSONValue interface{}
type JSONObject map[string]JSONValue

func (p *Parser) parseObject() (JSONObject, error) {
	obj := make(JSONObject)
	for {
		p.nextToken()
		if p.cur.Type == TokenRightCurlyBracket {
			return obj, nil
		}
		// if p.cur.Type != TokenDoubleQuote {
		// 	return nil, fmt.Errorf("Expected double quotes bug got %v", string(p.cur.Value))
		// }
		// p.nextToken()
		if p.cur.Type != TokenIdentifier {
			return nil, fmt.Errorf("Expected identifier bug got %v", string(p.cur.Value))
		}
		key := p.cur.Value
		p.nextToken()
		if p.cur.Type != TokenColon {
			return nil, fmt.Errorf("Expected colon bug got %v", string(p.cur.Value))
		}
		// p.nextToken()
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = value
		// fmt.Println("CUR: ", string(p.cur.Value))
		p.nextToken()
		// fmt.Println("CUR: ", string(p.cur.Value))
		if p.cur.Type == TokenRightCurlyBracket {
			return obj, nil
		}
		if p.cur.Type != TokenComma {
			return nil, fmt.Errorf("Expected Comma but got %v", string(p.cur.Value))
		}
	}
}

func (p *Parser) parseValue() (JSONValue, error) {
	// stack := Stack{}

	p.nextToken()
	switch p.cur.Type {
	case TokenIdentifier:
		return string(p.cur.Value), nil
	case TokenBoolean:
		return string(p.cur.Value), nil
	case TokenNull:
		return string(p.cur.Value), nil
	case TokenNumber:
		val, _ := strconv.Atoi(p.cur.Value)
		return val, nil
	case TokenLeftCurlyBracket:
		return p.parseObject()
	// case TokenLeftSquareBracket:
	// 	return parseArray(p)
	case TokenEOF:
		return TokenEOF, nil
	default:
		return nil, fmt.Errorf("Unexpected token %v", string(p.cur.Value))
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
	// fmt.Println("CUR: ", string(l.curChar), l.pos)
	switch {
	case l.curChar == '{':
		token := Token{Type: TokenLeftCurlyBracket, Value: string(l.curChar)}
		l.readChar()
		return token
	case l.curChar == '}':
		token := Token{Type: TokenRightCurlyBracket, Value: string(l.curChar)}
		l.readChar()
		return token
	case l.curChar == '"':
		tokenType, value := l.readString()
		token := Token{Type: tokenType, Value: value}
		return token
	case l.curChar == '(':
		token := Token{Type: TokenLeftRoundBracket, Value: string(l.curChar)}
		l.readChar()
		return token
	case l.curChar == ')':
		token := Token{Type: TokenRightRoundBracket, Value: string(l.curChar)}
		l.readChar()
		return token
	case l.curChar == '[':
		token := Token{Type: TokenLeftSquareBracket, Value: string(l.curChar)}
		l.readChar()
		return token
	case l.curChar == ']':
		token := Token{Type: TokenRightSquareBracket, Value: string(l.curChar)}
		l.readChar()
		return token
	case l.curChar == ',':
		token := Token{Type: TokenComma, Value: string(l.curChar)}
		l.readChar()
		return token
	case l.curChar == 't':
		token := Token{Type: TokenBoolean, Value: string("true")}
		for i := 0; i < 4; i++ {
			l.readChar()
		}
		return token
	case l.curChar == 'f':
		token := Token{Type: TokenBoolean, Value: string("false")}
		for i := 0; i < 5; i++ {
			l.readChar()
		}
		return token
	case l.curChar == 'n':
		token := Token{Type: TokenNull, Value: string("null")}
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
		token := Token{Type: TokenNumber, Value: fmt.Sprint(val)}
		return token

	case l.curChar == ':':
		token := Token{Type: TokenColon, Value: string(l.curChar)}
		l.readChar()
		return token
	// case unicode.IsLetter(l.curChar) || unicode.IsDigit(l.curChar):
	// 	var identifier string
	// 	for unicode.IsLetter(l.curChar) || unicode.IsDigit(l.curChar) {
	// 		identifier += string(l.curChar)
	// 		l.readChar()
	// 	}
	// 	token := Token{Type: TokenIdentifier, Value: identifier}
	// 	// l.readChar()
	// 	return token
	case l.curChar == 0:
		return Token{Type: TokenEOF, Value: ""}
		// default:
		// 	l.readChar()
	default:
		// l.readChar()
		return Token{Type: TokenUnknown, Value: string(l.curChar)}
	}
}

func checkValidity(fileContent string) {
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
		return
	}
	if parsedValue == TokenEOF {
		fmt.Println("Invalid JSON")
		return
	}
	fmt.Println("Valid JSON")
	fmt.Println(parsedValue)
	// valid := parser.parseJson()
	// if valid {
	// 	fmt.Println("Valid JSON file")
	// } else {
	// 	fmt.Println("Invalid JSON file")
	// }
}
