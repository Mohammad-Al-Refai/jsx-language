package lexer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Tokenized struct {
	Pos     Position `json:"position"`
	Literal string   `json:"literal"`
	Token   Token    `json:"token"`
}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Lexer struct {
	Pos    Position
	Reader *bufio.Reader
}

func (l *Lexer) LoadFileReader(reader io.Reader) []Tokenized {
	tokenized := []Tokenized{}
	l.Pos = Position{Line: 1, Column: 0}
	l.Reader = bufio.NewReader(reader)
	for {
		pos, tok, lit := l.Lex()
		if tok == EOF {
			tokenized = append(tokenized, Tokenized{Pos: pos, Literal: "", Token: EOF})
			break
		}
		tokenized = append(tokenized, Tokenized{Pos: pos, Literal: lit, Token: tok})
	}
	return tokenized
}
func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.Pos, EOF, ""
			}
			panic(err)
		}
		l.Pos.Column++

		switch r {
		case '\n':
			l.resetPosition()
		case '>':
			return l.Pos, CLOSE_OPEN_TAG, ">"
		case '<':
			return l.Pos, OPEN_TAG, "<"
		case '+':
			return l.Pos, ADD, "+"
		case '[':
			return l.Pos, LBRACK, "["
		case ']':
			return l.Pos, RBRACK, "]"
		case '.':
			return l.Pos, DOT, "."
		case '%':
			return l.Pos, MOD, "%"
		case '/':
			next, err := l.Reader.Peek(1)
			if err != nil {
				if err == io.EOF {
					panic(err)
				}
			}
			if string(next) == ">" {
				l.Reader.ReadRune()
				return l.Pos, CLOSE_TAG, "/>"
			} else {
				return l.Pos, DIV, "/"
			}
		case '-':
			return l.Pos, SUB, "-"
		case '*':
			return l.Pos, MUL, "*"
		case '{':
			return l.Pos, LBRACE, "{"
		case '}':
			return l.Pos, RBRACE, "}"
		case '(':
			return l.Pos, LPAREN, "("
		case ')':
			return l.Pos, RPAREN, ")"
		case ',':
			return l.Pos, COMMA, ","
		case '=':
			next, err := l.Reader.Peek(1)
			if err != nil {
				if err == io.EOF {
					return l.Pos, EQUAL, "="
				}
			}
			if string(next) == "=" {
				l.Reader.ReadRune()
				return l.Pos, EQUAL_EQUAL, "=="
			} else {
				return l.Pos, EQUAL, "="
			}
		case '!':
			next, err := l.Reader.Peek(1)
			if err != nil {
				if err == io.EOF {
					l.threwError(fmt.Sprintf("unknown token '%v'", "!"))
				}
			}
			if string(next) == "=" {
				l.Reader.ReadRune()
				return l.Pos, NOT_EQUAL, "!="
			}
		case '#':
			// ignore anything after # until detect a new line
			for {
				x, _, err := l.Reader.ReadRune()
				if err != nil {
					break
				}
				if x == '\n' {
					l.resetPosition()
					break
				}
			}
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.Pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				startPos := l.Pos
				l.backup()
				lit := l.lexIdent()
				if ok, token := IsKeyword(lit); ok {
					return startPos, token, lit
				}
				return startPos, IDENT, lit
			} else if r == '"' {
				lit := l.lexQuotation()
				return l.Pos, STRING, lit
			} else {
				return l.Pos, ILLEGAL, string(r)
			}
		}
	}

}
func (l *Lexer) resetPosition() {
	l.Pos.Line++
	l.Pos.Column = 0
}
func (l *Lexer) backup() {
	if err := l.Reader.UnreadRune(); err != nil {
		panic(err)
	}
	l.Pos.Column--
}

// lexInt scans the input until the end of an integer and then returns the
// literal.
func (l *Lexer) lexInt() string {
	var literal string
	for {
		r, _, err := l.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return literal
			}
		}

		l.Pos.Column++
		if unicode.IsDigit(r) {
			literal = literal + string(r)
		} else {
			// scanned something not in the integer
			l.backup()
			return literal
		}
	}
}
func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.Pos.Column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}
func (l *Lexer) threwError(message string) {
	fmt.Println(fmt.Errorf(fmt.Sprintf("[LexerError] %v", message)))
	os.Exit(1)
}
func (l *Lexer) lexQuotation() string {
	var str string
	for {
		r, _, err := l.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.threwError(fmt.Sprintf(`Missing '"' at %v:%v`, l.Pos.Line, l.Pos.Column))
				return str
			}
		}
		l.Pos.Column++
		if r == '\n' {
			l.threwError(fmt.Sprintf(`Missing '"' at %v:%v`, l.Pos.Line, l.Pos.Column))
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || unicode.IsSymbol(r) || string(r) == "!" {
			str = str + string(r)
		} else if r == '"' {
			return str
		} else {
			return str
		}
	}
}
