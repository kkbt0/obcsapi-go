package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Element represents a parsed Markdown element
type Element struct {
	Type     string
	Content  string
	Url      string
	Children []Element
}

// parseMarkdown parses the Markdown string into a list of elements
func parseMarkdown(markdown string) []Element {
	lines := strings.Split(markdown, "\n")
	elements := []Element{}
	stack := []*Element{}

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "#") {
			level := strings.Count(trimmedLine, "#")
			title := strings.TrimSpace(trimmedLine[level:])
			element := Element{Type: fmt.Sprintf("h%d", level), Content: title}
			elements = append(elements, element)
		} else if match := regexp.MustCompile(`(^- \[.\].*)`).FindStringSubmatch(trimmedLine); match != nil {
			text := strings.TrimSpace(match[1]) // ul
			element := Element{Type: "check", Content: text}
			elements = append(elements, element)
		} else if match := regexp.MustCompile(`^- (.*)`).FindStringSubmatch(trimmedLine); match != nil {
			text := strings.TrimSpace(match[1]) // ul
			element := Element{Type: "ul-li", Content: text}
			elements = append(elements, element)
		} else if match := regexp.MustCompile(`^\d+\.(.*)`).FindStringSubmatch(trimmedLine); match != nil {
			text := strings.TrimSpace(match[1]) // ol
			element := Element{Type: "ol-li", Content: text}
			if len(stack) > 0 {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, element)
			} else {
				elements = append(elements, element)
			}
		} else if match := regexp.MustCompile(`!\[([^\]]+)\]\(([^)]+)\)`).FindStringSubmatch(trimmedLine); match != nil {
			altText := match[1]
			imageURL := match[2]
			element := Element{Type: "img", Content: altText, Url: imageURL}
			elements = append(elements, element)
		} else if match := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`).FindStringSubmatch(trimmedLine); match != nil {
			altText := match[1]
			imageURL := match[2]
			element := Element{Type: "url", Content: altText, Url: imageURL}
			elements = append(elements, element)

		} else if trimmedLine == "" {
			// elements = append(elements, Element{Type: "br"})
		} else if trimmedLine == "" && len(stack) > 0 {
			stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, Element{Type: "br"})
		} else {
			element := Element{Type: "p", Content: trimmedLine}
			if len(stack) > 0 {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, element)
			} else {
				elements = append(elements, element)
			}
		}
	}
	// ol ul children
	var ol []Element
	var ul []Element
	var ans []Element
	for _, element := range elements {
		if element.Type == "ol-li" {
			ol = append(ol, Element{
				Type:    "li",
				Content: element.Content,
			})
		} else if element.Type == "ul-li" {
			ul = append(ul, Element{
				Type:    "li",
				Content: element.Content,
			})

		} else if len(ol) != 0 {
			ans = append(ans, Element{
				Type:     "ol",
				Children: ol,
			})
			ol = nil
		} else if len(ul) != 0 {
			ans = append(ans, Element{
				Type:     "ul",
				Children: ul,
			})
			ul = nil
		} else {
			ans = append(ans, element)
		}
	}
	return ans
}

func main() {
	markdown := `
- 20:05 也就是20号
# Heading 1

Some text.

## Heading 2

1. Ordered list item 1
2. Ordered list item 2

- xxx
- xxx

 - xxx1
 - xxx2

Paragraph with **bold** and *italic* text.

![OpenAI Logo](https://openai.com/assets/openai-logo.svg)

[OpenAI Logo](https://openai.com/assets/openai-logo.svg)

- [ ] x1
- [x] x2
 - [ ] x3
 - [x] x4
`
	elements := parseMarkdown(markdown)

	for _, element := range elements {
		printElement(element, 0)
	}
}

func printElement(element Element, indent int) {
	indentStr := strings.Repeat("  ", indent)
	fmt.Printf("---\nType: %s\n", element.Type)
	fmt.Printf("%sContent: %s\n", indentStr, element.Content)
	fmt.Printf("%sChild: %s\n", indentStr, element.Children)
	if element.Type == "img" || element.Type == "url" {
		fmt.Printf("%sURL: %s\n", indentStr, element.Url)
	}
	// for _, child := range element.Children {
	// 	printElement(child, indent+1)
	// }
}
