package analyzer

import (
	"fmt"
	"golang.org/x/tools/go/analysis/analysistest"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestValidCodeAnalysis(t *testing.T) {
	analysistest.Run(t, testdataDir(), New(), "valid")
}

func TestInvalidCodeAnalysis(t *testing.T) {
	wantErrs := []string{
		"invalid/function_private.go:13:2: unexpected diagnostic: call to private method",
		"invalid/new_func.go:19:3: unexpected diagnostic: private field \"bio\" should not be initialized with a composite literal",
		"invalid/new_func.go:28:3: unexpected diagnostic: private field \"bio\" should not be initialized with a composite literal",
		"invalid/new_func.go:37:2: unexpected diagnostic: private field \"n.bio\" should not be initialized with a composite literal",
		"invalid/not_own_fields.go:20:13: unexpected diagnostic: initialization of private field: l",
		"invalid/not_own_fields.go:21:2: unexpected diagnostic: access to private field",
		"invalid/not_own_fields.go:22:2: unexpected diagnostic: access to private field",
		"invalid/not_own_fields.go:23:9: unexpected diagnostic: access to private field",
		"invalid/not_own_method.go:10:2: unexpected diagnostic: call to private method",
	}
	var gotErrs []string

	analysistest.Run(lintErrors{&gotErrs}, testdataDir(), New(), "invalid")

	slices.Sort(gotErrs)
	if !slices.Equal(wantErrs, gotErrs) {
		t.Fatalf("want: %v, got: %v", wantErrs, gotErrs)
	}
}

func testdataDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(filepath.Dir(wd), "testdata")
}

type lintErrors struct {
	Msgs *[]string
}

func (a lintErrors) Errorf(format string, args ...interface{}) {
	*a.Msgs = append(*a.Msgs, fmt.Sprintf(format, args...))
}
