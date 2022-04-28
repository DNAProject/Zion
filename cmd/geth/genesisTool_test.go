package main

import (
	"errors"
	"testing"
)

func TestGenesisTool(t *testing.T){
	testCases := []struct {
		AllArgs []string
		Expect error
	}{
		{
			[]string{"genesisTool", "generate", "3"},
			errors.New("Fatal: got 3 nodes, but hotstuff BFT requires at least 4 nodes"),
		},
		{
			[]string{"genesisTool", "generate", "5"},
			nil,
		},
		{
			[]string{"genesisTool", "generate", "5", "-basePath", "./../../build/bin/testGenesis/"},
			nil,
		},
	}

	for _, testCase := range testCases {
		geth := runGeth(t, testCase.AllArgs...)
		if testCase.Expect != nil {
			geth.ExpectRegexp(testCase.Expect.Error())
		}
		geth.ExpectRegexp("")
	}
}
