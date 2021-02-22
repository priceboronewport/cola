package tempfile

import (
	"fmt"
	"os"
)

type TempFile struct {
	Filename string
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func New(id string) *TempFile {
	path := os.TempDir()
	seq := 0
	filename := fmt.Sprintf("%s/~tmp_%s_%d%d", path, id, os.Getpid(), seq)
	for FileExists(filename) {
		seq++
		filename = fmt.Sprintf("%s/~tmp_%s_%d%d", path, id, os.Getpid(), seq)
	}
	tf := TempFile{Filename: filename}
	return &tf
}

func (tf *TempFile) Close() {
	os.Remove(tf.Filename)
}
