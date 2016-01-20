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

func TestCopyIntegerAndIntegerPtr(t *testing.T) {
	type SampleStruct struct {
		Int      int
		IntPtr   *int
		Int64    int64
		Int64Ptr *int64
	}

	intPtr := int(1001)
	int64Ptr := int64(1002)

	src := SampleStruct{
		Int:      2001,
		IntPtr:   &intPtr,
		Int64:    2002,
		Int64Ptr: &int64Ptr,
	}

	dst := SampleStruct{}

	errs := Copy(&dst, src, false)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.Int, dst.Int)
	assertEqual(t, src.Int64, dst.Int64)

	assertEqual(t, true, src.IntPtr != dst.IntPtr)
	assertEqual(t, *src.IntPtr, *dst.IntPtr)
	assertEqual(t, *src.Int64Ptr, *dst.Int64Ptr)
}

func TestCopyStringAndStringPtr(t *testing.T) {
	type SampleStruct struct {
		String    string
		StringPtr *string
	}

	strPtr := "This is string for pointer test"
	src := SampleStruct{
		String:    "This is string for test",
		StringPtr: &strPtr,
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src, false)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.String, dst.String)
	assertEqual(t, *src.StringPtr, *dst.StringPtr)
	assertEqual(t, true, src.StringPtr != dst.StringPtr)
}

//
// TODO for Boolean, float, etc
//

func TestCopySliceStringAndSliceStringPtr(t *testing.T) {
	type SampleStruct struct {
		SliceString    []string
		SliceStringPtr *[]string
	}

	sliceStrPtr := []string{"This is slice string test pointer."}
	src := SampleStruct{
		SliceString:    []string{"This is slice string test."},
		SliceStringPtr: &sliceStrPtr,
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src, false)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.SliceString, dst.SliceString)
	assertEqual(t, *src.SliceStringPtr, *dst.SliceStringPtr)
	assertEqual(t, true, src.SliceStringPtr != dst.SliceStringPtr)
}

func TestCopySliceElementsPtr(t *testing.T) {
	type SampleStruct struct {
		SliceIntPtr    []*int
		SliceInt64Ptr  []*int64
		SliceStringPtr []*string
		SliceFloat32   []*float32
		SliceFloat64   []*float64
	}

	i1 := int(1)
	i2 := int(2)
	i3 := int(3)

	i11 := int64(11)
	i12 := int64(12)
	i13 := int64(13)

	str1 := "This is string pointer 1"
	str2 := "This is string pointer 2"
	str3 := "This is string pointer 3"

	f1 := float32(0.1)
	f2 := float32(0.2)
	f3 := float32(0.3)

	f11 := float64(0.11)
	f12 := float64(0.12)
	f13 := float64(0.13)

	src := SampleStruct{
		SliceIntPtr:    []*int{&i1, &i2, &i3},
		SliceInt64Ptr:  []*int64{&i11, &i12, &i13},
		SliceStringPtr: []*string{&str1, &str2, &str3},
		SliceFloat32:   []*float32{&f1, &f2, &f3},
		SliceFloat64:   []*float64{&f11, &f12, &f13},
	}

	dst := SampleStruct{}

	errs := Copy(&dst, src, false)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.SliceIntPtr, dst.SliceIntPtr)
	assertEqual(t, src.SliceInt64Ptr, dst.SliceInt64Ptr)
	assertEqual(t, src.SliceStringPtr, dst.SliceStringPtr)
	assertEqual(t, src.SliceFloat32, dst.SliceFloat32)
	assertEqual(t, src.SliceFloat64, dst.SliceFloat64)
}

//
// TODO slice with interface{}, etc.
//

func TestCopyMapElements(t *testing.T) {
	type SampleSubInfo struct {
		Name string
		Year int
	}

	type SampleStruct struct {
		MapIntInt       map[int]int
		MapStringInt    map[string]int
		MapStringString map[string]string
		MapStruct       map[string]SampleSubInfo
		MapInterfaces   map[string]interface{}
	}

	src := SampleStruct{
		MapIntInt:       map[int]int{1: 1001, 2: 1002, 3: 1003, 4: 1004},
		MapStringInt:    map[string]int{"first": 1001, "second": 1002, "third": 1003, "forth": 1004},
		MapStringString: map[string]string{"first": "1001", "second": "1002", "third": "1003"},
		MapStruct: map[string]SampleSubInfo{
			"struct1": SampleSubInfo{Name: "struct 1 value", Year: 2001},
			"struct2": SampleSubInfo{Name: "struct 2 value", Year: 2002},
			"struct3": SampleSubInfo{Name: "struct 3 value", Year: 2003},
		},
		MapInterfaces: map[string]interface{}{
			"inter1": 100001,
			"inter2": "This is my interface string",
			"inter3": SampleSubInfo{Name: "struct 3 value", Year: 2003},
			"inter4": float32(1.6546565),
			"inter5": float64(1.6546565),
			"inter6": &SampleSubInfo{Name: "struct 3 value", Year: 2006},
		},
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src, false)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.MapIntInt, dst.MapIntInt)
	assertEqual(t, src.MapStringInt, dst.MapStringInt)
	assertEqual(t, src.MapStringString, dst.MapStringString)
	assertEqual(t, src.MapStruct, dst.MapStruct)
	assertEqual(t, src.MapInterfaces, dst.MapInterfaces)
}

func TestCopyStructEmbededAndAttribute(t *testing.T) {
	type SampleSubInfo struct {
		Name string
		Year int
	}

	type SampleStruct struct {
		SampleSubInfo
		Level1Struct     SampleSubInfo
		Level1StructPtr  *SampleSubInfo
		Level1StructOmit *SampleSubInfo `model:",omitnested"`
		CreatedTime      time.Time
	}

	src := SampleStruct{
		SampleSubInfo:    SampleSubInfo{Name: "This embeded struct", Year: 2016},
		Level1Struct:     SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:  &SampleSubInfo{Name: "This level 2 struct", Year: 2014},
		Level1StructOmit: &SampleSubInfo{Name: "This omit nested traverse struct", Year: 2013},
		CreatedTime:      time.Now(),
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src, false)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.Name, dst.Name)
	assertEqual(t, src.Year, dst.Year)

	assertEqual(t, src.Level1Struct.Name, dst.Level1Struct.Name)
	assertEqual(t, src.Level1Struct.Year, dst.Level1Struct.Year)

	assertEqual(t, src.Level1StructPtr.Name, dst.Level1StructPtr.Name)
	assertEqual(t, src.Level1StructPtr.Year, dst.Level1StructPtr.Year)

	assertEqual(t, true, src.CreatedTime == dst.CreatedTime)
	assertEqual(t, src.Level1StructOmit.Year, dst.Level1StructOmit.Year)
}

func TestCopyZeroInput(t *testing.T) {
	errs := Copy(&SampleStruct{}, SampleStruct{}, false)

	assertEqual(t, "Source struct is empty", errs[0].Error())
}

func TestCopyDestinationIsNotPointer(t *testing.T) {
	type SampleStruct struct {
		Name string
	}
	errs := Copy(SampleStruct{}, SampleStruct{Name: "Not a pointer"}, false)

	assertEqual(t, "Destination struct is not a pointer", errs[0].Error())
}

func TestCopyInputIsNotStruct(t *testing.T) {
	type SampleStruct struct {
		Name string
	}
	errs := Copy(&SampleStruct{}, map[string]string{"1": "2001"}, false)

	assertEqual(t, "Source or Destination is not a struct", errs[0].Error())
}

func TestCopyStructElementKindDiff(t *testing.T) {
	type Source struct {
		Name string
	}

	type Destination struct {
		Name int
	}

	errs := Copy(&Destination{}, Source{Name: "This struct element kind is different"}, false)

	assertEqual(t, "Field: 'Name', src & dst kind doesn't match", errs[0].Error())
}

func TestCopyStructElementIsNotValidInDst(t *testing.T) {
	type Source struct {
		Name string
		Year int
	}

	type Destination struct {
		Name string
	}

	src := Source{Year: 2016}
	dst := Destination{Name: "Value is gonna disappear"}

	errs := Copy(&dst, src, true)

	assertEqual(t, "Field: 'Year', dst is not valid", errs[0].Error())
}

func TestCopyStructZeroValToDst(t *testing.T) {
	type Source struct {
		Name string
		Year int
	}

	type Destination struct {
		Name string
		Year int
	}

	src := Source{Year: 2016}
	dst := Destination{Name: "Value is gonna disappear"}

	errs := Copy(&dst, src, true)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, "", dst.Name)
	assertEqual(t, 2016, dst.Year)
}

func TestAddNoTraverseType(t *testing.T) {
	if !isNoTraverseType(valueOf(os.File{})) {
		t.Errorf("Given type not found in omit list")
	}

	// Already registered
	AddNoTraverseType(os.File{})
}

func TestRemoveNoTraverseType(t *testing.T) {
	RemoveNoTraverseType(os.File{})

	if isNoTraverseType(valueOf(os.File{})) {
		t.Errorf("Type should not exists in the NoTraverseTypeList")
	}

	AddNoTraverseType(os.File{})

	// test again
	if !isNoTraverseType(valueOf(os.File{})) {
		t.Errorf("Type should exists in the NoTraverseTypeList")
	}
}

func TestIsZero(t *testing.T) {
	IsZero(nil) // nil check

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
	case reflect.Slice, reflect.Map:
		r = reflect.DeepEqual(e, g)
	}

	return
}

func logSrcDst(t *testing.T, src, dst interface{}) {
	logIt(t, "Source", src)
	fmt.Println()
	logIt(t, "Destination", dst)
}

func logIt(t *testing.T, str string, v interface{}) {
	t.Logf("%v: %#v", str, v)
}
