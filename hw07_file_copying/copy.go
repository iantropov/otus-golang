package main

import (
	"errors"
	"io"
	"os"
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
		} else {
			return err
		}
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

	if limit == 0 {
		limit = size
	}

	bufferSize := limit
	if bufferSize > CHUNK_SIZE {
		bufferSize = CHUNK_SIZE
	}
	buffer := make([]byte, bufferSize)

	curPos := offset
	for curPos < offset+limit {
		read, readErr := inputFile.ReadAt(buffer, curPos)
		if readErr != nil && readErr != io.EOF {
			return readErr
		}
		if curPos+int64(read) > offset+limit {
			read = int(offset + limit - curPos)
		}
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
		curPos += int64(read)
	}
	return nil
}
