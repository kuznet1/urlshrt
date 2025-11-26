package main

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
	"strings"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	res := analysistest.Run(t, testdata, structAnalyzer, "main")
	resets := strings.Join(res[0].Result.([]string), "\n")
	expected := `
func (s *Test) Reset() {
	if s == nil {
		return
	}

	s.I = 0
	s.F32 = 0
	s.Str = ""
	s.Flag = false
	// unsupported value field 'Custom'

	if s.PI != nil {
		*s.PI = 0
	}

	if s.PStr != nil {
		*s.PStr = ""
	}

	if s.PBool != nil {
		*s.PBool = false
	}

	if s.PCustom != nil {
		if resetter, ok := interface{}(s.PCustom).(interface{ Reset() }); ok {
			resetter.Reset()
		}
	}
	s.Sl = s.Sl[:0]
	// unsupported array field 'Arr'
	clear(s.M)

	if resetter, ok := s.Any.(interface{ Reset() }); ok {
		resetter.Reset()
	}
	// unsupported pointer field 'SomeStr'
}`
	assert.Equal(t, expected, resets)
}
