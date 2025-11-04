// Package aac provides AAC encoding functionality for beep audio streams.
package aac

import (
	"fmt"
	"io"

	"github.com/gen2brain/aac-go"
	"github.com/gopxl/beep/v2"
)

// EncodeOptions defines the options for AAC encoding.
type EncodeOptions struct {
	// Embed the source audio format. Precision must be 2 (16-bit PCM) to match the AAC encoder requirements.
	beep.Format
	// BitRate in bits/sec.
	// If 0, defaults to 64000 bits/sec.
	BitRate int
}

// Encode writes a [beep.Streamer] to an [io.Writer] in AAC format.
// Returns an error if encoding fails, if the format has an unsupported precision, or if the format has more than 2 channels.
func Encode(w io.Writer, s beep.Streamer, opts EncodeOptions) error {
	// Validate sample precision (encoder expects 16-bit PCM input)
	if opts.Precision != 2 {
		return fmt.Errorf(
			"unsupported precision: %d (must be 2 bytes per sample for 16-bit PCM)",
			opts.Precision,
		)
	}

	// Validate number of channels (beep supports stereo only)
	if opts.NumChannels < 1 || opts.NumChannels > 2 {
		return fmt.Errorf("unsupported number of channels: %d (must be 1 or 2)", opts.NumChannels)
	}

	enc, err := aac.NewEncoder(w, &aac.Options{
		SampleRate:  int(opts.SampleRate),
		NumChannels: opts.NumChannels,
		BitRate:     opts.BitRate,
	})
	if err != nil {
		return fmt.Errorf("failed to create AAC encoder: %w", err)
	}

	if err := enc.Encode(newStreamReader(s, opts.Format)); err != nil {
		return fmt.Errorf("failed to encode audio: %w", err)
	}

	if err := enc.Close(); err != nil {
		return fmt.Errorf("failed to close encoder: %w", err)
	}

	return nil
}
