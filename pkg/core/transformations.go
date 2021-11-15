//go:generate mockgen -package mocks -destination=mocks/mock_transformationport.go go-secretshelper/pkg/core TransformationPort
package core

import "context"

// Transformations is an array of Transformation structs
type Transformations []*Transformation

// Transformation describe a single transformation
type Transformation struct {
	// Input is the list of input variables for this transformation. These must have
	// been defined as secrets or must have been processed before by other transformations
	Input []string `yaml:"in" validate:"required"`

	// Output is the name of output variable. The result of the transformation will go here.
	Output string `yaml:"out" validate:"required"`

	// Type is the type of transformation
	Type string `yaml:"type" validate:"required"`

	// Spec is the generic specification for a transformation of a given type
	Spec map[interface{}]interface{} `yaml:"spec" validate:""`
}

// TransformationPort is the interface for a single transformation
type TransformationPort interface {
	// ProcessSecret applies the Transformation to the Secret and returns an updated Secret
	ProcessSecret(context.Context, *Defaults, *Secret, *Transformation) (*Secret, error)
}
