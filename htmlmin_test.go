package htmlmin

import (
	"fmt"
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
   <footer>
      Copyright &copy;    Decent    People
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
<footer>
Copyright &copy; Decent People
</footer>
<script>
     (function()
       alert("Please   leave   this   unchanged! Thanks");
     )();
   </script>
</body>
</html>`

func TestMinify(t *testing.T) {
	result, err := Minify([]byte(doc))
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if string(result) != miniDoc {
		t.Errorf("incorrect result of minifying")
	}
	fmt.Printf("---\n%s\n---\n", result)
}
