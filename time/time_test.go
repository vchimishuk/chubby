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

import "testing"

func TestString(t *testing.T) {
	assertStrEq(t, Time(60*60).String(), "1:00:00")
	assertStrEq(t, Time(60*60+60*2+3).String(), "1:02:03")
	assertStrEq(t, Time(60).String(), "1:00")
	assertStrEq(t, Time(60+59).String(), "1:59")
	assertStrEq(t, Time(59).String(), "0:59")
	assertStrEq(t, Time(0).String(), "0:00")
}

func TestParse(t *testing.T) {
	tm, err := Parse("0")
	assertTimeEq(t, Time(0), tm)
	assertErrNil(t, err)
	tm, err = Parse("30")
	assertTimeEq(t, Time(30), tm)
	assertErrNil(t, err)
	tm, err = Parse("59")
	assertTimeEq(t, Time(59), tm)
	assertErrNil(t, err)
	tm, err = Parse("1:09")
	assertTimeEq(t, Time(69), tm)
	assertErrNil(t, err)
	tm, err = Parse("59:59")
	assertTimeEq(t, Time(59*60+59), tm)
	assertErrNil(t, err)
	tm, err = Parse("01:02:03")
	assertTimeEq(t, Time(1*60*60+2*60+3), tm)
	assertErrNil(t, err)
	tm, err = Parse("00:00:00")
	assertTimeEq(t, Time(0), tm)
	assertErrNil(t, err)

	_, err = Parse("")
	assertStrEq(t, "seconds: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
	_, err = Parse("FF")
	assertStrEq(t, "seconds: strconv.Atoi: parsing \"FF\": invalid syntax", err.Error())
	_, err = Parse("00:00:00:00")
	assertStrEq(t, "bad format", err.Error())
}

func assertStrEq(t *testing.T, a, b string) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func assertTimeEq(t *testing.T, a, b Time) {
	if a != b {
		t.Fatalf("%d != %d", a, b)
	}
}

func assertErrNil(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%s != nil", err)
	}
}
