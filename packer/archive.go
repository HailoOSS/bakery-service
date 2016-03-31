package packer

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func UnzipReader(r io.Reader, dst string) error {
	archive := filepath.Join(dst, "archive.zip")
	arcF, err := os.Create(archive)
	if err != nil {
		return err
	}

	if _, err := io.Copy(arcF, r); err != nil {
		return err
	}

	zipR, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	for _, f := range zipR.File {
		srcF, err := f.Open()
		if err != nil {
			return err
		}

		dstF, err := os.Create(filepath.Join(dst, f.Name))
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
