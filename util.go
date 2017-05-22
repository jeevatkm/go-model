// Copyright (c) Jeevanandam M. (https://github.com/jeevatkm).
// go-model source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"errors"
	"fmt"
	"reflect"
)

var errFieldNotExists = errors.New("Field does not exists")

func isFieldZero(f reflect.Value) bool {
	// zero value of the given field
	// For example: reflect.Zero(reflect.TypeOf(42)) returns a Value with Kind Int and value 0
	zero := reflect.Zero(f.Type()).Interface()

	return reflect.DeepEqual(f.Interface(), zero)
}

func isNoTraverseType(v reflect.Value) bool {
	if !isStruct(v) {
		return false
	}

	t := deepTypeOf(v)

	_, found := noTraverseTypeList[t]
	return found
}

func validateCopyField(f reflect.StructField, sfv, dfv reflect.Value) error {
	// check dst field is exists, if not valid move on
	if !dfv.IsValid() {
		return errFieldNotExists
		//return fmt.Errorf("Field does not exists in dst", f.Name)
	}

	if conversionExists(sfv.Type(), dfv.Type()) {
		return nil
	}

	// check kind of src and dst, if doesn't match move on
	if (sfv.Kind() != dfv.Kind()) && !isInterface(dfv) {
		return fmt.Errorf("Field: '%v', src [%v] & dst [%v] kind didn't match",
			f.Name,
			sfv.Kind(),
			dfv.Kind(),
		)
	}

	// check type of src and dst, if doesn't match move on
	sfvt := deepTypeOf(sfv)
	dfvt := deepTypeOf(dfv)

	if (sfvt.Kind() == reflect.Slice || sfvt.Kind() == reflect.Map) && sfvt.Kind() == dfvt.Kind() && conversionExists(sfvt.Elem(), dfvt.Elem()) {
		return nil
	}

	if (sfvt != dfvt) && !isInterface(dfv) {
		return fmt.Errorf("Field: '%v', src [%v] & dst [%v] type didn't match",
			f.Name,
			sfvt,
			dfvt,
		)
	}

	return nil
}

func modelFields(v reflect.Value) []reflect.StructField {
	v = indirect(v)
	t := v.Type()

	var fs []reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// Only exported fields of a struct can be accessed.
		// So, non-exported fields will be ignored
		if f.PkgPath == "" {
			fs = append(fs, f)
		}
	}

	return fs
}

func structValue(s interface{}) (reflect.Value, error) {
	if s == nil {
		return reflect.Value{}, errors.New("Invalid input <nil>")
	}

	sv := indirect(valueOf(s))

	if !isStruct(sv) {
		return reflect.Value{}, errors.New("Input is not a struct")
	}

	return sv, nil
}

func getField(sv reflect.Value, name string) (reflect.Value, error) {
	field := sv.FieldByName(name)
	if !field.IsValid() {
		return reflect.Value{}, fmt.Errorf("Field: '%v', does not exists", name)
	}

	return field, nil
}

func zeroOf(f reflect.Value) reflect.Value {

	// get zero value for type
	ftz := reflect.Zero(f.Type())

	if f.Kind() == reflect.Ptr {
		return ftz
	}

	// if not a pointer then get zero value for interface
	return indirect(valueOf(ftz.Interface()))
}

func deepTypeOf(v reflect.Value) reflect.Type {
	if isInterface(v) {

		// check zero or not
		if !isFieldZero(v) {
			v = valueOf(v.Interface())
		}

	}

	return v.Type()
}

func valueOf(i interface{}) reflect.Value {
	return reflect.ValueOf(i)
}

func indirect(v reflect.Value) reflect.Value {
	return reflect.Indirect(v)
}

func isPtr(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr
}

func isStruct(v reflect.Value) bool {
	if isInterface(v) {
		v = valueOf(v.Interface())
	}

	pv := indirect(v)

	// struct is not yet initialized
	if pv.Kind() == reflect.Invalid {
		return false
	}

	return pv.Kind() == reflect.Struct
}

func isInterface(v reflect.Value) bool {
	return v.Kind() == reflect.Interface
}

func extractType(x interface{}) reflect.Type {
	return reflect.TypeOf(x).Elem()
}

func conversionExists(srcType reflect.Type, destType reflect.Type) bool {
	if _, ok := converterMap[srcType]; !ok {
		return false
	}
	if _, ok := converterMap[srcType][destType]; !ok {
		return false
	}
	return true
}
