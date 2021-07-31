package analyzer

import (
	"testing"

	"github.com/Jesse-Cameron/unassignederr/testdata/src/a"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, UnassignedErrAnalyzer, "a")
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

	err = a.ErrorAssigned()
	assert.True(t, err != nil)
	assert.True(t, err.Error() == "error")
}
