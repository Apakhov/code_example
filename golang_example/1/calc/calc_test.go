package main

import (
	"bytes"
	"strings"
	"testing"
)

var basicOperationsTests = []struct {
	in     string
	expRes string
}{
	{"1 1 + =", "2\n"},
	{"3 5 - =", "-2\n"},
	{"15 15 * =", "225\n"},
	{"18 7 / =", "2\n"},
}

func TestBasicOperations(t *testing.T) {
	for _, test := range basicOperationsTests {
		out := new(bytes.Buffer)
		err := calc(strings.NewReader(test.in), out)
		if err != nil {
			t.Errorf("test for OK Failed - error")
		}
		result := out.String()
		if result != test.expRes {
			t.Errorf("test for OK Failed - results not match\nChecking:\n%v\nGot:\n%v\nExpected:\n%v", test.in, result, test.expRes)
		}
	}
}

var errorsTests = []struct {
	in     string
	expErr string
}{
	{"666 0 / =", "division by zero"},
	{"666 / =", "not enough arguments for /"},
	{"666 * =", "not enough arguments for *"},
	{"666 + =", "not enough arguments for +"},
	{"666 - =", "not enough arguments for -"},
	{"666 12 + 2394 329 - - as =", "illigal symbol: a"},
}

func TestErrIllegalSymb(t *testing.T) {
	for _, test := range errorsTests {
		out := new(bytes.Buffer)
		err := calc(strings.NewReader(test.in), out)
		if err == nil {
			t.Errorf("test for OK Failed - expected error")
		}
		if err.Error() != test.expErr {
			t.Errorf("test for OK Failed - errors not match\nChecking:\n%v\nGot:\n%v\nExpected:\n%v", test.in, err, test.expErr)
		}
	}
}

func TestComplicatedTask(t *testing.T) {
	out := new(bytes.Buffer)
	err := calc(strings.NewReader("15 15 * = 25 - = 8 4 / = / = 78 - ="), out)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	if result != "225\n200\n2\n100\n22\n" {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, "225\n200\n2\n100\n22\n")
	}
}

func TestStack(t *testing.T) {
	stack := newStack()
	_, err := stack.top()
	if err == nil {
		t.Errorf("test for OK Failed - expected error")
	}
	if err.Error() != "no values in stack" {
		t.Errorf("test for OK Failed - errors not match\nGot:\n%v\nExpected:\n%v", err, "no values in stack")
	}
	stack = newStack()
	_, err = stack.pop()
	if err == nil {
		t.Errorf("test for OK Failed - expected error")
	}
	if err.Error() != "no values in stack" {
		t.Errorf("test for OK Failed - errors not match\nGot:\n%v\nExpected:\n%v", err, "no values in stack")
	}
}
