# errctx

A Go package that extends error handling with contextual metadata.

## Overview

Tired of logging bad status, found "active" but wanted "deleted", and then having a nightmare filtering through the logs? errctx enables you to store data, keeping it through all the wrapping, so you can access it at the very end. Once logged, you can easily filter on the main error message or some metadata.

## Example

```Go
package main

import (
	"context"
	"errctx"
	"errors"
	"os"

	"github.com/rs/zerolog"
)

var ErrMultiOccurrenceChar = errctx.New("multi occurrence char")

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	ctx := logger.WithContext(context.Background())

	err := checkMultiOccurrenceChar("helloooo")
	if err != nil {
		zeroLog(ctx, "error in main", err)
	}
}

func zeroLog(ctx context.Context, msg string, e error) {
	var errctx errctx.ErrCtx
	if errors.As(e, &errctx) {
		zerolog.Ctx(ctx).Error().Err(errctx).Fields(errctx.Values()).Msg(msg)

		return
	}

	zerolog.Ctx(ctx).Err(e).Msg(msg)
}

func checkMultiOccurrenceChar(s string) error {
	chars := make(map[rune]int)
	multiChar := make([]string, 0)

	for _, char := range s {
		chars[char]++
		if chars[char] == 2 {
			multiChar = append(multiChar, string(char))
		}
	}

	if len(multiChar) > 0 {
		return errctx.NewFromErr(ErrMultiOccurrenceChar).With("chars", multiChar)
	}

	return nil
}
```

Output with errctx:
{"level":"error","error":"multi occurrence char","chars":["l","o"],"time":"2025-07-13T23:07:44+02:00","message":"error in main"}

Output  without errctx:
{"level":"error","error":"multi occurrence char: \"l\", \"o\","time":"2025-07-13T23:07:44+02:00","message":"error in main"}
