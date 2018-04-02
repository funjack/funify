package funscript

import (
	"testing"

	"github.com/funjack/launchcontrol/protocol/funscript"
)

func TestFunscriptAction(t *testing.T) {
	raw := []byte{1, 1, 1, 1, 1, 1, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 100, 100, 100, 100, 100}
	action, err := funscriptAction(raw, 0, 5.0, 5)
	if err != nil {
		t.Errorf("no error expected: %v", err)
	}
	want := funscript.Action{At: 0, Pos: 0}
	if action != want {
		t.Errorf("action did not match, want %+v, got %+v", want, action)
	}
	action, err = funscriptAction(raw, 5, 5.0, 5)
	if err != nil {
		t.Errorf("no error expected: %v", err)
	}
	want = funscript.Action{At: 1000, Pos: 0}
	if action != want {
		t.Errorf("action did not match, want %+v, got %+v", want, action)
	}
	action, err = funscriptAction(raw, 10, 5.0, 5)
	if err != nil {
		t.Errorf("no error expected: %v", err)
	}
	want = funscript.Action{At: 2000, Pos: 50}
	if action != want {
		t.Errorf("action did not match, want %+v, got %+v", want, action)
	}
}
