package htmlquery

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

const htmlSample = `<!DOCTYPE html><html lang="en-US">
<head>
<title>Hello,World!</title>
</head>
<body>
<div class="container">
<header>
	<!-- Logo -->
   <h1>City Gallery</h1>
</header>  
<nav>
  <ul>
    <li><a href="#">London</a></li>
    <li><a href="#">Paris</a></li>
    <li><a href="#">Tokyo</a></li>
  </ul>
</nav>
<article>
  <h1>London</h1>
  <img src="pic_mountain.jpg" alt="Mountain View" style="width:304px;height:228px;">
  <p>London is the capital city of England. It is the most populous city in the  United Kingdom, with a metropolitan area of over 13 million inhabitants.</p>
  <p>Standing on the River Thames, London has been a major settlement for two millennia, its history going back to its founding by the Romans, who named it Londinium.</p>
</article>
<footer>Copyright &copy; W3Schools.com</footer>
</div>
</body>
</html>
`

var testDoc, _ = loadHTML(htmlSample)

func TestHttpLoad(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	doc, err := LoadURL(ctx, "http://www.bing.com")
	if err != nil {
		t.Fatal(err)
	}
	if doc == nil {
		t.Fatal("doc is nil")
	}
}

func TestNavigator(t *testing.T) {
	top, _ := FindOne(testDoc, "//html")
	nav := &NodeNavigator{curr: top, root: top, attr: -1}
	nav.MoveToChild() // HEAD
	nav.MoveToNext()
	if nav.NodeType() != xpath.TextNode {
		t.Fatalf("expectd node type is TextNode,but got %v ", nav.NodeType())
	}
	nav.MoveToNext() // <BODY>
	bodyNode, _ := FindOne(testDoc, "//body")
	if nav.Value() != InnerText(bodyNode) {
		t.Fatal("body not equal")
	}
	nav.MoveToPrevious() //
	nav.MoveToParent()   //<HTML>
	if nav.curr != top {
		t.Fatal("current node is not html node")
	}
	nav.MoveToNextAttribute()
	if nav.LocalName() != "lang" {
		t.Fatal("node not move to lang attribute")
	}

	nav.MoveToParent()
	nav.MoveToFirst() // <!DOCTYPE html>
	if nav.curr.Type != html.DoctypeNode {
		t.Fatalf("expected node type is DoctypeNode,but got %d", nav.curr.Type)
	}
}

func TestXPath(t *testing.T) {
	node, _ := FindOne(testDoc, "//html")
	if SelectAttr(node, "lang") != "en-US" {
		t.Fatal("//html[@lang] != en-Us")
	}

	var c int
	FindEach(testDoc, "//li", func(i int, node *html.Node) {
		c++
	})
	length, _ := Find(testDoc, "//li")
	if c != len(length) {
		t.Fatal("li node count != 3")
	}
	node, _ = FindOne(testDoc, "//header")
	if strings.Index(InnerText(node), "Logo") > 0 {
		t.Fatal("InnerText() have comment node text")
	}
	if strings.Index(OutputHTML(node, true), "Logo") == -1 {
		t.Fatal("OutputHTML() shoud have comment node text")
	}
}

func TestXPathCdUp(t *testing.T) {
	doc, _ := loadHTML(`<html><b attr="1"></b></html>`)
	node, _ := FindOne(doc, "//b/@attr/..")
	t.Logf("node = %#v", node)
	if node == nil || node.Data != "b" {
		t.Fatal("//b/@id/.. != <b></b>")
	}
}

func loadHTML(str string) (*html.Node, error) {
	node, err := Parse(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	return node, nil
}
