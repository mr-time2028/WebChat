package validators

// Errors keeps all errors while validating
type Errors struct {
	MessageMap map[string][]string
	Code       int
}

// Add simply adds a key value error to the MessageMap
func (e Errors) Add(field, message string) {
	e.MessageMap[field] = append(e.MessageMap[field], message)
}

// Get simply gets a value from MessageMap by given key
func (e Errors) Get(field string) string {
	es := e.MessageMap[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
