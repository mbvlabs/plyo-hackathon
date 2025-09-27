package services

import (
	"bytes"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func newParser() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("gruvbox"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.TabWidth(4),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
}

func cleanMarkdownBlock(report string) string {
	report = strings.TrimSpace(report)

	// // Remove opening markdown block
	// if strings.HasPrefix(report, "```markdown\n") {
	// 	report = report[12:]
	// } else if strings.HasPrefix(report, "```\n") {
	// 	report = report[4:]
	// }
	//
	// // Remove closing markdown block
	// if strings.HasSuffix(report, "\n```") {
	// 	report = report[:len(report)-4]
	// } else if strings.HasSuffix(report, "```") {
	// 	report = report[:len(report)-3]
	// }
	//
	return strings.TrimSpace(report)
}

func ParseMarkdownToHTML(content string) (string, error) {
	parser := newParser()

	var htmlOutput bytes.Buffer
	if err := parser.Convert([]byte(cleanMarkdownBlock(content)), &htmlOutput); err != nil {
		return "", err
	}

	return htmlOutput.String(), nil
}
