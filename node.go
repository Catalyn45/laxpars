package main

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

const (
	NODE_TYPE_ROOT = "ROOT"
	NODE_TYPE_ELEMENT = "ELEMENT"
	NODE_TYPE_TEXT = "TEXT"
	NODE_TYPE_COMMENT = "COMMENT"
	NODE_TYPE_DECLARATION = "DECLARATION"
)

type Node struct {
	Type       string
	Value      string
	Attributes map[string]string
	Children   []*Node
	Parent     *Node
}

func newTextNode(text string) *Node {
	return &Node{
		Type:  NODE_TYPE_TEXT,
		Value: text,
	}
}

func (this *Node) AppendChild(child *Node) {
	this.Children = append(this.Children, child)
	child.Parent = this
}

func (this *Node) RemoveChild(child *Node) {
	for index, childNode := range this.Children {
		if childNode == child {
			this.Children = append(this.Children[0:index], this.Children[:index+1]...)
			return
		}
	}
}

func (this *Node) Remove() {
	this.Parent.RemoveChild(this)
}

func (this *Node) Before(node *Node) {
	this.Children = append([]*Node{node}, this.Children...)
}

func (this *Node) After(node *Node) {
	this.Children = append(this.Children, node)
}

func (this *Node) ReplaceChild(newChild *Node, oldChild *Node) {
	for index, childNode := range this.Children {
		if childNode == oldChild {
			this.Children[index] = newChild
			return
		}
	}
}

func (this *Node) ReplaceWith(newChild *Node) {
	this.Parent.ReplaceChild(newChild, this)
}

func (this *Node) CloneNode(deepClone bool) *Node {
	clone := *this

	clone.Attributes = maps.Clone(clone.Attributes)
	clone.Children = slices.Clone(clone.Children)

	if deepClone {
		for index, child := range clone.Children {
			clone.Children[index] = child.CloneNode(deepClone)
		}
	}

	return &clone
}

func (this *Node) GetElementById(id string) *Node {
	idValue, ok := this.Attributes["id"]
	if ok {
		idValue = strings.Trim(idValue, "'\"")
		if idValue == id {
			return this
		}
	}

	for _, child := range this.Children {
		node := child.GetElementById(id)
		if node != nil {
			return node
		}
	}

	return nil
}

func (this *Node) GetElementsByClassName(className string) []*Node {
	var nodes []*Node

	classNameValue, ok := this.Attributes["class"]
	if ok {
		classNameValue = strings.Trim(classNameValue, "'\"")
		if classNameValue == className {
			nodes = append(nodes, this)
		}
	}

	for _, child := range this.Children {
		childClassNodes := child.GetElementsByClassName(className)
		if childClassNodes != nil {
			nodes = append(nodes, childClassNodes...)
		}
	}

	return nodes
}

func (this *Node) GetElementsByTagName(tagName string) []*Node {
	if this.Type != NODE_TYPE_ELEMENT && this.Type != NODE_TYPE_ROOT {
		return nil
	}

	var nodes []*Node
	if this.Value == tagName {
		nodes = append(nodes, this)
	}

	for _, child := range this.Children {
		childTagNodes := child.GetElementsByTagName(tagName)
		if childTagNodes != nil {
			nodes = append(nodes, childTagNodes...)
		}
	}

	return nodes
}

func (this *Node) TextContent() string {
	content := ""
	for _, child := range this.Children {
		if  child.Type == NODE_TYPE_TEXT {
			content += child.Value
		}
	}

	return content
}

func (this *Node) StringHtml() string {
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
		repr += "<" + this.Value

		for name, value := range this.Attributes {
			repr += " " + name

			if value != "" {
				repr += "=" + value
			}
		}

		repr += ">"
	}

	for _, child := range this.Children {
		repr += child.StringHtml()
	}

	if this.Type == NODE_TYPE_ELEMENT && !isVoidTag(this.Value) {
		repr += "</" + this.Value + ">"
	}

	return repr
}

func (this *Node) String() string {
	return fmt.Sprint("{ Type: ", this.Type, " Value: ", strconv.Quote(this.Value), " Attributes: ", this.Attributes, " }")
}

func (this *Node) stringTreeDepth(depth int) string {
	repr := ""
	for i := 0; i < depth; i++ {
		repr += " "
	}

	repr += fmt.Sprintln(this)

	for _, child := range this.Children {
		repr += child.stringTreeDepth(depth + 4)
	}

	return repr
}

func (this *Node) StringTree() string {
	return this.stringTreeDepth(0)
}

