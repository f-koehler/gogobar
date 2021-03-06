package gogobar_test

import (
	gogobar "github.com/f-koehler/gogobar/gogobar"
	"testing"
)

func TestPadLeft(t *testing.T) {
	strs := []string{
		"a",
		"a",
		"a",
		"hello",
	}

	lengths := []int{
		0,
		1,
		2,
		8,
	}

	outputs := []string{
		"a",
		"a",
		"_a",
		"___hello",
	}

	for i := 0; i < len(strs); i++ {
		output := gogobar.PadLeft(strs[i], '_', lengths[i])
		if output != outputs[i] {
			t.Error(output, "!=", outputs[i])
		}
	}
}

func TestPadRight(t *testing.T) {
	strs := []string{
		"a",
		"a",
		"a",
		"hello",
	}

	lengths := []int{
		0,
		1,
		2,
		8,
	}

	outputs := []string{
		"a",
		"a",
		"a_",
		"hello___",
	}

	for i := 0; i < len(strs); i++ {
		output := gogobar.PadRight(strs[i], '_', lengths[i])
		if output != outputs[i] {
			t.Error(output, "!=", outputs[i])
		}
	}
}

func TestPadBoth(t *testing.T) {
	strs := []string{
		"a",
		"a",
		"a",
		"a",
		"a",
		"a",
	}

	lengths := []int{
		0,
		1,
		2,
		3,
		4,
		5,
	}

	outputs := []string{
		"a",
		"a",
		"a_",
		"_a_",
		"_a__",
		"__a__",
	}

	for i := 0; i < len(strs); i++ {
		output := gogobar.PadBoth(strs[i], '_', lengths[i])
		if output != outputs[i] {
			t.Error(output, "!=", outputs[i])
		}
	}
}
