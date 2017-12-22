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

import "testing"

func TestString(t *testing.T) {
	assertEq(t, Time(60*60).String(), "1:00:00")
	assertEq(t, Time(60*60+60*2+3).String(), "1:02:03")
	assertEq(t, Time(60).String(), "1:00")
	assertEq(t, Time(60+59).String(), "1:59")
	assertEq(t, Time(59).String(), "0:59")
	assertEq(t, Time(0).String(), "0:00")
}

func assertEq(t *testing.T, a, b string) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
