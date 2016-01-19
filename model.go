// Copyright (c) 2016 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	TagName    = "model"
	OmitField  = "-"
	OmitNested = "omitnested"
)

var (
	Version            = "0.1"
	omitNestedTypeList map[reflect.Type]bool
)

func init() {
	omitNestedTypeList = map[reflect.Type]bool{}

	// Auto omit nested struct list for not traversing DeepLevel
	// However, variable value will be evaluated
	AddToOmitNested(
		time.Time{},
		&time.Time{},
		os.File{},
		&os.File{},
		http.Request{},
		&http.Request{},
		http.Response{},
		&http.Response{},
	)
}

func AddToOmitNested(i ...interface{}) {
	for _, v := range i {
		t := typeOf(v)
		if _, ok := omitNestedTypeList[t]; ok {
			continue
		}

		omitNestedTypeList[t] = true
	}
}

func IsZero(s interface{}) bool {
	if s == nil {
		return true
	}

	sv := valueOf(s)
	fields := Fields(sv)
	sv = indirect(sv)

	for _, f := range fields {
		fv := sv.FieldByName(f.Name)

		// embeded or nested struct
		if isStruct(fv) {
			if isOmitNestedType(fv) || isOmitNested(f.Tag.Get(TagName)) {
				if !IsFieldZero(fv) {
					return false
				}

				continue
			}

			if !IsZero(fv.Interface()) {
				return false
			}

			continue
		}

		if !IsFieldZero(fv) {
			return false
		}
	}

	return true
}

func IsFieldZero(f reflect.Value) bool {
	// zero value of the given field
	// For example: Zero(TypeOf(42)) returns a Value with Kind Int and value 0
	zero := reflect.Zero(f.Type()).Interface()

	return reflect.DeepEqual(f.Interface(), zero)
}

func Fields(v reflect.Value) []reflect.StructField {
	v = indirect(v)
	t := v.Type()

	var fs []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// Only exported fields of a struct can be accessed,
		// non-exported fields will be ignored
		if f.PkgPath == "" {
			// `model="-"` attributes will be omited
			if tag := f.Tag.Get(TagName); tag == OmitField {
				continue
			}

			fs = append(fs, f)
		}
	}

	return fs
}

func Copy(dst, src interface{}, zero bool) []error {
	var errs []error

	if IsZero(src) {
		return append(errs, errors.New("Source struct is empty"))
	}

	dv := valueOf(dst)
	if !isPtr(dv) {
		return append(errs, errors.New("Destination struct is not a pointer"))
	}

	sv := valueOf(src)
	if !isStruct(sv) || !isStruct(dv) {
		return append(errs, errors.New("Source or Destination is not a struct"))
	}

	// processing copy value(s)
	errs = doCopy(dv, sv, zero)
	if errs != nil {
		return errs
	}

	return nil
}

func doCopy(dv, sv reflect.Value, zero bool) []error {
	dv = indirect(dv)
	sv = indirect(sv)
	fields := Fields(sv)

	fmt.Println("No of src fields ready for use:", len(fields))
	fmt.Println("Copy only non-zero:", zero)

	var errs []error

	for _, f := range fields {
		sfv := sv.FieldByName(f.Name)

		fmt.Println("--------------------------")
		fmt.Println("Processing field:", f.Name)

		isVal := !IsFieldZero(sfv)
		if isVal || zero {
			dfv := dv.FieldByName(f.Name)

			// check dst field is exists, if not valid move on
			if !dfv.IsValid() {
				errs = append(errs, fmt.Errorf("Field: %v, dst is not valid", f.Name))
				continue
			}

			// check kind of src and dst, if doesn't match move on
			if sfv.Kind() != dfv.Kind() {
				errs = append(errs, fmt.Errorf("Field: %v, src & dst is different kind", f.Name))
				continue
			}

			fmt.Println("Pre-conditions is met")

			// check dst field settable or not
			if dfv.CanSet() {
				fmt.Println("Can set value:", f.Name)

				// if src is zero make dst also zero
				if !isVal {
					dfv.Set(zeroVal(dfv))

					continue // move on to next attribute
				}

				if isStruct(sfv) {
					fmt.Println("This is struct kind:", typeOf(sfv))
					if isOmitNestedType(sfv) || isOmitNested(f.Tag.Get(TagName)) {
						fmt.Println("We are not going to traverse")
						// This is struct kind, but we are not going to traverse
						// since its in omitNestedTypeList
						// however we will take care of current value
						dfv.Set(val(sfv))
					} else {
						fmt.Println("We are going to traverse inside")

						ndv := reflect.New(indirect(sfv).Type())
						innerErrs := doCopy(ndv, sfv, zero)
						if innerErrs != nil {
							errs = append(errs, innerErrs...)
						}

						// handle based value ptr/non-ptr
						if isPtr(sfv) {
							dfv.Set(ndv)
						} else {
							dfv.Set(indirect(ndv))
						}
					}

					continue // move on to next attribute
				}

				dfv.Set(val(sfv))
			}
		}

	}

	return nil
}

// Non-exported methods of model library

func isOmitNestedType(v reflect.Value) bool {
	t := indirect(v).Type()

	if _, ok := omitNestedTypeList[t]; ok {
		return true
	}

	return false
}

func val(f reflect.Value) reflect.Value {
	if isPtr(f) {
		fmt.Println("Value is Ptr:", f.Interface(), f.Elem().Interface())

		fe := f.Elem()
		nf := reflect.New(fe.Type())
		nf.Elem().Set(fe)
		return nf
	}

	fmt.Println("Value is not a Ptr:", f.Interface())
	// regular attribute may hold pointer
	// for eg.: map, slice
	switch f.Kind() {
	case reflect.Map:
		fmt.Println("This is map kind: ", len(f.MapKeys()))
		if len(f.MapKeys()) > 0 {
			fv := f.MapIndex(f.MapKeys()[0])
			fmt.Println("map first value is:", fv.Interface())
			if isPtr(fv) {
				fmt.Println("Found pointer element, start processing")

				nf := reflect.MakeMap(f.Type())
				for _, key := range f.MapKeys() {
					ov := f.MapIndex(key)
					cv := reflect.New(ov.Type()).Elem()
					cv.Set(val(ov))
					nf.SetMapIndex(key, cv)
				}

				return nf
			}
		}
	case reflect.Slice:
		if f.Len() > 0 {
			// handle if slice element is pointer else return as-is
			if isPtr(f.Index(0)) {
				nf := reflect.MakeSlice(f.Type(), f.Len(), f.Cap())

				for i := 0; i < f.Len(); i++ {
					nf.Index(i).Set(val(f.Index(i)))
				}

				return indirect(nf)
			}
		}
	}

	return f // return as-is
}

func zeroVal(f reflect.Value) reflect.Value {
	// get zero value for type
	ftz := reflect.Zero(f.Type())

	if f.Kind() == reflect.Ptr {
		return ftz
	}

	// if not a pointer then get zero value for interface
	return indirect(valueOf(ftz.Interface()))
}

func isOmitNested(tag string) bool {
	return strings.Contains(tag, OmitNested)
}

func typeOf(i interface{}) reflect.Type {
	return reflect.TypeOf(i) // indirect(valueOf(i)).Type()
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
	return indirect(v).Kind() == reflect.Struct
}
