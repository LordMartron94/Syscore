package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type windowSpecManifest struct {
	Submodules map[string]windowSpecSubmodule `json:"submodules"`
	Inputs     map[string]windowSpecInput     `json:"inputs"`
}

type windowSpecSubmodule struct {
	Path   string `json:"path"`
	URL    string `json:"url"`
	Tag    string `json:"tag"`
	Commit string `json:"commit"`
}

type windowSpecInput struct {
	Submodule    string `json:"submodule"`
	RelativePath string `json:"relativePath"`
	SHA256       string `json:"sha256"`
}

func windowSpecManifestLoad(specRoot string) (windowSpecManifest, error) {
	manifestPath := filepath.Join(specRoot, "spec_manifest.json")
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return windowSpecManifest{}, fmt.Errorf("read spec manifest: %w", err)
	}

	var manifest windowSpecManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return windowSpecManifest{}, fmt.Errorf("parse spec manifest: %w", err)
	}
	return manifest, nil
}

func windowSpecInputPathResolve(specRoot string, manifest windowSpecManifest, inputKey string) (string, error) {
	input, ok := manifest.Inputs[inputKey]
	if !ok {
		return "", fmt.Errorf("spec manifest: unknown input %q", inputKey)
	}

	var absPath string
	if input.Submodule != "" {
		submodule, ok := manifest.Submodules[input.Submodule]
		if !ok {
			return "", fmt.Errorf("spec manifest: unknown submodule %q for input %q", input.Submodule, inputKey)
		}
		absPath = filepath.Join(specRoot, submodule.Path, input.RelativePath)
	} else {
		absPath = filepath.Join(specRoot, input.RelativePath)
	}

	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("spec input %q missing at %s (run: git submodule update --init): %w", inputKey, absPath, err)
	}

	if input.SHA256 != "" {
		if err := windowSpecFileSHA256Verify(absPath, input.SHA256); err != nil {
			return "", err
		}
	}

	return absPath, nil
}

func windowSpecFileSHA256Verify(absPath string, expected string) error {
	content, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", absPath, err)
	}
	sum := sha256.Sum256(content)
	actual := hex.EncodeToString(sum[:])
	if actual != expected {
		return fmt.Errorf("checksum mismatch for %s: got %s want %s", absPath, actual, expected)
	}
	return nil
}
