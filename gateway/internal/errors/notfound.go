package errors

import "fmt"

type NotFoundError struct {
	entityName string
	id         any
}

func NewNotFoundError(entityName string, id any) NotFoundError {
	return NotFoundError{
		entityName: entityName,
		id:         id,
	}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%q not found by identificator %v", e.entityName, e.id)
}
