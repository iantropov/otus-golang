package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptyFromFileName     = errors.New("empty from filename")
	ErrEmptyToFileName       = errors.New("empty to filename")
	ErrNegativeOffset        = errors.New("negative offset")
	ErrNegativeLimit         = errors.New("negative limit")
	ErrFromFileNotFound      = errors.New("from file not found")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return ErrEmptyFromFileName
	}
	if toPath == "" {
		return ErrEmptyToFileName
	}
	if offset < 0 {
		return ErrNegativeOffset
	}
	if limit < 0 {
		return ErrNegativeLimit
	}

	fileStat, err := os.Stat(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFromFileNotFound
		}
		return err
	}
	size := fileStat.Size()
	if size == 0 {
		return ErrUnsupportedFile
	}
	if offset >= size {
		return ErrOffsetExceedsFileSize
	}

	var inputFile *os.File
	inputFile, err = os.Open(fromPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	var outputFile *os.File
	outputFile, err = os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileStat.Mode())
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if limit == 0 || limit > size-offset {
		limit = size - offset
	}

	return copyFileContent(inputFile, outputFile, offset, limit)
}

func copyFileContent(inputFile, outputFile *os.File, offset, limit int64) error {
	const ChunkSize = 102400

	bufferSize := limit
	if bufferSize > ChunkSize {
		bufferSize = ChunkSize
	}
	buffer := make([]byte, bufferSize)

	bar := pb.Start64(limit)

	curPos := offset
	for curPos < limit+offset {
		read, readErr := inputFile.ReadAt(buffer, curPos)
		if readErr != nil && readErr != io.EOF {
			return readErr
		}
		if curPos+int64(read) > limit+offset {
			read = int(limit + offset - curPos)
		}

		curPos += int64(read)
		bar.Add64(int64(read))

		written := 0
		for written < read {
			nextWritten, writeErr := outputFile.Write(buffer[written:read])
			if writeErr != nil {
				return writeErr
			}
			written += nextWritten
		}
		if readErr == io.EOF {
			break
		}
	}
	bar.Finish()

	return nil
}
