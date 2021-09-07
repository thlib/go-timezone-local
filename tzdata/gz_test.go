package tzdata

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestExtractTarGz(t *testing.T) {

	file, err := os.Open(`output.tar.gz`)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer file.Close()

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err = ExtractTarGz(file, w, "backward")
	if err != nil {
		t.Errorf("%v", err)
	}

	// fmt.Printf("Data: %v", b.String())
}
