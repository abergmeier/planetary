package collection

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type Installer struct {
	installPath string
}

func NewInstaller() *Installer {
	return &Installer{
		installPath: "~/.ansible/collections/planetary",
	}
}

func (i *Installer) Run(downloadPrefix string) error {

	suffix := fmt.Sprintf("_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)

	foo := downloadPrefix + suffix

	out, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(foo)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Get of `%s` failed: %s", foo, resp.Status)
	}

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	defer gr.Close()

	tr := tar.NewReader(gr)

	// get the next file entry
	for {
		h, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Next failed when processing %s: %s", foo, err)
		}

		p := filepath.Join(i.installPath, h.Name)

		err = extractFile(tr, h.Size, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFile(tr *tar.Reader, size int64, path string) error {

	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := [1024]byte{}
	for {

		read, err := tr.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		_, err = f.Write(buf[:read])
		if err != nil {
			return err
		}
	}
}
