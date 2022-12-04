//go:build mage

package main

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

const VENV_PATH = ".venv"

var (
	PIP     string
	PYTHON  string
	JUPYTER string
)

func init() {
	PIP = filepath.Join(VENV_PATH, "bin", "pip")
	PYTHON = filepath.Join(VENV_PATH, "bin", "python")
	JUPYTER = filepath.Join(VENV_PATH, "bin", "jupyter")
}

func ensureVirtualenv() error {
	if err := sh.Run("python", "-mvirtualenv", VENV_PATH); err != nil {
		return errors.Wrap(err, "failed to create virtualenv")
	}
	return sh.Run(PIP, "install", "nbconvert")
}

func renderNotebook(notebook string) error {
	// jupyter nbconvert --to markdown --output-dir=content/posts/najemnine2/ ~/Downloads/ETN_SLO_CSV_A_NAJ/najemnine.ipynb
	outputDir := filepath.Dir(notebook)
	fmt.Println("outputDir", outputDir)
	return sh.Run(JUPYTER, "nbconvert", "--to", "markdown", "--output-dir", outputDir, notebook)
}

func walkNotebooks() error {
	return filepath.Walk("content", func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".ipynb" {
			return renderNotebook(path)
		}
		return nil
	})
}

// Runs go mod download and then installs the binary.
func Build() error {
	if err := ensureVirtualenv(); err != nil {
		return err
	}
	return walkNotebooks()
}
