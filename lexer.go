package easylex

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Lexer is a struct that holds all private state
// necessare for lexing.
type Lexer struct {
	input  string
	state  StateFn
	start  int
	pos    int
	width  int
	tokens chan Token
}

// Lex returns a new Lexer instance that will lex the provided
// input starting from the provided state function.
func Lex(input string, state StateFn) *Lexer {
	return &Lexer{
		input:  input,
		state:  state,
		tokens: make(chan Token, 2),
	}
}

// NextToken returns the next token in the input
// currently being lexed.
func (l *charMatchLexer) NextToken() Token {
	for {
		select {
		case tok := <-l.tokens:
			return tok
		default:
			l.state = l.state(l)
		}
	}
}

func (l *Lexer) Emit(t tokenType) {
	l.tokens <- Token{
		t,
		l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *Lexer) Errorf(format string, args ...interface{}) StateFn {
	l.tokens <- Token{
		tokenError,
		fmt.Sprintf(format, args),
	}
	return nil
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// backup decrements l.pos by the width of the last rune
// processed. backup can only be called once per call to
// next().
func (l *Lexer) backup() {
	l.pos -= l.width
}

// TODO: currently these two are not used, but they may be in the future.

// ignore resets l.start to the current value of l.pos.
// this ignores all the runes processed since the last
// call to ignore() or emit().
func (l *Lexer) ignore() {
	l.start = l.pos
}

// peek returns the value of the rune at l.pos + 1, but
// does not mutate lexer state.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}
