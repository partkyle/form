package form

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ParamType int

const (
	Query ParamType = iota
	FormValue
)

type Form struct {
	Fields []*Field
}

func NewForm() *Form {
	return &Form{Fields: make([]*Field, 0)}
}

func (f *Form) AddField(name string, opts ...FieldOpt) *Field {
	field := &Field{Name: name}
	for _, opt := range opts {
		opt(field)
	}

	f.Fields = append(f.Fields, field)

	return field
}

func (f *Form) Parse(r *http.Request) error {
	multiErr := make(MultiError, 0)

fieldsLoop:
	for _, field := range f.Fields {
		if field.Required {
			if _, ok := r.URL.Query()[field.Name]; !ok {
				multiErr.Add(fmt.Errorf("%s: required field missing", field.Name))
				continue fieldsLoop
			}
		}

		var value string

		switch field.Type {
		case Query:
			value = r.URL.Query().Get(field.Name)
		case FormValue:
			value = r.FormValue(field.Name)
		default:
			multiErr.Add(fmt.Errorf("%s: ParamType %q not valid", field.Name, field.Type))
			continue fieldsLoop
		}

		err := field.Value.Set(value)
		if err != nil {
			multiErr.Add(fmt.Errorf("%s: %s", field.Name, err))
		}
	}

	if len(multiErr) > 0 {
		return &multiErr
	}

	return nil
}

type Field struct {
	Name     string
	Value    Value
	Type     ParamType
	Required bool
}

type FieldOpt func(*Field)

func FieldValue(v Value) FieldOpt {
	return func(f *Field) {
		f.Value = v
	}
}

func QueryParam(f *Field) {
	f.Type = Query
}

func FormValueParam(f *Field) {
	f.Type = FormValue
}

func Required(f *Field) {
	f.Required = true
}

func NotRequired(f *Field) {
	f.Required = false
}

type Value interface {
	Set(string) error
}

type IntField struct {
	Field *int
}

func (intField *IntField) Set(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	*(intField.Field) = val

	return nil
}

type StringField struct {
	Field *string
}

func (stringField *StringField) Set(s string) error {
	*(stringField.Field) = s

	return nil
}

type MultiError []error

func (m *MultiError) Add(err error) {
	*m = append(*m, err)
}

func (m *MultiError) Error() string {
	errString := ""
	for _, err := range *m {
		errString += err.Error() + "; "
	}

	return strings.TrimSuffix(errString, "; ")
}
