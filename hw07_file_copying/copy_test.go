package main

import (
	"os"
	"testing"
)

// вспомогательная функция — сравнение файлов.
func equalFiles(t *testing.T, path1, path2 string) {
	t.Helper()

	data1, err := os.ReadFile(path1)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path1, err)
	}
	data2, err := os.ReadFile(path2)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path2, err)
	}

	if string(data1) != string(data2) {
		t.Errorf("files %s and %s differ", path1, path2)
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		name   string
		offset int64
		limit  int64
		want   string // путь к эталонному файлу.
		err    error
	}{
		{"offset0_limit0", 0, 0, "testdata/out_offset0_limit0.txt", nil},
		{"offset0_limit10", 0, 10, "testdata/out_offset0_limit10.txt", nil},
		{"offset0_limit1000", 0, 1000, "testdata/out_offset0_limit1000.txt", nil},
		{"offset0_limit10000", 0, 10000, "testdata/out_offset0_limit10000.txt", nil},
		{"offset100_limit1000", 100, 1000, "testdata/out_offset100_limit1000.txt", nil},
		{"offset6000_limit1000", 6000, 1000, "testdata/out_offset6000_limit1000.txt", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создаём временный файл-результат.
			dst, err := os.CreateTemp("", "copy-*.txt")
			if err != nil {
				t.Fatal(err)
			}
			dst.Close()
			defer os.Remove(dst.Name())

			// вызов функции Copy.
			err = Copy("testdata/input.txt", dst.Name(), tt.offset, tt.limit)
			if err != tt.err {
				t.Fatalf("unexpected error: got %v, want %v", err, tt.err)
			}

			// сравнение результата с эталонным файлом.
			equalFiles(t, dst.Name(), tt.want)
		})
	}
}

func TestCopyOffsetTooLarge(t *testing.T) {
	dst, _ := os.CreateTemp("", "copy-*.txt")
	dst.Close()
	defer os.Remove(dst.Name())

	err := Copy("testdata/input.txt", dst.Name(), 999999, 0)
	if err != ErrOffsetExceedsFileSize {
		t.Errorf("expected %v, got %v", ErrOffsetExceedsFileSize, err)
	}
}
