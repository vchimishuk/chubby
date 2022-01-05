// Copyright 2018 Viacheslav Chimishuk <vchimishuk@yandex.ru>
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

import "fmt"

type State string

const (
	StatePaused  State = "paused"
	StatePlaying State = "playing"
	StateStopped State = "stopped"
)

func parseState(s string) (State, error) {
	var st State = StateStopped
	var err error

	switch s {
	case "paused":
		st = StatePaused
	case "playing":
		st = StatePlaying
	case "stopped":
		st = StateStopped
	default:
		err = fmt.Errorf("invalid state: %s", s)
	}

	return st, err
}
