//go:generate mockgen -package mocks -destination=mocks/mock_transformationport.go go-secretshelper/pkg/core TransformationPort
package core

import "context"

// Transformations is an array of Transformation structs
type Transformations []*Transformation

// Transformation describe a single transformation
type Transformation struct {
	// Input is the list of input variables for this transformation. These must have
	// been defined as secrets or must have been processed before by other transformations
	Input []string `yaml:"in" validate:"required,dive,required"`

	// Output is the name of output variable. The result of the transformation will go here.
	Output string `yaml:"out" validate:"required"`

	// Type is the type of transformation
	Type string `yaml:"type" validate:"required"`

	// Spec is the generic specification for a transformation of a given type
	Spec TransformationSpec `yaml:"spec" validate:""`
}

type TransformationSpec map[interface{}]interface{}

// TransformationPort is the interface for a single transformation
type TransformationPort interface {
	// ProcessSecret applies the Transformation, using given Secrets and returns an updated Secret
	ProcessSecret(context.Context, *Defaults, *Secrets, *Transformation) (*Secret, error)
}
