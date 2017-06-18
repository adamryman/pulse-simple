// generate some simple tones and play them on the default audio output.
//
// to run this example, simply `go run examples/sinewave.go`.
package main

import (
	"encoding/binary"
	"fmt"
	"github.com/mesilliac/pulse-simple" // pulse-simple
	"github.com/mjibson/go-dsp/wav"
	"math"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w, err := wav.New(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(w.SampleRate)
	fmt.Println(w.BitsPerSample)
	ss := pulse.SampleSpec{pulse.SAMPLE_FLOAT32LE, 192000, 1}
	pb, err := pulse.Playback("pulse-simple test", "playback test", &ss)
	defer pb.Free()
	defer pb.Drain()
	if err != nil {
		fmt.Printf("Could not create playback stream: %s\n", err)
		return
	}
	playwav(pb, &ss, w)
}

func playwav(s *pulse.Stream, ss *pulse.SampleSpec, w *wav.Wav) {
	data := make([]byte, 4*ss.Rate)
	for {
		for i := 0; i < int(ss.Rate); i++ {
			sample, err := w.ReadFloats(1)
			if err != nil {
				fmt.Println(err)
				return
			}
			bits := math.Float32bits(sample[0])

			binary.LittleEndian.PutUint32(data[4*i:4*i+4], bits)
		}
		s.Write(data)
	}
}
