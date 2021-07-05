package analyzer

import (
	"testing"

	"github.com/Jesse-Cameron/golang-nil-error-struct/testdata/src/a"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, NilAnalyzer, "a")
}

func TestErrorFuncs(t *testing.T) {
	err := a.MakePointerError()
	assert.False(t, err == nil)

	err = a.MakeError()
	assert.False(t, err == nil)

	err = a.MakeErrorParenDecl()
	assert.False(t, err == nil)

	err = a.MakeErrorListDecl()
	assert.False(t, err == nil)
}
