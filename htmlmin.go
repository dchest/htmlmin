// Copyright 2013 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package htmlmin minifies HTML.
package htmlmin

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
	"github.com/dchest/cssmin"
	"github.com/dchest/jsmin"
)

type Options struct {
	MinifyScripts bool // if true, use jsmin to minify contents of script tags.
	MinifyStyles  bool // if true, use cssmin to minify contents of style tags and inline styles.
	UnquoteAttrs  bool // if true, remove quotes from HTML attributes where possible.
}

var DefaultOptions = &Options{
	MinifyScripts: false,
	MinifyStyles:  false,
	UnquoteAttrs:  false,
}

// Minify returns minified version of the given HTML data.
// If passed options is nil, uses default options.
func Minify(data []byte, options *Options) (out []byte, err error) {
	if options == nil {
		options = DefaultOptions
	}
	var b bytes.Buffer
	z := html.NewTokenizer(bytes.NewReader(data))
	raw := 0
	javascript := false
	style := false
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
			switch string(tagName) {
			case "script":
				javascript = true
				raw++
			case "style":
				style = true
				raw++
			case "pre", "code", "textarea":
				raw++
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
				if string(k) == "style" && options.MinifyStyles {
					v = []byte("a{" + string(v) + "}") // simulate "full" CSS
					v = cssmin.Minify(v)
					v = v[2 : len(v)-1] // strip simulation
				}
				if isFirst {
					b.WriteByte(' ')
					isFirst = false
				}
				b.Write(k)
				b.WriteByte('=')
				qv := html.EscapeString(string(v))
				if !options.UnquoteAttrs || shouldQuote(v) {
					// Quoted value.
					b.WriteByte('"')
					b.WriteString(qv)
					b.WriteByte('"')
				} else {
					// Unquoted value.
					b.WriteString(qv)
				}
				if hasAttr {
					b.WriteByte(' ')
				}
			}
			b.WriteByte('>')
		case html.EndTagToken:
			tagName, _ := z.TagName()
			switch string(tagName) {
			case "script":
				javascript = false
				raw--
			case "style":
				style = false
				raw--
			case "pre", "code", "textarea":
				raw--
			}
			b.Write([]byte("</"))
			b.Write(tagName)
			b.WriteByte('>')
		case html.CommentToken:
			if bytes.HasPrefix(z.Raw(), []byte("<!--[if")) ||
				bytes.HasPrefix(z.Raw(), []byte("<!--//")) {
				// Preserve IE conditional and special style comments.
				b.Write(z.Raw())
			}
			// ... otherwise, skip.
		case html.TextToken:
			if javascript && options.MinifyScripts {
				min, err := jsmin.Minify(z.Raw())
				if err != nil {
					// Just write it as is.
					b.Write(z.Raw())
				} else {
					b.Write(min)
				}
			} else if style && options.MinifyStyles {
				b.Write(cssmin.Minify(z.Raw()))
			} else if raw > 0 {
				b.Write(z.Raw())
			} else {
				b.Write(trimTextToken(z.Raw()))
			}
		default:
			b.Write(z.Raw())
		}

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

func shouldQuote(b []byte) bool {
	if len(b) == 0 || bytes.IndexAny(b, "\"'`=<> \n\r\t\b") != -1 {
		return true
	}
	return false
}
