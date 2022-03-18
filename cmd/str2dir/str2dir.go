package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: str2dir <string>")
	}

	a := strings.Split(os.Args[1], "#")
	if len(a) != 2 {
		log.Fatal("invalid string")
	}

	dn, err := base64.StdEncoding.DecodeString(a[0])
	if err != nil {
		log.Fatal(err)
	}
	dirname := string(dn)
	b, err := base64.StdEncoding.DecodeString(a[1])
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	target := filepath.Join(wd, dirname)
	_, err = os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(target, 0755)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("target %s exists", target)
	}

	err = extract(target, bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("unpacked files to '%s'\n", dirname)
}

func extract(dst string, r io.Reader) error {
	reader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tarReader); err != nil {
				return err
			}

			err = f.Close()
			if err != nil {
				return err
			}
		}
	}
}
