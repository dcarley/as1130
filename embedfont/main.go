package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	fontFilename := flag.String("f", "font.ttf", "Font filename")
	outFilename := flag.String("o", "font.go", "Go filename")
	packageName := flag.String("p", "main", "Package name")
	flag.Parse()

	file, err := os.Open(*fontFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var zbuf bytes.Buffer
	zw := gzip.NewWriter(&zbuf)
	io.Copy(zw, file)
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	fmt.Fprintf(&out, `// DO NOT EDIT
	// rendered by go generate and as1130/embedfont

	package %s

	import (
		"bytes"
		"compress/gzip"
		"io"
	)

	// fontBytes returns an uncompressed byte slice of a truetype font file.
	var fontBytes = func() ([]byte, error) {
		var buf bytes.Buffer

		zbuf := bytes.NewBuffer(fontCompressed)
		zr, err := gzip.NewReader(zbuf)
		if err != nil {
			return buf.Bytes(), err
		}

		io.Copy(&buf, zr)

		return buf.Bytes(), zr.Close()
	}
	`, *packageName)

	fmt.Fprintf(&out, "var fontCompressed = []byte{")
	for i, x := range zbuf.Bytes() {
		if i&15 == 0 {
			out.WriteByte('\n')
		}
		fmt.Fprintf(&out, "%#02x,", x)
	}
	fmt.Fprintf(&out, "\n}\n")

	dst, err := format.Source(out.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(*outFilename, dst, 0666); err != nil {
		log.Fatal(err)
	}
}
