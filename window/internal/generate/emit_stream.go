package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codegen"
	gocode "codegen/go"
)

type bindingFileStream struct {
	path       string
	file       *os.File
	stream     *gocode.FileStream
	flushEvery int
}

func bindingFileStreamOpen(outputDir, fileName, packageName, buildTag string, flushEvery int) (*bindingFileStream, error) {
	return bindingFileStreamCreate(outputDir, fileName, bindingFileStreamPreamble(packageName, buildTag), flushEvery)
}

func bindingFileStreamCreate(outputDir, fileName string, initialElements []codegen.FileElement, flushEvery int) (*bindingFileStream, error) {
	absPath := filepath.Join(outputDir, fileName)
	file, err := os.Create(absPath)
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", absPath, err)
	}

	stream := gocode.FileStreamCreate(file, gocode.RenderOptionsEnumBindings(), gocode.FileStreamConfig{
		FlushEveryElements: flushEvery,
	})

	out := &bindingFileStream{
		path:       absPath,
		file:       file,
		stream:     stream,
		flushEvery: flushEvery,
	}

	if len(initialElements) > 0 {
		if err := bindingFileStreamWriteElements(out, initialElements...); err != nil {
			file.Close()
			os.Remove(absPath)
			return nil, err
		}
	}
	return out, nil
}

func bindingFileStreamPreamble(packageName, buildTag string) []codegen.FileElement {
	return bindingFileHeaderElements(packageName, buildTag)
}

func bindingFileStreamWriteElements(stream *bindingFileStream, elements ...codegen.FileElement) error {
	if stream == nil || stream.stream == nil {
		return nil
	}
	return gocode.FileStreamWriteElements(stream.stream, elements...)
}

func bindingFileStreamClose(stream *bindingFileStream) error {
	if stream == nil {
		return nil
	}
	if err := gocode.FileStreamClose(stream.stream); err != nil {
		return fmt.Errorf("flush %s: %w", stream.path, err)
	}
	if err := stream.file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", stream.path, err)
	}
	return nil
}

func bindingFileStreamFormat(stream *bindingFileStream) error {
	if stream == nil {
		return nil
	}
	return bindingFilesFormat(filepath.Dir(stream.path), filepath.Base(stream.path))
}
