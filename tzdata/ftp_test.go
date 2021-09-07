package tzdata

import (
	"io/ioutil"
	"testing"
)

func TestFTPDownload(t *testing.T) {
	buf, err := FTPDownload(`ftp://ftp.iana.org/tz/tzdata-latest.tar.gz`)
	if err != nil {
		t.Errorf("%v", err)
	}

	err = ioutil.WriteFile("output.tar.gz", buf.Bytes(), 0644)
	if err != nil {
		t.Errorf("%v", err)
	}
}

// func TestFTPDownload2(t *testing.T) {
// 	err := FTPDownload2()
// 	if err != nil {
// 		t.Errorf("%v", err)
// 	}
// }

// func TestFiles(t *testing.T) {
// 	b1, err := ioutil.ReadFile(`output.tar.gz`)
// 	if err != nil {
// 		panic(err)
// 	}
// 	b2, err := ioutil.ReadFile(`tzdata-latest.tar.gz`)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for i := 0; i < len(b1); i++ {
// 		if b1[i] != b2[i] {
// 			t.Errorf("The wrong value: %v vs %v", b1[i], b2[i])
// 			t.Errorf(string(b1[i-1 : 30]))

// 			t.Errorf(string(b2[i-1 : 30]))
// 			break
// 		}
// 	}
// }

// func BenchmarkFTPDownload(b *testing.B) {

// 	server, client := net.Pipe()
// 	go func() {
// 		// Do some stuff
// 		server.Close()
// 	}()

// 	// Do some stuff
// 	client.Close()

// 	for i := 0; i < b.N; i++ {

// 	}
// }
