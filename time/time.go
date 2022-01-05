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

package time

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func Parse(s string) (Time, error) {
	pts := reverse(strings.Split(s, ":"))
	i, err := parseSecMin(pts[0])
	if err != nil {
		return 0, fmt.Errorf("seconds: %w", err)
	}
	t := i

	if len(pts) > 1 {
		i, err := parseSecMin(pts[1])
		if err != nil {
			return 0, fmt.Errorf("minutes: %w", err)
		}
		t += i * 60
	}
	if len(pts) > 2 {
		i, err := parseSecMin(pts[2])
		if err != nil {
			return 0, fmt.Errorf("hours: %w", err)
		}
		t += i * 60 * 60
	}
	if len(pts) > 3 {
		return 0, errors.New("bad format")
	}

	return Time(t), nil
}

func parseSecMin(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, errors.New("out of range")
	}
	if i > 59 {
		return 0, errors.New("out of range")
	}

	return i, nil
}

func parseHour(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, errors.New("out of range")
	}

	return i, nil
}

func reverse(s []string) []string {
	var ss []string
	for i := len(s) - 1; i >= 0; i-- {
		ss = append(ss, s[i])
	}

	return ss
}
