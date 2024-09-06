package process

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/actord/actord/pkg/process/schema"
)

func Parse(processFolderPath string) (*schema.Schema, error) {
	packageName, data, err := ReadPackage(processFolderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package: %w", err)
	}

	file, diags := hclsyntax.ParseConfig(data, packageName, hcl.Pos{Line: 1, Column: 1})
	if diags != nil && diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse config: %w", diags)
	}

	if file == nil {
		return nil, fmt.Errorf("failed to parse config: file is nil")
	}
	var s schema.Schema

	diags = gohcl.DecodeBody(file.Body, nil, &s)
	if diags.HasErrors() {
		showError(data, diags[0])
		return nil, fmt.Errorf("failed to decode body: %w", diags)
	}

	log.Println("Schema:", s)

	return &s, nil
}

func showError(data []byte, diag *hcl.Diagnostic) {
	if diag == nil {
		panic("diag nil")
	}
	var rng *hcl.Range
	if diag.Subject != nil {
		rng = diag.Subject
	} else if diag.Context != nil {
		rng = diag.Context
	} else {
		fmt.Println(diag.Error())
		return
	}
	line := rng.Start.Line
	dataLines := strings.Split(string(data), "\n")

	start := line - 5
	if start < 0 {
		start = 0
	}

	end := line + 5
	if end > len(dataLines) {
		end = len(dataLines)
	}

	shownLineWithError := line - 1

	for i := start; i < end; i++ {
		fmt.Println(dataLines[i])
		if i == shownLineWithError {
			for j := 0; j < rng.Start.Column-1; j++ {
				fmt.Print(" ")
			}
			for j := rng.Start.Column - 1; j < rng.End.Column; j++ {
				fmt.Print("^")
			}
			fmt.Println(" ->", diag.Summary)
		}
	}
}

func ReadPackage(processFolderPath string) (string, []byte, error) {
	files, err := os.ReadDir(processFolderPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read process directory: %w", err)
	}

	var allFiles []io.Reader

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".hcl") {
			continue
		}
		f, err := os.Open(path.Join(processFolderPath, file.Name()))
		if err != nil {
			return "", nil, fmt.Errorf("failed to open file: %w", err)
		}
		allFiles = append(allFiles, f)
	}

	data, err := io.ReadAll(io.MultiReader(allFiles...))
	if err != nil {
		return "", nil, fmt.Errorf("failed to read all files: %w", err)
	}

	processFolderPathItems := strings.Split(processFolderPath, "/")
	packageName := processFolderPathItems[len(processFolderPathItems)-1]

	return packageName, data, nil
}
