package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type sliceValidateError []error

func (err sliceValidateError) Error() string {
	var errMsgs []string
	for i, e := range err {
		if e == nil {
			continue
		}
		errMsgs = append(errMsgs, fmt.Sprintf("[%d]: %s", i, e.Error()))
	}
	return strings.Join(errMsgs, "\n")
}

type Validator struct {
	once sync.Once
	*validator.Validate
}

func New() *Validator {
	v := &Validator{
		Validate: validator.New(),
	}

	return v.Engine().(*Validator)
}

func (v *Validator) lazyinit() {
	v.once.Do(func() {
		v.SetTagName("validate")
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})
		v.RegisterValidation("str_gt", v.validateStringGraterThan)
		v.RegisterValidation("str_lt", v.validateStringLessThan)
		v.RegisterValidation("has_lowercase", v.validateStringHasLowercase)
		v.RegisterValidation("has_uppercase", v.validateStringHasUppercase)
		v.RegisterValidation("has_special", v.validateStringHasSpecial)
	})
}

func (v *Validator) Engine() interface{} {
	v.lazyinit()
	return v
}

func (v *Validator) kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func (v *Validator) validateStringGraterThan(fl validator.FieldLevel) bool {
	l, _ := strconv.Atoi(fl.Param())
	return len(strings.TrimSpace(fl.Field().String())) > l
}

func (v *Validator) validateStringLessThan(fl validator.FieldLevel) bool {
	l, _ := strconv.Atoi(fl.Param())
	return len(strings.TrimSpace(fl.Field().String())) < l
}

func (v *Validator) validateStringHasLowercase(fl validator.FieldLevel) bool {
	for _, s := range fl.Field().String() {
		if unicode.IsLower(s) {
			return true
		}
	}
	return false
}

func (v *Validator) validateStringHasUppercase(fl validator.FieldLevel) bool {
	for _, s := range fl.Field().String() {
		if unicode.IsUpper(s) {
			return true
		}
	}
	return false
}

func (v *Validator) validateStringHasSpecial(fl validator.FieldLevel) bool {
	for _, s := range fl.Field().String() {
		if unicode.IsPunct(s) || unicode.IsSymbol(s) {
			return true
		}
	}
	return false
}

func CheckValidationErrors(err error) (e []error) {

	if _, ok := err.(*validator.InvalidValidationError); ok {
		e = append(e, err)
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		e = append(e, err)
		return
	}

	for _, validationError := range errs {
		message := errorMessages[validationError.Tag()]

		switch strings.Count(message, "%s") {
		case 0:
			e = append(e, errors.New(fmt.Sprintf(message)))
		case 1:
			e = append(e, errors.New(fmt.Sprintf(message, validationError.Field())))
		case 2:
			e = append(e, errors.New(fmt.Sprintf(message, validationError.Field(), validationError.Param())))
		}
	}

	return
}

func (v *Validator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.Struct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(sliceValidateError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}
