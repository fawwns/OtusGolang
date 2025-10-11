package main

import (
	"errors"
	"io"
	"os"

	//nolint:depguard
	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer file.Close()

	copyFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer copyFile.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	size := fileInfo.Size()

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	bytesToCopy := limit
	if limit == 0 || offset+limit > size {
		bytesToCopy = size - offset
	}

	if bytesToCopy == 0 {
		return nil
	}

	bar := pb.Full.Start64(bytesToCopy)
	barReader := bar.NewProxyReader(file)

	_, err = io.CopyN(copyFile, barReader, bytesToCopy)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	bar.Finish()
	return nil
}
