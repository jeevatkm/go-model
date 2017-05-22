// Copyright (c) Jeevanandam M. (https://github.com/jeevatkm).
// go-model source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"reflect"
	"strings"
)

type tag struct {
	Name    string
	Options string
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

	if fv, ok := sv.Type().FieldByName(name); ok {
		return fv.Tag, nil
	}

	return "", fmt.Errorf("Field: '%v', does not exists", name)
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

func newTag(modelTag string) *tag {
	t := tag{}
	values := strings.Split(modelTag, ",")

	t.Name = values[0]
	t.Options = strings.Join(values[1:], ",")

	return &t
}

func (t *tag) isOmitField() bool {
	return t.Name == OmitField
}

func (t *tag) isOmitEmpty() bool {
	return t.isExists(OmitEmpty)
}

func (t *tag) isNoTraverse() bool {
	return t.isExists(NoTraverse)
}

func (t *tag) isExists(opt string) bool {
	return strings.Contains(t.Options, opt)
}

func isStringEmpty(str string) bool {
	return (len(strings.TrimSpace(str)) == 0)
}
