package main

import (
	"bufio"
	"encoding/base32"
	"encoding/binary"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"strings"
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

	decode(input, output)
}

// Input 8 * 2 = 16 byte
// => 16 + 4 (CRC32) = 20 byte
// => 5 * 4 => 8 * 4 char
const bufSize = 16

var crcTable = crc32.MakeTable(crc32.IEEE)

func decode(r io.Reader, w io.Writer) {
	sc := bufio.NewScanner(r)

	lineNum := 0
	for sc.Scan() {
		lineNum++

		line := sc.Text()
		line = strings.Replace(line, " ", "", -1)
		bs, err := base32.StdEncoding.DecodeString(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n")
			log.Fatalf("Decode error: Line %d: %v\n", lineNum, err)
		}

		if len(bs) < 5 {
			fmt.Fprintf(os.Stderr, "\n")
			log.Fatalf("Fatal: Line %d is too short.\n", lineNum)
		}

		crcBs := bs[0:4]
		contentBs := bs[4:]

		tobeCRC := binary.LittleEndian.Uint32(crcBs)
		acctualCRC := crc32.ChecksumIEEE(contentBs)

		if tobeCRC != acctualCRC {
			fmt.Fprintf(os.Stderr, "\n")
			log.Fatalf("Detect input error at line %d.\n", lineNum)
		}

		w.Write(contentBs)
	}
}
