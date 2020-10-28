package mapper

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

func ConvertMapToModel(model interface{}, data interface{}) error {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr {
		return errors.New("model must be a pointer")
	}

	modelValue := reflect.ValueOf(model).Elem()
	return setValue(modelValue, data)
}

func setValue(value reflect.Value, data interface{}) error {
	switch value.Kind() {
	case reflect.Struct:
		return setStructValue(value, data)
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		return setBasicTypeValue(value, data)
	case reflect.Slice:
		return setSliceValue(value, data)
	case reflect.Ptr:
		return setPointerValue(value, data)
	case reflect.Interface:
		value.Set(reflect.ValueOf(data))
		return nil
	}

	return nil
}

func setPointerValue(ptr reflect.Value, data interface{}) error {
	if ptr.Kind() != reflect.Ptr {
		return errors.New("this is not a pointer")
	}
	if data == nil {
		return nil
	}
	newValue := reflect.New(ptr.Type().Elem())
	err := setValue(newValue.Elem(), data)
	if err != nil {
		return err
	}
	ptr.Set(newValue)
	return nil
}

func setStructValue(structValue reflect.Value, data interface{}) error {
	if structValue.Type().Name() == "Time" {
		if data == nil {
			return nil
		}
		var err error
		var tm time.Time
		t, ok := data.(string)
		if ok && len(t) > 0 {
			tm, err = time.Parse(time.RFC3339, t)
			if err == nil {
				structValue.Set(reflect.ValueOf(tm))
			}
		} else {
			return errors.New("this is not a datetime")
		}

		return err
	}

	d, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("data must be map[string]interface{}")
	}

	if structValue.Kind() != reflect.Struct {
		return errors.New("this is not a struct")
	}

	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		tag := structField.Tag.Get("json")
		v, ok := d[strings.Split(tag, ",")[0]]
		if !ok {
			continue
		}

		fieldValue := structValue.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		err := setValue(fieldValue, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func setBasicTypeValue(field reflect.Value, data interface{}) error {
	dataValue := reflect.ValueOf(data)
	dataKind := dataValue.Kind()

	switch field.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		if isInt(dataKind) {
			field.SetInt(dataValue.Int())
		} else if isFloat(dataKind) {
			field.SetInt(int64(dataValue.Float()))
		}
	case reflect.Float32, reflect.Float64:
		if isInt(dataKind) {
			field.SetInt(int64(dataValue.Float()))
		} else if isFloat(dataKind) {
			field.SetFloat(dataValue.Float())
		}
	case reflect.String:
		if isString(dataKind) {
			field.SetString(dataValue.String())
		}
	case reflect.Bool:
		if dataKind == reflect.Bool {
			field.SetBool(dataValue.Bool())
		}
	}

	return nil
}

func setSliceValue(value reflect.Value, data interface{}) error {
	if value.Kind() != reflect.Slice {
		return errors.New("this is not a slice")
	}

	if data == nil {
		return nil
	}

	d, ok := data.([]interface{})
	if !ok {
		return errors.New("data must be slice of interface{}")
	}

	for _, v := range d {
		elem := reflect.New(value.Type().Elem()).Elem()
		err := setValue(elem, v)
		if err != nil {
			return err
		}
		value.Set(reflect.Append(value, elem))
	}

	return nil
}

func isInt(kind reflect.Kind) bool {
	return kind == reflect.Int || kind == reflect.Int32 || kind == reflect.Int64
}

func isFloat(kind reflect.Kind) bool {
	return kind == reflect.Float32 || kind == reflect.Float64
}

func isString(kind reflect.Kind) bool {
	return kind == reflect.String
}
