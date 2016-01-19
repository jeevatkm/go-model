// Copyright (c) 2016 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

type SampleStruct struct {
	Integer            int
	IntegerPtr         *int
	String             string
	StringPtr          *string
	Boolean            bool
	BooleanPtr         *bool
	BooleanOmit        bool `model:"-"`
	SliceString        []string
	SliceStringPtr     *[]string
	SliceStringPtrOmit *[]string `model:"-"`
	SliceStringPtrStr  []*string
	SliceStruct        []SubInfo
	SliceStructPtr     []*SubInfo
	SliceInt           []int
	SliceIntPtr        []*int
	Time               time.Time
	TimePtr            *time.Time
	Struct             SubInfo
	StructPtr          *SubInfo
	StructOmit         SubInfo  `model:",omitnested"`
	StructPtrOmit      *SubInfo `model:",omitnested"`
	StructDeep         SubInfoDeep
	StructDeepPtr      *SubInfoDeep
	SubInfo
}

type SubInfo struct {
	Name string
	Year int
}

type SubInfoDeep struct {
	Name          string
	NamePtr       *string `model:"-"`
	Year          int     `model:"-"`
	YearPtr       *int
	Struct        SubInfo
	StructPtr     *SubInfo
	StructOmit    SubInfo  `model:",omitnested"`
	StructPtrOmit *SubInfo `model:",omitnested"`
	SubInfo
}

func TestAddOmitNested(t *testing.T) {
	if !isOmitNestedType(valueOf(os.File{})) {
		t.Errorf("Given type not found in omit list")
	}

	// Already registered
	AddToOmitNested(os.File{})
}

func TestIsZero(t *testing.T) {
	if !IsZero(SampleStruct{}) {
		t.Error("SampleStruct - supposed to be zero")
	}

	if !IsZero(&SampleStruct{}) {
		t.Error("SampleStruct Ptr - supposed to be zero")
	}

	if !IsZero(&SampleStruct{Struct: SubInfo{}, StructPtr: &SubInfo{}}) {
		t.Error("SampleStruct with sub struct 1 - supposed to be zero")
	}

	if !IsZero(&SampleStruct{Struct: SubInfo{Name: "go-model"}, StructPtr: &SubInfo{}}) {
		t.Log("SampleStruct with sub struct 2 - supposed to be zero")
	} else {
		t.Error("SampleStruct with sub struct 2 - supposed to be zero")
	}

	deepStruct := SampleStruct{
		StructDeepPtr: &SubInfoDeep{
			StructPtr: &SubInfo{
				Name: "I'm here",
			},
			StructOmit: SubInfo{
				Year: 2005,
			},
		},
	}
	if IsZero(deepStruct) {
		t.Error("SampleStruct deep level - supposed to be non-zero")
	}
}

func TestNonZeroCheck(t *testing.T) {
	if IsZero(&SampleStruct{Time: time.Now()}) {
		t.Error("SampleStruct omitnested - supposed to be zero")
	}

	if IsZero(SampleStruct{SubInfo: SubInfo{Year: 2010}}) {
		t.Error("SampleStruct embeded struct - supposed to be non-zero")
	}
}

func TestCopyOnlyNonZeroPtr(t *testing.T) {
	intPtr := 1002
	deepIntPtr := 2006
	deepNamePtr := "Deep: go-model struct ptr"
	src := SampleStruct{
		Integer:    1001,
		IntegerPtr: &intPtr,
		Struct:     SubInfo{Name: "go-model struct", Year: 2014},
		StructPtr:  &SubInfo{Name: "go-model struct ptr", Year: 2015},
		StructDeep: SubInfoDeep{
			Name:      "Deep: go-model struct",
			Year:      1996,
			NamePtr:   &deepNamePtr,
			YearPtr:   &deepIntPtr,
			Struct:    SubInfo{Name: "Deep: go-model struct", Year: 1994},
			StructPtr: &SubInfo{Name: "Deep: go-model struct ptr", Year: 1995},
			SubInfo:   SubInfo{Name: "Deep: go-model embeded", Year: 1996},
		},
		StructDeepPtr: &SubInfoDeep{
			Name:      "Deep: go-model struct ptr",
			Year:      2006,
			NamePtr:   &deepNamePtr,
			YearPtr:   &deepIntPtr,
			Struct:    SubInfo{Name: "Deep: go-model struct", Year: 2004},
			StructPtr: &SubInfo{Name: "Deep: go-model struct ptr", Year: 2005},
			SubInfo:   SubInfo{Name: "Deep: go-model embeded", Year: 2006},
		},
		SubInfo: SubInfo{Name: "go-model embeded", Year: 2016},
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src, false)
	fmt.Println("Errors:", errs)

	fmt.Printf("\nSource     : %#v\n", src)
	fmt.Printf("\nDestination: %#v\n", dst)
	fmt.Println()

	assertEqual(t, src.Integer, dst.Integer)
	assertEqual(t, true, src.IntegerPtr != dst.IntegerPtr)
	assertEqual(t, *src.IntegerPtr, *dst.IntegerPtr)

	// Level 1 struct Assertion
	assertEqual(t, src.Struct.Name, dst.Struct.Name)
	assertEqual(t, src.Struct.Year, dst.Struct.Year)

	assertEqual(t, true, src.StructPtr != dst.StructPtr)
	assertEqual(t, src.StructPtr.Name, dst.StructPtr.Name)
	assertEqual(t, src.StructPtr.Year, dst.StructPtr.Year)

	// Level 2 strcut Assertion
	assertEqual(t, 0, dst.StructDeep.Year)
	assertEqual(t, true, dst.StructDeep.NamePtr == nil)
	assertEqual(t, src.StructDeep.Struct.Name, dst.StructDeep.Struct.Name)
	assertEqual(t, true, src.StructDeep.StructPtr != dst.StructDeep.StructPtr)
	assertEqual(t, src.StructDeep.StructPtr.Year, dst.StructDeep.StructPtr.Year)
	assertEqual(t, src.StructDeep.SubInfo.Year, dst.StructDeep.SubInfo.Year)

	assertEqual(t, 0, dst.StructDeepPtr.Year)
	assertEqual(t, true, dst.StructDeepPtr.NamePtr == nil)
	assertEqual(t, true, src.StructDeepPtr != dst.StructDeepPtr)
	assertEqual(t, src.StructDeepPtr.Struct.Name, dst.StructDeepPtr.Struct.Name)
	assertEqual(t, true, src.StructDeepPtr.StructPtr != dst.StructDeepPtr.StructPtr)
	assertEqual(t, src.StructDeepPtr.StructPtr.Year, dst.StructDeepPtr.StructPtr.Year)
	assertEqual(t, src.StructDeepPtr.SubInfo.Year, dst.StructDeepPtr.SubInfo.Year)
}

func assertError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error occurred [%v]", err)
	}
}

func assertEqual(t *testing.T, e, g interface{}) (r bool) {
	r = compare(e, g)
	if !r {
		t.Errorf("Expected [%v], got [%v]", e, g)
	}

	return
}

func assertNotEqual(t *testing.T, e, g interface{}) (r bool) {
	if compare(e, g) {
		t.Errorf("Expected [%v], got [%v]", e, g)
	} else {
		r = true
	}

	return
}

func compare(e, g interface{}) (r bool) {
	ev := reflect.ValueOf(e)
	gv := reflect.ValueOf(g)

	if ev.Kind() != gv.Kind() {
		return
	}

	switch ev.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r = (ev.Int() == gv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r = (ev.Uint() == gv.Uint())
	case reflect.Float32, reflect.Float64:
		r = (ev.Float() == gv.Float())
	case reflect.String:
		r = (ev.String() == gv.String())
	case reflect.Bool:
		r = (ev.Bool() == gv.Bool())
	}

	return
}
