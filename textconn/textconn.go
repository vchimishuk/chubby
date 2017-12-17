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

package textconn

import (
	"bufio"
	"net"
	"net/textproto"
)

type TextConn struct {
	conn   net.Conn
	reader *textproto.Reader
	writer *bufio.Writer
}

func New(conn net.Conn) *TextConn {
	return &TextConn{
		conn:   conn,
		reader: textproto.NewReader(bufio.NewReader(conn)),
		writer: bufio.NewWriter(conn),
	}
}

func (c *TextConn) ReadLine() (string, error) {
	return c.reader.ReadLine()
}

func (c *TextConn) WriteLine(line string) (int, error) {
	n, err := c.writer.WriteString(line)
	if err != nil {
		return n, err
	}
	n, err = c.writer.WriteString("\n")

	return n, err
}

func (c *TextConn) Flush() error {
	return c.writer.Flush()
}

func (c *TextConn) Close() error {
	return c.conn.Close()
}
