package aac

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/generators"
)

func TestEncode(t *testing.T) {
	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}

	// Generate 1 second of sine wave
	sine, err := generators.SineTone(format.SampleRate, 440)
	if err != nil {
		t.Fatalf("Failed to create sine tone: %v", err)
	}
	limited := beep.Take(format.SampleRate.N(1*time.Second), sine)

	buf := &bytes.Buffer{}
	err = Encode(buf, limited, EncodeOptions{
		Format:  format,
		BitRate: 128000,
	})
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("Expected non-empty output")
	}
	if buf.Len() < 1000 {
		t.Fatalf("Output too small: %d bytes", buf.Len())
	}
}

func TestEncodeInvalidPrecision(t *testing.T) {
	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   1,
	}

	sine, err := generators.SineTone(format.SampleRate, 440)
	if err != nil {
		t.Fatalf("Failed to create sine tone: %v", err)
	}

	buf := &bytes.Buffer{}
	err = Encode(buf, sine, EncodeOptions{
		Format: format,
	})
	if err == nil {
		t.Fatal("expected error for unsupported precision")
	}
}

func TestEncodeMonoAudio(t *testing.T) {
	format := beep.Format{
		SampleRate:  22050,
		NumChannels: 1,
		Precision:   2,
	}

	sine, err := generators.SineTone(format.SampleRate, 440)
	if err != nil {
		t.Fatalf("Failed to create sine tone: %v", err)
	}
	limited := beep.Take(format.SampleRate.N(1*time.Second), sine)

	buf := &bytes.Buffer{}
	err = Encode(buf, limited, EncodeOptions{
		Format:  format,
		BitRate: 64000,
	})
	if err != nil {
		t.Fatalf("Mono encode failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("Expected non-empty output")
	}
	if buf.Len() < 500 {
		t.Fatalf("Output too small: %d bytes", buf.Len())
	}
}

func TestEncodeDifferentSampleRates(t *testing.T) {
	sampleRates := []beep.SampleRate{8000, 22050, 44100, 48000}

	for _, sr := range sampleRates {
		t.Run(fmt.Sprintf("%d", sr), func(t *testing.T) {
			format := beep.Format{
				SampleRate:  sr,
				NumChannels: 2,
				Precision:   2,
			}

			sine, err := generators.SineTone(format.SampleRate, 440)
			if err != nil {
				t.Fatalf("Failed to create sine tone: %v", err)
			}
			limited := beep.Take(format.SampleRate.N(500*time.Millisecond), sine)

			buf := &bytes.Buffer{}
			err = Encode(buf, limited, EncodeOptions{
				Format:  format,
				BitRate: 96000,
			})
			if err != nil {
				t.Fatalf("Encode failed for sample rate %d: %v", sr, err)
			}

			if buf.Len() == 0 {
				t.Fatal("Expected non-empty output")
			}
		})
	}
}

func TestEncodeDifferentBitRates(t *testing.T) {
	bitRates := []int{64000, 96000, 128000, 192000, 256000}

	for _, br := range bitRates {
		t.Run(fmt.Sprintf("%d", br), func(t *testing.T) {
			format := beep.Format{
				SampleRate:  44100,
				NumChannels: 2,
				Precision:   2,
			}

			sine, err := generators.SineTone(format.SampleRate, 440)
			if err != nil {
				t.Fatalf("Failed to create sine tone: %v", err)
			}
			limited := beep.Take(format.SampleRate.N(1*time.Second), sine)

			buf := &bytes.Buffer{}
			err = Encode(buf, limited, EncodeOptions{
				Format:  format,
				BitRate: br,
			})
			if err != nil {
				t.Fatalf("Encode failed for bit rate %d: %v", br, err)
			}

			if buf.Len() == 0 {
				t.Fatal("Expected non-empty output")
			}
		})
	}
}

func TestEncodeDefaultBitRate(t *testing.T) {
	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}

	sine, err := generators.SineTone(format.SampleRate, 440)
	if err != nil {
		t.Fatalf("Failed to create sine tone: %v", err)
	}
	limited := beep.Take(format.SampleRate.N(1*time.Second), sine)

	buf := &bytes.Buffer{}
	// Don't specify BitRate, should use default
	err = Encode(buf, limited, EncodeOptions{
		Format: format,
	})
	if err != nil {
		t.Fatalf("Encode with default bit rate failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("Expected non-empty output")
	}
}

func TestEncodeEmptyStream(t *testing.T) {
	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}

	// Create an empty streamer
	sine, err := generators.SineTone(format.SampleRate, 440)
	if err != nil {
		t.Fatalf("Failed to create sine tone: %v", err)
	}
	empty := beep.Take(0, sine)

	buf := &bytes.Buffer{}
	err = Encode(buf, empty, EncodeOptions{
		Format:  format,
		BitRate: 128000,
	})
	if err != nil {
		t.Fatalf("Encode empty stream failed: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatal("Expected empty output")
	}
}

func BenchmarkEncode(b *testing.B) {
	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}

	for b.Loop() {
		sine, err := generators.SineTone(format.SampleRate, 440)
		if err != nil {
			b.Fatalf("Failed to create sine tone: %v", err)
		}
		limited := beep.Take(format.SampleRate.N(1*time.Second), sine)

		buf := &bytes.Buffer{}
		err = Encode(buf, limited, EncodeOptions{
			Format:  format,
			BitRate: 128000,
		})
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
	}
}
