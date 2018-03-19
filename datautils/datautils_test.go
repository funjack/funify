package datautils

import "testing"

func TestSmooth(t *testing.T) {
	cases := []struct {
		Input      []int
		WindowSize int
		Output     []int
	}{
		{
			Input:      []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			WindowSize: 1,
			Output:     []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			Input:      []int{1, 10, 1, 10, 1, 10, 1, 10, 1, 10},
			WindowSize: 1,
			Output:     []int{1, 4, 7, 4, 7, 4, 7, 4, 7, 10},
		},
		{
			Input:      []int{1, 10, 10, 10, 1, 10, 10, 10, 1, 10},
			WindowSize: 1,
			Output:     []int{1, 7, 10, 7, 7, 7, 10, 7, 7, 10},
		},
	}

	for cid, c := range cases {
		got := Smooth(c.Input, c.WindowSize)
		if len(c.Output) != len(got) {
			t.Errorf("case %d: output length did not match, want %d, got %d", cid, len(c.Output), len(got))
			continue
		}
		for i, v := range c.Output {
			if v != got[i] {
				t.Errorf("case %d: output did not match at position %d, want %v, got %v", cid, i, c.Output, got)
				break
			}
		}
	}
}
