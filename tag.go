// Copyright (c) 2016 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "strings"

type tag struct {
	Name    string
	Options string
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
