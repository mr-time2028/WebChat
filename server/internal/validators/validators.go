package validators

import "net/http"

// Validation contains information of validation (final result)
type Validation struct {
	Errors Errors
}

// New create an instance of Validation struct
func New() Validation {
	return Validation{
		Errors{
			map[string][]string{},
			http.StatusBadRequest,
		},
	}
}

// Valid check if validation was successful or not
func (v *Validation) Valid() bool {
	return len(v.Errors.MessageMap) == 0
}
