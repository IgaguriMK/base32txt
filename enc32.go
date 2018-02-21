package main

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"flag"
	"hash/crc32"
	"io"
	"log"
	"os"
)

func main() {
	var inputName string
	flag.StringVar(&inputName, "i", "", "Input file name.")
	var outputName string
	flag.StringVar(&outputName, "o", "", "Output file name.")

	flag.Parse()

	var input io.Reader
	if inputName == "" {
		input = os.Stdin
	} else {
		inf, err := os.Open(inputName)
		if err != nil {
			log.Fatalf("Can't open input file %q\n", inputName)
		}
		defer inf.Close()
		input = inf
	}

	var output io.Writer
	if outputName == "" {
		output = os.Stdout
	} else {
		outf, err := os.Create(outputName)
		if err != nil {
			log.Fatalf("Can't open output file %q\n", outputName)
		}
		defer outf.Close()
		output = outf
	}

	encode(input, output)
}

// Input 8 * 2 = 16 byte
// => 16 + 4 (CRC32) = 20 byte
// => 5 * 4 => 8 * 4 char
const bufSize = 16

var crcTable = crc32.MakeTable(crc32.IEEE)

func encode(r io.Reader, w io.Writer) {
	newline := []byte("\n")

	for {
		buf := make([]byte, bufSize)
		size := fillRead(r, buf)
		if size == 0 {
			return
		}

		bytes := buf[:size]
		crc := crc32.ChecksumIEEE(bytes)
		crcBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(crcBytes, crc)
		bytes = append(crcBytes, bytes...)

		encStr := base32.StdEncoding.EncodeToString(bytes)
		res := insertSpace([]byte(encStr))

		w.Write(res)
		w.Write(newline)
	}
}

func fillRead(r io.Reader, buf []byte) int {
	b := buf

	readsize := 0
	for {
		bs := len(b)
		size, err := r.Read(b)
		readsize += size
		if err == io.EOF {
			return readsize
		}
		if err != nil {
			log.Fatalf("Fatal error: %v\n", err)
		}

		if size == bs {
			return readsize
		}

		b = b[size:]
	}
}

const (
	spacePerChars = 4
)

func insertSpace(bs []byte) []byte {
	slices := make([][]byte, 0)

	for {
		if len(bs) <= spacePerChars {
			slices = append(slices, bs)
			break
		}

		b := bs[:spacePerChars]
		bs = bs[spacePerChars:]

		slices = append(slices, b)
	}

	return bytes.Join(slices, []byte(" "))
}
