package laxpars

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type NodeVisitor interface {
	VisitNode(node *Node)
}

type Parser struct {
	lexer   *Lexer
	current *Token
	nodeVisitor NodeVisitor

	Root *Node
}

func NewParser(input io.Reader) *Parser {
	return &Parser {
		lexer: NewLexer(input),
	}
}

func (this *Parser) SetVisitor(nodeVisitor NodeVisitor) {
	this.nodeVisitor = nodeVisitor
}

func (this *Parser) next() {
	this.current = this.lexer.Next()
}

func (this *Parser) consume(tokenType string) string {
	if this.current.Type != tokenType {
		panic(fmt.Errorf("Unexpected token %s, expected %s", this.current, tokenType))
	}

	token := this.current

	this.next()

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
			this.next()

			if this.current.Type != TOKEN_TYPE_STRING && this.current.Type != TOKEN_TYPE_NAME {
				panic(fmt.Errorf("string or name expected for attribute value"))
			}

			attributeValue = this.current.Value
			this.next()
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

func (this *Parser) visitNode(node *Node) {
	if this.nodeVisitor != nil {
		this.nodeVisitor.VisitNode(node)
	}
}

func (this *Parser) parseElement() *Node {
	if this.current.Type == TOKEN_TYPE_COMMENT {
		return &Node{
			Type: TOKEN_TYPE_COMMENT,
			Value: this.current.Value,
		}
	}

	tag := this.parseStartTag()
	defer this.visitNode(tag)

	if isVoidTag(tag.Value) {
		return tag
	}

	if this.isRawTextTag(tag.Value) {
		text := this.lexer.GetRawText(tag.Value)
		if text != "" {
			tag.AppendChild(NewTextNode(text))
		}

		return tag
	}

	for {
		text := this.lexer.GetText()
		if text != "" {
			tag.AppendChild(NewTextNode(text))
		}
		this.next()

		if this.current.Type != TOKEN_TYPE_START_TAG_OPEN && this.current.Type != TOKEN_TYPE_COMMENT {
			break
		}

		child := this.parseElement()
		tag.AppendChild(child)
	}

	this.parseEndTag(tag.Value)

	return tag
}

func (this *Parser) Parse() *Node {
	if this.Root != nil {
		return this.Root
	}

	this.next()

	root := Node {
		Type: NODE_TYPE_ROOT,
	}

	for ; this.current.Type == TOKEN_TYPE_COMMENT; this.next() {
		root.AppendChild(&Node {
			Type: NODE_TYPE_COMMENT,
			Value: this.current.Value,
		})
	}

	if this.current.Type == TOKEN_TYPE_DECLARATION {
		root.AppendChild(&Node {
			Type: NODE_TYPE_DECLARATION,
			Value: strings.ToLower(this.current.Value),
		})

		this.next()
	}

	for ; this.current.Type != TOKEN_TYPE_EOF; this.next() {
		element := this.parseElement()
		root.AppendChild(element)
	}

	this.Root = &root

	return this.Root
}