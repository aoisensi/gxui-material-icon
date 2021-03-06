// +build ignore

package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	url  = "https://github.com/google/material-design-icons/raw/master/iconfont/MaterialIcons-Regular.ttf"
	name = "font_gen.go"
)

func main() {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	data, err = deflate(data)
	if err != nil {
		log.Fatal(err)
	}

	w := new(bytes.Buffer)
	fmt.Fprintln(w, "// GENERATED BY generate.go")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "package materialicon")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "var data = []byte{")

	for len(data) > 0 {
		n := 16
		if n > len(data) {
			n = len(data)
		}
		for _, c := range data[:n] {
			fmt.Fprintf(w, "0x%02x,", c)
		}
		fmt.Fprintln(w)
		data = data[n:]
	}

	fmt.Fprintln(w, "}")
	wbytes := w.Bytes()

	b, err := format.Source(wbytes)
	if err != nil {
		os.Stderr.Write(wbytes)
		log.Fatal(err)
	}

	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(b); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func deflate(src []byte) (dst []byte, err error) {
	buf := new(bytes.Buffer)
	w, err := flate.NewWriter(buf, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(src); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
