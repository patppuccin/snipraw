package helpers

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/patppuccin/snipraw/src/web"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var (
	ChromaLightCSS string
	ChromaDarkCSS  string
)

var sniprawLight = chroma.MustNewStyle("snipraw-light", chroma.StyleEntries{
	chroma.Background:          "#1a1a1a bg:#f5f4f2",
	chroma.Comment:             "italic #6b7280",
	chroma.Keyword:             "bold #3d7a5a",
	chroma.KeywordDeclaration:  "bold #3d7a5a",
	chroma.KeywordNamespace:    "bold #3d7a5a",
	chroma.NameFunction:        "#5a8c6e",
	chroma.NameClass:           "#7aab8a",
	chroma.NameBuiltin:         "#5a8c6e",
	chroma.LiteralString:       "#16a34a",
	chroma.LiteralStringDouble: "#16a34a",
	chroma.LiteralStringSingle: "#16a34a",
	chroma.LiteralNumber:       "#f87171",
	chroma.Operator:            "#4a7a5e",
	chroma.Punctuation:         "#1a1a1a",
	chroma.LineNumbers:         "#9ca3af",
})

var sniprawDark = chroma.MustNewStyle("snipraw-dark", chroma.StyleEntries{
	chroma.Background:          "#e5e5e5 bg:#1a1a1a",
	chroma.Comment:             "italic #6b7280",
	chroma.Keyword:             "bold #6fa783",
	chroma.KeywordDeclaration:  "bold #6fa783",
	chroma.KeywordNamespace:    "bold #6fa783",
	chroma.NameFunction:        "#8dc8a8",
	chroma.NameClass:           "#a3c9b0",
	chroma.NameBuiltin:         "#8dc8a8",
	chroma.LiteralString:       "#4ade80",
	chroma.LiteralStringDouble: "#4ade80",
	chroma.LiteralStringSingle: "#4ade80",
	chroma.LiteralNumber:       "#f87171",
	chroma.Operator:            "#5a9470",
	chroma.Punctuation:         "#e5e5e5",
	chroma.LineNumbers:         "#6b7280",
})

func init() {
	formatter := chromahtml.New(chromahtml.WithClasses(true))
	var light, dark strings.Builder
	if err := formatter.WriteCSS(&light, sniprawLight); err == nil {
		ChromaLightCSS = strings.ReplaceAll(light.String(), ".chroma", "[data-theme='light'] .chroma")
	}
	if err := formatter.WriteCSS(&dark, sniprawDark); err == nil {
		ChromaDarkCSS = strings.ReplaceAll(dark.String(), ".chroma", "[data-theme='dark'] .chroma")
	}
}

var calloutTypes = []string{"NOTE", "TIP", "WARNING", "IMPORTANT", "CAUTION"}

func rewriteCallouts(html string) string {
	for _, t := range calloutTypes {
		old := fmt.Sprintf("<blockquote>\n<p>[!%s]\n", t)
		new := fmt.Sprintf(`<blockquote class="callout callout-%s"><p class="callout-title">%s</p><p>`, strings.ToLower(t), calloutTitle(t))
		html = strings.ReplaceAll(html, old, new)
	}
	return html
}

func calloutTitle(t string) string {
	switch t {
	case "NOTE":
		return "📝 Note"
	case "TIP":
		return "💡 Tip"
	case "WARNING":
		return "⚠️ Warning"
	case "IMPORTANT":
		return "🔔 Important"
	case "CAUTION":
		return "🚨 Caution"
	default:
		return t
	}
}

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Table,
		extension.Strikethrough,
		extension.TaskList,
		highlighting.NewHighlighting(
			highlighting.WithCSSWriter(io.Discard),
			highlighting.WithFormatOptions(
				chromahtml.WithClasses(true),
				chromahtml.WithLineNumbers(false),
			),
			highlighting.WithStyle("fallback"),
			highlighting.WithGuessLanguage(true),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

func RenderMarkdown(content []byte, project, dir string) string {
	reader := text.NewReader(content)
	doc := md.Parser().Parse(reader)

	// Rewrite relative image srcs to raw endpoint
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		img, ok := n.(*ast.Image)
		if !ok {
			return ast.WalkContinue, nil
		}
		dest := string(img.Destination)
		if strings.HasPrefix(dest, "http://") || strings.HasPrefix(dest, "https://") || strings.HasPrefix(dest, "/") {
			return ast.WalkContinue, nil
		}
		// Resolve relative to current dir within project
		resolved := path.Join(dir, dest)
		img.Destination = []byte("/project/" + project + "/blob/" + resolved)
		return ast.WalkContinue, nil
	})

	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, content, doc); err != nil {
		return ""
	}

	return rewriteCallouts(buf.String())

}

func ListProjects(root string) ([]web.Project, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var projects []web.Project
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		if e.Name() == ".git" {
			continue
		}

		count := 0
		filepath.WalkDir(filepath.Join(root, e.Name()), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() && d.Name() == ".git" {
				return filepath.SkipDir
			}
			if !d.IsDir() {
				count++
			}
			return nil
		})

		projects = append(projects, web.Project{
			Name:      e.Name(),
			FileCount: count,
		})
	}

	return projects, nil
}

func CodeHighlight(content []byte, filename string) (string, error) {
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Analyse(string(content))
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	formatter := html.New(
		html.WithLineNumbers(true),
		html.LineNumbersInTable(true),
		html.WithClasses(true),
		html.TabWidth(2),
	)

	var buf strings.Builder
	iterator, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		return "", err
	}

	if err := formatter.Format(&buf, styles.Fallback, iterator); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func IsImage(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return true
	}
	return false
}

func IsPDF(name string) bool {
	return strings.ToLower(filepath.Ext(name)) == ".pdf"
}

func IsMarkdown(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".md", ".markdown", ".mkd":
		return true
	}
	return false
}

var readmeNames = []string{
	"README.md",
	"README.markdown",
	"README.mkd",
	"README.txt",
	"README",
	"readme.md",
	"readme.markdown",
	"readme.txt",
	"readme",
}

func FindReadme(dir string) []byte {
	for _, name := range readmeNames {
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err == nil {
			return content
		}
	}
	return nil
}
