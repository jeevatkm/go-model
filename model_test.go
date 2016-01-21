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

func TestCopyBooleanAndBooleanPtr(t *testing.T) {
	type SampleStruct struct {
		Boolean    bool
		BooleanPtr *bool
	}

	boolPtr := true
	src := SampleStruct{
		Boolean:    true,
		BooleanPtr: &boolPtr,
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src, false)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.Boolean, dst.Boolean)
	assertEqual(t, *src.BooleanPtr, *dst.BooleanPtr)
	assertEqual(t, true, src.BooleanPtr != dst.BooleanPtr)
}

func TestCopyFloatAndFloatPtr(t *testing.T) {
	type SampleStruct struct {
		Float32    float32
		Float64    float64
		Float32Ptr *float32
		Float64Ptr *float64
	}

	f32 := float32(0.1)
	f64 := float64(0.2)

	src := SampleStruct{
		Float32:    float32(0.11),
		Float32Ptr: &f32,
		Float64:    float64(0.22),
		Float64Ptr: &f64,
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src, false)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.Float32, dst.Float32)
	assertEqual(t, *src.Float32Ptr, *dst.Float32Ptr)

	assertEqual(t, src.Float64, dst.Float64)
	assertEqual(t, *src.Float64Ptr, *dst.Float64Ptr)

	assertEqual(t, true, src.Float32Ptr != dst.Float32Ptr)
	assertEqual(t, true, src.Float64Ptr != dst.Float64Ptr)
}

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
	type SampleSubInfo2 struct {
		SliceIntPtr    []*int
		SliceInt64Ptr  []*int64
		SliceStringPtr []*string
		SliceFloat32   []*float32
		SliceFloat64   []*float64
		SliceInterface []interface{}
	}

	type SampleSubInfo1 struct {
		SliceIntPtr    []*int
		SliceInt64Ptr  []*int64
		SliceStringPtr []*string
		SliceFloat32   []*float32
		SliceFloat64   []*float64
		SliceInterface []interface{}
		Level2         SampleSubInfo2
	}

	type SampleStruct struct {
		SliceIntPtr    []*int
		SliceInt64Ptr  []*int64
		SliceStringPtr []*string
		SliceFloat32   []*float32
		SliceFloat64   []*float64
		SliceInterface []interface{}
		Level1         SampleSubInfo1
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
		SliceInterface: []interface{}{&i1, i11, &str1, &f1, f11},
		Level1: SampleSubInfo1{
			SliceIntPtr:    []*int{&i1, &i2, &i3},
			SliceInt64Ptr:  []*int64{&i11, &i12, &i13},
			SliceStringPtr: []*string{&str1, &str2, &str3},
			SliceFloat32:   []*float32{&f1, &f2, &f3},
			SliceFloat64:   []*float64{&f11, &f12, &f13},
			SliceInterface: []interface{}{&i1, i11, &str1, &f1, f11},
			Level2: SampleSubInfo2{
				SliceIntPtr:    []*int{&i1, &i2, &i3},
				SliceInt64Ptr:  []*int64{&i11, &i12, &i13},
				SliceStringPtr: []*string{&str1, &str2, &str3},
				SliceFloat32:   []*float32{&f1, &f2, &f3},
				SliceFloat64:   []*float64{&f11, &f12, &f13},
				SliceInterface: []interface{}{&i1, i11, &str1, &f1, f11},
			},
		},
	}

	dst := SampleStruct{}

	errs := Copy(&dst, src, false)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	// Level 0 assertion
	assertEqual(t, true, src.SliceIntPtr[0] != dst.SliceIntPtr[0])
	assertEqual(t, src.SliceIntPtr, dst.SliceIntPtr)

	assertEqual(t, true, src.SliceInt64Ptr[0] != dst.SliceInt64Ptr[0])
	assertEqual(t, src.SliceInt64Ptr, dst.SliceInt64Ptr)

	assertEqual(t, true, src.SliceStringPtr[0] != dst.SliceStringPtr[0])
	assertEqual(t, src.SliceStringPtr, dst.SliceStringPtr)

	assertEqual(t, true, src.SliceFloat32[0] != dst.SliceFloat32[0])
	assertEqual(t, src.SliceFloat32, dst.SliceFloat32)

	assertEqual(t, true, src.SliceFloat64[0] != dst.SliceFloat64[0])
	assertEqual(t, src.SliceFloat64, dst.SliceFloat64)

	assertEqual(t, true, src.SliceInterface[0] != dst.SliceInterface[0])
	assertEqual(t, src.SliceInterface, dst.SliceInterface)

	// Level 1 assertion
	assertEqual(t, true, src.Level1.SliceIntPtr[0] != dst.Level1.SliceIntPtr[0])
	assertEqual(t, src.Level1.SliceIntPtr, dst.Level1.SliceIntPtr)

	assertEqual(t, true, src.Level1.SliceInt64Ptr[0] != dst.Level1.SliceInt64Ptr[0])
	assertEqual(t, src.Level1.SliceInt64Ptr, dst.Level1.SliceInt64Ptr)

	assertEqual(t, true, src.Level1.SliceStringPtr[0] != dst.Level1.SliceStringPtr[0])
	assertEqual(t, src.Level1.SliceStringPtr, dst.Level1.SliceStringPtr)

	assertEqual(t, true, src.Level1.SliceFloat32[0] != dst.Level1.SliceFloat32[0])
	assertEqual(t, src.Level1.SliceFloat32, dst.Level1.SliceFloat32)

	assertEqual(t, true, src.Level1.SliceFloat64[0] != dst.Level1.SliceFloat64[0])
	assertEqual(t, src.Level1.SliceFloat64, dst.Level1.SliceFloat64)

	assertEqual(t, true, src.Level1.SliceInterface[0] != dst.Level1.SliceInterface[0])
	assertEqual(t, src.Level1.SliceInterface, dst.Level1.SliceInterface)

	// Level 2 assertion
	assertEqual(t, true, src.Level1.Level2.SliceIntPtr[0] != dst.Level1.Level2.SliceIntPtr[0])
	assertEqual(t, src.Level1.SliceIntPtr, dst.Level1.SliceIntPtr)

	assertEqual(t, true, src.Level1.Level2.SliceInt64Ptr[0] != dst.Level1.Level2.SliceInt64Ptr[0])
	assertEqual(t, src.Level1.Level2.SliceInt64Ptr, dst.Level1.Level2.SliceInt64Ptr)

	assertEqual(t, true, src.Level1.Level2.SliceStringPtr[0] != dst.Level1.Level2.SliceStringPtr[0])
	assertEqual(t, src.Level1.Level2.SliceStringPtr, dst.Level1.Level2.SliceStringPtr)

	assertEqual(t, true, src.Level1.Level2.SliceFloat32[0] != dst.Level1.Level2.SliceFloat32[0])
	assertEqual(t, src.Level1.Level2.SliceFloat32, dst.Level1.Level2.SliceFloat32)

	assertEqual(t, true, src.Level1.Level2.SliceFloat64[0] != dst.Level1.Level2.SliceFloat64[0])
	assertEqual(t, src.Level1.Level2.SliceFloat64, dst.Level1.Level2.SliceFloat64)

	assertEqual(t, true, src.Level1.Level2.SliceInterface[0] != dst.Level1.Level2.SliceInterface[0])
	assertEqual(t, src.Level1.Level2.SliceInterface, dst.Level1.Level2.SliceInterface)
}

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
		Level1Struct           SampleSubInfo `model:",notraverse"`
		Level1StructPtr        *SampleSubInfo
		Level1StructNoTraverse *SampleSubInfo `model:",notraverse"`
		CreatedTime            time.Time
		SampleSubInfo
	}

	src := SampleStruct{
		SampleSubInfo:          SampleSubInfo{Name: "This embeded struct", Year: 2016},
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 2 struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
		CreatedTime:            time.Now(),
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
	assertEqual(t, src.Level1StructNoTraverse.Year, dst.Level1StructNoTraverse.Year)
}

func TestCopyStructEmbededAndAttributeDstPtr(t *testing.T) {
	type SampleSubInfo struct {
		Name string
		Year int
	}

	type SampleStruct struct {
		Level1Struct           SampleSubInfo `model:",notraverse"`
		Level1StructPtr        *SampleSubInfo
		Level1StructPtrZero    *SampleSubInfo
		Level1StructNoTraverse *SampleSubInfo `model:",notraverse"`
		CreatedTime            time.Time
		SampleSubInfo
	}

	src := SampleStruct{
		SampleSubInfo:          SampleSubInfo{Name: "This embeded struct", Year: 2016},
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 2 struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
		CreatedTime:            time.Now(),
	}

	dst := SampleStruct{
		Level1StructPtrZero: &SampleSubInfo{Name: "This level 1 struct ptr zero", Year: 2015},
	}

	errs := Copy(&dst, &src, true)
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
	assertEqual(t, src.Level1StructNoTraverse.Year, dst.Level1StructNoTraverse.Year)

	assertEqual(t, true, dst.Level1StructPtrZero == nil)
}

func TestCopyStructEmbededAndAttributeMakeZeroInDst(t *testing.T) {
	type SampleSubInfo struct {
		Name string
		Year int
	}

	type SampleStruct struct {
		Level1Struct           SampleSubInfo `model:",notraverse"`
		Level1StructPtr        *SampleSubInfo
		Level1StructNoTraverse *SampleSubInfo `model:",notraverse"`
		CreatedTime            time.Time
		SampleSubInfo
	}

	src := SampleStruct{CreatedTime: time.Now()}

	dst := SampleStruct{
		SampleSubInfo:          SampleSubInfo{Name: "This embeded struct", Year: 2016},
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 2 struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
	}

	errs := Copy(&dst, &src, true)
	if errs != nil {
		fmt.Println(errs)
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, true, src.CreatedTime == dst.CreatedTime)

	assertEqual(t, true, IsZero(dst.Level1Struct))
	assertEqual(t, true, IsZero(dst.SampleSubInfo))

	assertEqual(t, true, dst.Level1StructPtr == nil)
	assertEqual(t, true, dst.Level1StructNoTraverse == nil)
}

type SampleStruct struct {
	Integer             int
	IntegerPtr          *int
	String              string
	StringPtr           *string
	Boolean             bool
	BooleanPtr          *bool
	BooleanOmit         bool `model:"-"`
	SliceString         []string
	SliceStringOmit     []string `model:"-"`
	SliceStringPtr      *[]string
	SliceStringPtrOmit  *[]string `model:"-"`
	SliceStringPtrStr   []*string
	Float32             float32
	Float32Ptr          *float32
	Float32Omit         float32  `model:"-"`
	Float32PtrOmit      *float32 `model:"-"`
	Float64             float64
	Float64Ptr          *float64
	Float64Omit         float64  `model:"-"`
	Float64PtrOmit      *float64 `model:"-"`
	SliceStruct         []SampleSubInfo
	SliceStructPtr      []*SampleSubInfo
	SliceInt            []int
	SliceIntPtr         []*int
	Time                time.Time
	TimePtr             *time.Time
	Struct              SampleSubInfo
	StructPtr           *SampleSubInfo
	StructNoTraverse    SampleSubInfo  `model:",notraverse"`
	StructPtrNoTraverse *SampleSubInfo `model:",notraverse"`
	StructDeep          SampleSubInfoDeep
	StructDeepPtr       *SampleSubInfoDeep
	SampleSubInfo
}

type SampleSubInfo struct {
	Name string
	Year int
}

type SampleSubInfoDeep struct {
	Name                string
	NamePtr             *string `model:"-"`
	Year                int     `model:"-"`
	YearPtr             *int
	Struct              SampleSubInfo
	StructPtr           *SampleSubInfo
	StructNoTraverse    SampleSubInfo  `model:",notraverse"`
	StructPtrNoTraverse *SampleSubInfo `model:",notraverse"`
	SampleSubInfo
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

	assertEqual(t, "Field: 'Name', src [string] & dst [int] kind doesn't match", errs[0].Error())
}

func TestCopyStructElementTypeDiffOnLevel1(t *testing.T) {
	type SampleLevelSrc struct {
		Name string
	}

	type SampleLevelDst struct {
		Name int
	}

	type Source struct {
		Name   string
		Level1 SampleLevelSrc
	}

	type Destination struct {
		Name   int
		Level1 SampleLevelDst
	}

	src := Source{
		Name: "This struct element kind is different",
		Level1: SampleLevelSrc{
			Name: "Level1: This struct element kind is different",
		},
	}

	dst := Destination{}

	errs := Copy(&dst, src, false)

	logSrcDst(t, src, dst)

	assertEqual(t, "Field: 'Name', src [string] & dst [int] kind doesn't match", errs[0].Error())
	assertEqual(t,
		"Field: 'Level1', src [model.SampleLevelSrc] & dst [model.SampleLevelDst] type doesn't match",
		errs[1].Error(),
	)
}

func TestCopyStructTypeDiffOnLevel1Interface(t *testing.T) {
	type SampleLevelSrc struct {
		Name string
	}

	type Source struct {
		Name   string
		Level1 SampleLevelSrc
	}

	type Destination struct {
		Name   int
		Level1 interface{}
	}

	src := Source{
		Name: "This struct element kind is different",
		Level1: SampleLevelSrc{
			Name: "Level1: This struct element kind is different",
		},
	}

	dst := Destination{}

	errs := Copy(&dst, src, false)

	logSrcDst(t, src, dst)

	assertEqual(t, "Field: 'Name', src [string] & dst [int] kind doesn't match", errs[0].Error())
	assertEqual(t, 0, dst.Name)
	assertEqual(t, src.Level1.Name, dst.Level1.(SampleLevelSrc).Name)
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

	if !IsZero(&SampleStruct{Struct: SampleSubInfo{}, StructPtr: &SampleSubInfo{}}) {
		t.Error("SampleStruct with sub struct 1 - supposed to be zero")
	}

	if !IsZero(&SampleStruct{Struct: SampleSubInfo{Name: "go-model"}, StructPtr: &SampleSubInfo{}}) {
		t.Log("SampleStruct with sub struct 2 - supposed to be zero")
	} else {
		t.Error("SampleStruct with sub struct 2 - supposed to be zero")
	}

	deepStruct := SampleStruct{
		StructDeepPtr: &SampleSubInfoDeep{
			StructPtr: &SampleSubInfo{
				Name: "I'm here",
			},
			StructNoTraverse: SampleSubInfo{
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
		t.Error("SampleStruct notraverse - supposed to be zero")
	}

	if IsZero(SampleStruct{SampleSubInfo: SampleSubInfo{Year: 2010}}) {
		t.Error("SampleStruct embeded struct - supposed to be non-zero")
	}
}

func TestIsStructMethod(t *testing.T) {
	src := map[string]interface{}{
		"struct": &SampleStruct{Time: time.Now()},
	}

	mv := valueOf(src)
	keys := mv.MapKeys()

	assertEqual(t, true, isStruct(mv.MapIndex(keys[0])))
}

func TestIsZeroNotAStructInput(t *testing.T) {
	result1 := IsZero(10001)
	assertEqual(t, false, result1)

	result2 := IsZero(map[string]int{"1": 101, "2": 102, "3": 103})
	assertEqual(t, false, result2)

	floatVar := float64(1.7367643)
	result3 := IsZero(&floatVar)
	assertEqual(t, false, result3)

	str := "This is not a struct"
	result4 := IsZero(&str)
	assertEqual(t, false, result4)
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
	t.Log()
	logIt(t, "Destination", dst)
}

func logIt(t *testing.T, str string, v interface{}) {
	t.Logf("%v: %#v", str, v)
}
