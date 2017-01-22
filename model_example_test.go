// Copyright (c) Jeevanandam M. (https://github.com/jeevatkm).
// go-model source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"reflect"
	"strconv"
)

//
// Examples
//

// Register a custom `Converter` to allow conversions from `int` to `string`.
func ExampleAddConversion() {
	AddConversion((*int)(nil), (*string)(nil), func(in reflect.Value) (reflect.Value, error) {
		return reflect.ValueOf(strconv.FormatInt(in.Int(), 10)), nil
	})
	type StructA struct {
		Mixed string
	}

	type StructB struct {
		Mixed int
	}
	src := StructB{Mixed: 123}
	dst := StructA{}

	errs := Copy(&dst, &src)
	if errs != nil {
		panic(errs)
	}
	fmt.Printf("%v", dst)
	// Output: {123}
}

// Register a custom `Converter` to allow conversions from `*int` to `string`.
func ExampleAddConversion_sourcePointer() {
	AddConversion((**int)(nil), (*string)(nil), func(in reflect.Value) (reflect.Value, error) {
		return reflect.ValueOf(strconv.FormatInt(in.Elem().Int(), 10)), nil
	})
	type StructA struct {
		Mixed string
	}

	type StructB struct {
		Mixed *int
	}
	val := 123
	src := StructB{Mixed: &val}
	dst := StructA{}

	errs := Copy(&dst, &src)
	if errs != nil {
		panic(errs[0])
	}
	fmt.Printf("%v", dst)
	// Output: {123}
}

// Register a custom `Converter` to allow conversions from `int` to `*string`.
func ExampleAddConversion_destinationPointer() {
	AddConversion((*int)(nil), (**string)(nil), func(in reflect.Value) (reflect.Value, error) {
		str := strconv.FormatInt(in.Int(), 10)
		return reflect.ValueOf(&str), nil
	})
	type StructA struct {
		Mixed *string
	}

	type StructB struct {
		Mixed int
	}
	src := StructB{Mixed: 123}
	dst := StructA{}

	errs := Copy(&dst, &src)
	if errs != nil {
		panic(errs[0])
	}
	fmt.Printf("%v", *dst.Mixed)
	// Output: 123
}

// Register a custom `Converter` to allow conversions from `int` to `*string` by passing types.
func ExampleAddConversion_destinationPointerByType() {
	srcType := reflect.TypeOf((*int)(nil)).Elem()
	targetType := reflect.TypeOf((**string)(nil)).Elem()
	AddConversionByType(srcType, targetType, func(in reflect.Value) (reflect.Value, error) {
		str := strconv.FormatInt(in.Int(), 10)
		return reflect.ValueOf(&str), nil
	})
	type StructA struct {
		Mixed *string
	}

	type StructB struct {
		Mixed int
	}
	src := StructB{Mixed: 123}
	dst := StructA{}

	errs := Copy(&dst, &src)
	if errs != nil {
		panic(errs[0])
	}
	fmt.Printf("%v", *dst.Mixed)
	// Output: 123
}
