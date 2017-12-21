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
	"reflect"
	"testing"
)

func TestNumber(t *testing.T) {
	err := testMap(`foo: 0, bar: 123, baz: 123456`,
		map[string]interface{}{
			"foo": 0,
			"bar": 123,
			"baz": 123456,
		})
	if err != nil {
		t.Fatal(err)
	}
}

func TestString(t *testing.T) {
	err := testMap(`aaa: "foo", bbb: "foo bar baz", ccc: "foo\"bar'baz", ddd: "абвгд"`,
		map[string]interface{}{
			"aaa": "foo",
			"bbb": "foo bar baz",
			"ccc": `foo"bar'baz`,
			"ddd": "абвгд",
		})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBool(t *testing.T) {
	err := testMap(`foo: true, bar: false`,
		map[string]interface{}{
			"foo": true,
			"bar": false,
		})
	if err != nil {
		t.Fatal(err)
	}
}

func testMap(s string, expected map[string]interface{}) error {
	m, err := Parse(s)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(expected, m) {
		return fmt.Errorf("%+v != %+v", expected, m)
	}

	return nil
}
