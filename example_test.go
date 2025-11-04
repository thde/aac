package aac_test

import (
	"log"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/generators"
	"github.com/thde/aac"
)

// Example demonstrates encoding generated audio (sine wave).
func ExampleEncode() {
	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}

	// Generate a 440 Hz sine wave (A4 note) for 3 seconds
	sine, err := generators.SineTone(format.SampleRate, 440)
	if err != nil {
		log.Fatal(err)
	}
	limited := beep.Take(format.SampleRate.N(3*time.Second), sine)

	// Create output file
	out, err := os.Create("sine.aac")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Encode to AAC
	if err = aac.Encode(out, limited, aac.EncodeOptions{
		Format:  format,
		BitRate: 128000,
	}); err != nil {
		log.Fatal(err)
	}
}
