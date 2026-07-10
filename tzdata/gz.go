package tzdata

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

func ExtractTarGz(src io.Reader, dst io.Writer, name string) error {
	uncompressedStream, err := gzip.NewReader(src)
	if err != nil {
		return fmt.Errorf("ExtractTarGz: NewReader failed: %w", err)
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			return fmt.Errorf("ExtractTarGz: %q not found", name)
		}

		if err != nil {
			return fmt.Errorf("ExtractTarGz: Next() failed: %w", err)
		}

		if header.Typeflag == tar.TypeReg || header.Typeflag == tar.TypeRegA {
			if header.Name == name {
				if _, err := io.Copy(dst, tarReader); err != nil {
					return fmt.Errorf("ExtractTarGz: Copy() failed: %w", err)
				}
				return nil
			}
		}
	}
}
