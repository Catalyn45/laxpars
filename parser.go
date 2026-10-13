package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Parser struct {
	lexer   *Lexer
	current *Token
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

func (this *Parser) Next() {
	this.current = this.lexer.Next()
	fmt.Println("current token: ", this.current)
}

func (this *Parser) consume(tokenType string) string {
	if this.current.Type != tokenType {
		panic(fmt.Errorf("Unexpected token %s, expected %s", this.current, tokenType))
	}

	token := this.current

	this.Next()

	return token.Value
}

func (this *Parser) endTag() {
	if this.current.Type != TOKEN_TYPE_TAG_CLOSE {
		panic(fmt.Errorf("Unexpected token %s, expected %s", this.current, TOKEN_TYPE_TAG_CLOSE))
	}
}

func (this *Parser) parseStartTag() *Node {
	this.consume(TOKEN_TYPE_START_TAG_OPEN)
	tag := strings.ToLower(this.consume(TOKEN_TYPE_NAME))

	node := &Node {
		Type: NODE_TYPE_ELEMENT,
		Value: tag,
		Attributes: make(map[string]string),
	}

	for {
		if this.current.Type == TOKEN_TYPE_TAG_CLOSE {
			break
		}

		attributeName := strings.ToLower(this.consume(TOKEN_TYPE_NAME))

		attributeValue := ""
		if this.current.Type == TOKEN_TYPE_EQUAL {
			this.Next()

			if this.current.Type != TOKEN_TYPE_STRING && this.current.Type != TOKEN_TYPE_NAME {
				panic(fmt.Errorf("string or name expected for attribute value"))
			}

			attributeValue = this.current.Value
			this.Next()
		}

		node.Attributes[attributeName] = attributeValue
	}

	this.endTag()

	return node
}

func (this *Parser) parseEndTag(tagName string) {
	this.consume(TOKEN_TYPE_END_TAG_OPEN)

	name := strings.ToLower(this.consume(TOKEN_TYPE_NAME))

	if name != tagName {
		panic(fmt.Errorf("unmatched closing tag, expected %s, got %s", tagName, name))
	}

	this.endTag()
}

var voidTags = []string{ 
	"area",
	"base",
	"br",
	"col",
	"embed",
	"hr",
	"img",
	"input",
	"link",
	"meta",
	"source",
	"track",
	"wbr",
}

var escapableRawTextTags = []string{
	"textarea",
	"title",
}

var rawTextTags = []string{
	"script",
	"style",
}

func isVoidTag(tag string) bool {
	return slices.Contains(voidTags, tag)
}

func (this *Parser) isRawTextTag(tag string) bool {
	if slices.Contains(rawTextTags, tag) || slices.Contains(escapableRawTextTags, tag) {
		return true
	}
	
	return false
}

func (this *Parser) parseElement() *Node {
	if this.current.Type == TOKEN_TYPE_COMMENT {
		return &Node{
			Type: TOKEN_TYPE_COMMENT,
			Value: this.current.Value,
		}
	}

	tag := this.parseStartTag()
	if isVoidTag(tag.Value) {
		return tag
	}

	if this.isRawTextTag(tag.Value) {
		text := this.lexer.GetRawText(tag.Value)
		if text != "" {
			tag.AddChild(newTextNode(text))
		}

		return tag
	}

	for {
		text := this.lexer.GetText()
		if text != "" {
			tag.AddChild(newTextNode(text))
		}
		this.Next()

		if this.current.Type != TOKEN_TYPE_START_TAG_OPEN && this.current.Type != TOKEN_TYPE_COMMENT {
			break
		}

		child := this.parseElement()
		tag.AddChild(child)
	}

	this.parseEndTag(tag.Value)

	return tag
}

const (
	NODE_TYPE_ROOT = "ROOT"
	NODE_TYPE_ELEMENT = "ELEMENT"
	NODE_TYPE_TEXT = "TEXT"
	NODE_TYPE_COMMENT = "COMMENT"
	NODE_TYPE_DECLARATION = "DECLARATION"
)

type Node struct {
	Type string
	Value string
	Attributes map[string]string
	Children []*Node
	Parent *Node
}

func newTextNode(text string) *Node {
	return &Node {
		Type: NODE_TYPE_TEXT,
		Value: text,
	}
}

func (this *Node) AddChild(child *Node) {
	this.Children = append(this.Children, child)
	child.Parent = this
}

func (this *Node) ToHtml() string {
	if this.Type == NODE_TYPE_TEXT {
		return this.Value
	}

	if this.Type == NODE_TYPE_COMMENT {
		return "<!--" + this.Value + "-->"
	}

	if this.Type == NODE_TYPE_DECLARATION {
		return "<!" + this.Value + ">"
	}

	repr := ""
	if this.Type == NODE_TYPE_ELEMENT {
		repr += "<" + this.Value;

		for name, value := range this.Attributes {
			repr += " " + name

			if value != "" {
				repr += "=" + value
			}
		}

		repr += ">"
	}

	for _, child := range this.Children {
		repr += child.ToHtml()
	}

	if this.Type == NODE_TYPE_ELEMENT && !isVoidTag(this.Value) {
		repr += "</" + this.Value + ">"
	}

	return repr
}

func (this *Node) String() string {
	return this.StringDepth(0)
}

func (this *Node) StringDepth(depth int) string {
	repr := ""
	for i := 0; i < depth + 4; i++ {
		repr += " "
	}

	repr += fmt.Sprintln("{Type:", this.Type, "Value:", strconv.Quote(this.Value), "Attributes", this.Attributes, "}")

	for _, child := range this.Children {
		repr += child.StringDepth(depth + 4)
	}

	return repr
}

func (this *Parser) Parse() {
	this.Next()

	root := Node {
		Type: NODE_TYPE_ROOT,
	}

	if this.current.Type == TOKEN_TYPE_DECLARATION {
		root.AddChild(&Node {
			Type: NODE_TYPE_DECLARATION,
			Value: strings.ToLower(this.current.Value),
		})

		this.Next()
	}

	for ; this.current.Type != TOKEN_TYPE_EOF; this.Next() {
		element := this.parseElement()
		root.AddChild(element)
	}

	fmt.Println(root.String())
	fmt.Println(root.ToHtml())
}