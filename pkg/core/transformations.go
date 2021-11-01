//go:generate mockgen -package mocks -destination=mocks/mock_transformationport.go go-secretshelper/pkg/core TransformationPort
package core

import "context"

type Transformations []Transformation

type Transformation struct {
	Var string `yaml:"var" validate:"required"`
	ToVar string `yaml:"toVar" validate:""`
}

// TransformationPort is the interface for a single transformation
type TransformationPort interface {

	// ProcessSecret applies the Transformation to the Secret and returns an updated Secret
	ProcessSecret(context.Context, *Defaults, *Secret, *Transformation) (*Secret, error)
}