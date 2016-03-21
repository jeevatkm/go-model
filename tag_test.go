package model

import "testing"

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
