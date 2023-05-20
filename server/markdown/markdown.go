package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

type Element struct {
	Type     string    `json:"type"`
	Content  string    `json:"content"`
	Url      string    `json:"url"`
	Children []Element `json:"children"`
}

// parseMarkdown parses the Markdown string into a list of elements
func ParseMarkdown(markdown string) []Element {
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
		} else if match := regexp.MustCompile(`!\[\]\(([^)]+)\)`).FindStringSubmatch(trimmedLine); match != nil {
			imageURL := match[1]
			element := Element{Type: "img", Content: "", Url: imageURL}
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
	for _, ele := range elements {
		if ele.Type == "ol-li" {
			ol = append(ol, Element{
				Type:    "li",
				Content: ele.Content,
			})
		} else if ele.Type == "ul-li" {
			ul = append(ul, Element{
				Type:    "li",
				Content: ele.Content,
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
			ans = append(ans, ele)
		}
	}
	if len(ol) != 0 {
		ans = append(ans, Element{
			Type:     "ol",
			Children: ol,
		})
		ol = nil
	}
	if len(ul) != 0 {
		ans = append(ans, Element{
			Type:     "ul",
			Children: ul,
		})
		ul = nil
	}
	return ans
}

func ParseMemos(memos []string) [][]Element {
	var ans [][]Element
	for _, m := range memos {
		ans = append(ans, ParseMarkdown(m))
	}
	return ans
}
