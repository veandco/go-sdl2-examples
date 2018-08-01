package main

import (
	"fmt"
	"io"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

var rwops *sdl.RWops
var writer io.Writer
var reader io.Reader
var seeker io.Seeker
var closer io.Closer
var data = []byte{0, 10, 20, 30}
var err error

func size() {
	if n, err := rwops.Size(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("size:", n)
	}
}

func read() {
	buf := make([]byte, 4)
	reader.Read(buf[:])
	fmt.Println("read:", buf[0], buf[1], buf[2], buf[3])
}

func tell() {
	if n, err := rwops.Tell(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("tell:", n)
	}
}

func seek() {
	// Seek
	if n, err := seeker.Seek(0, io.SeekStart); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("seek: offset is", n)
	}
}

func write() {
	if n, err := writer.Write([]byte{40, 50, 60, 70}); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("write:", n, "bytes")
	}
}

func main() {
	rwops, err = sdl.RWFromMem(data)
	if err != nil {
		log.Println(err)
	}

	writer = rwops
	reader = rwops
	seeker = rwops
	closer = rwops

	defer closer.Close()

	size()  // Print RWops data size
	tell()  // Print current data offset
	read()  // Read data via RWops
	tell()  // Print current data offset
	seek()  // Seek to start of the data
	write() // Write some data via RWops
	tell()  // Print current data offset
	seek()  // Seek to start of the data
	tell()  // Print current data offset
	read()  // Read data via RWops again
}
