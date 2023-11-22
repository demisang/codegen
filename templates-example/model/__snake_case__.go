package models

import (
	"errors"
	"time"
)

var Err__PascalCase__NotFound = errors.New("__camelCase__ not found")

type __PascalCase__ID int

type __PascalCase__ struct {
	ID   __PascalCase__ID `json:"id"`
	Name string           `json:"name"`
	// @TODO database fields here
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type __PascalCase__CreateDto struct {
	Name string `json:"name"`
	// @TODO database fields here
}

type __PascalCase__UpdateDto struct {
	Name string `json:"name"`
	// @TODO database fields here
}

type __PascalCase__Filter struct {
	Limit  uint64
	Offset uint64
	IDs    *[]__PascalCase__ID
	// @TODO database fields here
	SearchString *string
}

// New__PascalCase__UpdateDtoFromModel helps construct pre-filled update DTO
func New__PascalCase__UpdateDtoFromModel(model __PascalCase__) __PascalCase__UpdateDto {
	return __PascalCase__UpdateDto{
		Name: model.Name,
		// @TODO database fields here
	}
}
