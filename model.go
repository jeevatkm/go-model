// Copyright (c) Jeevanandam M. (https://github.com/jeevatkm).
// go-model source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

// Package model provides robust and easy-to-use model mapper and utility methods for Go.
// These typical methods increase productivity and make Go development more fun :)
package model

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"
)

// Converter is used to provide custom mappers for a datatype pair.
type Converter func(in reflect.Value) (reflect.Value, error)

const (
	// TagName is used to mention field options for go-model library.
	//
	// Example:
	// --------
	// BookCount	int		`model:"bookCount"`
	// ArchiveInfo	StoreInfo	`model:"archiveInfo,notraverse"`
	TagName = "model"

	// OmitField value is used to omit field(s) from processing
	OmitField = "-"

	// OmitEmpty option is used skip field(s) from output if it's zero value
	OmitEmpty = "omitempty"

	// NoTraverse option makes sure the go-model library to not to traverse inside the struct object.
	// However, the field value will be evaluated or processed by library.
	NoTraverse = "notraverse"
)

var (
	// Version # of go-model library
	Version = "1.1.0"

	// NoTraverseTypeList keeps track of no-traverse type list at library level
	noTraverseTypeList map[reflect.Type]bool

	// Type conversion functions at library level
	converterMap map[reflect.Type]map[reflect.Type]Converter

	typeOfBytes     = reflect.TypeOf([]byte(nil))
	typeOfInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

// AddNoTraverseType method adds the Go Lang type into global `NoTraverseTypeList`.
// The type(s) from list is considered as "No Traverse" type by go-model library
// for model mapping process. See also `RemoveNoTraverseType()` method.
// 		model.AddNoTraverseType(time.Time{}, &time.Time{}, os.File{}, &os.File{})
//
// Default NoTraverseTypeList: time.Time{}, &time.Time{}, os.File{}, &os.File{},
// http.Request{}, &http.Request{}, http.Response{}, &http.Response{}
//
func AddNoTraverseType(i ...interface{}) {
	for _, v := range i {
		t := reflect.TypeOf(v)
		if _, ok := noTraverseTypeList[t]; ok {

			// already registered for no traverse, move on
			continue
		}

		// not found, add it
		noTraverseTypeList[t] = true
	}
}

// RemoveNoTraverseType method is used to remove Go Lang type from the `NoTraverseTypeList`.
// See also `AddNoTraverseType()` method.
// 		model.RemoveNoTraverseType(http.Request{}, &http.Request{})
//
func RemoveNoTraverseType(i ...interface{}) {
	for _, v := range i {
		t := reflect.TypeOf(v)
		if _, ok := noTraverseTypeList[t]; ok {

			// found, delete it
			delete(noTraverseTypeList, t)
		}
	}
}

// AddConversion mothod allows registering a custom `Converter` into the global `converterMap`
// by supplying pointers of the target types.
func AddConversion(in interface{}, out interface{}, converter Converter) {
	srcType := extractType(in)
	targetType := extractType(out)
	AddConversionByType(srcType, targetType, converter)
}

// AddConversionByType allows registering a custom `Converter` into golbal `converterMap` by types.
func AddConversionByType(srcType reflect.Type, targetType reflect.Type, converter Converter) {
	if _, ok := converterMap[srcType]; !ok {
		converterMap[srcType] = map[reflect.Type]Converter{}
	}
	converterMap[srcType][targetType] = converter
}

// RemoveConversion registered conversions
func RemoveConversion(in interface{}, out interface{}) {
	srcType := extractType(in)
	targetType := extractType(out)
	if _, ok := converterMap[srcType]; !ok {
		return
	}
	if _, ok := converterMap[srcType][targetType]; !ok {
		return
	}
	delete(converterMap[srcType], targetType)
}

// IsZero method returns `true` if all the exported fields in a given `struct`
// are zero value otherwise `false`. If input is not a struct, method returns `false`.
//
// A "model" tag with the value of "-" is ignored by library for processing.
// 		Example:
//
// 		// Field is ignored by go-model processing
// 		BookCount	int	`model:"-"`
// 		BookCode	string	`model:"-"`
//
// A "model" tag value with the option of "notraverse"; library will not traverse
// inside the struct object. However, the field value will be evaluated whether
// it's zero value or not.
// 		Example:
//
// 		// Field is not traversed but value is evaluated/processed
// 		ArchiveInfo	BookArchive	`model:"archiveInfo,notraverse"`
// 		Region		BookLocale	`model:",notraverse"`
//
func IsZero(s interface{}) bool {
	if s == nil {
		return true
	}

	sv, err := structValue(s)
	if err != nil {
		return false
	}

	fields := modelFields(sv)

	for _, f := range fields {
		fv := sv.FieldByName(f.Name)
		tag := newTag(f.Tag.Get(TagName))

		if tag.isOmitField() {
			continue
		}

		// embedded or nested struct
		if isStruct(fv) {
			// check type is in NoTraverseTypeList or has 'notraverse' tag option
			if isNoTraverseType(fv) || tag.isNoTraverse() {

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

// IsZeroInFields method verifies the value for the given list of field names against
// given struct. Method returns `Field Name` and `true` for the zero value field.
// Otherwise method returns empty `string` and `false`.
//
// Note:
// [1] This method doesn't traverse nested and embedded `struct`, instead it just evaluates that `struct`.
// [2] If given field is not exists in the struct, method moves on to next field
//
// A "model" tag with the value of "-" is ignored by library for processing.
// 		Example:
//
// 		// Field is ignored by go-model processing
// 		BookCount	int	`model:"-"`
// 		BookCode	string	`model:"-"`
//
func IsZeroInFields(s interface{}, names ...string) (string, bool) {
	if s == nil || len(names) == 0 {
		return "", true
	}

	sv, err := structValue(s)
	if err != nil {
		return "", false
	}

	for _, name := range names {
		fv := sv.FieldByName(name)

		// if given field is not exists then continue
		if !fv.IsValid() {
			continue
		}

		if isFieldZero(fv) {
			return name, true
		}
	}

	return "", false
}

// HasZero method returns `true` if any one of the exported fields in a given
// `struct` is zero value otherwise `false`. If input is not a struct, method
// returns `false`.
//
// A "model" tag with the value of "-" is ignored by library for processing.
// 		Example:
//
// 		// Field is ignored by go-model processing
// 		BookCount	int	`model:"-"`
// 		BookCode	string	`model:"-"`
//
// A "model" tag value with the option of "notraverse"; library will not traverse
// inside the struct object. However, the field value will be evaluated whether
// it's zero value or not.
// 		Example:
//
// 		// Field is not traversed but value is evaluated/processed
// 		ArchiveInfo	BookArchive	`model:"archiveInfo,notraverse"`
// 		Region		BookLocale	`model:",notraverse"`
//
func HasZero(s interface{}) bool {
	if s == nil {
		return true
	}

	sv, err := structValue(s)
	if err != nil {
		return false
	}

	fields := modelFields(sv)

	for _, f := range fields {
		fv := sv.FieldByName(f.Name)
		tag := newTag(f.Tag.Get(TagName))

		if tag.isOmitField() {
			continue
		}

		// embedded or nested struct
		if isStruct(fv) {
			// check type is in NoTraverseTypeList or has 'notraverse' tag option
			if isNoTraverseType(fv) || tag.isNoTraverse() {

				// not traversing inside, but evaluating a value
				if isFieldZero(fv) {
					return true
				}

				continue
			}

			if HasZero(fv.Interface()) {
				return true
			}

			continue
		}

		if isFieldZero(fv) {
			return true
		}
	}

	return false
}

// Copy method copies all the exported field values from source `struct` into destination `struct`.
// The "Name", "Type" and "Kind" is should match to qualify a copy. One exception though;
// if the destination field type is "interface{}" then "Type" and "Kind" doesn't matter,
// source value gets copied to that destination field.
//
// 		Example:
//
// 		src := SampleStruct { /* source struct field values go here */ }
// 		dst := SampleStruct {}
//
// 		errs := model.Copy(&dst, src)
// 		if errs != nil {
// 			fmt.Println("Errors:", errs)
// 		}
//
// Note:
// [1] Copy process continues regardless of the case it qualifies or not. The non-qualified field(s)
// gets added to '[]error' that you will get at the end.
// [2] Two dimensional slice type is not supported yet.
//
// A "model" tag with the value of "-" is ignored by library for processing.
// 		Example:
//
// 		// Field is ignored while processing
// 		BookCount	int	`model:"-"`
// 		BookCode	string	`model:"-"`
//
// A "model" tag value with the option of "omitempty"; library will not copy those values
// into destination struct object. It may be handy for partial put or patch update
// request scenarios; if you don't want to copy empty/zero value into destination object.
// 		Example:
//
// 		// Field is not copy into 'dst' if it's empty/zero value
// 		ArchiveInfo	BookArchive	`model:"archiveInfo,omitempty"`
// 		Region		BookLocale	`model:",omitempty,notraverse"`
//
// A "model" tag value with the option of "notraverse"; library will not traverse
// inside the struct object. However, the field value will be evaluated whether
// it's zero value or not, and then copied to the destination object accordingly.
// 		Example:
//
// 		// Field is not traversed but value is evaluated/processed
// 		ArchiveInfo	BookArchive	`model:"archiveInfo,notraverse"`
// 		Region		BookLocale	`model:",notraverse"`
//
func Copy(dst, src interface{}) []error {
	var errs []error

	if src == nil || dst == nil {
		return append(errs, errors.New("Source or Destination is nil"))
	}

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

	// processing, copy field value(s)
	errs = doCopy(dv, sv)
	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Clone method creates a clone of given `struct` object. As you know go-model does, deep processing.
// So all field values you get in the result.
//
// 		Example:
// 		input := SampleStruct { /* input struct field values go here */ }
//
// 		clonedObj := model.Clone(input)
//
// 		fmt.Printf("\nCloned Object: %#v\n", clonedObj)
//
// Note:
// [1] Two dimensional slice type is not supported yet.
//
// A "model" tag with the value of "-" is ignored by library for processing.
// 		Example:
//
// 		// Field is ignored while processing
// 		BookCount	int	`model:"-"`
// 		BookCode	string	`model:"-"`
//
// A "model" tag value with the option of "omitempty"; library will not clone those values
// into result struct object.
// 		Example:
//
// 		// Field is not cloned into 'result' if it's empty/zero value
// 		ArchiveInfo	BookArchive	`model:"archiveInfo,omitempty"`
// 		Region		BookLocale	`model:",omitempty,notraverse"`
//
// A "model" tag value with the option of "notraverse"; library will not traverse
// inside the struct object. However, the field value will be evaluated whether
// it's zero value or not, and then cloned to the result accordingly.
// 		Example:
//
// 		// Field is not traversed but value is evaluated/processed
// 		ArchiveInfo	BookArchive	`model:"archiveInfo,notraverse"`
// 		Region		BookLocale	`model:",notraverse"`
//
func Clone(s interface{}) (interface{}, error) {
	sv, err := structValue(s)
	if err != nil {
		return nil, err
	}

	// figure out target type
	st := deepTypeOf(sv)

	// create a target type
	dv := reflect.New(st)

	// apply copy to target
	doCopy(dv, sv)

	return dv.Interface(), nil
}

// Map method converts all the exported field values from the given `struct`
// into `map[string]interface{}`. In which the keys of the map are the field names
// and the values of the map are the associated values of the field.
//
// 		Example:
//
// 		src := SampleStruct { /* source struct field values go here */ }
//
// 		err := model.Map(src)
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 		}
//
// Note:
// [1] Two dimensional slice type is not supported yet.
//
// The default 'Key Name' string is the struct field name. However, it can be
// changed in the struct field's tag value via "model" tag.
// 		Example:
//
// 		// Now field 'Key Name' is customized
// 		BookTitle	string	`model:"bookTitle"`
//
// A "model" tag with the value of "-" is ignored by library for processing.
// 		Example:
//
// 		// Field is ignored while processing
// 		BookCount	int	`model:"-"`
// 		BookCode	string	`model:"-"`
//
// A "model" tag value with the option of "omitempty"; library will not include those values
// while converting to map[string]interface{}. If it's empty/zero value.
// 		Example:
//
// 		// Field is not included in result map if it's empty/zero value
// 		ArchivedDate	time.Time	`model:"archivedDate,omitempty"`
// 		Region		BookLocale	`model:",omitempty,notraverse"`
//
// A "model" tag value with the option of "notraverse"; library will not traverse
// inside the struct object. However, the field value will be evaluated whether
// it's zero value or not, and then added to the result map accordingly.
// 		Example:
//
// 		// Field is not traversed but value is evaluated/processed
// 		ArchivedDate	time.Time	`model:"archivedDate,notraverse"`
// 		Region		BookLocale	`model:",notraverse"`
//
func Map(s interface{}) (map[string]interface{}, error) {
	sv, err := structValue(s)
	if err != nil {
		return nil, err
	}

	// processing, field value(s) into map
	return doMap(sv), nil
}

// Fields method returns the exported struct fields from the given `struct`.
// 		Example:
//
// 		src := SampleStruct { /* source struct field values go here */ }
//
// 		fields, _ := model.Fields(src)
// 		for _, f := range fields {
// 			tag := newTag(f.Tag.Get("model"))
// 			fmt.Println("Field Name:", f.Name, "Tag Name:", tag.Name, "Tag Options:", tag.Options)
// 		}
//
func Fields(s interface{}) ([]reflect.StructField, error) {
	sv, err := structValue(s)
	if err != nil {
		return nil, err
	}

	return modelFields(sv), nil
}

// Kind method returns `reflect.Kind` for the given field name from the `struct`.
// 		Example:
//
// 		src := SampleStruct {
// 			BookCount      int         `json:"-"`
// 			BookCode       string      `json:"-"`
// 			ArchiveInfo    BookArchive `json:"archive_info,omitempty"`
// 			Region         BookLocale  `json:"region,omitempty"`
// 		}
//
// 		fieldKind, _ := model.Kind(src, "ArchiveInfo")
// 		fmt.Println("Field kind:", fieldKind)
//
func Kind(s interface{}, name string) (reflect.Kind, error) {
	sv, err := structValue(s)
	if err != nil {
		return reflect.Invalid, err
	}

	fv, err := getField(sv, name)
	if err != nil {
		return reflect.Invalid, err
	}

	return fv.Type().Kind(), nil
}

// Get method returns a field value from `struct` by field name.
// 		Example:
//
// 		src := SampleStruct {
// 			BookCount      int         `json:"-"`
// 			BookCode       string      `json:"-"`
// 			ArchiveInfo    BookArchive `json:"archive_info,omitempty"`
// 			Region         BookLocale  `json:"region,omitempty"`
// 		}
//
// 		value, err := model.Get(src, "ArchiveInfo")
// 		fmt.Println("Field Value:", value)
// 		fmt.Println("Error:", err)
//
// Note: Get method does not honor model tag annotations. Get simply access
// value on exported fields.
//
func Get(s interface{}, name string) (interface{}, error) {
	sv, err := structValue(s)
	if err != nil {
		return nil, err
	}

	fv, err := getField(sv, name)
	if err != nil {
		return nil, err
	}

	return fv.Interface(), nil
}

// Set method sets a value into field on struct by field name.
// 		Example:
//
// 		src := SampleStruct {
// 			BookCount      int         `json:"-"`
// 			BookCode       string      `json:"-"`
// 			ArchiveInfo    BookArchive `json:"archive_info,omitempty"`
// 			Region         BookLocale  `json:"region,omitempty"`
// 		}
//
// 		bookLocale := BookLocale {
//			Locale: "en-US",
//			Language: "en",
//			Region: "US",
// 		}
//
// 		err := model.Set(&src, "Region", bookLocale)
// 		fmt.Println("Error:", err)
//
// Note: Set method does not honor model tag annotations. Set simply given
// value by field name on exported fields.
//
func Set(s interface{}, name string, value interface{}) error {
	if s == nil {
		return errors.New("Invalid input <nil>")
	}

	sv := valueOf(s)
	if isPtr(sv) {
		sv = sv.Elem()
	} else {
		return errors.New("Destination struct is not a pointer")
	}

	fv, err := getField(sv, name)
	if err != nil {
		return err
	}

	if !fv.CanSet() {
		return fmt.Errorf("Field: %v, cannot be settable", name)
	}

	tv := valueOf(value)
	if isPtr(tv) {
		tv = tv.Elem()
	}

	if (fv.Kind() != tv.Kind()) || fv.Type() != tv.Type() {
		return fmt.Errorf("Field: %v, type/kind did not match", name)
	}

	// assign the given value
	fv.Set(tv)

	return nil
}

//
// go-model init
//

func init() {
	noTraverseTypeList = map[reflect.Type]bool{}
	converterMap = map[reflect.Type]map[reflect.Type]Converter{}

	// Default NoTraverseTypeList
	// --------------------------
	// Auto No Traverse struct list for not traversing Deep Level
	// However, field value will be evaluated/processed by go-model library
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

func doCopy(dv, sv reflect.Value) []error {
	dv = indirect(dv)
	sv = indirect(sv)
	fields := modelFields(sv)

	var errs []error

	for _, f := range fields {
		sfv := sv.FieldByName(f.Name)
		tag := newTag(f.Tag.Get(TagName))

		if tag.isOmitField() {
			continue
		}

		// check type is in NoTraverseTypeList or has 'notraverse' tag option
		noTraverse := (isNoTraverseType(sfv) || tag.isNoTraverse())

		// check whether field is zero or not
		var isVal bool
		if isStruct(sfv) && !noTraverse {
			isVal = !IsZero(sfv.Interface())
		} else {
			isVal = !isFieldZero(sfv)
		}

		// get dst field by name
		dfv := dv.FieldByName(f.Name)

		// validate field - exists in dst, kind and type
		err := validateCopyField(f, sfv, dfv)
		if err != nil {
			if err != errFieldNotExists {
				errs = append(errs, err)
			}

			continue
		}

		// if value is not exists
		if !isVal {
			// field value is zero and check 'omitempty' option present
			// then don't copy into destination struct
			// otherwise copy to dst
			if !tag.isOmitEmpty() {
				dfv.Set(zeroOf(dfv))
			}
			continue
		}

		// check dst field settable or not
		if dfv.CanSet() {
			if isStruct(sfv) {
				// handle embedded or nested struct
				v, innerErrs := copyVal(dfv.Type(), sfv, noTraverse)

				// add errors to main stream
				errs = append(errs, innerErrs...)

				// handle based on ptr/non-ptr value
				dfv.Set(v)
			} else {
				v, err := copyVal(dfv.Type(), sfv, false)
				errs = append(errs, err...)
				dfv.Set(v)
			}
		}
	}

	return errs
}

func doMap(sv reflect.Value) map[string]interface{} {
	sv = indirect(sv)
	fields := modelFields(sv)
	m := map[string]interface{}{}

	for _, f := range fields {
		fv := sv.FieldByName(f.Name)
		tag := newTag(f.Tag.Get(TagName))

		if tag.isOmitField() {
			continue
		}

		// map key name
		keyName := f.Name
		if !isStringEmpty(tag.Name) {
			keyName = tag.Name
		}

		// check type is in NoTraverseTypeList or has 'notraverse' tag option
		noTraverse := (isNoTraverseType(fv) || tag.isNoTraverse())

		// check whether field is zero or not
		var isVal bool
		if isStruct(fv) && !noTraverse {
			isVal = !IsZero(fv.Interface())
		} else {
			isVal = !isFieldZero(fv)
		}

		if !isVal {
			// field value is zero and has 'omitempty' option present
			// then not include in the Map
			if !tag.isOmitEmpty() {
				m[keyName] = zeroOf(fv).Interface()
			}

			continue
		}

		// handle embedded or nested struct
		if isStruct(fv) {

			if noTraverse {
				// This is struct kind and it's present in NoTraverseTypeList or
				// has 'notraverse' tag option. So go-model is not gonna traverse inside.
				// however will take care of field value
				m[keyName] = mapVal(fv, true).Interface()
			} else {

				// embedded struct values gets mapped at embedded level
				// as represented by Go instead of object
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
	}

	return m
}

func copyVal(dt reflect.Type, f reflect.Value, notraverse bool) (reflect.Value, []error) {
	var (
		ptr  bool
		nf   reflect.Value
		errs []error
	)

	if conversionExists(f.Type(), dt) && !notraverse {
		// handle custom converters
		res, err := converterMap[f.Type()][dt](f)
		if err != nil {
			errs = append(errs, err)
		}
		return res, errs
	}

	// take care interface{} and its actual value
	if isInterface(f) {
		f = valueOf(f.Interface())
	}

	// if ptr, let's take a note
	if isPtr(f) {
		ptr = true
		f = f.Elem()
	}

	// two dimensional slice is not yet supported by this library
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
		if dt.Kind() == reflect.Ptr {
			dt = dt.Elem()
		}
		nf = reflect.MakeMap(dt)

		for _, key := range f.MapKeys() {
			ov := f.MapIndex(key)

			cv := reflect.New(dt.Elem()).Elem()
			v, err := copyVal(dt.Elem(), ov, isNoTraverseType(ov))
			if len(err) > 0 {
				errs = append(errs, err...)
			} else {
				cv.Set(v)
				nf.SetMapIndex(key, cv)
			}
		}
	case reflect.Slice:
		if f.Type() == typeOfBytes {
			nf = f
		} else {
			if dt.Kind() == reflect.Ptr {
				dt = dt.Elem()
			}
			nf = reflect.MakeSlice(dt, f.Len(), f.Cap())

			for i := 0; i < f.Len(); i++ {
				ov := f.Index(i)

				cv := reflect.New(dt.Elem()).Elem()
				v, err := copyVal(dt.Elem(), ov, isNoTraverseType(ov))
				if len(err) > 0 {
					errs = append(errs, err...)
				} else {
					cv.Set(v)
					nf.Index(i).Set(cv)
				}
			}
		}
	default:
		nf = f
	}

	if ptr {
		// wrap
		o := reflect.New(nf.Type())
		o.Elem().Set(nf)

		return o, errs
	}

	return nf, errs
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

	// if ptr, let's take a note
	if isPtr(f) {
		ptr = true
		f = f.Elem()
	}

	// two dimensional slice is not yet supported by this library
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
			nv := mapVal(mv, isNoTraverseType(mv))
			nmv[skey] = nv.Interface()
		}

		nf = valueOf(nmv)
	case reflect.Slice:
		if f.Type() == typeOfBytes {
			nf = f
		} else {
			if f.Len() > 0 {
				fsv := f.Index(0)

				// figure out target slice type
				if isStruct(fsv) {
					nf = reflect.MakeSlice(reflect.SliceOf(typeOfInterface), f.Len(), f.Cap())
				} else {
					nf = reflect.MakeSlice(f.Type(), f.Len(), f.Cap())
				}

				for i := 0; i < f.Len(); i++ {
					sv := f.Index(i)

					var dv reflect.Value
					if isStruct(sv) {
						dv = reflect.New(typeOfInterface).Elem()
					} else {
						dv = reflect.New(sv.Type()).Elem()
					}

					dv.Set(mapVal(sv, isNoTraverseType(sv)))
					nf.Index(i).Set(dv)
				}
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
