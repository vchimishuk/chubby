// Copyright 2017 Viacheslav Chimishuk <vchimishuk@yandex.ru>
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

package chubby

import (
	"fmt"
	"strconv"
)

type Time int

func (t Time) Hour() int {
	return int(t) / (60 * 60)
}

func (t Time) Minute() int {
	return int(t) % (60 * 60) / 60
}

func (t Time) Second() int {
	return int(t) % 60
}

func (t Time) String() string {
	r := ""
	h := t.Hour()
	m := t.Minute()
	s := t.Second()

	if h > 0 {
		r += strconv.Itoa(h) + ":"
	}
	if h > 0 {
		r += fmt.Sprintf("%02d:", m)
	} else {
		r += strconv.Itoa(m) + ":"
	}
	r += fmt.Sprintf("%02d", s)

	return r
}
