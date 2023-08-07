// Copyright 2018-2023 Viacheslav Chimishuk <vchimishuk@yandex.ru>
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
	"errors"
	"fmt"

	"github.com/vchimishuk/chubby/parser"
	"github.com/vchimishuk/chubby/time"
)

type Event interface {
	Event() string
	Serialize() string
}

type CreatePlaylistEvent struct {
	s    string
	Name string
}

func (e *CreatePlaylistEvent) Event() string {
	return "create-playlist"
}

func (e *CreatePlaylistEvent) Serialize() string {
	return e.s
}

type DeletePlaylistEvent struct {
	s    string
	Name string
}

func (e *DeletePlaylistEvent) Event() string {
	return "delete-playlist"
}

func (e *DeletePlaylistEvent) Serialize() string {
	return e.s
}

type StatusEvent struct {
	s           string
	State       State
	PlaylistPos int
	TrackPos    time.Time
	Playlist    *Playlist
	Track       *Track
}

func (e *StatusEvent) Event() string {
	return "status"
}

func (e *StatusEvent) Serialize() string {
	return e.s
}

func parseEvent(name string, lines []string) (Event, error) {
	if len(lines) != 1 {
		return nil, errors.New("protocol error")
	}
	s := lines[0]

	p, err := parser.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("protocol: %w", err.Error())
	}

	switch name {
	case "create-playlist":
		return createCreatePlaylist(s, p)
	case "delete-playlist":
		return createDeletePlaylist(s, p)
	case "status":
		return createStatus(s, p)
	default:
		return nil, fmt.Errorf("protocol: invalid event: %s", name)
	}
}

func createCreatePlaylist(s string, m map[string]any) (Event, error) {
	return &CreatePlaylistEvent{
		s:    s,
		Name: m["name"].(string),
	}, nil
}

func createDeletePlaylist(s string, m map[string]any) (Event, error) {
	return &DeletePlaylistEvent{
		s:    s,
		Name: m["name"].(string),
	}, nil
}

func createStatus(s string, m map[string]any) (Event, error) {
	state, err := parseState(m["state"].(string))
	if err != nil {
		return nil, err
	}

	if state == StateStopped {
		return &StatusEvent{
			s:     s,
			State: state,
		}, nil
	} else {
		return &StatusEvent{
			s:           s,
			State:       state,
			PlaylistPos: m["playlist-position"].(int),
			TrackPos:    time.Time(m["track-position"].(int)),
			Playlist: &Playlist{
				Name:     m["playlist-name"].(string),
				Duration: time.Time(m["playlist-duration"].(int)),
				Length:   m["playlist-length"].(int),
			},
			Track: &Track{
				Path:   m["track-path"].(string),
				Artist: m["track-artist"].(string),
				Album:  m["track-album"].(string),
				Year:   m["track-year"].(int),
				Title:  m["track-title"].(string),
				Number: m["track-number"].(int),
				Length: time.Time(m["track-length"].(int)),
			},
		}, nil
	}
}
