// Copyright 2016 Viacheslav Chimishuk <vchimishuk@yandex.ru>
//
// This file is part of chubby.
//
// Chub is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Chub is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Chub. If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type impl struct {
	s   string
	pos int
}

func Parse(s string) (map[string]interface{}, error) {
	p := &impl{s: s, pos: 0}
	m := make(map[string]interface{})

	for !p.eol() {
		p.skipSpaces()

		key, err := p.key()
		if err != nil {
			return nil, err
		}
		p.skipSpaces()
		err = p.consume(":")
		if err != nil {
			return nil, err
		}
		p.skipSpaces()
		val, err := p.val()
		if err != nil {
			return nil, err
		}
		p.skipSpaces()
		if !p.eol() {
			if err := p.consume(","); err != nil {
				return nil, err
			}
		}

		m[key] = val
	}

	return m, nil
}

func (p *impl) key() (string, error) {
	first := true
	k := 0

	for {
		r, n := utf8.DecodeRuneInString(p.s[p.pos+k:])
		if (first && unicode.IsLetter(r)) || (unicode.IsLetter(r) ||
			unicode.IsNumber(r) || r == '_') {
			k += n
		} else {
			break
		}
		first = false
	}

	if k == 0 {
		return "", newError(p.pos, "identifier expected")
	}

	tok := p.s[p.pos : p.pos+k]
	p.pos += k

	return tok, nil
}

func (p *impl) val() (interface{}, error) {
	var val interface{}
	var err error

	r, _ := utf8.DecodeRuneInString(p.s[p.pos:])
	if r == '"' {
		val, err = p.string()
	} else if unicode.IsNumber(r) {
		val, err = p.number()
	} else if r == 't' || r == 'f' {
		val, err = p.boolean()
	} else {
		err = newError(p.pos, "value expected")
	}

	return val, err
}

func (p *impl) string() (string, error) {
	var buf []rune
	k := 0

	if err := p.consume(`"`); err != nil {
		return "", newError(p.pos, "\" expected1")
	}

	for !p.eol() {
		r, n := utf8.DecodeRuneInString(p.s[p.pos+k:])
		if r == '"' {
			break
		}
		if r == '\\' {
			if p.eol() {
				return "", newError(p.pos+k, "end of string expected")
			}
			rr, nn := utf8.DecodeRuneInString(p.s[p.pos+k+n:])
			r = rr
			n += nn
		}

		buf = append(buf, r)
		k += n
	}
	p.pos += k

	if err := p.consume(`"`); err != nil {
		return "", newError(p.pos, "\" expected2")
	}

	return string(buf), nil
}

func (p *impl) number() (int, error) {
	k := 0

	for !p.eol() {
		r, n := utf8.DecodeRuneInString(p.s[p.pos+k:])
		if unicode.IsNumber(r) {
			k += n
		} else {
			break
		}
	}
	if k == 0 {
		return 0, newError(p.pos, "number expected")
	}
	s := p.s[p.pos : p.pos+k]
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, newError(p.pos, "invalid number")
	}
	p.pos += k

	return n, nil
}

func (p *impl) boolean() (bool, error) {
	s := p.s[p.pos:]
	b := false
	n := 0

	if strings.HasPrefix(s, "true") {
		b, n = true, 4
	} else if strings.HasPrefix(s, "false") {
		b, n = false, 5
	} else {
		return false, newError(p.pos, "invalid boolean")
	}

	p.pos += n

	return b, nil
}

func (p *impl) consume(s string) error {
	var err error

	if p.eol() {
		err = newError(p.pos, "EOL")
	} else if strings.HasPrefix(p.s[p.pos:], s) {
		p.pos += len([]byte(s))
	} else {
		err = newError(p.pos, fmt.Sprintf("'%s' expected", s))
	}

	return err
}

func (p *impl) skipSpaces() {
	for p.consume(" ") == nil {

	}
}

func (p *impl) eol() bool {
	return p.pos >= len(p.s)
}
