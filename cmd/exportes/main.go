package main

import (
	"github.com/akrovv/exportes/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.New()}, nil
}

func main() {
	anal, err := New(true)
	if err != nil {
		return
	}

	singlechecker.Main(anal[0])
	//anal[0].Run()
}
