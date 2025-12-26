package parser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type SpecFile struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	HTMLBody string `json:"html_body"`
	IsDir    bool   `json:"is_dir"`
}

type SpecTree struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"is_dir"`
	Children []*SpecTree `json:"children,omitempty"`
	File     *SpecFile   `json:"file,omitempty"`
}

func ParseDirectory(rootPath string) (*SpecTree, error) {
	root := &SpecTree{
		Name:  filepath.Base(rootPath),
		Path:  rootPath,
		IsDir: true,
	}

	err := walkDir(rootPath, root)
	if err != nil {
		return nil, err
	}

	return root, nil
}

func ParseMultipleDirectories(paths []string) ([]*SpecTree, error) {
	var trees []*SpecTree

	for _, path := range paths {
		tree, err := ParseDirectory(path)
		if err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}

	return trees, nil
}

func walkDir(dirPath string, node *SpecTree) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// Skip hidden files and directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			child := &SpecTree{
				Name:  entry.Name(),
				Path:  fullPath,
				IsDir: true,
			}
			node.Children = append(node.Children, child)
			if err := walkDir(fullPath, child); err != nil {
				return err
			}
		} else if strings.HasSuffix(entry.Name(), ".md") {
			specFile, err := ParseMarkdownFile(fullPath)
			if err != nil {
				return err
			}

			child := &SpecTree{
				Name:  entry.Name(),
				Path:  fullPath,
				IsDir: false,
				File:  specFile,
			}
			node.Children = append(node.Children, child)
		}
	}

	// Sort children: directories first, then files alphabetically
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		return node.Children[i].Name < node.Children[j].Name
	})

	return nil
}

func ParseMarkdownFile(filePath string) (*SpecFile, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var htmlBuf strings.Builder
	if err := md.Convert(content, &htmlBuf); err != nil {
		return nil, err
	}

	return &SpecFile{
		Path:     filePath,
		Name:     filepath.Base(filePath),
		Content:  string(content),
		HTMLBody: htmlBuf.String(),
		IsDir:    false,
	}, nil
}

func GetFileContent(filePath string) (*SpecFile, error) {
	return ParseMarkdownFile(filePath)
}
