package htmlmin

import (
	"testing"
)

var doc = `<!doctype html>
<head>
  <meta charset="utf-8">
  <title>Sample document</title>
</head>
<body>
   <!-- This is a comment -->
   <p  class="quoted value" data-something="x">
     Hello, this is a <b>document</b>.<br/>About
     something.
   </p>
   <img alt="" width="100">
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
<head>
<meta charset=utf-8>
<title>Sample document</title>
</head>
<body>

<p class="quoted value" data-something=x>
Hello, this is a <b>document</b>.<br>About
something.
</p>
<img alt="" width=100>
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

var miniDocJS = `<!doctype html>
<head>
<meta charset=utf-8>
<title>Sample document</title>
</head>
<body>

<p class="quoted value" data-something=x>
Hello, this is a <b>document</b>.<br>About
something.
</p>
<img alt="" width=100>
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
	if string(result) != miniDoc {
		t.Errorf("1. incorrect result of minifying\n---\n%s\n---\n", result)
	}

	result, err = Minify([]byte(doc), &Options{MinifyScripts: true})
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if string(result) != miniDocJS {
		t.Errorf("2. incorrect result of minifying\n---\n%s\n---\n", result)
	}
}
