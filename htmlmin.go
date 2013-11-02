// Copyright 2013 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package htmlmin minifies HTML.
package htmlmin

import (
	"bytes"
	"io"

	"code.google.com/p/go.net/html"
	"github.com/dchest/jsmin"
)

type Options struct {
	MinifyScripts bool // if true, use jsmin to minify contents of script tags.
}

var DefaultOptions = &Options{
	MinifyScripts: false,
}

// Minify returns minified version of the given HTML data.
// If passed options is nil, uses default options.
func Minify(data []byte, options *Options) (out []byte, err error) {
	if options == nil {
		options = DefaultOptions
	}
	var b bytes.Buffer
	z := html.NewTokenizer(bytes.NewReader(data))
	raw := false
	javascript := false
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				return b.Bytes(), nil
			}
			return nil, err
		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, hasAttr := z.TagName()
			raw = isRawTagName(tagName)
			if string(tagName) == "script" {
				javascript = true
			}
			b.WriteByte('<')
			b.Write(tagName)
			var k, v []byte
			isFirst := true
			for hasAttr {
				k, v, hasAttr = z.TagAttr()
				if javascript && string(k) == "type" && string(v) != "text/javascript" {
					javascript = false
				}
				if isFirst {
					b.WriteByte(' ')
					isFirst = false
				}
				b.Write(k)
				b.WriteByte('=')
				if quoteChar := valueQuoteChar(v); quoteChar != 0 {
					// Quoted value.
					b.WriteByte(quoteChar)
					b.WriteString(html.EscapeString(string(v)))
					b.WriteByte(quoteChar)
				} else {
					// Unquoted value.
					b.Write(v)
				}
				if hasAttr {
					b.WriteByte(' ')
				}
			}
			b.WriteByte('>')
		case html.EndTagToken:
			tagName, _ := z.TagName()
			raw = false
			if javascript && string(tagName) == "script" {
				javascript = false
			}
			b.Write([]byte("</"))
			b.Write(tagName)
			b.WriteByte('>')
		case html.CommentToken:
			// skip
		case html.TextToken:
			if javascript && options.MinifyScripts {
				min, err := jsmin.Minify(z.Raw())
				if err != nil {
					// Just write it as is.
					b.Write(z.Raw())
				} else {
					b.Write(min)
				}
			} else if raw {
				b.Write(z.Raw())
			} else {
				b.Write(trimTextToken(z.Raw()))
			}
		default:
			b.Write(z.Raw())
		}

	}
}

func isRawTagName(tagName []byte) bool {
	switch string(tagName) {
	case "script", "pre", "code", "textarea":
		return true
	default:
		return false
	}
}

func trimTextToken(b []byte) (out []byte) {
	out = make([]byte, 0)
	seenSpace := false
	for _, c := range b {
		switch c {
		case ' ', '\n', '\r', '\t':
			if !seenSpace {
				out = append(out, c)
				seenSpace = true
			}
		default:
			out = append(out, c)
			seenSpace = false
		}
	}
	return out
}

func valueQuoteChar(b []byte) byte {
	if len(b) == 0 || bytes.IndexAny(b, "'`=<> \n\r\t\b") != -1 {
		return '"' // quote with quote mark
	}
	if bytes.IndexByte(b, '"') != -1 {
		return '\'' // quote with apostrophe
	}
	return 0 // do not quote
}
