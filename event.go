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
	"fmt"

	"github.com/vchimishuk/chubby/parser"
	"github.com/vchimishuk/chubby/time"
)

type StatusEvent struct {
	State       State
	PlaylistPos int
	TrackPos    time.Time
	Playlist    *Playlist
	Track       *Track
}

type kv map[string]any

func parseEvent(name string, lines []string) (any, error) {
	m := make([]kv, 0, len(lines))
	for _, l := range lines {
		p, err := parser.Parse(l)
		if err != nil {
			return nil, fmt.Errorf("protocol: %w", err.Error())
		}
		m = append(m, p)
	}

	switch name {
	case "status":
		return parseStatus(m)
	default:
		return nil, fmt.Errorf("protocol: invalid event: %s", name)
	}
}

func parseStatus(lines []kv) (interface{}, error) {
	e := &StatusEvent{Playlist: &Playlist{}, Track: &Track{}}

	for _, l := range lines {
		for k, v := range l {
			switch k {
			case "state":
				st, err := parseState(v.(string))
				if err != nil {
					return nil, err
				}
				e.State = st
			case "playlist-position":
				e.PlaylistPos = v.(int)
			case "track-position":
				e.TrackPos = time.Time(v.(int))
			case "playlist-name":
				e.Playlist.Name = v.(string)
			case "playlist-duration":
				e.Playlist.Duration = time.Time(v.(int))
			case "playlist-length":
				e.Playlist.Length = v.(int)
			case "track-path":
				e.Track.Path = v.(string)
			case "track-artist":
				e.Track.Artist = v.(string)
			case "track-album":
				e.Track.Album = v.(string)
			case "track-title":
				e.Track.Title = v.(string)
			case "track-number":
				e.Track.Number = v.(int)
			case "track-length":
				e.Track.Length = time.Time(v.(int))
			}
		}
	}

	return e, nil
}
