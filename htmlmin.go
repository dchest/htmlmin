// Copyright 2013 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package htmlmin minifies HTML.
package htmlmin

import (
	"bytes"
	"io"

	"code.google.com/p/go.net/html"
)

// Minify returns minified version of the given HTML data.
func Minify(data []byte) (out []byte, err error) {
	var b bytes.Buffer
	z := html.NewTokenizer(bytes.NewReader(data))
	raw := false
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
			b.WriteByte('<')
			b.Write(tagName)
			var k, v []byte
			isFirst := true
			for hasAttr {
				k, v, hasAttr = z.TagAttr()
				if isFirst {
					b.WriteByte(' ')
					isFirst = false
				}
				b.Write(k)
				if len(v) == 0 {
					// Empty attribute values can be skipped.
					continue
				}
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
			b.Write([]byte("</"))
			b.Write(tagName)
			b.WriteByte('>')
		case html.CommentToken:
			// skip
		case html.TextToken:
			if raw {
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
	if bytes.IndexAny(b, "'`=<> \n\r\t\b") != -1 {
		return '"' // quote with quote mark
	}
	if bytes.IndexByte(b, '"') != -1 {
		return '\'' // quote with apostrophe
	}
	return 0 // do not quote
}
