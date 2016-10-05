package archiver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkCompress(b *testing.B) {
	tmp, err := ioutil.TempDir("", "archiver")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	// Test creating archive
	outfile := filepath.Join(tmp, "test")

	for _, c := range []struct {
		name string
		cf   CompressFunc
	}{
		{"Zip", Zip},
		{"Rar", Rar},
		{"TarBz2", TarBz2},
		{"TarGz", TarGz},
	} {
		cf := c.cf
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cf(outfile, []string{"testdata"})
			}
		})
	}
}
