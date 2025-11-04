# aac

[![Go Reference](https://pkg.go.dev/badge/github.com/thde/aac.svg)](https://pkg.go.dev/github.com/thde/aac) [![tests](https://github.com/thde/aac/actions/workflows/test.yml/badge.svg)](https://github.com/thde/aac/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/thde/aac)](https://goreportcard.com/report/github.com/thde/aac)

This package provides a simple way to encode audio from [gopxl/beep](https://github.com/gopxl/beep)'s [`Streamer`](https://pkg.go.dev/github.com/gopxl/beep/v2#Streamer) interface directly to AAC format using the [gen2brain/aac-go](https://github.com/gen2brain/aac-go) encoder.

## Installation

```sh
go get -u github.com/thde/aac
```

## Usage

```go
package main

import (
    "os"

    "github.com/gopxl/beep/v2"
    "github.com/gopxl/beep/v2/wav"
    "github.com/thde/aac"
)

func main() {
    f, err := os.Open("input.wav")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    streamer, format, err := wav.Decode(f)
    if err != nil {
        panic(err)
    }
    defer streamer.Close()

    out, err := os.Create("output.aac")
    if err != nil {
        panic(err)
    }
    defer out.Close()

    err = aac.Encode(out, streamer, aac.EncodeOptions{
        Format:  format,
        BitRate: 128000, // 128 kbps
    })
    if err != nil {
        panic(err)
    }
}
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
