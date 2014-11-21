package form

import (
	"errors"
	"net/http"
	"strconv"
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
	for _, field := range f.Fields {
		if field.Required {
			if _, ok := r.URL.Query()[field.Name]; !ok {
				return errors.New("Value was not available")
			}
		}

		var value string

		switch field.Type {
		case Query:
			value = r.URL.Query().Get(field.Name)
		case FormValue:
			// TODO(partkyle): figure out why the for is not posting correctly
			err := r.ParseForm()
			if err != nil {
				return err
			}
			value = r.PostForm.Get(field.Name)
		default:
			return errors.New("Unsupported ParamType")
		}

		err := field.Value.Set(value)
		if err != nil {
			return err
		}
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
