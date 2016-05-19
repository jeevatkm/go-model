// Copyright (c) 2016 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
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
	Version = "0.4"

	// NoTraverseTypeList keeps track of no-traverse type list at library level
	NoTraverseTypeList map[reflect.Type]bool

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

// Tag method returns the exported struct field `Tag` value from the given struct.
// 		Example:
//
// 		src := SampleStruct {
// 			BookCount      int         `json:"-"`
// 			BookCode       string      `json:"-"`
// 			ArchiveInfo    BookArchive `json:"archive_info,omitempty"`
// 			Region         BookLocale  `json:"region,omitempty"`
// 		}
//
// 		tag, _ := model.Tag(src, "ArchiveInfo")
// 		fmt.Println("Tag Value:", tag.Get("json"))
//
// 		// Output:
// 		Tag Value: archive_info,omitempty
//
func Tag(s interface{}, name string) (reflect.StructTag, error) {
	sv, err := structValue(s)
	if err != nil {
		return "", err
	}

	fv, ok := sv.Type().FieldByName(name)
	if !ok {
		return "", fmt.Errorf("Field: '%v', does not exists", name)
	}

	return fv.Tag, nil
}

// Tags method returns the exported struct fields `Tag` value from the given struct.
// 		Example:
//
// 		src := SampleStruct {
// 			BookCount      int         `json:"-"`
// 			BookCode       string      `json:"-"`
// 			ArchiveInfo    BookArchive `json:"archive_info,omitempty"`
// 			Region         BookLocale  `json:"region,omitempty"`
// 		}
//
// 		tags, _ := model.Tags(src)
// 		fmt.Println("Tags:", tags)
//
func Tags(s interface{}) (map[string]reflect.StructTag, error) {
	sv, err := structValue(s)
	if err != nil {
		return nil, err
	}

	tags := map[string]reflect.StructTag{}

	fields := modelFields(sv)
	for _, f := range fields {
		tags[f.Name] = f.Tag
	}

	return tags, nil
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
	fv, err := getField(s, name)
	if err != nil {
		return reflect.Invalid, err
	}

	return fv.Type().Kind(), nil
}

//
// go-model init
//

func init() {
	NoTraverseTypeList = map[reflect.Type]bool{}

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

		// validate field - exists in dst, kind and type
		err := valiadateCopyField(f, sfv, dfv)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		// check dst field settable or not
		if dfv.CanSet() {

			// handle embedded or nested struct
			if isStruct(sfv) {

				if noTraverse {
					// This is struct kind and it's present in NoTraverseTypeList or
					// has 'notraverse' tag option. So go-model is not gonna traverse inside.
					// however will take care of field value
					dfv.Set(copyVal(sfv, true))
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

			dfv.Set(copyVal(sfv, false))
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

func copyVal(f reflect.Value, notraverse bool) reflect.Value {
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
			cv.Set(copyVal(ov, isNoTraverseType(ov)))

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
				cv.Set(copyVal(ov, isNoTraverseType(ov)))

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

	if _, ok := NoTraverseTypeList[t]; ok {
		return true
	}

	return false
}

func valiadateCopyField(f reflect.StructField, sfv, dfv reflect.Value) error {
	// check dst field is exists, if not valid move on
	if !dfv.IsValid() {
		return fmt.Errorf("Field: '%v', does not exists in dst", f.Name)
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
		// TODO Go 1.6 changes -> f.PkgPath != "" && !f.Anonymous
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

func getField(s interface{}, name string) (reflect.Value, error) {
	sv, err := structValue(s)
	if err != nil {
		return reflect.Value{}, err
	}

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
		// vt := deepTypeOf(v)
		// return vt.Elem().Kind() == reflect.Struct
		return false
	}

	return pv.Kind() == reflect.Struct
}

func isInterface(v reflect.Value) bool {
	return v.Kind() == reflect.Interface
}
