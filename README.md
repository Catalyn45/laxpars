# LaxPars

LaxPars is a simple HTML parser written in `go` that allows parsing custom HTML elements. Basically it will parse
the document as-is as long as the tags are closed correctly.

## LaxPars will
- Parse any custom tag name
- Fail if any non-void element tag doesn't have the corresponding closing tag

## LaxPars will NOT
- Do any modification to the elements order and structure
- Validate if certain elements are allowed to be inside other elements

## Who is this for?

This parser is useful for doing preprocessing on html files known to be valid or creating templating frameworks.
In theory it can be used for HTML scrapping, but it will fail if the HTML is not valid.

It can also be used as a learning exercise since it's just a few `.go` files.

## Examples

parse a html returned by a website and find a certain element
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/Catalyn45/laxpars"
	"os"
)

func main() {
	resp, err := http.Get(os.Args[1])
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status:", resp.Status)
		return
	}
  
	parser := laxpars.NewParser(resp.Body)

	root := parser.Parse()

	el := root.GetElementById("id123")

	fmt.Println(el.StringHtml())
}
```

<hr>
  
Insert a new node somewhere in the document
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/Catalyn45/laxpars"
	"os"
)

func main() {
	resp, err := http.Get(os.Args[1])
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status:", resp.Status)
		return
	}
  
	parser := laxpars.NewParser(resp.Body)

	root := parser.Parse()

	el := root.GetElementById("id123")

	el.AppendChild(&laxpars.Node {
		Type: NODE_TYPE_ELEMENT,
		Value: "p",
		Children: []*laxpars.Node { laxpars.NewTextNode("Hello") },
	})

	fmt.Println(root.StringHtml())
}
```
