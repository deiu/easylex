package easylex

import (
	"strings"
)

type Matcher struct {
	*unionMatcher
}

// NewMatcher creates a new instance of Matcher with default
// behavior (in other words, it will not match anything).
func NewMatcher() *Matcher {
	return &Matcher{
		&unionMatcher{
			[]textMatcher{},
		},
	}
}

// AcceptRunes modifies a Matcher to accept any runes that
// are contained withing the provided string.
// The modified Matcher is returned to the caller.
func (m *Matcher) AcceptRunes(valid string) *Matcher {
	// TODO: check up on the implementation details of the rune vs byte slice thing
	r := &runeMatcher{valid}
	m.add(r)
	return m
}

// RejectRunes modifies a Matcher to accept any runes
// that are not contained withing the provided string.
// The modified Matcher is returned to the caller.
func (m *Matcher) RejectRunes(invalid string) *Matcher {
	r := &runeFilter{invalid}
	m.add(r)
	return m
}

// AcceptUnicodeRange modifies a Matcher to accept any runes
// that fall between the two provided runes (inclusive).
// The modified Matcher is returned to the caller.
func (m *Matcher) AcceptUnicodeRange(first, last rune) *Matcher {
	u := &unicodeRangeMatcher{first, last}
	m.add(u)
	return m
}

// AcceptString modifies a Matcher to accept a string of
// characters in the input that exactly matches the provided
// string. The modified Matcher is returned to the caller.
func (m *Matcher) AcceptString(exact string) *Matcher {
	p := &prefixMatcher{exact}
	m.add(p)
	return m
}

// Union modifies a Matcher to accept the set of characters
// equal to the union between the current Matcher's set of
// accepted characters and another Matcher's set of
// accepted characters.
// The modified Matcher is returned to the caller.
func (m *Matcher) Union(other *Matcher) *Matcher {
	u := other
	m.add(u)
	return m
}

// MatchOne accepts the next input rune if that rune conforms
// to the rules currently represented by this Matcher.
// If the next character was accepted, MatchOne returns true.
// If the next character was not accepted, MatchOne returns
// false and the state of the Lexer is left unmodified.
func (m *Matcher) MatchOne(l *Lexer) bool {
	return m.match(l)
}

// MatchRun accepts as many consecutive input characters as
// fit the rules currently represented by this Matcher.
// If at least one character was accepted, MatchRun returns true.
// If no characters were accepted, MatchRun returns false and
// the state of the Lexer is left unmodified.
func (m *Matcher) MatchRun(l *Lexer) bool {
	success := false
	for m.match(l) {
		if !success {
			success = true
		}
	}
	return success
}

// TODO: make textMatcher an exported interface and allow the
// addition of custom matcher modules to a Matcher
type textMatcher interface {
	match(*Lexer) bool
}

type runeMatcher struct {
	valid string
}

func (r *runeMatcher) match(l *Lexer) bool {
	if strings.IndexRune(r.valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

type runeFilter struct {
	invalid string
}

func (r *runeFilter) match(l *Lexer) bool {
	if strings.IndexRune(r.invalid, l.Next()) >= 0 {
		l.Backup()
		return false
	}
	return true
}

type unicodeRangeMatcher struct {
	first rune
	last  rune
}

func (u *unicodeRangeMatcher) match(l *Lexer) bool {
	next := l.Next()
	if next >= u.first && next <= u.last {
		return true
	}
	l.Backup()
	return false
}

type prefixMatcher struct {
	prefix string
}

func (p *prefixMatcher) match(l *Lexer) bool {
	if strings.HasPrefix(l.input[l.pos:], p.prefix) {
		l.pos += len(p.prefix)
		return true
	}
	return false
}

type unionMatcher struct {
	matchers []textMatcher
}

func (u *unionMatcher) match(l *Lexer) bool {
	for _, m := range u.matchers {
		if m.match(l) {
			return true
		}
	}
	return false
}

func (u *unionMatcher) add(t textMatcher) {
	u.matchers = append(u.matchers, t)
}
