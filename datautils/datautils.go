// Package datautils provides helpers for dealing with sampled data.
package datautils

// Smooth values using a x-point (size*2+1) window moving average filter.
func Smooth(v []int, size int) []int {
	if size < 1 {
		return v
	}
	window := size*2 + 1
	if len(v) < window {
		return v
	}
	var r = make([]int, len(v))
	for i := 0; i < size; i++ {
		r[i] = v[i]
		r[len(v)-1-i] = v[len(v)-1-i]
	}
	for i := size; i < len(v)-size; i++ {
		var total int64
		for j := -size; j <= size; j++ {
			total += int64(v[i+j])
		}
		r[i] = int(total / int64(window))
	}
	return r
}

// Extrema returns the lowest and highest value in the given slice.
func Extrema(v []int) (min, max int) {
	if len(v) < 1 {
		return
	}
	min, max = v[0], v[0]
	for _, value := range v {
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}
	return
}

// PercentValues returns a slice with all values in v converted to a percentage.
func PercentValues(v []int, invert bool) []byte {
	min, max := Extrema(v)
	size := max - min
	out := make([]byte, len(v))
	for i := range v {
		out[i] = byte((max - v[i]) * 100 / size)
		if invert {
			out[i] = 100 - out[i]
		}
	}
	return out
}
