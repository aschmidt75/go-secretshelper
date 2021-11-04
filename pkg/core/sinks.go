package core

import (
	"context"
	"fmt"
)

type Sinks []*Sink

type Sink struct {
	Type string `yaml:"type" validate:"required"`
	Path string `yaml:"path" validate:"required"`
	Var string `yaml:"var" validate:"required"`
	Spec SinkSpec `yaml:"spec" validate:""`
}

type SinkSpec map[interface{}]interface{}

func (s Sink) String() string {
	return fmt.Sprintf("Sink:[Var=%s, Type=%s, Path=%s]", s.Var, s. Type, s.Path)
}

// SinkWriterPort is able to write a secret into a defined sink
type SinkWriterPort interface {

	// Write takes the raw content of given secret and writes it to the sink using the defaults
	Write(context.Context, *Defaults, *Secret, *Sink) error
}
