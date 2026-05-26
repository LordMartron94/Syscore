package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"codegen"
	gocode "codegen/go"
)

const syscoreWindowBindingsGeneratorTool = "syscore window bindings generator"

func bindingFilePreamble(packageName, buildTag string) []codegen.FileElement {
	return bindingFileHeaderElements(packageName, buildTag)
}

func bindingFileHeaderElements(packageName, buildTag string) []codegen.FileElement {
	elements := make([]codegen.FileElement, 0, 8)
	if buildTag != "" {
		elements = append(elements, gocode.GoBuildConstraint(buildTag))
		elements = append(elements, gocode.GoBlankLine())
	}
	elements = append(elements, gocode.GoGeneratedFileHeader(syscoreWindowBindingsGeneratorTool, time.Now())...)
	elements = append(elements, gocode.FileElementFrom(gocode.DeclPackage(packageName)))
	elements = append(elements, gocode.GoBlankLine())
	return elements
}

func writeBindingFile(outputDir string, fileName string, elements []codegen.FileElement) error {
	file := gocode.DeclFile(elements...)
	content, err := gocode.GoFileRenderWithOptions(file, gocode.RenderOptionsEnumBindings())
	if err != nil {
		return fmt.Errorf("render %s: %w", fileName, err)
	}

	absPath := filepath.Join(outputDir, fileName)
	if err := os.WriteFile(absPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", absPath, err)
	}
	return nil
}

func bindingFilesFormat(outputDir string, fileNames ...string) error {
	if len(fileNames) == 0 {
		return bindingDirFormat(outputDir)
	}
	args := []string{"gofmt", "-w"}
	for _, name := range fileNames {
		args = append(args, filepath.Join(outputDir, name))
	}
	cmd := exec.Command(args[0], args[1:]...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("gofmt: %w: %s", err, string(out))
	}
	return nil
}

func bindingDirFormat(outputDir string) error {
	cmd := exec.Command("gofmt", "-w", outputDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("gofmt %s: %w: %s", outputDir, err, string(out))
	}
	return nil
}

func bindingBlankLine(elements *[]codegen.FileElement) {
	*elements = append(*elements, gocode.GoBlankLine())
}

func loaderGenHeaderElements(buildTag string) []codegen.FileElement {
	elements := make([]codegen.FileElement, 0, 4)
	if buildTag != "" {
		elements = append(elements, gocode.GoBuildConstraint(buildTag))
		elements = append(elements, gocode.GoBlankLine())
	}
	elements = append(elements, gocode.GoGeneratedFileHeader(syscoreWindowBindingsGeneratorTool, time.Now())...)
	return elements
}

func bindingConstSpecString(name string, value string, doc string) gocode.ConstSpec {
	return gocode.ConstSpecNew(name, nil, fmt.Sprintf("%q", value), doc)
}

func bindingElementsWithImport(elements []codegen.FileElement, importPath string) []codegen.FileElement {
	insertAt := bindingElementsIndexAfterPackage(elements)
	if insertAt < 0 {
		return elements
	}
	out := make([]codegen.FileElement, 0, len(elements)+2)
	out = append(out, elements[:insertAt]...)
	out = append(out, gocode.FileElementFrom(gocode.DeclImportBlock(importPath)))
	out = append(out, gocode.GoBlankLine())
	out = append(out, elements[insertAt:]...)
	return out
}

func bindingElementsIndexAfterPackage(elements []codegen.FileElement) int {
	for i, element := range elements {
		if element.Node.NodeKind() == gocode.KindPackageDecl {
			return i + 1
		}
	}
	return -1
}
