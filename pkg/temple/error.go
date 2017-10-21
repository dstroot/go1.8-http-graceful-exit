package temple

import (
	"fmt"
)

// TemplateError describes template errors
type TemplateError struct {
	Format     string
	Parameters []interface{}
}

func (e *TemplateError) Error() string {
	return fmt.Sprintf(e.Format, e.Parameters...)
}

// Errf returns a TemplateError
func Errf(format string, parameters ...interface{}) error {
	return &TemplateError{
		Format:     format,
		Parameters: parameters,
	}
}
