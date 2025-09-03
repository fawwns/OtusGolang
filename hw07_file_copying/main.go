package main

import (
	"flag"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" || to == "" {
		println("Usage: -from <source file> -to <destination file> [-offset N] [-limit N]")
		return
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		println("Error:", err.Error())
		return
	}

	println("File copied successfully!")
}
