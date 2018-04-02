// Package funscript provides tools the generate and manipulate funscripts.
package funscript

import (
	"errors"

	"github.com/funjack/launchcontrol/protocol/funscript"
)

// Cleanup returns a copy of a list of actions with all redundant actions
// removed.
func Cleanup(f []funscript.Action) []funscript.Action {
	if len(f) < 2 {
		return f
	}
	cleanActions := make([]funscript.Action, 1, len(f))
	cleanActions[0] = f[0]
	for i := 1; i < len(f)-1; i++ {
		prev, cur, next := f[i-1], f[i], f[i+1]
		if cur.Pos != prev.Pos || cur.Pos != next.Pos {
			cleanActions = append(cleanActions, cur)
		}
	}
	cleanActions = append(cleanActions, f[len(f)-1])
	return cleanActions
}

// GeneratetOpts bundles the script are used during the conversion of raw data
// into a funscript.
type GeneratetOpts struct {
	// SampleRate in Hz
	SampleRate float64
	// Tolerance that is used to round values during comparisons.
	Tolerance byte
	// MinInterval is the minimum amount of time in miliseconds between actions.
	MinInterval int
}

// Generate a funscript from raw data in percent (0-100).
func Generate(p []byte, opts GeneratetOpts) funscript.Script {
	// default options when not specified
	if opts.SampleRate <= 0 {
		opts.SampleRate = 1
	}
	if opts.Tolerance == 0 {
		opts.Tolerance = 5
	}
	if opts.MinInterval < 1 {
		opts.MinInterval = 150
	}

	var actions = make([]funscript.Action, 0, len(p))

	minInterval := int(opts.SampleRate * float64(opts.MinInterval) / 1e3)
	if minInterval < 1 {
		minInterval = 1
	}
	maxima := localminmax(p, opts.Tolerance)
	for i := 1; i < len(maxima); i++ {
		ppos, pos := maxima[i-1], maxima[i]
		pausepositions := pauses(p[ppos:pos], minInterval, opts.Tolerance)
		for _, pausepos := range pausepositions {
			action, err := funscriptAction(p, ppos+pausepos, opts.SampleRate, opts.Tolerance)
			if err != nil {
				continue
			}
			if len(actions) < 1 || actions[len(actions)-1] != action {
				actions = append(actions, action)
			}
		}
		action, err := funscriptAction(p, pos, opts.SampleRate, opts.Tolerance)
		if err != nil {
			continue
		}
		if len(actions) < 1 || actions[len(actions)-1] != action {
			actions = append(actions, action)
		}
	}
	return funscript.Script{
		Actions:  Cleanup(actions),
		Version:  "1.0",
		Inverted: false,
		Range:    funscript.Range(90),
	}
}

// funscriptAction returns a funscript action for the given sample in raw data
// using the specified sample rate in Hz. The pos is rounded to tolerance.
func funscriptAction(p []byte, sample int, rate float64, tolerance byte) (funscript.Action, error) {
	if sample < 0 || sample > len(p)-1 {
		return funscript.Action{}, errors.New("invalid sample")
	}
	return funscript.Action{
		At:  sampleinms(sample, rate),
		Pos: int(round(p[sample], tolerance)),
	}, nil
}

// localminmax finds all the local maxima locations in p with all the values
// rounded to tolerance.
func localminmax(p []byte, tolerance byte) []int {
	var (
		last byte
		up   bool
		f    = make([]int, 0, len(p)/10)
	)
	if len(p) < 2 {
		return f
	}
	last = p[0]
	if p[0] < p[1] {
		up = true
	}
	for i, v := range p[1:] {
		rounded := round(v, tolerance)
		if up {
			if rounded < last {
				up = false
				f = append(f, i-1)
			}
		} else {
			if rounded > last {
				up = true
				f = append(f, i-1)
			}
		}
		last = rounded
	}
	return f
}

// round value b to the given tolerance.
func round(b, tolerance byte) byte {
	return b / tolerance * tolerance
}

// pauses returns all the start and stop sample positions in p with pauses
// within the given tolerance and have the length of the specified interval.
func pauses(p []byte, interval int, tolerance byte) []int {
	var f []int
	last := p[0]
	for i := interval; i <= (len(p) - interval); i += interval {
		if round(p[i], tolerance) == round(last, tolerance) {
			f = append(f, i-interval, i)
		}
		last = p[i]
	}
	return f
}

// sampleinms returns the time in ms of a sample number given the samplerate in
// Hz.
func sampleinms(sample int, rate float64) int64 {
	return int64(1e3 / rate * float64(sample))
}
