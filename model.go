// Copyright (c) 2016 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

/*
Package model provides robust & easy to use model mapper and model utility methods for Go.
These typical methods increases productivity and makes Go developement more fun :)
*/
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
	// TagName is used to mention field options for go-model library.
	//
	// For Example:
	// ------------
	// BookCount	int		`model:"bookCount"`
	// ArchiveInfo	StoreInfo	`model:"archiveInfo,notraverse"`
	TagName = "model"

	// OmitField value is used omit field(s) from go-model processing
	OmitField = "-"

	// OmitEmpty option is used skip field(s) from output if zero value
	OmitEmpty = "omitempty"

	// NoTraverse means go-model library will not traverse inside those struct object.
	// However, attribute value will be evaluated/processed by library.
	NoTraverse = "notraverse"
)

var (
	// Version # of go-model library
	Version = "0.1-beta"

	// NoTraverseTypeList keeps track of no-traverse type list at library level
	NoTraverseTypeList map[reflect.Type]bool

	typeOfBytes = reflect.TypeOf([]byte(nil))
)

// AddNoTraverseType method adds the Go Lang type into global `NoTraverseTypeList`.
// Those type(s) from list is considered as "No Traverse" type by go-model library
// for model mapping process. See also `RemoveNoTraverseType()` method.
// 		model.AddNoTraverseType(time.Time{}, &time.Time{}, os.File{}, &os.File{})
//
// Default NoTraverseTypeList: time.Time{}, &time.Time{}, os.File{}, &os.File{},
// http.Request{}, &http.Request{}, http.Response{}, &http.Response{}
//
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

// RemoveNoTraverseType method is used to remove Go Lang type from the `NoTraverseTypeList`.
// See also `AddNoTraverseType()` method.
// 		model.RemoveNoTraverseType(http.Request{}, &http.Request{})
//
func RemoveNoTraverseType(i ...interface{}) {
	for _, v := range i {
		t := reflect.TypeOf(v)
		if _, ok := NoTraverseTypeList[t]; ok {

			// found, delete it
			delete(NoTraverseTypeList, t)
		}
	}
}

// IsZero method returns true if all the exported fields in a given struct
// is a zero value. If its not struct then method returns false.
//
// A "model" tag with the value of "-" is ignored by library for processing.
// 		For Example:
//
// 		// Field is ignored by go-model processing
// 		BookCount	int	`model:"-"`
// 		BookCode	string	`model:"-"`
//
// A "model" tag value with the option of "notraverse"; library will not traverse
// inside those struct object. However field value will be evaluated whether
// its zero or not.
// 		For Example:
//
// 		// Field is not traversed but value is evaluated/processed
// 		ArchiveInfo	BookArchive	`model:"archiveInfo,notraverse"`
// 		Region		BookLocale	`model:",notraverse"`
//
func IsZero(s interface{}) bool {
	if s == nil {
		return true
	}

	sv := indirect(valueOf(s))

	if !isStruct(sv) {
		return false
	}

	fields := getFields(sv)

	for _, f := range fields {
		fv := sv.FieldByName(f.Name)

		// embeded or nested struct
		if isStruct(fv) {

			if isNoTraverseType(fv) || isNoTraverse(f.Tag.Get(TagName)) {

				// not traversing inside, but evaluating a value
				if !isFieldZero(fv) {
					return false
				}

				continue
			}

			if !IsZero(fv.Interface()) {
				return false
			}

			continue
		}

		if !isFieldZero(fv) {
			return false
		}
	}

	return true
}

// Copy method copies all the exported values from source struct into destination struct.
// The "Name", "Type" and "Kind" is should match to qualify a copy. One exception though;
// if the destination Field type is "interface{}" then "Type" and "Kind" doesn't matter,
// source value gets copied to that destination Field.
//
// Copy method handles:
// 		// List out what copy method can do
//
// Sample Snippet
// 		src := SampleStruct{ /* source values goes here */ }
// 		dst := SampleStruct{}
//
// 		errs := model.Copy(&dst, src)
// 		if errs != nil {
// 			fmt.Println("Errors:", errs)
// 		}
//
// Note:
// [1] Copy process continues regardless of the case it qualify or not. Not qualified field(s)
// gets added to '[]error' that you will get at the end.
// [2] Slice3 type is not yet supported by Copy method.
//
// A "model" tag with the value of "-" is ignored by library for processing.
// 		For Example:
//
// 		// Field is ignored by go-model processing
// 		BookCount	int	`model:"-"`
// 		BookCode	string	`model:"-"`
//
// A "model" tag value with the option of "omitempty"; library will not copy those values
// into destination struct object. It may be very handy for partial put or patch update
// request scenarios; you don't want to copy empty/zero value into destination object
// 		For Example:
//
// 		// Field is not copy into 'dst' if it is empty/zero value
// 		ArchiveInfo	BookArchive	`model:"archiveInfo,omitempty"`
// 		Region		BookLocale	`model:",omitempty,notraverse"`
//
// A "model" tag value with the option of "notraverse"; library will not traverse
// inside those struct object. However field value will be evaluated whether
// its zero or not.
// 		For Example:
//
// 		// Field is not traversed but value is evaluated/processed
// 		ArchiveInfo	BookArchive	`model:"archiveInfo,notraverse"`
// 		Region		BookLocale	`model:",notraverse"`
//
func Copy(dst, src interface{}) []error {
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
	errs = doCopy(dv, sv)
	if errs != nil {
		return errs
	}

	return nil
}

func Map(s interface{}) (map[string]interface{}, error) {
	if s == nil {
		return nil, errors.New("Invalid input <nil>")
	}

	sv := indirect(valueOf(s))

	if !isStruct(sv) {
		return nil, errors.New("Input is not a struct")
	}

	// processing map value(s)
	return doMap(sv), nil
}

//
// go-model init
//

func init() {
	NoTraverseTypeList = map[reflect.Type]bool{}

	// Default NoTraverseTypeList
	// --------------------------
	// Auto No Traverse struct list for not traversing Deep Level
	// However, attribute value will be evaluated by go-model library
	AddNoTraverseType(
		time.Time{},
		&time.Time{},
		os.File{},
		&os.File{},
		http.Request{},
		&http.Request{},
		http.Response{},
		&http.Response{},

		// it's better to add it to the list for appropriate type(s)
	)
}

//
// Non-exported methods of model library
//

func getFields(v reflect.Value) []reflect.StructField {
	v = indirect(v)
	t := v.Type()

	var fs []reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// Only exported fields of a struct can be accessed,
		// non-exported fields will be ignored
		if f.PkgPath == "" {

			// `model="-"` attributes will be omitted
			if tag := f.Tag.Get(TagName); tag == OmitField {
				continue
			}

			fs = append(fs, f)
		}
	}

	return fs
}

func isFieldZero(f reflect.Value) bool {
	// zero value of the given field
	// For example: reflect.Zero(reflect.TypeOf(42)) returns a Value with Kind Int and value 0
	zero := reflect.Zero(f.Type()).Interface()

	return reflect.DeepEqual(f.Interface(), zero)
}

func isNoTraverseType(v reflect.Value) bool {
	t := dTypeOf(v)

	if _, ok := NoTraverseTypeList[t]; ok {
		return true
	}

	return false
}

func doCopy(dv, sv reflect.Value) []error {
	dv = indirect(dv)
	sv = indirect(sv)
	fields := getFields(sv)

	var errs []error

	for _, f := range fields {
		sfv := sv.FieldByName(f.Name)
		tag := newTag(f.Tag.Get(TagName))

		// compute no-traverse scope
		noTraverse := (isNoTraverseType(sfv) || tag.isNoTraverse())

		// check whether field is zero or not
		var isVal bool
		if isStruct(sfv) && !noTraverse {
			isVal = !IsZero(sfv.Interface())
		} else {
			isVal = !isFieldZero(sfv)
		}

		dfv := dv.FieldByName(f.Name)

		if isVal {

			// check dst field is exists, if not valid move on
			if !dfv.IsValid() {
				errs = append(errs, fmt.Errorf("Field: '%v', dst is not valid", f.Name))
				continue
			}

			// check kind of src and dst, if doesn't match move on
			if (sfv.Kind() != dfv.Kind()) && !isInterface(dfv) {
				errs = append(errs, fmt.Errorf("Field: '%v', src [%v] & dst [%v] kind doesn't match",
					f.Name,
					sfv.Kind(),
					dfv.Kind(),
				))
				continue
			}

			// check type of src and dst, if doesn't match move on
			sfvt := dTypeOf(sfv)
			dfvt := dTypeOf(dfv)
			if (sfvt != dfvt) && !isInterface(dfv) {
				errs = append(errs, fmt.Errorf("Field: '%v', src [%v] & dst [%v] type doesn't match",
					f.Name,
					sfvt,
					dfvt,
				))
				continue
			}

			// check dst field settable or not
			if dfv.CanSet() {

				// handle embeded or nested struct
				if isStruct(sfv) {

					if noTraverse {
						// This is struct kind and its present in NoTraverseTypeList,
						// so go-model is not gonna traverse inside.
						// however will take care of field value
						dfv.Set(val(sfv, true))
					} else {
						ndv := reflect.New(indirect(sfv).Type())
						innerErrs := doCopy(ndv, sfv)

						// add errors to main stream
						errs = append(errs, innerErrs...)

						// handle based on ptr/non-ptr value
						if isPtr(sfv) {
							dfv.Set(ndv)
						} else {
							dfv.Set(indirect(ndv))
						}
					}

					continue
				}

				dfv.Set(val(sfv, false))
			}
		} else {

			// field value is zero and 'omitempty' option present
			// then not copying into destination struct
			if !tag.isOmitEmpty() {
				dfv.Set(zeroVal(dfv))
			}
		}
	}

	return errs
}

func doMap(sv reflect.Value) map[string]interface{} {
	sv = indirect(sv)
	fields := getFields(sv)
	m := map[string]interface{}{}

	for _, f := range fields {
		fv := sv.FieldByName(f.Name)
		tag := newTag(f.Tag.Get(TagName))

		// map key name
		keyName := f.Name
		if !isStringEmpty(tag.Name) {
			keyName = tag.Name
		}

		// compute no-traverse scope
		noTraverse := (isNoTraverseType(fv) || tag.isNoTraverse())

		// check whether field is zero or not
		var isVal bool
		if isStruct(fv) && !noTraverse {
			isVal = !IsZero(fv.Interface())
		} else {
			isVal = !isFieldZero(fv)
		}

		// field value is not zero
		if isVal {

			// handle embeded or nested struct
			if isStruct(fv) {

				if noTraverse {
					// This is struct kind and its present in NoTraverseTypeList,
					// so go-model is not gonna traverse inside.
					// however will take care of field value
					m[keyName] = mapVal(fv, true).Interface()
				} else {

					// embeded struct values gets mapped at
					// Field level instead of object
					fmv := doMap(fv)
					if f.Anonymous {
						for k, v := range fmv {
							m[k] = v
						}
					} else {
						m[keyName] = fmv
					}
				}

				continue
			}

			m[keyName] = mapVal(fv, false).Interface()
		} else {

			// field value is zero and 'omitempty' option present
			// then not including into Map
			if !tag.isOmitEmpty() {
				m[keyName] = zeroVal(fv).Interface()
			}
		}
	}

	return m
}

func val(f reflect.Value, notraverse bool) reflect.Value {
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

	// reflect.Slice3 is not yet supported by this library
	switch f.Kind() {
	case reflect.Struct:

		if notraverse {
			nf = f
		} else {
			nf = reflect.New(f.Type())

			// currently, struct within map/slice errors doesn't get propagated
			doCopy(nf, f)

			// unwrap
			nf = nf.Elem()
		}
	case reflect.Map:
		nf = reflect.MakeMap(f.Type())

		for _, key := range f.MapKeys() {
			ov := f.MapIndex(key)

			cv := reflect.New(ov.Type()).Elem()
			traverse := isNoTraverseType(ov)
			cv.Set(val(ov, traverse))

			nf.SetMapIndex(key, cv)
		}
	case reflect.Slice:
		if f.Type() == typeOfBytes {
			nf = f
		} else {
			nf = reflect.MakeSlice(f.Type(), f.Len(), f.Cap())

			for i := 0; i < f.Len(); i++ {
				ov := f.Index(i)

				cv := reflect.New(ov.Type()).Elem()
				traverse := isNoTraverseType(ov)
				cv.Set(val(ov, traverse))

				nf.Index(i).Set(cv)
			}
		}
	default:
		nf = f
	}

	if ptr {
		// wrap
		o := reflect.New(nf.Type())
		o.Elem().Set(nf)

		return o
	}

	return nf
}

func mapVal(f reflect.Value, notraverse bool) reflect.Value {
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

	// reflect.Slice3 is not yet supported by this library
	switch f.Kind() {
	case reflect.Struct:
		if notraverse {
			nf = f
		} else {
			nf = valueOf(doMap(f))
		}
	case reflect.Map:
		nmv := map[string]interface{}{}

		for _, key := range f.MapKeys() {
			skey := fmt.Sprintf("%v", key.Interface())
			mv := f.MapIndex(key)

			var nv reflect.Value
			if isStruct(mv) {
				nv = mapVal(mv, isNoTraverseType(mv))
			} else {
				nv = mapVal(mv, false)
			}

			nmv[skey] = nv.Interface()
		}

		nf = valueOf(nmv)
	case reflect.Slice:
		if f.Type() == typeOfBytes {
			nf = f
		} else {
			nf = reflect.MakeSlice(f.Type(), f.Len(), f.Cap())
			for i := 0; i < f.Len(); i++ {
				sv := f.Index(i)
				dv := reflect.New(sv.Type()).Elem()

				if isStruct(sv) {
					// TODO
				} else {
					dv.Set(mapVal(sv, false))
				}

				nf.Index(i).Set(dv)
			}
		}
	default:
		nf = f
	}

	if ptr {
		// wrap
		o := reflect.New(nf.Type())
		o.Elem().Set(nf)

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

	v = indirect(v)

	// struct is not yet initialized
	if v.Kind() == reflect.Invalid {
		return false
	}

	return v.Kind() == reflect.Struct
}

func isInterface(v reflect.Value) bool {
	return v.Kind() == reflect.Interface
}
