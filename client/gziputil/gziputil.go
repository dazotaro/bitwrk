//  BitWrk - A Bitcoin-friendly, anonymous marketplace for computing power
//  Copyright (C) 2013-2019  Jonas Eschenburg <jonas@bitwrk.net>
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package gziputil enables transparent GZIP compression in HTTP connections.
package gziputil

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
)

// Given a handler, returns a handler with transparent support for receiving gzip-compressed POST data
func WithCompression(handle func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.Header.Get("Content-Encoding") == "gzip" {
			log.Printf("Handling GZIP-compressed POST.\n")
			if gz, err := newGZIPBody(r.Body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				// copy request data, substitute body
				r2 := *r
				r2.Body = gz
				// Call original handler
				handle(w, &r2)
			}
		} else {
			handle(w, r)
		}
	}
}

type gzipBody struct {
	compressed, uncompressed io.ReadCloser
}

func newGZIPBody(compressed io.ReadCloser) (*gzipBody, error) {
	if gz, err := gzip.NewReader(compressed); err != nil {
		return nil, err
	} else {
		return &gzipBody{compressed, gz}, nil
	}
}

func (gz *gzipBody) Read(data []byte) (int, error) {
	return gz.uncompressed.Read(data)
}

func (gz *gzipBody) Close() error {
	if err := gz.uncompressed.Close(); err != nil {
		gz.compressed.Close()
		return err
	} else if err := gz.compressed.Close(); err != nil {
		return err
	} else {
		return nil
	}
}

// When streaming chunk data in compressed mode, we must make sure that the gzip stream is
// flushed every couple of writes in order to avoid deadlocking.
type flushingCompressor struct {
	w   *gzip.Writer
	n   int
	err error
}

// Flush after every write. This wastes some compression but prevents deadlocks.
func (c *flushingCompressor) Write(buf []byte) (int, error) {
	if c.err != nil {
		return 0, c.err
	} else if written, err := c.w.Write(buf); err != nil {
		c.err = err
		return written, err
	} else {
		c.err = c.w.Flush()
		return written, nil
	}
}

func (c *flushingCompressor) Close() error {
	return c.w.Close()
}

// Function NewFlushingCompressor returns a WriteCloser that compresses as gzip and flushes on every Write.
func NewFlushingCompressor(w io.Writer) io.WriteCloser {
	c := gzip.NewWriter(w)
	return &flushingCompressor{w: c}
}

type nopCompressor struct {
	w io.Writer
}

func (c nopCompressor) Write(p []byte) (n int, err error) {
	return c.Write(p)
}

func (c nopCompressor) Close() error {
	return nil
}

// Function NewNopCompressor returns an implementation of WriteCloser that just passes writes to the
// underlying Writer and ignores Close.
func NewNopCompressor(w io.Writer) io.WriteCloser {
	return nopCompressor{w}
}
