package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	opened, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	openedStat, err := opened.Stat()
	if err != nil {
		return err
	}

	if err := validate(openedStat, offset); err != nil {
		return err
	}

	copied, err := os.Create(toPath)
	if err != nil {
		return err
	}

	if _, err = opened.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	limit = defineLimit(openedStat.Size(), offset, limit)

	if _, err = io.CopyN(copied, opened, limit); err != nil {
		return err
	}

	if err := opened.Close(); err != nil {
		return err
	}
	if err := copied.Close(); err != nil {
		return err
	}

	return nil
}

func defineLimit(fileSize int64, offset int64, limit int64) int64 {
	if fileSize-offset < limit || limit == 0 {
		return fileSize - offset
	}
	return limit
}

func validate(stat os.FileInfo, offset int64) error {
	if !stat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if stat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	return nil
}
