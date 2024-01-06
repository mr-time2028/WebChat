package validators

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	MinErrMsg               = `this field must be a minimum length of %d characters`
	MaxErrMsg               = `the maximum length for this field must be %d characters`
	RequiredErrMsg          = `this field is required`
	JSONEmptyBodyErrMsg     = `request body is empty`
	JSONSyntaxErrMsg        = `request body contains invalid JSON syntax at position %s`
	JSONUnmarshalTypeErrMsg = `request body contains invalid data type at position %s`
	JSONValueErrMsg         = `body must have only a single json value`
	JSONUnknownField        = `unknown field in JSON: %s`
)

// MinLengthTag check min length of a json field
func (v *Validation) MinLengthTag(field reflect.StructField, fieldValue reflect.Value) {
	jsonFieldName := field.Tag.Get("json")
	minTag := field.Tag.Get("min")

	if minTag != "" && fieldValue.Kind() == reflect.String && len(fieldValue.String()) > 0 {
		minLength, _ := strconv.Atoi(minTag)
		if len(fieldValue.String()) < minLength {
			v.Errors.Add(jsonFieldName, fmt.Sprintf(MinErrMsg, minLength))
		}
	}
}

// MaxLengthTag check max length of a json field
func (v *Validation) MaxLengthTag(field reflect.StructField, fieldValue reflect.Value) {
	jsonFieldName := field.Tag.Get("json")
	maxTag := field.Tag.Get("max")

	if maxTag != "" && fieldValue.Kind() == reflect.String && len(fieldValue.String()) > 0 {
		maxLength, _ := strconv.Atoi(maxTag)
		if len(fieldValue.String()) > maxLength {
			v.Errors.Add(jsonFieldName, fmt.Sprintf(MaxErrMsg, maxLength))
		}
	}
}

// RequiredTag force client to send a specific json field
func (v *Validation) RequiredTag(field reflect.StructField, fieldValue reflect.Value) {
	fieldType := fieldValue.Type()

	jsonFieldName := field.Tag.Get("json")
	requiredTag := field.Tag.Get("required")

	if requiredTag == "true" {
		zeroValue := reflect.Zero(fieldType)
		isZero := reflect.DeepEqual(fieldValue.Interface(), zeroValue.Interface())
		if isZero {
			v.Errors.Add(jsonFieldName, RequiredErrMsg)
		}
	}
}

// JsonValidation decode json and do json validation
func (v *Validation) JsonValidation(r *http.Request, data interface{}) {
	if r.Body == nil {
		v.Errors.Add("json", JSONEmptyBodyErrMsg)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// validation for json
	err := dec.Decode(data)
	if err != nil {
		if err == io.EOF {
			v.Errors.Add("json", JSONEmptyBodyErrMsg)
		} else if syntaxErr, ok := err.(*json.SyntaxError); ok {
			v.Errors.Add("json", fmt.Sprintf(JSONSyntaxErrMsg, strconv.Itoa(int(syntaxErr.Offset))))
		} else if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			v.Errors.Add("json", fmt.Sprintf(JSONUnmarshalTypeErrMsg, strconv.Itoa(int(unmarshalErr.Offset))))
		} else if strings.HasPrefix(err.Error(), "json: unknown field") {
			v.Errors.Add("json", fmt.Sprintf(JSONUnknownField, strings.TrimPrefix(err.Error(), "json: unknown field ")))
		} else {
			v.Errors.Add("json", err.Error())
		}
		return
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		v.Errors.Add("json", JSONValueErrMsg)
		return
	}

	// validation for each field
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()

	for i := 0; i < dataType.NumField(); i++ {
		//validator := New() // we want to do validation for each json field

		field := dataType.Field(i)
		fieldValue := dataValue.Field(i)

		// required field validation
		v.RequiredTag(field, fieldValue)

		// minimum length validation
		v.MinLengthTag(field, fieldValue)

		// maximum length validation
		v.MaxLengthTag(field, fieldValue)
	}
}
