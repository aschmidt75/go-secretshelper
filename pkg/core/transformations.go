//go:generate mockgen -package mocks -destination=mocks/mock_transformationport.go go-secretshelper/pkg/core TransformationPort
package core

import "context"

// Transformations is an array of Transformation structs
type Transformations []Transformation

// Transformation describe a single transformation spec, from
// a named secred (Var) to a new, transformed one (ToVar)
type Transformation struct {
	Var string `yaml:"var" validate:"required"`
	ToVar string `yaml:"toVar" validate:""`
}

// TransformationPort is the interface for a single transformation
type TransformationPort interface {

	// ProcessSecret applies the Transformation to the Secret and returns an updated Secret
	ProcessSecret(context.Context, *Defaults, *Secret, *Transformation) (*Secret, error)
}