package packer

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
)

func UnzipReader(r io.Reader, dst string) error {
	var buf bytes.Buffer

	buf.ReadFrom(r)

	zipR, err := zip.NewReader(buf, buf.Len())
	if err != nil {
		return err
	}

	for _, f := range zipR.File {
		srcF, err := f.Open()
		if err != nil {
			return err
		}

		dstF, err := os.Create(f.Name)
		if err != nil {
			srcF.Close()
			return err
		}
		_, err = io.Copy(dstF, srcF)
		srcF.Close()
		dstF.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
