package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/funjack/funify/datautils"
	"github.com/funjack/funify/funscript"
)

func readIntsFromFile(filename string) ([]int, error) {
	f, err := os.Open(filepath.Clean(filename))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	values := make([]int, 0, 1e6)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		values = append(values, value)
	}
	return values, scanner.Err()
}

func main() {
	var (
		filename   = flag.String("input", "", "input file")
		invert     = flag.Bool("invert", false, "invert positions")
		samplerate = flag.Float64("rate", 0, "sample rate in Hz (samples per second)")
		duration   = flag.Int64("duration", 0, "duration in milliseconds")
		speed      = flag.Float64("speed", 1, "speed factor")
		tolerance  = flag.Int("tolerance", 5, "rounding tolerance in percent")
		smooth     = flag.Int("smooth", 1, "smoothing window size in number samples")
	)
	flag.Parse()

	rawvalues, err := readIntsFromFile(*filename)
	if err != nil {
		log.Fatal(err)
	}

	var rate float64
	if *samplerate > 0 {
		rate = *samplerate * *speed
	} else if *duration > 0 {
		rate = float64(len(rawvalues)*1e3) / float64(*duration) * *speed
		fmt.Fprintf(os.Stderr, "detected sample rate: %.2f Hz\n", rate)
	} else {
		log.Fatal("no samplerate or duration specified")
	}

	if *smooth > 0 {
		rawvalues = datautils.Smooth(rawvalues, *smooth)
	}
	values := datautils.PercentValues(rawvalues, *invert)
	script := funscript.Generate(values,
		funscript.GeneratetOpts{
			SampleRate: rate,
			Tolerance:  byte(*tolerance),
		})
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(&script); err != nil {
		log.Fatal(err)
	}
}
