// Copyright 2017-2023 Viacheslav Chimishuk <vchimishuk@yandex.ru>
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
	"io"
	"net"
	"strings"

	"github.com/vchimishuk/chubby/parser"
	"github.com/vchimishuk/chubby/textconn"
	"github.com/vchimishuk/chubby/time"
)

const (
	cmdCreatePlaylist = "create-playlist"
	cmdDeletePlaylist = "delete-playlist"
	cmdEvents         = "events"
	cmdKill           = "kill"
	cmdList           = "list"
	cmdNext           = "next"
	cmdPause          = "pause"
	cmdPing           = "ping"
	cmdPlay           = "play"
	cmdPlaylists      = "playlists"
	cmdPrev           = "prev"
	cmdRenamePlaylist = "rename-playlist"
	cmdSeek           = "seek"
	cmdStatus         = "status"
	cmdStop           = "stop"
)

type SeekMode int

const (
	SeekModeAbs     SeekMode = 0
	SeekModeForward SeekMode = 1
	SeekModeRewind  SeekMode = -1
)

const (
	eventsChSize = 10
)

type Playlist struct {
	Name     string
	Duration time.Time
	Length   int
}

type Status struct {
	State       State
	PlaylistPos int
	TrackPos    time.Time
	Playlist    *Playlist
	Track       *Track
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
	Year   int
	Title  string
	Number int
	Length time.Time
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
	conn   *textconn.TextConn
	resps  chan []string
	events chan Event
	err    chan error
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

	c.resps = make(chan []string, 1)
	c.events = make(chan Event, eventsChSize)
	c.err = make(chan error, 1)

	go c.read()

	return nil
}

func (c *Chubby) Close() error {
	if c.conn == nil {
		return errors.New("not connected")
	}
	err := c.conn.Close()
	c.conn = nil

	// Wait for read() goroutine to exit.
	<-c.err

	return err
}

func (c *Chubby) CreatePlaylist(name string) error {
	_, err := c.cmd(cmdCreatePlaylist, name)

	return err
}

func (c *Chubby) DeletePlaylist(name string) error {
	_, err := c.cmd(cmdDeletePlaylist, name)

	return err
}

func (c *Chubby) Events(enable bool) (<-chan Event, error) {
	_, err := c.cmd(cmdEvents, enable)

	return c.events, err
}

func (c *Chubby) Kill() error {
	_, err := c.cmd(cmdKill)

	return err
}

func (c *Chubby) List(path string) ([]Entry, error) {
	lines, err := c.cmd(cmdList, path)
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
	_, err := c.cmd(cmdNext)

	return err
}

func (c *Chubby) Pause() error {
	_, err := c.cmd(cmdPause)

	return err
}

func (c *Chubby) Ping() error {
	_, err := c.cmd(cmdPing)

	return err
}

func (c *Chubby) Play(pth string) error {
	_, err := c.cmd(cmdPlay, pth)

	return err
}

func (c *Chubby) Playlists() ([]*Playlist, error) {
	lines, err := c.cmd(cmdPlaylists)
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
	_, err := c.cmd(cmdPrev)

	return err
}

func (c *Chubby) RenamePlaylist(from, to string) error {
	_, err := c.cmd(cmdRenamePlaylist, from, to)

	return err
}

func (c *Chubby) Seek(time time.Time, mode SeekMode) error {
	var t int
	var rel bool

	switch mode {
	case SeekModeAbs:
		t = int(time)
		rel = false
	case SeekModeForward:
		t = int(time)
		rel = true
	case SeekModeRewind:
		t = -int(time)
		rel = true
	default:
		panic("unsupported SeekMode")
	}

	_, err := c.cmd(cmdSeek, t, rel)

	return err
}

func (c *Chubby) Status() (*Status, error) {
	lines, err := c.cmd(cmdStatus)

	if len(lines) != 1 {
		return nil, err
	}

	m, err := parser.Parse(lines[0])
	if err != nil {
		return nil, err
	}

	st, err := parseState(m["state"].(string))
	if err != nil {
		return nil, err
	}

	s := &Status{State: st, Playlist: &Playlist{}, Track: &Track{}}

	if st != StateStopped {
		s.PlaylistPos = m["playlist-position"].(int)
		s.TrackPos = time.Time(m["track-position"].(int))
		s.Playlist.Name = m["playlist-name"].(string)
		s.Playlist.Duration = time.Time(m["playlist-duration"].(int))
		s.Playlist.Length = m["playlist-length"].(int)
		s.Track.Path = m["track-path"].(string)
		s.Track.Artist = m["track-artist"].(string)
		s.Track.Album = m["track-album"].(string)
		s.Track.Year = m["track-year"].(int)
		s.Track.Title = m["track-title"].(string)
		s.Track.Number = m["track-number"].(int)
		s.Track.Length = time.Time(m["track-length"].(int))
	}

	return s, nil
}

func (c *Chubby) Stop() error {
	_, err := c.cmd(cmdStop)

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

	select {
	case r := <-c.resps:
		return r, nil
	case err = <-c.err:
		return nil, err
	}
}

func (c *Chubby) read() {
	var err error

	for {
		var event string
		var resp []string
		var nerr net.Error
		event, resp, err = c.readResp()
		if err != nil {
			if errors.As(err, &nerr) || errors.Is(err, io.EOF) {
				break
			} else {
				c.err <- err
			}
		} else if event != "" {
			if len(c.events) < eventsChSize {
				var e Event
				e, err = parseEvent(event, resp)
				// Ignore invalid/unknown events.
				if err == nil {
					c.events <- e
				}
			}
		} else {
			c.resps <- resp
		}
	}

	c.conn.Close()
	c.err <- err
	close(c.events)
	close(c.resps)
	close(c.err)
}

func (c *Chubby) readResp() (string, []string, error) {
	line, err := c.conn.ReadLine()
	if err != nil {
		return "", nil, err
	}

	event := ""
	pts := strings.SplitN(line, " ", 2)
	if pts[0] == "OK" {
		// Do nothing.
	} else if pts[0] == "EVENT" {
		if len(pts) != 2 {
			return "", nil,
				fmt.Errorf("protocol: invalid header")
		}
		event = pts[1]
	} else if pts[0] == "ERR" {
		if len(pts) != 2 {
			return "", nil,
				fmt.Errorf("protocol: invalid header")
		}
		return "", nil, newServerError(pts[1])
	} else {
		return "", nil, fmt.Errorf("protocol: invalid header")
	}

	lines := make([]string, 0, 8)
	for {
		line, err := c.conn.ReadLine()
		if err != nil {
			return "", nil, err
		}
		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}

	return event, lines, nil
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
				Year:   m["year"].(int),
				Title:  m["title"].(string),
				Number: m["number"].(int),
				Length: time.Time(m["length"].(int))},
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
		Duration: time.Time(m["duration"].(int)),
		Length:   m["length"].(int),
	}, nil
}
