package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func addToTar(fi os.FileInfo, file string, root string, tarWriter *tar.Writer, isDir bool) error {
	header, err := tar.FileInfoHeader(fi, fi.Name())
	if err != nil {
		return err
	}

	header.Name = strings.TrimPrefix(strings.Replace(file, root, "", -1), string(filepath.Separator))
	if isDir {
		header.Typeflag = tar.TypeDir
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if isDir {
		return nil
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if _, err := io.Copy(tarWriter, f); err != nil {
		return err
	}

	return nil
}

func archivePath(root string, tarWriter *tar.Writer) error {
	return filepath.Walk(root, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.Mode().IsDir() {
			return addToTar(fi, file, root, tarWriter, true)
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		return addToTar(fi, file, root, tarWriter, false)
	})
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error get current dir: %v", err)
	}

	file, err := os.CreateTemp("", "dir2str*.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Remove(file.Name()) }()

	gzipWriter := gzip.NewWriter(file)
	//defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	//defer tarWriter.Close()

	err = archivePath(wd, tarWriter)
	if err != nil {
		log.Fatal(err)
	}

	err = tarWriter.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = gzipWriter.Close()
	if err != nil {
		log.Fatal(err)
	}

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	content := make([]byte, fi.Size())
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	n, err := file.ReadAt(content, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s#%s\n", base64.StdEncoding.EncodeToString([]byte(filepath.Base(wd))), base64.StdEncoding.EncodeToString(content[:n]))
}
