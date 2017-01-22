// Copyright (c) Jeevanandam M. (https://github.com/jeevatkm).
// go-model source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

//
// Copy test cases
//

func TestConverter(t *testing.T) {
	type SampleStructA struct {
		Int    int
		String string
		Mixed  string
	}

	type SampleStructB struct {
		Int    int
		String string
		Mixed  int
	}

	AddConversion((*int)(nil), (*string)(nil), func(in reflect.Value) (reflect.Value, error) {
		return reflect.ValueOf(strconv.FormatInt(in.Int(), 10) + "lala"), nil
	})

	src := SampleStructB{Mixed: 123, Int: 5, String: "string"}
	dst := SampleStructA{}

	errs := Copy(&dst, src)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}
	assertEqual(t, "123lala", dst.Mixed)
	assertEqual(t, 5, dst.Int)
	assertEqual(t, "string", dst.String)
}

func TestMissingConverter(t *testing.T) {
	type SampleStructA struct {
		Int    int
		String string
		Mixed  string
	}

	type SampleStructB struct {
		Int    int
		String string
		Mixed  int
	}

	RemoveConversion((*int)(nil), (*string)(nil))

	src := SampleStructB{Mixed: 123, Int: 5, String: "string"}
	dst := SampleStructA{}

	errs := Copy(&dst, src)
	if errs == nil {
		t.Error("The conversion between int and string should have failed.")
	}
	assertEqual(t, "", dst.Mixed)
	assertEqual(t, 5, dst.Int)
	assertEqual(t, "string", dst.String)
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

	errs := Copy(&dst, src)
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

	errs := Copy(&dst, &src)
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
		Boolean             bool
		BooleanPtr          *bool
		BooleanOmitEmpty    bool  `model:",omitempty"`
		BooleanOmitEmptyPtr *bool `model:",omitempty"`
	}

	boolPtr := true
	src := SampleStruct{
		Boolean:    true,
		BooleanPtr: &boolPtr,
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.Boolean, dst.Boolean)
	assertEqual(t, *src.BooleanPtr, *dst.BooleanPtr)
	assertEqual(t, true, src.BooleanPtr != dst.BooleanPtr)
	assertEqual(t, false, dst.BooleanOmitEmpty)
	assertEqual(t, true, dst.BooleanOmitEmptyPtr == nil)
}

func TestCopyFloatAndFloatPtr(t *testing.T) {
	type SampleStruct struct {
		Float32        float32
		Float64        float64
		Float32PtrOmit *float32 `model:"-"`
		Float32Ptr     *float32
		Float64Ptr     *float64
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

	errs := Copy(&dst, &src)
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

	errs := Copy(&dst, &src)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.SliceString, dst.SliceString)
	assertEqual(t, *src.SliceStringPtr, *dst.SliceStringPtr)
	assertEqual(t, true, src.SliceStringPtr != dst.SliceStringPtr)
}

func TestCopyByteAndByteSlice(t *testing.T) {
	type SampleStruct struct {
		Byte          byte
		SliceBytes    []byte
		SliceBytesPtr *[]byte
	}

	bytesPtr := []byte("This is byte pointer value")

	src := SampleStruct{
		Byte:          byte('A'),
		SliceBytes:    []byte("This is byte value"),
		SliceBytesPtr: &bytesPtr,
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, src.Byte, dst.Byte)
	assertEqual(t, true, &src.Byte != &dst.Byte)

	assertEqual(t, src.SliceBytes, dst.SliceBytes)
	assertEqual(t, true, &src.SliceBytes != &dst.SliceBytes)

	assertEqual(t, *src.SliceBytesPtr, *dst.SliceBytesPtr)
	assertEqual(t, true, src.SliceBytesPtr != dst.SliceBytesPtr)
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

	errs := Copy(&dst, src)
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
			"struct1": {Name: "struct 1 value", Year: 2001},
			"struct2": {Name: "struct 2 value", Year: 2002},
			"struct3": {Name: "struct 3 value", Year: 2003},
		},
		MapInterfaces: map[string]interface{}{
			"inter1": 100001,
			"inter2": "This is my interface string",
			"inter3": SampleSubInfo{Name: "inter3: struct 1 value", Year: 2003},
			"inter4": float32(1.6546565),
			"inter5": float64(1.6546565),
			"inter6": &SampleSubInfo{Name: "inter6: struct 2 value", Year: 2006},
		},
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src)
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
		SampleSubInfo:          SampleSubInfo{Name: "This embedded struct", Year: 2016},
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 1 ptr struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
		CreatedTime:            time.Now(),
	}

	dst := SampleStruct{}

	errs := Copy(&dst, &src)
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
		SampleSubInfo:          SampleSubInfo{Name: "This embedded struct", Year: 2016},
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 1 ptr struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
		CreatedTime:            time.Now(),
	}

	dst := SampleStruct{
		Level1StructPtrZero: &SampleSubInfo{Name: "This level 1 struct ptr zero", Year: 2015},
	}

	errs := Copy(&dst, &src)
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
		SampleSubInfo:          SampleSubInfo{Name: "This embedded struct", Year: 2016},
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 1 ptr struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
	}

	errs := Copy(&dst, &src)
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

func TestCopyStructEmbededAndAttributeOmitEmpty(t *testing.T) {
	type SampleSubInfo struct {
		Name string
		Year int
	}

	type SampleStruct struct {
		Level1Struct           SampleSubInfo  `model:",omitempty,notraverse"`
		Level1StructPtr        *SampleSubInfo `model:",omitempty"`
		Level1StructNoTraverse *SampleSubInfo `model:",omitempty,notraverse"`
		CreatedTime            time.Time
		SampleSubInfo          `model:",omitempty"`
	}

	src := SampleStruct{CreatedTime: time.Now()}

	dst := SampleStruct{
		SampleSubInfo:          SampleSubInfo{Name: "This embedded struct", Year: 2016},
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 1 ptr struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
	}

	errs := Copy(&dst, src)
	if errs != nil {
		fmt.Println(errs)
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, true, src.CreatedTime == dst.CreatedTime)

	assertEqual(t, 2016, dst.Year)
	assertEqual(t, "This embedded struct", dst.Name)

	assertEqual(t, 2013, dst.Level1StructNoTraverse.Year)
	assertEqual(t, "This nested no traverse struct", dst.Level1StructNoTraverse.Name)

	assertEqual(t, 2015, dst.Level1Struct.Year)
	assertEqual(t, "This level 1 struct", dst.Level1Struct.Name)

	assertEqual(t, 2014, dst.Level1StructPtr.Year)
	assertEqual(t, "This level 1 ptr struct", dst.Level1StructPtr.Name)
}

func TestCopyDestinationIsNotPointer(t *testing.T) {
	type SampleStruct struct {
		Name string
	}
	errs := Copy(SampleStruct{}, SampleStruct{Name: "Not a pointer"})

	assertEqual(t, "Destination struct is not a pointer", errs[0].Error())
}

func TestCopyInputIsNotStruct(t *testing.T) {
	type SampleStruct struct {
		Name string
	}
	errs := Copy(&SampleStruct{}, map[string]string{"1": "2001"})

	assertEqual(t, "Source or Destination is not a struct", errs[0].Error())
}

func TestCopyStructElementKindDiff(t *testing.T) {
	type Source struct {
		Name string
	}

	type Destination struct {
		Name int
	}

	errs := Copy(&Destination{}, Source{Name: "This struct element kind is different"})

	assertEqual(t, "Field: 'Name', src [string] & dst [int] kind didn't match", errs[0].Error())
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

	errs := Copy(&dst, src)

	logSrcDst(t, src, dst)

	assertEqual(t, "Field: 'Name', src [string] & dst [int] kind didn't match", errs[0].Error())
	assertEqual(t,
		"Field: 'Level1', src [model.SampleLevelSrc] & dst [model.SampleLevelDst] type didn't match",
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

	errs := Copy(&dst, src)

	logSrcDst(t, src, dst)

	assertEqual(t, "Field: 'Name', src [string] & dst [int] kind didn't match", errs[0].Error())
	assertEqual(t, 0, dst.Name)
	assertEqual(t, src.Level1.Name, dst.Level1.(SampleLevelSrc).Name)
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

	errs := Copy(&dst, src)
	if errs != nil {
		t.Error("Error occurred while copying.")
	}

	logSrcDst(t, src, dst)

	assertEqual(t, "", dst.Name)
	assertEqual(t, 2016, dst.Year)
}

//
// NoTraverseTypeLis test cases
//

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

//
// Zero test cases
//

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
	errs := Copy(&SampleStruct{}, SampleStruct{})
	assertEqual(t, "Source struct is empty", errs[0].Error())

	errs = Copy(nil, nil)
	assertEqual(t, "Source or Destination is nil", errs[0].Error())
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
		t.Error("SampleStruct embedded struct - supposed to be non-zero")
	}
}

func TestNonHasZeroCheck(t *testing.T) {
	type SampleSubInfo struct {
		Name string
		Year int
	}

	type SampleStruct struct {
		Level1Struct           SampleSubInfo
		Level1StructPtr        *SampleSubInfo
		Level1StructNoTraverse *SampleSubInfo `model:",notraverse"`
		CreatedTime            time.Time
		SampleSubInfo
	}

	src1 := SampleStruct{
		SampleSubInfo:          SampleSubInfo{Name: "This embedded struct", Year: 2016},
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 1 ptr struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
		CreatedTime:            time.Now(),
	}

	if HasZero(src1) {
		t.Error("SampleStruct supposed to be non-zero")
	}

	src2 := SampleStruct{
		SampleSubInfo:          SampleSubInfo{Name: "This embedded struct", Year: 2016},
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 1 ptr struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
	}

	if !HasZero(src2) {
		t.Error("SampleStruct supposed to have one-zero i.e. CreatedTime field")
	}

	src3 := SampleStruct{
		Level1Struct:           SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleSubInfo{Name: "This level 1 ptr struct", Year: 2014},
		Level1StructNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2013},
		SampleSubInfo:          SampleSubInfo{Name: "This embedded struct"},
	}

	if !HasZero(src3) {
		t.Error("SampleStruct supposed to have one-zero i.e. SampleSubInfo -> Year field")
	}
}

func TestHasZeroForField(t *testing.T) {
	type SampleSubInfo struct {
		OmitThisField string `model:"-"`
		Name          string
		Year          int
	}

	if !HasZero(&SampleSubInfo{Name: "only I have populated"}) {
		t.Error("Supposed to have one-zero filed i.e. Year")
	}

	type SampleSrcStruct1 struct {
		Level1StructOmit       SampleSubInfo `model:"-"`
		Level1StructPtr        *SampleSubInfo
		Level1Struct           SampleSubInfo
		Level1StructNoTraverse *SampleSubInfo `model:",notraverse"`
		CreatedTime            time.Time
		SampleSubInfo
	}

	src1 := SampleSrcStruct1{}
	if !HasZero(src1) {
		t.Error("Suppose to be empty")
	}

	type SampleSrcStruct2 struct {
		Level1StructOmit       SampleSubInfo
		Level1StructPtr        *SampleSubInfo
		Level1Struct           SampleSubInfo
		Level1StructNoTraverse *SampleSubInfo `model:",notraverse"`
		CreatedTime            time.Time
		SampleSubInfo
	}

	src2 := SampleSrcStruct2{}
	if !HasZero(src2) {
		t.Error("Suppose to be empty")
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

func TestHasZeroNotAStructInput(t *testing.T) {
	result1 := HasZero(10001)
	assertEqual(t, false, result1)

	result2 := HasZero(map[string]int{"1": 101, "2": 102, "3": 103})
	assertEqual(t, false, result2)

	floatVar := float64(1.7367643)
	result3 := HasZero(&floatVar)
	assertEqual(t, false, result3)

	str := "This is not a struct"
	result4 := HasZero(&str)
	assertEqual(t, false, result4)

	assertEqual(t, true, HasZero(nil))
}

//
// Map test cases
//

func TestMapMethodValidation(t *testing.T) {
	_, err1 := Map(nil)
	assertEqual(t, "Invalid input <nil>", err1.Error())

	_, err2 := Map("not struct")
	assertEqual(t, "Input is not a struct", err2.Error())
}

func TestMapIntegerAndIntegerPtrWithDefaultKeyName(t *testing.T) {
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

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	// Assertion

	value1, found1 := result["Int"]
	assertEqual(t, true, found1)
	assertEqual(t, src.Int, value1)

	value2, found2 := result["Int64Ptr"]
	assertEqual(t, true, found2)
	assertEqual(t, *src.Int64Ptr, *value2.(*int64))
}

func TestMapIntegerAndIntegerPtrWithCustomKeyName(t *testing.T) {
	type SampleStruct struct {
		Int      int    `model:"int"`
		IntPtr   *int   `model:"32Pointer"`
		Int64    int64  `model:"int64"`
		Int64Ptr *int64 `model:"64Pointer"`
	}

	intPtr := int(1001)
	int64Ptr := int64(1002)

	src := SampleStruct{
		Int:      2001,
		IntPtr:   &intPtr,
		Int64:    2002,
		Int64Ptr: &int64Ptr,
	}

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	// Assertion

	value1, found1 := result["int"]
	assertEqual(t, true, found1)
	assertEqual(t, src.Int, value1)

	value2, found2 := result["64Pointer"]
	assertEqual(t, true, found2)
	assertEqual(t, *src.Int64Ptr, *value2.(*int64))
}

func TestMapStringAndStringPtr(t *testing.T) {
	type SampleStruct struct {
		String        string `model:"myStringKey"`
		StringPtr     *string
		StringZero    string `model:",omitempty"`
		StringPtrZero string `model:",omitempty"`
	}

	strPtr := "Map: This is string for pointer test"
	src := SampleStruct{
		String:    "Map: This is string for test",
		StringPtr: &strPtr,
	}

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	// Assertion

	value1, found1 := result["myStringKey"]
	assertEqual(t, src.String, value1)
	assertEqual(t, true, found1)

	vaule2, found2 := result["StringPtr"]
	assertEqual(t, *src.StringPtr, *vaule2.(*string))
	assertEqual(t, true, found2)

	_, notFound1 := result["StringZero"]
	assertEqual(t, false, notFound1)

	_, notFound2 := result["StringPtrZero"]
	assertEqual(t, false, notFound2)
}

func TestMapBooleanAndBooleanPtr(t *testing.T) {
	type SampleStruct struct {
		Boolean    bool
		BooleanPtr *bool
	}

	boolPtr := true
	src := SampleStruct{
		Boolean:    true,
		BooleanPtr: &boolPtr,
	}

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	// Assertion

	value1, found1 := result["Boolean"]
	assertEqual(t, true, found1)
	assertEqual(t, src.Boolean, value1)

	value2, found2 := result["BooleanPtr"]
	assertEqual(t, true, found2)
	assertEqual(t, *src.BooleanPtr, *value2.(*bool))
}

func TestMapByteAndByteSlice(t *testing.T) {
	type SampleStruct struct {
		Byte          byte
		SliceBytes    []byte
		SliceBytesPtr *[]byte
	}

	bytesPtr := []byte("This is byte pointer value")

	src := SampleStruct{
		Byte:          byte('A'),
		SliceBytes:    []byte("This is byte value"),
		SliceBytesPtr: &bytesPtr,
	}

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	value1, found1 := result["Byte"]
	assertEqual(t, true, found1)
	assertEqual(t, src.Byte, value1.(byte))

	value2, found2 := result["SliceBytes"]
	assertEqual(t, true, found2)
	assertEqual(t, src.SliceBytes, value2.([]byte))
	assertEqual(t, string(src.SliceBytes), string(value2.([]byte)))

	value3, found3 := result["SliceBytesPtr"]
	assertEqual(t, true, found3)
	assertEqual(t, true, reflect.DeepEqual(src.SliceBytesPtr, value3.(*[]byte)))
}

func TestMapSliceStringAndSliceStringPtr(t *testing.T) {
	type SampleStruct struct {
		SliceString    []string
		SliceStringPtr *[]string
	}

	sliceStrPtr := []string{
		"Val1: This is slice string test pointer.",
		"Val2: This is slice string test pointer.",
	}

	src := SampleStruct{
		SliceString: []string{
			"Val1: This is slice string test.",
			"Val2: This is slice string test.",
		},
		SliceStringPtr: &sliceStrPtr,
	}

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	// Assertion

	value1, found1 := result["SliceString"]
	assertEqual(t, true, found1)
	assertEqual(t, src.SliceString, value1)

	value2, found2 := result["SliceStringPtr"]
	assertEqual(t, true, found2)
	assertEqual(t, *src.SliceStringPtr, *value2.(*[]string))
}

func TestMapSliceElementsPtr(t *testing.T) {
	type SampleSubInfo2 struct {
		SliceIntPtr      []*int
		SliceInt64Ptr    []*int64
		SliceStringPtr   []*string  `model:"stringPtr"`
		SliceFloat32Omit []*float32 `model:"-"`
		SliceFloat32     []*float32
		SliceFloat64     []*float64
		SliceInterface   []interface{} `model:"interface"`
	}

	type SampleSubInfo1 struct {
		SliceIntPtr    []*int
		SliceInt64Ptr  []*int64 `model:"int64Ptr"`
		SliceStringPtr []*string
		SliceFloat32   []*float32
		SliceFloat64   []*float64     `model:"float64Ptr"`
		SliceInterface []interface{}  `model:"interface"`
		Level2         SampleSubInfo2 `model:"level2"`
	}

	type SampleStruct struct {
		SliceIntPtr    []*int `model:"intPtr"`
		SliceInt64Ptr  []*int64
		SliceStringPtr []*string  `model:"stringPtr"`
		SliceFloat32   []*float32 `model:"float32"`
		SliceFloat64   []*float64
		SliceInterface []interface{}  `model:"interface"`
		Level1         SampleSubInfo1 `model:"level1"`
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

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	// Assertion

	// Level 0 assertion
	value1, _ := result["intPtr"]
	assertEqual(t, src.SliceIntPtr, value1.([]*int))

	value2, _ := result["SliceInt64Ptr"]
	assertEqual(t, src.SliceInt64Ptr, value2.([]*int64))

	value3, _ := result["stringPtr"]
	assertEqual(t, src.SliceStringPtr, value3.([]*string))

	value4, _ := result["float32"]
	assertEqual(t, src.SliceFloat32, value4.([]*float32))

	value5, _ := result["SliceFloat64"]
	assertEqual(t, src.SliceFloat64, value5.([]*float64))

	value6, _ := result["interface"]
	assertEqual(t, src.SliceInterface, value6.([]interface{}))

	// // Level 1 assertion
	l1, _ := result["level1"].(map[string]interface{})

	l1value1, _ := l1["SliceIntPtr"]
	assertEqual(t, src.Level1.SliceIntPtr, l1value1.([]*int))

	l1value2, _ := l1["int64Ptr"]
	assertEqual(t, src.Level1.SliceInt64Ptr, l1value2.([]*int64))

	l1value3, _ := l1["SliceStringPtr"]
	assertEqual(t, src.Level1.SliceStringPtr, l1value3.([]*string))

	l1value4, _ := l1["SliceFloat32"]
	assertEqual(t, src.Level1.SliceFloat32, l1value4.([]*float32))

	l1value5, _ := l1["float64Ptr"]
	assertEqual(t, src.Level1.SliceFloat64, l1value5.([]*float64))

	l1value6, _ := l1["interface"]
	assertEqual(t, src.Level1.SliceInterface, l1value6.([]interface{}))

	// // Level 2 assertion
	l2, _ := l1["level2"].(map[string]interface{})

	l2value1, _ := l2["SliceIntPtr"]
	assertEqual(t, src.Level1.Level2.SliceIntPtr, l2value1.([]*int))

	l2value2, _ := l2["SliceInt64Ptr"]
	assertEqual(t, src.Level1.Level2.SliceInt64Ptr, l2value2.([]*int64))

	l2value3, _ := l2["stringPtr"]
	assertEqual(t, src.Level1.Level2.SliceStringPtr, l2value3.([]*string))

	l2value4, _ := l2["SliceFloat32"]
	assertEqual(t, src.Level1.Level2.SliceFloat32, l2value4.([]*float32))

	l2value5, _ := l2["SliceFloat64"]
	assertEqual(t, src.Level1.Level2.SliceFloat64, l2value5.([]*float64))

	l2value6, _ := l2["interface"]
	assertEqual(t, src.Level1.Level2.SliceInterface, l2value6.([]interface{}))
}

func TestMapMapElements(t *testing.T) {
	type SampleSubInfo struct {
		Name string
		Year int
	}

	type SampleStruct struct {
		MapIntInt       map[int]int
		MapStringInt    map[string]int `model:"stringInt"`
		MapStringString map[string]string
		MapStruct       map[string]SampleSubInfo
		MapInterfaces   map[string]interface{}
	}

	src := SampleStruct{
		MapIntInt:       map[int]int{1: 1001, 2: 1002, 3: 1003, 4: 1004},
		MapStringInt:    map[string]int{"first": 1001, "second": 1002, "third": 1003, "forth": 1004},
		MapStringString: map[string]string{"first": "1001", "second": "1002", "third": "1003"},
		MapStruct: map[string]SampleSubInfo{
			"struct1": {Name: "struct 1 value", Year: 2001},
			"struct2": {Name: "struct 2 value", Year: 2002},
			"struct3": {Name: "struct 3 value", Year: 2003},
		},
		MapInterfaces: map[string]interface{}{
			"inter1": 100001,
			"inter2": "This is my interface string",
			"inter3": SampleSubInfo{Name: "inter3: struct 1 value", Year: 2003},
			"inter4": float32(1.6546565),
			"inter5": float64(1.6546565),
			"inter6": &SampleSubInfo{Name: "inter6: struct 2 value", Year: 2006},
			"l1map1": map[int]int{1: 1001, 2: 1002, 3: 1003, 4: 1004},
			"l1map2": map[string]int{"first": 1001, "second": 1002, "third": 1003, "forth": 1004},
			"l2map1": map[string]interface{}{
				"struct1": SampleSubInfo{Name: "l2map1: struct 1 value", Year: 2001},
				"struct2": SampleSubInfo{Name: "l2map1: struct 2 value", Year: 2002},
				"struct3": SampleSubInfo{Name: "l2map1: struct 3 value", Year: 2003},
				"l3map1":  map[string]string{"first": "1001", "second": "1002", "third": "1003"},
			},
		},
	}

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	// Assertion

	// Field: MapIntInt
	value1, found1 := result["MapIntInt"].(map[string]interface{})
	assertEqual(t, true, found1)

	value11, found11 := value1["2"]
	assertEqual(t, true, found11)
	assertEqual(t, 1002, value11)

	// Field: MapStringString
	value2, found2 := result["MapStringString"].(map[string]interface{})
	assertEqual(t, true, found2)

	value21, found21 := value2["third"]
	assertEqual(t, true, found21)
	assertEqual(t, "1003", value21)

	// Field: MapStruct -> struct2 -> Name
	value3, found3 := result["MapStruct"].(map[string]interface{})
	assertEqual(t, true, found3)

	value31, found31 := value3["struct2"].(map[string]interface{})
	assertEqual(t, true, found31)

	value32, found32 := value31["Name"]
	assertEqual(t, true, found32)
	assertEqual(t, "struct 2 value", value32)

	// Field: MapInterfaces -> inter4
	value4, found4 := result["MapInterfaces"].(map[string]interface{})
	assertEqual(t, true, found4)

	value41, found41 := value4["inter4"]
	assertEqual(t, true, found41)
	assertEqual(t, float32(1.6546565), value41.(float32))

	// Field: MapInterfaces -> inter6 -> Name
	value42, found42 := value4["inter6"].(*map[string]interface{})
	assertEqual(t, true, found42)
	assertEqual(t, true, value42 != nil)

	// Field: MapInterfaces -> l1map1 -> 4
	value43, found43 := value4["l1map1"].(map[string]interface{})
	assertEqual(t, true, found43)

	value431, found431 := value43["4"]
	assertEqual(t, true, found431)
	assertEqual(t, 1004, value431.(int))

	// Field: MapInterfaces -> l2map1 -> struct3 -> Name
	value44, found44 := value4["l2map1"].(map[string]interface{})
	assertEqual(t, true, found44)

	value441, found441 := value44["struct3"].(map[string]interface{})
	assertEqual(t, true, found441)

	value4411, found4411 := value441["Name"]
	assertEqual(t, true, found4411)
	assertEqual(t, "l2map1: struct 3 value", value4411.(string))

	// Field: MapInterfaces -> l2map1 -> l3map1 -> first
	value442, found442 := value44["l3map1"].(map[string]interface{})
	assertEqual(t, true, found442)

	value4421, found4421 := value442["first"]
	assertEqual(t, true, found4421)
	assertEqual(t, "1001", value4421.(string))
}

func TestMapStructEmbededAndAttribute(t *testing.T) {
	type SampleSubInfo struct {
		Name string
		Year int `model:"year"`
		Goal string
	}

	type SampleStruct struct {
		Level1Struct              SampleSubInfo `model:"level1Struct"`
		Level1StructNoTraverse    SampleSubInfo `model:",notraverse"`
		Level1StructPtr           *SampleSubInfo
		Level1StructPtrNoTraverse *SampleSubInfo `model:",notraverse"`
		Level1StructEmpty         SampleSubInfo  `model:",omitempty"`
		Level1StructPtrEmpty      *SampleSubInfo `model:",omitempty"`
		CreatedTime               time.Time      `model:"created_time"`
		CreatedTimePtr            *time.Time     `model:"created_time_ptr"`
		UpdateTimeOmitEmpty       time.Time      `model:"update_time,omitempty"`
		SampleSubInfo
	}

	timePtr := time.Now()
	src := SampleStruct{
		SampleSubInfo:             SampleSubInfo{Name: "This embedded struct", Year: 2016},
		Level1Struct:              SampleSubInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructNoTraverse:    SampleSubInfo{Name: "This level 1 struct no traverse", Year: 2014},
		Level1StructPtr:           &SampleSubInfo{Name: "This level 1 struct pointer", Year: 2013},
		Level1StructPtrNoTraverse: &SampleSubInfo{Name: "This nested no traverse struct", Year: 2012},
		CreatedTime:               time.Now(),
		CreatedTimePtr:            &timePtr,
	}

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	// Assertion

	// Embedded struct assertion
	// Field: Name
	value1, found1 := result["Name"]
	assertEqual(t, true, found1)
	assertEqual(t, "This embedded struct", value1.(string))

	// Field: year
	value2, found2 := result["year"]
	assertEqual(t, true, found2)
	assertEqual(t, 2016, value2.(int))

	// Field: level1Struct -> Name
	value3, found3 := result["level1Struct"].(map[string]interface{})
	assertEqual(t, true, found3)

	value31, found31 := value3["Name"]
	assertEqual(t, true, found31)
	assertEqual(t, "This level 1 struct", value31.(string))

	// Field: level1Struct -> Goal (should be empty)
	value32, found32 := value3["Goal"]
	assertEqual(t, true, found32)
	assertEqual(t, "", value32.(string))

	// Field: created_time
	value4, found4 := result["created_time"]
	assertEqual(t, true, found4)
	assertEqual(t, true, src.CreatedTime == value4.(time.Time))

	value5, found5 := result["created_time_ptr"]
	assertEqual(t, true, found5)
	assertEqual(t, true, src.CreatedTimePtr != value5.(*time.Time))

	// Field should not exists: Level1StructEmpty, Level1StructPtrEmpty, UpdateTimeOmitEmpty
	_, notfound1 := result["Level1StructEmpty"]
	assertEqual(t, false, notfound1)

	_, notfound2 := result["Level1StructPtrEmpty"]
	assertEqual(t, false, notfound2)

	_, notfound3 := result["UpdateTimeOmitEmpty"]
	assertEqual(t, false, notfound3)
}

func TestMapSliceStructAndSliceStructPtr(t *testing.T) {
	type SampleSubInfo struct {
		Name string
		Year int `model:"year"`
		Goal string
	}
	type SampleStruct struct {
		SliceStruct    []SampleSubInfo
		SliceStructPtr *[]SampleSubInfo
	}

	sliceStructPtr := []SampleSubInfo{
		{Name: "Struct: Slice Ptr 1", Year: 2016},
		{Name: "Struct: Slice Ptr 2", Year: 2015},
		{Name: "Struct: Slice Ptr 3", Year: 2014},
	}
	src := SampleStruct{
		SliceStruct: []SampleSubInfo{
			{Name: "Struct: Slice 1", Year: 2006},
			{Name: "Struct: Slice 2", Year: 2005},
			{Name: "Struct: Slice 3", Year: 2004},
		},
		SliceStructPtr: &sliceStructPtr,
	}

	result, err := Map(src)
	if err != nil {
		t.Error("Error occurred while Map export.")
	}

	logSrcDst(t, src, result)

	value1 := result["SliceStruct"].([]interface{})[0].(map[string]interface{})

	value11, found11 := value1["Name"]
	assertEqual(t, true, found11)
	assertEqual(t, src.SliceStruct[0].Name, value11.(string))

	value12, found12 := value1["year"]
	assertEqual(t, true, found12)
	assertEqual(t, src.SliceStruct[0].Year, value12.(int))

	value13, found13 := value1["Goal"]
	assertEqual(t, true, found13)
	assertEqual(t, src.SliceStruct[0].Goal, value13.(string))
}

func TestCloneInputNil(t *testing.T) {
	result, err := Clone(nil)

	assertEqual(t, "Invalid input <nil>", err.Error())
	assertEqual(t, true, result == nil)
}

func TestCloneNotAStruct(t *testing.T) {
	result, err := Clone("I'm not a struct")

	assertEqual(t, "Input is not a struct", err.Error())
	assertEqual(t, true, result == nil)
}

func TestCloneStruct(t *testing.T) {
	type SampleInfo struct {
		Name string
		Year int `model:"year"`
		Goal string
	}

	src := SampleInfo{Name: "My name is go-model", Year: 2016}

	result, err := Clone(src)

	assertEqual(t, true, result != nil)
	assertEqual(t, true, err == nil)
	assertEqual(t, src.Name, result.(*SampleInfo).Name)
	assertEqual(t, src.Year, result.(*SampleInfo).Year)
}

func TestCloneStructPtr(t *testing.T) {
	type SampleInfo struct {
		Name string
		Year int `model:"year"`
		Goal string
	}

	src := SampleInfo{Name: "My name is go-model ptr", Year: 2015}

	result, err := Clone(&src)

	assertEqual(t, true, result != nil)
	assertEqual(t, true, err == nil)
	assertEqual(t, src.Name, result.(*SampleInfo).Name)
	assertEqual(t, src.Year, result.(*SampleInfo).Year)
}

//
// IsZeroInFields test case
//

func TestIsZeroInFields(t *testing.T) {
	type SampleInfo struct {
		Name   string
		Year   int
		Level2 float32
	}

	type SampleStruct struct {
		Name                   string
		Year                   int `model:"year"`
		Goal                   string
		Level1Struct           SampleInfo
		Level1StructPtr        *SampleInfo
		Level1StructNoTraverse *SampleInfo `model:",notraverse"`
		CreatedTime            time.Time
	}

	_, basic1 := IsZeroInFields(nil, "TestField")
	assertEqual(t, true, basic1)

	_, basic2 := IsZeroInFields(SampleStruct{})
	assertEqual(t, true, basic2)

	_, basic3 := IsZeroInFields("I'm not a struct", "TestField")
	assertEqual(t, false, basic3)

	src1 := SampleStruct{
		Name:                   "I'm Name",
		Year:                   2016,
		Goal:                   "To test IsZeroInFields",
		Level1Struct:           SampleInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleInfo{Name: "This level 1 ptr struct", Year: 2014},
		Level1StructNoTraverse: &SampleInfo{Name: "This nested no traverse struct", Year: 2013},
		CreatedTime:            time.Now(),
	}

	name1, zero1 := IsZeroInFields(src1, "Name", "Year", "Level1StructNoTraverse")
	assertEqual(t, false, zero1)
	assertEqual(t, "", name1)

	src2 := SampleStruct{
		Name:                   "I'm Name",
		Year:                   2016,
		Goal:                   "To test IsZeroInFields",
		Level1Struct:           SampleInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructNoTraverse: &SampleInfo{Name: "This nested no traverse struct", Year: 2013},
		CreatedTime:            time.Now(),
	}
	name2, zero2 := IsZeroInFields(src2, "Goal", "Level1Struct", "Level1StructNoTraverse", "Level1StructPtr")
	assertEqual(t, true, zero2)
	assertEqual(t, "Level1StructPtr", name2)

	src3 := SampleStruct{
		Name:                   "I'm Name",
		Year:                   2016,
		Goal:                   "To test IsZeroInFields",
		Level1Struct:           SampleInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr:        &SampleInfo{Name: "This level 1 ptr struct", Year: 2014},
		Level1StructNoTraverse: &SampleInfo{Name: "This nested no traverse struct"},
		CreatedTime:            time.Now(),
	}
	name3, zero3 := IsZeroInFields(src3, "Year1")
	assertEqual(t, false, zero3)
	assertEqual(t, "", name3)
}

func TestIsZeroInFieldsEmbedded(t *testing.T) {
	type SampleInfo struct {
		Name   string
		Year   int
		Level2 float32
	}

	type SampleStruct struct {
		Level1Struct           SampleInfo
		Level1StructPtr        *SampleInfo
		Level1StructNoTraverse SampleInfo `model:",notraverse"`
		CreatedTime            time.Time
		SampleInfo
	}

	src1 := SampleStruct{
		Level1Struct:    SampleInfo{Name: "This level 1 struct", Year: 2015},
		Level1StructPtr: &SampleInfo{Name: "This level 1 ptr struct", Year: 2014},
		CreatedTime:     time.Now(),
		SampleInfo: SampleInfo{
			Name: "I'm Name",
			Year: 2016,
		},
	}

	name1, zero1 := IsZeroInFields(src1, "SampleInfo")
	assertEqual(t, false, zero1)
	assertEqual(t, "", name1)

	name2, zero2 := IsZeroInFields(SampleStruct{}, "SampleInfo")
	assertEqual(t, true, zero2)
	assertEqual(t, "SampleInfo", name2)
}

func TestFields(t *testing.T) {
	type SampleInfo struct {
		Name   string
		Year   int
		Level2 float32
	}

	type SampleStruct struct {
		Level1Struct           SampleInfo
		Level1StructPtr        *SampleInfo
		Level1StructNoTraverse SampleInfo `model:",notraverse"`
		CreatedTime            time.Time
		SampleInfo
	}

	fields1, err1 := Fields(nil)
	assertEqual(t, true, err1.Error() == "Invalid input <nil>")
	assertEqual(t, true, fields1 == nil)

	fields2, err2 := Fields(&SampleStruct{})

	assertError(t, err2)
	assertEqual(t, true, len(fields2) > 0)
}

func TestKind(t *testing.T) {
	type SampleInfo struct {
		MapIntInt       map[int]int
		MapStringInt    map[string]int `model:"stringInt"`
		MapStringString map[string]string
	}

	type SampleStruct struct {
		Name                   string      `json:"name,omitempty"`
		Year                   int         `json:"year"`
		Level2                 float32     `json:"level2"`
		Struct                 SampleInfo  `json:"struct"`
		StructPtr              *SampleInfo `json:"struct_ptr"`
		Level1StructNoTraverse SampleInfo  `model:",notraverse"`
		CreatedTime            time.Time   `json:"created_time,omitempty"`
	}

	s := SampleStruct{}

	kind1, err1 := Kind(s, "Name")
	assertError(t, err1)
	assertEqual(t, true, reflect.String == kind1)

	kind2, err2 := Kind(s, "StructPtr")
	assertError(t, err2)
	assertEqual(t, true, reflect.Ptr == kind2)

	kind3, err3 := Kind(s, "CreatedTime")
	assertError(t, err3)
	assertEqual(t, true, reflect.Struct == kind3)

	kind4, err4 := Kind(s, "Level2")
	assertError(t, err4)
	assertEqual(t, true, reflect.Float32 == kind4)

	kind5, err5 := Kind(nil, "NoExists")
	assertEqual(t, "Invalid input <nil>", err5.Error())
	assertEqual(t, true, reflect.Invalid == kind5)

	kind6, err6 := Kind(s, "NoExists")
	assertEqual(t, "Field: 'NoExists', does not exists", err6.Error())
	assertEqual(t, true, reflect.Invalid == kind6)
}

func TestNestedStructToStructMapping(t *testing.T) {
	type C struct {
		X string
	}

	type A struct {
		V C
	}

	type B struct {
		V C
	}

	a := A{V: C{"1"}}
	b := B{}

	Copy(&b, &a)

	assertEqual(t, a.V.X, b.V.X)
}

func TestStructToStructPtrWithConverter(t *testing.T) {
	type C struct {
		X string
	}

	type D struct {
		X *string
	}

	type A struct {
		V C
	}

	type B struct {
		V D
	}

	a := A{V: C{"1"}}
	b := B{}

	AddConversion(&C{}, &D{}, func(in reflect.Value) (reflect.Value, error) {
		x := in.Interface().(C).X
		d := D{X: &x}
		return reflect.ValueOf(d), nil
	})

	Copy(&b, &a)

	assertEqual(t, a.V.X, *b.V.X)
}

func TestStructWithConverter(t *testing.T) {

	type C struct {
		X string
	}

	type D struct {
		X string
	}

	type A struct {
		V C
	}

	type B struct {
		V D
	}

	a := A{V: C{"1"}}
	b := B{}

	AddConversion(&C{}, &D{}, func(in reflect.Value) (reflect.Value, error) {
		x := in.Interface().(C).X
		d := D{X: x}
		return reflect.ValueOf(d), nil
	})

	Copy(&b, &a)

	assertEqual(t, a.V.X, b.V.X)
}

func TestSliceWithConverter(t *testing.T) {

	type C struct {
		X string
	}

	type D struct {
		X string
	}

	type A struct {
		V []C
	}

	type B struct {
		V []D
	}

	a := A{V: []C{{"1"}, {"2"}}}
	b := B{}

	AddConversion(&C{}, &D{}, func(in reflect.Value) (reflect.Value, error) {
		x := in.Interface().(C).X
		d := D{X: x}
		return reflect.ValueOf(d), nil
	})

	Copy(&b, &a)

	assertEqual(t, a.V[0].X, b.V[0].X)
	assertEqual(t, a.V[1].X, b.V[1].X)

}

func TestMapWithConverter(t *testing.T) {
	type C struct {
		X string
	}

	type D struct {
		X string
	}

	type A struct {
		M map[string]C
	}

	type B struct {
		M map[string]D
	}

	a := A{M: map[string]C{"1": {"1"}, "2": {"2"}, "3": {"error"}}}
	b := B{}

	AddConversion(&C{}, &D{}, func(in reflect.Value) (reflect.Value, error) {
		x := in.Interface().(C).X
		d := D{X: x}
		if x == "error" {
			return reflect.ValueOf(d), errors.New("Custom conversion failed.")
		}
		return reflect.ValueOf(d), nil
	})

	errs := Copy(&b, &a)
	assertEqual(t, a.M["1"].X, b.M["1"].X)
	assertEqual(t, a.M["2"].X, b.M["2"].X)
	assertEqual(t, "Custom conversion failed.", errs[0].Error())
}

func TestGetField(t *testing.T) {
	type SampleStruct struct {
		Int    int
		String string
	}

	src := SampleStruct{
		Int:    10,
		String: "go-model",
	}

	// scenario 1 int
	value1, err1 := Get(src, "Int")
	assertEqual(t, 10, value1)
	assertError(t, err1)

	// scenario 2 string
	value2, err2 := Get(src, "String")
	assertEqual(t, "go-model", value2)
	assertError(t, err2)

	// scenario 3 field not exists
	_, err := Get(src, "NotExists")
	assertEqual(t, "Field: 'NotExists', does not exists", err.Error())

	// scenario 4 struct is nil
	_, err = Get(nil, "Int")
	assertEqual(t, "Invalid input <nil>", err.Error())
}

func TestSetField(t *testing.T) {
	type SampleStruct struct {
		Int    int
		String string
	}

	src := SampleStruct{
		Int:    10,
		String: "go-model",
	}

	// scenario 1 direct value int
	err := Set(&src, "Int", 20)
	assertError(t, err)

	value1, err1 := Get(src, "Int")
	assertEqual(t, 20, value1)
	assertError(t, err1)

	// scenario 2 direct value string
	err = Set(&src, "String", "go-model set")
	assertError(t, err)

	value2, err2 := Get(src, "String")
	assertEqual(t, "go-model set", value2)
	assertError(t, err2)

	// scenario 3 value is pointer
	newVal := "go-model set ptr"
	err = Set(&src, "String", &newVal)
	assertError(t, err)

	value3, err3 := Get(src, "String")
	assertEqual(t, "go-model set ptr", value3)
	assertError(t, err3)

	// scenario 4 struct is not pointer
	err = Set(src, "Int", 10)
	assertEqual(t, "Destination struct is not a pointer", err.Error())

	// scenario 5 field not exists
	err = Set(&src, "NotExists", "test value")
	assertEqual(t, "Field: 'NotExists', does not exists", err.Error())

	// scenario 6 struct is nil
	err = Set(nil, "Int", 30)
	assertEqual(t, "Invalid input <nil>", err.Error())

	// scenario 7 different type
	err = Set(&src, "String", 30)
	assertEqual(t, "Field: String, type/kind did not match", err.Error())
}

func TestImprovedCopy(t *testing.T) {
	type DomainObject struct {
		Name    string
		Address string
		Phone   string
	}

	type LoginGreeterDTO struct {
		Name string
	}

	src := DomainObject{
		Name:    "go-model",
		Address: "123 sample street",
		Phone:   "000-000-0000",
	}

	dst := LoginGreeterDTO{}

	errs := Copy(&dst, src)
	assertEqual(t, 0, len(errs))
	assertEqual(t, "go-model", dst.Name)
}

//
// helper test methods
//

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
	t.Log("")
	logIt(t, "Destination", dst)
}

func logIt(t *testing.T, str string, v interface{}) {
	t.Logf("%v: %#v", str, v)
}
