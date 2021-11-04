package adapters

import (
	"context"
	"errors"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/core"
	"log"
	"os"
	"strconv"
)

// FileSinkSpec is a specialisation of the SinkSpec interface for file sink
type FileSinkSpec struct {
	Mode *uint32 `yaml:"mode,omitempty"`
	UserID  *int `yaml:"user,omitempty"`
	GroupID *int `yaml:"group,omitempty"`
}

// FileSink is a file-based sink endpoint, where secrets are written to files
type FileSink struct {
	log *log.Logger
	fs afero.Fs
}

// NewFileSink creates a new FileSink, based on given Afero file system and a logger
func NewFileSink(log *log.Logger, fs afero.Fs) *FileSink {
	return &FileSink{
		log: log,
		fs: fs,
	}
}

// NewFileSinkSpec creates a FileSinkSpec struct from abstract map
func NewFileSinkSpec(in map[interface{}]interface{}) (FileSinkSpec, error) {
	var res FileSinkSpec

	var defaultMode uint32 = 400
	res.Mode = &defaultMode

	v, ex := in["mode"]
	if ex {
		vn, err := stringOrIntToI(v)
		if err != nil {
			return res, err
		}
		var vn2 uint32 = uint32(vn)
		res.Mode = &vn2
	}
	v, ex = in["user"]
	if ex {
		vn, err := stringOrIntToI(v)
		if err != nil {
			return res, err
		}
		res.UserID = &vn
	}
	v, ex = in["group"]
	if ex {
		vn, err := stringOrIntToI(v)
		if err != nil {
			return res, err
		}
		res.GroupID = &vn
	}

	return res, nil
}

func stringOrIntToI(v interface{}) (int, error) {
	var vn int
	var err error

	vs, ok := v.(string)
	if ok {
		vn, err = strconv.Atoi(vs)
		if err != nil {
			return vn, errors.New("mode parameter in file sink spec must be string with valid file mode")
		}
	} else {
		vn, ok = v.(int)
		if !ok {
			return vn, errors.New("mode parameter in file sink spec must be string or integer")
		}
	}

	return vn, nil
}

func (s *FileSink) Write(ctx context.Context, defaults *core.Defaults, secret *core.Secret, sink *core.Sink) error {

	spec, err := NewFileSinkSpec(sink.Spec)
	if err != nil {
		return err
	}

	f, err := s.fs.OpenFile(sink.Path, os.O_WRONLY|os.O_CREATE, os.FileMode(*spec.Mode))
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
	s.log.Printf("Written secret %s: %d bytes, to file %s, mode %d\n", secret.Name, n, sink.Path, os.FileMode(*spec.Mode))

	uid := -1
	gid := -1
	if spec.UserID != nil {
		uid = *spec.UserID
	}
	if spec.GroupID != nil {
		gid = *spec.GroupID
	}
	if err := s.fs.Chown(sink.Path, uid, gid); err != nil {
		return err
	}

	return nil
}

