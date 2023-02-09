package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

const CHUNK_SIZE = 102400

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	// from = "/dev/urandom"
	// to = "app2"
	// offset = 10
	// limit = 20
	fmt.Printf("FROM: %s, TO: %s, OFFSET: %d, LIMIT: %d\n", from, to, offset, limit)

	err := Copy(from, to, offset, limit)
	if err != nil {
		panic(err)
	}
}
