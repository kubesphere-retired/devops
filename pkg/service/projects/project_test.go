package projects

import (
	"testing"
)

func Test_checkJenkinsGoodName(t *testing.T) {
	goodValues := []string{
		"a", "ab", "abc", "a1", "a-1", "a--1--2--b",
		"0", "01", "012", "1a", "1-a", "1--a--b--2",
		"AB", "A.B", "A_B", "A_B", "A1a",
	}
	for _, val := range goodValues {
		if err := checkJenkinsGoodName(val); err != nil {
			t.Fatal(err)
		}
	}

	badValues := []string{
		".", "/", "a/b", "你好",
		" ", "a ", " a", "a b", "1 ", " 1", "1 2",
	}
	for _, val := range badValues {
		if err := checkJenkinsGoodName(val); err == nil {
			t.Fatalf("[%s] should be bad name", val)
		}
	}

}
