// Copyright (c) Jeevanandam M. (https://github.com/jeevatkm).
// go-model source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"testing"
	"time"
)

func TestTag(t *testing.T) {
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
		CreatedTime            *time.Time  `json:"created_time,omitempty"`
	}

	s := SampleStruct{}

	tag1, err1 := Tag(s, "StructPtr")
	assertError(t, err1)
	assertEqual(t, "struct_ptr", tag1.Get("json"))

	tag2, err2 := Tag(s, "CreatedTime")
	assertError(t, err2)
	assertEqual(t, "created_time,omitempty", tag2.Get("json"))

	tag3, err3 := Tag(s, "Level1StructNoTraverse")
	assertError(t, err3)
	assertEqual(t, "", tag3.Get("json"))

	_, err4 := Tag(nil, "Level1StructNoTraverse")
	assertEqual(t, "Invalid input <nil>", err4.Error())

	_, err5 := Tag(s, "NotExists")
	assertEqual(t, "Field: 'NotExists', does not exists", err5.Error())
}

func TestTags(t *testing.T) {
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
		CreatedTime            *time.Time  `json:"created_time,omitempty"`
	}

	s := SampleStruct{}

	tags, err1 := Tags(s)
	assertError(t, err1)
	assertEqual(t, "struct_ptr", tags["StructPtr"].Get("json"))
	assertEqual(t, "created_time,omitempty", tags["CreatedTime"].Get("json"))

	_, err2 := Tags(nil)
	assertEqual(t, "Invalid input <nil>", err2.Error())
}

func TestNewTag(t *testing.T) {
	tag := newTag("fieldName,omitempty,notraverse")

	logIt(t, "Model Tag", tag)

	assertEqual(t, "fieldName", tag.Name)
	assertEqual(t, "omitempty,notraverse", tag.Options)
}

func TestNewTagWithoutName(t *testing.T) {
	tag := newTag(",omitempty,notraverse")

	logIt(t, "Model Tag", tag)

	assertEqual(t, "", tag.Name)
	assertEqual(t, "omitempty,notraverse", tag.Options)
}

func TestNewTagNoValues(t *testing.T) {
	tag := newTag("")

	logIt(t, "Model Tag", tag)

	assertEqual(t, "", tag.Name)
	assertEqual(t, "", tag.Options)
}

func TestNewTagEmptyValues(t *testing.T) {
	tag := newTag(",")

	logIt(t, "Model Tag", tag)

	assertEqual(t, "", tag.Name)
	assertEqual(t, "", tag.Options)
}

func TestNewTagSkipField(t *testing.T) {
	tag := newTag("-")

	logIt(t, "Model Tag", tag)

	assertEqual(t, "-", tag.Name)
	assertEqual(t, "", tag.Options)
	assertEqual(t, true, tag.isOmitField())
}

func TestIsOmitEmpty(t *testing.T) {
	tag1 := newTag("fieldName,omitempty,notraverse")
	logIt(t, "Model Tag", tag1)
	assertEqual(t, true, tag1.isOmitEmpty())

	tag2 := newTag(",omitempty")
	logIt(t, "Model Tag", tag2)
	assertEqual(t, true, tag2.isOmitEmpty())

	tag3 := newTag(",omitempty,notraverse")
	logIt(t, "Model Tag", tag3)
	assertEqual(t, true, tag3.isOmitEmpty())

	tag4 := newTag(",notraverse")
	logIt(t, "Model Tag", tag4)
	assertEqual(t, false, tag4.isOmitEmpty())

	tag5 := newTag("fieldName")
	logIt(t, "Model Tag", tag5)
	assertEqual(t, false, tag5.isOmitEmpty())
}

func TestIsNoTraverse(t *testing.T) {
	tag1 := newTag("fieldName,omitempty,notraverse")
	logIt(t, "Model Tag", tag1)
	assertEqual(t, true, tag1.isNoTraverse())

	tag2 := newTag(",notraverse")
	logIt(t, "Model Tag", tag2)
	assertEqual(t, true, tag2.isNoTraverse())

	tag3 := newTag(",omitempty,notraverse")
	logIt(t, "Model Tag", tag3)
	assertEqual(t, true, tag3.isNoTraverse())

	tag4 := newTag(",omitempty")
	logIt(t, "Model Tag", tag4)
	assertEqual(t, false, tag4.isNoTraverse())

	tag5 := newTag("fieldName")
	logIt(t, "Model Tag", tag5)
	assertEqual(t, false, tag5.isNoTraverse())
}
