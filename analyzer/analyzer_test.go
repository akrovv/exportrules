package analyzer

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"os"
	"path/filepath"
	"testing"
)

func TestPositive(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	testdata := filepath.Join(filepath.Dir(wd), "testdata")
	analysistest.Run(t, testdata, New(), "tests")
}

func TestNegative(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	skipForTest = true
	testdata := filepath.Join(filepath.Dir(wd), "testdata")
	analysistest.Run(t, testdata, New(), "negative")
}
