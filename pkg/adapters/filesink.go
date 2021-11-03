package adapters

import (
	"context"
	"errors"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/core"
	"log"
	"os"
)

type FileSink struct {
	log *log.Logger
	fs afero.Fs
}

// NewFileSink creates a new FileSink, based on given Afero file system
func NewFileSink(log *log.Logger, fs afero.Fs) *FileSink {
	return &FileSink{
		log: log,
		fs: fs,
	}
}

func (s *FileSink) Write(ctx context.Context, defaults *core.Defaults, secret *core.Secret, sink *core.Sink) error {

	f, err := s.fs.OpenFile(sink.Path, os.O_WRONLY|os.O_CREATE, 0200)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := f.Write(secret.RawContent)
	if err != nil {
		return err
	}
	if n != len(secret.RawContent) {
		return errors.New("invalid number of bytes")
	}
	s.log.Printf("Written secret %s: %d bytes, to file %s\n", secret.Name, n, sink.Path)

	return nil
}

