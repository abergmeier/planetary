package collection

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInstall(t *testing.T) {
	prefix := "https://github.com/abergmeier/ansible-collection-binary-builtin/releases/download/v0.0.5/ansible-collection-binary-builtin_0.0.5"

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tempDir)

	i := &Installer{
		installPath: tempDir,
	}
	err = i.Run(prefix)
	if err != nil {
		t.Fatalf("Unexpected error from Install: %s", err)
	}

	expectedPaths := []string{
		"plugins/modules/assert",
		"plugins/modules/git",
		"plugins/modules/unarchive",
	}

	for _, expectedPath := range expectedPaths {
		_, err = os.Stat(filepath.Join(tempDir, expectedPath))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
	}
}
