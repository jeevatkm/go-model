// Copyright (c) 2016 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	TagName    = "model"
	OmitField  = "-"
	NoTraverse = "notraverse"
)

var (
	Version            = "0.1-beta"
	NoTraverseTypeList map[reflect.Type]bool
)

func AddNoTraverseType(i ...interface{}) {
	for _, v := range i {
		t := reflect.TypeOf(v)
		if _, ok := NoTraverseTypeList[t]; ok {
			// already registered for no traverse, move on
			continue
		}

		// not found, add it
		NoTraverseTypeList[t] = true
	}
}

func RemoveNoTraverseType(i ...interface{}) {
	for _, v := range i {
		t := reflect.TypeOf(v)
		if _, ok := NoTraverseTypeList[t]; ok {

			// found, delete it
			delete(NoTraverseTypeList, t)
		}
	}
}

func IsZero(s interface{}) bool {
	if s == nil {
		return true
	}

	sv := indirect(valueOf(s))
	fields := Fields(sv)

	for _, f := range fields {
		fv := sv.FieldByName(f.Name)

		// embeded or nested struct
		if isStruct(fv) {

			if isNoTraverseType(fv) || isNoTraverse(f.Tag.Get(TagName)) {
				// not traversing inside, but evaluating a value
				if !IsFieldZero(fv) {
					return false
				}

				continue // move on to next attribute
			}

			if !IsZero(fv.Interface()) {
				return false
			}

			continue // move on to next attribute
		}

		if !IsFieldZero(fv) {
			return false
		}
	}

	return true
}

func IsFieldZero(f reflect.Value) bool {
	// zero value of the given field
	// For example: reflect.Zero(reflect.TypeOf(42)) returns a Value with Kind Int and value 0
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

	sv := valueOf(src)
	dv := valueOf(dst)
	if !isStruct(sv) || !isStruct(dv) {
		return append(errs, errors.New("Source or Destination is not a struct"))
	}

	if !isPtr(dv) {
		return append(errs, errors.New("Destination struct is not a pointer"))
	}

	if IsZero(src) {
		return append(errs, errors.New("Source struct is empty"))
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
		fmt.Println("Interface: ", sfv.Interface(), "Type:", sfv.Type(), "isStruct:", isStruct(sfv))

		// check whether field is zero or not
		var isVal bool
		if isStruct(sfv) {
			isVal = !IsZero(sfv)
		} else {
			isVal = !IsFieldZero(sfv)
		}

		if isVal || zero {
			dfv := dv.FieldByName(f.Name)

			// check dst field is exists, if not valid move on
			if !dfv.IsValid() {
				errs = append(errs, fmt.Errorf("Field: '%v', dst is not valid", f.Name))
				continue
			}

			// check kind of src and dst, if doesn't match move on
			if sfv.Kind() != dfv.Kind() {
				errs = append(errs, fmt.Errorf("Field: '%v', src & dst kind doesn't match", f.Name))
				continue
			}

			fmt.Println("Pre-conditions is met")

			// check dst field settable or not
			if dfv.CanSet() {
				fmt.Println("Can set value:", f.Name)

				// if src is zero make dst also zero
				// since zero=true
				if !isVal {
					dfv.Set(zeroVal(dfv))

					continue // move on to next attribute
				}

				// handle embeded/nested struct
				if isStruct(sfv) {
					fmt.Println("This is struct kind:", dTypeOf(sfv))
					if isNoTraverseType(sfv) || isNoTraverse(f.Tag.Get(TagName)) {
						fmt.Println("We are not going to traverse")
						// This is struct kind, but we are not going to traverse
						// since its in NoTraverseTypeList
						// however we will take care of attribute value
						dfv.Set(val(sfv, zero, true))
					} else {
						fmt.Println("We are going to traverse inside")

						ndv := reflect.New(indirect(sfv).Type())
						innerErrs := doCopy(ndv, sfv, zero)
						if innerErrs != nil {
							errs = append(errs, innerErrs...)
						}

						// handle based on ptr/non-ptr value
						if isPtr(sfv) {
							dfv.Set(ndv)
						} else {
							dfv.Set(indirect(ndv))
						}
					}

					continue // move on to next attribute
				}

				dfv.Set(val(sfv, zero, false))
			} else {
				errs = append(errs, fmt.Errorf("Field: '%v', can't set value in the dst", f.Name))
			}
		}
	}

	return errs
}

//
// go-model init
//

func init() {
	NoTraverseTypeList = map[reflect.Type]bool{}

	// Default NoTraverseTypeList
	// --------------------------
	// Auto No Traverse struct list for not traversing DeepLevel
	// However, attribute value will be evaluated by go-model
	AddNoTraverseType(
		time.Time{},
		&time.Time{},
		os.File{},
		&os.File{},
		// it's better to add it to the list for appropriate type(s)
	)
}

//
// Non-exported methods of model library
//

func isNoTraverseType(v reflect.Value) bool {
	t := dTypeOf(v)

	if _, ok := NoTraverseTypeList[t]; ok {
		return true
	}

	return false
}

func val(f reflect.Value, zero, notraverse bool) reflect.Value {
	var (
		ptr bool
		nf  reflect.Value
	)

	// take care interface{} and its actual value
	if isInterface(f) {
		f = valueOf(f.Interface())
	}

	// ptr, let's take a note
	if isPtr(f) {
		ptr = true
		f = f.Elem()
	}

	fmt.Println("==> Pointer:", ptr)
	fmt.Println("==> Type:", f.Type(), "Value:", f.Interface())

	switch f.Kind() {
	case reflect.Struct:
		fmt.Println("VAL==>VAL: We got the struct here")
		if notraverse {
			fmt.Println("==> Ended up notraverse:", notraverse)
			nf = f
		} else {
			nf = reflect.New(f.Type())
			// TODO propagate errors
			doCopy(nf, f, zero)
			fmt.Printf("\n==> f type: %v, value: %#v\n", f.Type(), f)

			// unwrap
			nf = nf.Elem()
			fmt.Println("==> Type nf:", nf.Type())
		}
	case reflect.Map:
		fmt.Println("Map value count:", f.Len())
		if f.Len() > 0 {
			nf = reflect.MakeMap(f.Type())

			for _, key := range f.MapKeys() {
				ov := f.MapIndex(key)
				fmt.Println("||===> Map Type:", ov.Type(), isStruct(ov), dTypeOf(ov))
				cv := reflect.New(ov.Type()).Elem()
				traverse := isNoTraverseType(ov) // TODO No traverse tag needs to handled
				cv.Set(val(ov, zero, traverse))
				nf.SetMapIndex(key, cv)
			}
		}
	case reflect.Slice:
		fmt.Println("Slice value count:", f.Len())
		if f.Len() > 0 {
			nf = reflect.MakeSlice(f.Type(), f.Len(), f.Cap())

			for i := 0; i < f.Len(); i++ {
				ov := f.Index(i)
				fmt.Println("||===> Slice Type:", ov.Type(), isStruct(ov), dTypeOf(ov))
				cv := reflect.New(ov.Type()).Elem()
				traverse := isNoTraverseType(ov) // TODO No traverse tag needs to handled
				cv.Set(val(ov, zero, traverse))
				nf.Index(i).Set(cv)
			}
		}
	default:
		nf = f
	}

	if ptr {
		// wrap
		fmt.Println("==> PTR ==> Type:", nf.Type(), "f:", f.Type())

		o := reflect.New(nf.Type())
		o.Elem().Set(nf)

		fmt.Println("==> o type:", o.Type())
		return o
	}

	return nf
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

func isNoTraverse(tag string) bool {
	return strings.Contains(tag, NoTraverse)
}

func dTypeOf(v reflect.Value) reflect.Type {
	if isInterface(v) {
		v = valueOf(v.Interface())
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

	return indirect(v).Kind() == reflect.Struct
}

func isInterface(v reflect.Value) bool {
	return v.Kind() == reflect.Interface
}
