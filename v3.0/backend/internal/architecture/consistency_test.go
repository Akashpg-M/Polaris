package architecture

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestProductionCodeHasSingleSpatialRoutingAndCommandAuthorities(t *testing.T) {
	backend, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	forbiddenImports := []string{"/algo_/quadtree", "/algo_/graph", "/internal/core/actor"}
	err = filepath.WalkDir(backend, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".phase41-go-cache" || entry.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || strings.Contains(path, string(filepath.Separator)+"api"+string(filepath.Separator)+"proto") {
			return nil
		}
		parsed, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if parseErr != nil {
			return parseErr
		}
		for _, imported := range parsed.Imports {
			name, _ := strconv.Unquote(imported.Path.Value)
			for _, forbidden := range forbiddenImports {
				if strings.HasSuffix(name, forbidden) {
					t.Errorf("%s imports retired authority %s", path, name)
				}
			}
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(content), "telemetry:commands") || strings.Contains(string(content), "StartAutonomousLoop") {
			t.Errorf("%s retains a retired direct command path", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNodeTypeIsNotUsedAsBitmask(t *testing.T) {
	backend, _ := filepath.Abs(filepath.Join("..", ".."))
	_ = filepath.WalkDir(backend, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return walkErr
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(node ast.Node) bool {
			binary, ok := node.(*ast.BinaryExpr)
			if ok && (binary.Op == token.AND || binary.Op == token.OR) {
				text := expressionName(binary.X) + expressionName(binary.Y)
				if strings.Contains(strings.ToLower(text), "type") || strings.Contains(strings.ToLower(text), "class") {
					t.Errorf("%s contains a type/class bitmask expression", path)
				}
			}
			return true
		})
		return nil
	})
}

func expressionName(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return expressionName(value.X) + value.Sel.Name
	}
	return ""
}
