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
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/vchimishuk/chubby/parser"
	"github.com/vchimishuk/chubby/textconn"
)

const (
	CmdCreatePlaylist = "create-playlist"
	CmdDeletePlaylist = "delete-playlist"
	CmdKill           = "kill"
	CmdList           = "list"
	CmdNext           = "next"
	CmdPause          = "pause"
	CmdPing           = "ping"
	CmdPlay           = "play"
	CmdPlaylists      = "playlists"
	CmdPrev           = "prev"
	CmdRenamePlaylist = "rename-playlist"
	CmdStatus         = "status"
	CmdStop           = "stop"
)

type Playlist struct {
	Name     string
	Duration Time
	Length   int
}

type State int

const (
	StatePlaying State = iota
	StatePaused
	StateStopped
)

type Status struct {
	State       State
	Playlist    string
	PlaylistLen int
	PlaylistPos int
	Track       string
	TrackLen    Time
	TrackPos    Time
}

type Entry interface {
	IsDir() bool
	Dir() *Dir
	Track() *Track
}

type Dir struct {
	Path string
	Name string
}

func (d *Dir) IsDir() bool {
	return true
}

func (d *Dir) Dir() *Dir {
	return d
}

func (d *Dir) Track() *Track {
	panic("not a track")
}

type Track struct {
	Path   string
	Artist string
	Album  string
	Title  string
	Number int
	Length Time
}

func (t *Track) IsDir() bool {
	return false
}

func (t *Track) Dir() *Dir {
	panic("not a directory")
}

func (t *Track) Track() *Track {
	return t
}

type Chubby struct {
	conn *textconn.TextConn
}

func (c *Chubby) Connect(host string, port int) error {
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

	// Read server's greetings.
	if _, err := c.conn.ReadLine(); err != nil {
		c.conn.Close()
		return err
	}
	if _, err := c.conn.ReadLine(); err != nil {
		c.conn.Close()
		return err
	}

	return nil
}

func (c *Chubby) Close() error {
	if c.conn == nil {
		return errors.New("not connected")
	}
	err := c.conn.Close()
	c.conn = nil

	return err
}

func (c *Chubby) CreatePlaylist(name string) error {
	_, err := c.cmd(CmdCreatePlaylist, name)

	return err
}

func (c *Chubby) DeletePlaylist(name string) error {
	_, err := c.cmd(CmdDeletePlaylist, name)

	return err
}

func (c *Chubby) Kill() error {
	_, err := c.cmd(CmdKill)

	return err
}

func (c *Chubby) List(path string) ([]Entry, error) {
	lines, err := c.cmd(CmdList, path)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, len(lines))
	for i, line := range lines {
		entries[i], err = parseEntry(line)
		if err != nil {
			return nil, err
		}
	}

	return entries, nil
}

func (c *Chubby) Next() error {
	_, err := c.cmd(CmdNext)

	return err
}

func (c *Chubby) Pause() error {
	_, err := c.cmd(CmdPause)

	return err
}

func (c *Chubby) Ping() error {
	_, err := c.cmd(CmdPing)

	return err
}

func (c *Chubby) Play(pth string) error {
	_, err := c.cmd(CmdPlay, pth)

	return err
}

func (c *Chubby) Playlists() ([]*Playlist, error) {
	lines, err := c.cmd(CmdPlaylists)
	if err != nil {
		return nil, err
	}

	pls := make([]*Playlist, len(lines))
	for i, line := range lines {
		pls[i], err = parsePlaylist(line)
		if err != nil {
			return nil, err
		}
	}

	return pls, nil
}

func (c *Chubby) Prev() error {
	_, err := c.cmd(CmdPrev)

	return err
}

func (c *Chubby) RenamePlaylist(from, to string) error {
	_, err := c.cmd(CmdRenamePlaylist, from, to)

	return err
}

func (c *Chubby) Status() (*Status, error) {
	lines, err := c.cmd(CmdStatus)

	if len(lines) != 1 {
		return nil, errors.New("invalid server response")
	}

	m, err := parser.Parse(lines[0])
	if err != nil {
		return nil, err
	}

	var st State

	if m["state"].(string) == "playing" {
		st = StatePlaying
	} else if m["state"].(string) == "paused" {
		st = StatePaused
	} else if m["state"].(string) == "stopped" {
		st = StateStopped
	} else {
		return nil, fmt.Errorf("invalid status: %s", m["state"].(string))
	}

	s := &Status{}
	s.State = st

	if st != StateStopped {
		s.Playlist = m["playlist-name"].(string)
		s.PlaylistLen = m["playlist-length"].(int)
		s.PlaylistPos = m["playlist-position"].(int)
		s.Track = m["track-path"].(string)
		s.TrackLen = Time(m["track-length"].(int))
		s.TrackPos = Time(m["track-position"].(int))
	}

	return s, nil
}

func (c *Chubby) Stop() error {
	_, err := c.cmd(CmdStop)

	return err
}

func (c *Chubby) cmd(name string, args ...interface{}) ([]string, error) {
	buf := name
	for _, arg := range args {
		buf += fmt.Sprintf(" %#v", arg)
	}

	_, err := c.conn.WriteLine(buf)
	if err != nil {
		return nil, err
	}
	c.conn.Flush()
	line, err := c.conn.ReadLine()
	if err != nil {
		return nil, err
	}
	if err := parseRespStatus(line); err != nil {
		return nil, err
	}
	lines := make([]string, 0, 8)
	for {
		line, err := c.conn.ReadLine()
		if err != nil {
			return nil, err
		}
		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func parseRespStatus(line string) error {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) == 0 {
		return errors.New("invalid server response")
	}
	if parts[0] == "OK" {
		return nil
	} else if parts[0] == "ERR" {
		if len(parts) == 2 {
			return fmt.Errorf("server error: %s", parts[1])
		}
	}

	return errors.New("invalid server response")
}

func parseEntry(s string) (Entry, error) {
	m, err := parser.Parse(s)
	if err != nil {
		return nil, err
	}

	if tp, ok := m["type"].(string); ok && tp == "dir" {
		return &Dir{Path: m["path"].(string),
				Name: m["name"].(string)},
			nil
	} else {
		return &Track{Path: m["path"].(string),
				Artist: m["artist"].(string),
				Album:  m["album"].(string),
				Title:  m["title"].(string),
				Number: m["number"].(int),
				Length: Time(m["length"].(int))},
			nil
	}
}

func parsePlaylist(s string) (*Playlist, error) {
	m, err := parser.Parse(s)
	if err != nil {
		return nil, err
	}

	return &Playlist{
		Name:     m["name"].(string),
		Duration: Time(m["duration"].(int)),
		Length:   m["length"].(int),
	}, nil
}
