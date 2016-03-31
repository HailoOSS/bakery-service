package packer

import (
	"archive/zip"
	"io"
)

func UnzipReader(r io.Reader, dst string) error {
	zipR, err := zip.NewReader(rc, len(r))
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
