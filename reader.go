package aac

import (
	"encoding/binary"
	"io"

	"github.com/gopxl/beep/v2"
)

// streamerReader is an adapter that makes a [beep.Streamer] behave like an [io.Reader],
// providing audio data as 16-bit signed little-endian PCM as required by [aac.Encode].
// The source [beep.Format] precision is ignored because samples are always converted to 16-bit PCM.
type streamerReader struct {
	s      beep.Streamer
	format beep.Format
	buf    [][2]float64
}

// newStreamReader creates a new reader for the given streamer and format.
func newStreamReader(s beep.Streamer, format beep.Format) *streamerReader {
	return &streamerReader{
		s:      s,
		format: format,
		// The buffer size can be adjusted for performance.
		// 512 samples ≈ 2KB for stereo 16-bit PCM
		buf: make([][2]float64, 512),
	}
}

// Read implements the [io.Reader] interface. It streams audio samples, converts
// them to PCM, and writes them into the byte slice p.
func (sr *streamerReader) Read(p []byte) (n int, err error) {
	// Determine how many samples are needed to fill the byte buffer p.
	// Each sample is 2 bytes (int16) per channel.
	samplesNeeded := len(p) / (sr.format.NumChannels * 2)
	if samplesNeeded == 0 {
		return 0, nil
	}

	// Ensure our internal sample buffer is large enough.
	if len(sr.buf) < samplesNeeded {
		sr.buf = make([][2]float64, samplesNeeded)
	}

	// Stream audio data from the beep.Streamer.
	numStreamed, ok := sr.s.Stream(sr.buf[:samplesNeeded])
	if !ok && numStreamed == 0 {
		return 0, io.EOF
	}

	// Convert float64 samples to int16 PCM
	for i := range numStreamed {
		for ch := 0; ch < sr.format.NumChannels; ch++ {
			// Clamp the sample value to the [-1.0, 1.0] range.
			val := sr.buf[i][ch]
			if val < -1.0 {
				val = -1.0
			}
			if val > 1.0 {
				val = 1.0
			}
			// Convert float64 [-1.0, 1.0] to int16 [-32768, 32767]
			sample := int16(val * 32767)

			// Write the int16 sample as two bytes (little-endian) into the buffer p.
			offset := (i*sr.format.NumChannels + ch) * 2
			binary.LittleEndian.PutUint16(p[offset:], uint16(sample))
		}
	}

	return numStreamed * sr.format.NumChannels * 2, nil
}
