package htmlmin

import (
	"io/ioutil"
	"strings"
	"testing"
)

var doc = `<!doctype html>
<!--[if lt IE 7]><html class="ie6"><![endif]-->
<head>
  <meta charset="utf-8">
  <title>Sample document</title>
  <style>
  body {
	  color: #cccccc;
  }
  </style>
</head>
<body>
   <!-- This is a comment -->
   <p  class="quoted value" data-something="x">
     Hello, this is a <b>document</b>.<br/>About
     something.
   </p>
   <pre><code>
   Hello <b>world!</b>

   Nice.
   </code></pre>
   <img alt="" width="100" style="color: #aaaaaa; padding: 0px;">
   <footer>
	   Copyright &copy;    <A HREF="http://www.example.com/?q=1&amp;m=2">Decent</A>    People
   </footer>
   <script>
     (function()
       alert("Please   leave   this   unchanged! Thanks");
     )();
   </script>
</body>
</html>`

var miniDoc = `<!doctype html>
<!--[if lt IE 7]><html class="ie6"><![endif]-->
<head>
<meta charset=utf-8>
<title>Sample document</title>
<style>
  body {
	  color: #cccccc;
  }
  </style>
</head>
<body>

<p class="quoted value" data-something=x>
Hello, this is a <b>document</b>.<br>About
something.
</p>
<pre><code>
   Hello <b>world!</b>

   Nice.
   </code></pre>
<img alt="" width=100 style="color: #aaaaaa; padding: 0px;">
<footer>
Copyright &copy; <a href="http://www.example.com/?q=1&amp;m=2">Decent</a> People
</footer>
<script>
     (function()
       alert("Please   leave   this   unchanged! Thanks");
     )();
   </script>
</body>
</html>`

var miniDocFull = `<!doctype html>
<!--[if lt IE 7]><html class="ie6"><![endif]-->
<head>
<meta charset=utf-8>
<title>Sample document</title>
<style>body{color:#ccc}</style>
</head>
<body>

<p class="quoted value" data-something=x>
Hello, this is a <b>document</b>.<br>About
something.
</p>
<pre><code>
   Hello <b>world!</b>

   Nice.
   </code></pre>
<img alt="" width=100 style=color:#aaa;padding:0>
<footer>
Copyright &copy; <a href="http://www.example.com/?q=1&amp;m=2">Decent</a> People
</footer>
<script>(function()
alert("Please   leave   this   unchanged! Thanks");)();</script>
</body>
</html>`

func TestMinify(t *testing.T) {
	result, err := Minify([]byte(doc), nil)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	ioutil.WriteFile("expected.txt", []byte(miniDoc), 0644)
	ioutil.WriteFile("result.txt", []byte(result), 0644)

	if string(result) != miniDoc {
		t.Errorf("Incorrect result of minifying #1")
		diffLines(t, miniDoc, string(result))
	}

	result, err = Minify([]byte(doc), FullOptions)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if string(result) != miniDocFull {
		t.Errorf("Incorrect result of minifying #2")
		diffLines(t, miniDocFull, string(result))
	}
}

func diffLines(t *testing.T, expected, got string) {
	explines := strings.Split(expected, "\n")
	gotlines := strings.Split(got, "\n")
	for ei, el := range explines {
		if len(gotlines) <= ei {
			t.Errorf("result is shorter than expected")
			return
		}
		gl := gotlines[ei]
		if el != gl {
			t.Errorf("lines differ:")
			t.Errorf("%d: %s", ei, el)
			t.Errorf("%d: %s", ei, gl)
		}
	}
}
