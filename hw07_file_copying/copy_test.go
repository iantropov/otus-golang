package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		title, from, to string
		offset, limit   int64
		err             error
	}{
		{title: "empty FROM", from: "", to: "", offset: 0, limit: 0, err: ErrEmptyFromFileName},
		{title: "empty TO", from: "asd", to: "", offset: 0, limit: 0, err: ErrEmptyToFileName},
		{title: "negative OFFSET", from: "asd", to: "asd", offset: -1, limit: 0, err: ErrNegativeOffset},
		{title: "negative LIMIT", from: "asd", to: "asd", offset: 0, limit: -1, err: ErrNegativeLimit},
		{title: "unsupported file FROM", from: "/dev/urandom", to: "asd", offset: 0, limit: 0, err: ErrUnsupportedFile},
		{title: "file FROM not found", from: "/asdasdasd", to: "asd", offset: 0, limit: 0, err: ErrFromFileNotFound},
		{title: "OFFSET exceeds file size", from: "./testdata/out_offset0_limit10.txt", to: "asd", offset: 11, limit: 0, err: ErrOffsetExceedsFileSize},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.Truef(t, errors.Is(err, tc.err), "actual error %q", err)
		})
	}
}
