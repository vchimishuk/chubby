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

import (
	"errors"
	"fmt"
	"net"

	"github.com/vchimishuk/chubby/parser"
	"github.com/vchimishuk/chubby/textconn"
)

type StatusEvent struct {
	State       State
	PlaylistPos int
	TrackPos    Time
	Playlist    *Playlist
	Track       *Track
}

type mapLine map[string]interface{}

type protocolError string

func (err protocolError) Error() string {
	return "invalid server response: " + string(err)
}

type NotifClient struct {
	conn *textconn.TextConn
}

func (c *NotifClient) Connect(host string, port int) error {
	if c.conn != nil {
		return errors.New("already connected")
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	c.conn = textconn.New(conn)

	// TODO: Read server's greetings.
	// if _, err := c.conn.ReadLine(); err != nil {
	// 	c.conn.Close()
	// 	return err
	// }
	// if _, err := c.conn.ReadLine(); err != nil {
	// 	c.conn.Close()
	// 	return err
	// }

	return nil
}

func (c *NotifClient) Accept() (interface{}, error) {
	name, err := c.conn.ReadLine()
	if err != nil {
		return nil, err
	}

	var lines []mapLine
	for {
		l, err := c.conn.ReadLine()
		if err != nil {
			return nil, err
		}
		if l == "" {
			break
		}
		m, err := parser.Parse(l)
		if err != nil {
			return nil, protocolError(err.Error())
		}
		lines = append(lines, m)
	}

	var e interface{}
	switch name {
	case "status":
		e, err = status(lines)
	default:
		return nil, protocolError(fmt.Sprintf("invlid event %s", name))
	}

	return e, err
}

func (c *NotifClient) Close() error {
	if c.conn == nil {
		return errors.New("not connected")
	}
	err := c.conn.Close()
	c.conn = nil

	return err
}

func status(lines []mapLine) (interface{}, error) {
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
				e.TrackPos = Time(v.(int))
			case "playlist-name":
				e.Playlist.Name = v.(string)
			case "playlist-duration":
				e.Playlist.Duration = Time(v.(int))
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
				e.Track.Length = Time(v.(int))
			}
		}
	}

	return e, nil
}
