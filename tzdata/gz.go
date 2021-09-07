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
		return fmt.Errorf("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("ExtractTarGz: Next() failed: %w", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			return fmt.Errorf("ExtractTarGz: directories are not supported: %w", err)
		case tar.TypeReg:
			if header.Name == name {
				if _, err := io.Copy(dst, tarReader); err != nil {
					return fmt.Errorf("ExtractTarGz: Copy() failed: %w", err)
				}
				return nil
			}
		default:
			return fmt.Errorf("ExtractTarGz: uknown type: %v in %v", header.Typeflag, header.Name)
		}
	}
}
