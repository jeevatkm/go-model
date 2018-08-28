# go-model [![Stability: Sustained](https://masterminds.github.io/stability/sustained.svg)](https://masterminds.github.io/stability/sustained.html) [![Build Status](https://travis-ci.org/jeevatkm/go-model.svg?branch=master)](https://travis-ci.org/jeevatkm/go-model) [![codecov](https://codecov.io/gh/jeevatkm/go-model/branch/master/graph/badge.svg)](https://codecov.io/gh/jeevatkm/go-model/branch/master)  [![GoReport](https://goreportcard.com/badge/jeevatkm/go-model)](https://goreportcard.com/report/jeevatkm/go-model) [![Version](https://img.shields.io/badge/version-1.1.0-blue.svg)](https://github.com/jeevatkm/go-model/releases/latest)  [![GoDoc](https://godoc.org/github.com/jeevatkm/go-model?status.svg)](https://godoc.org/github.com/jeevatkm/go-model) [![License](https://img.shields.io/github/license/jeevatkm/go-model.svg)](LICENSE)

Robust & Easy to use model mapper and utility methods for Go `struct`. Typical methods increase productivity and make Go development more fun :smile:

***v1.1.0 [released](https://github.com/jeevatkm/go-model/releases/latest) and tagged on Aug 27, 2018***

go-model tested with Go `v1.2` and above.

## Features
go-model library provides [handy methods](#supported-methods) to process `struct` with below highlighted features. It's born from typical need while developing Go application or utility. I hope it's helpful to Go community!
* Embedded/Anonymous struct
* Multi-level nested struct/map/slice
* Pointer and non-pointer within struct/map/slice
* Struct within map and slice
* Embedded/Anonymous struct fields appear in map at same level as represented by Go
* Interface within struct/map/slice
* Get struct field `reflect.Kind` by field name
* Get all the struct field tags (`reflect.StructTag`) or selectively by field name
* Get all `reflect.StructField` for given struct instance
* Get or Set by individual field name on struct
* Add global no traverse type to the list or use `notraverse` option in the struct field
* Options to name map key, omit empty fields, and instruct not to traverse with struct/map/slice
* Conversions between mixed non-pointer types - add custom conversation method, refer to usage

## Installation

#### Stable Version - Production Ready
Please refer section [Versioning](#versioning) for detailed info.

**go.mod**
```sh
require gopkg.in/jeevatkm/go-model.v1 v1.1.0
```

**go get**
```sh
go get -u gopkg.in/jeevatkm/go-model.v1
```

#### Latest Version - Development Edge
```sh
# install the latest & greatest library
go get -u github.com/jeevatkm/go-model
```

## It might be beneficial for your project :smile:

go-model author also published following projects to Go Community.

  * [aah framework](https://aahframework.org) - A secure, flexible, rapid Go web framework.
  * [go-resty](https://github.com/go-resty/resty) - Simple HTTP and REST client for Go.


## Usage
Import go-model into your code and refer it as `model`. Have a look on [model test cases](model_test.go) to know more possibilities.
```go
import (
  "gopkg.in/jeevatkm/go-model.v1"
)
```

### Supported Methods
* Copy - [usage](#copy-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Copy)
* Map - [usage](#map-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Map)
* Clone - [usage](#clone-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Clone)
* IsZero - [usage](#iszero-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#IsZero)
* HasZero - [usage](#haszero-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#HasZero)
* IsZeroInFields - [usage](#iszeroinfields-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#IsZeroInFields)
* Fields - [usage](#fields-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Fields)
* Kind - [usage](#kind-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Kind)
* Tag - [usage](#tag-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Tag)
* Tags - [usage](#tags-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Tags)
* Get - [usage](#get-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Get)
* Set - [usage](#set-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Set)
* AddNoTraverseType - [usage](#addnotraversetype--removenotraversetype-methods), [godoc](https://godoc.org/github.com/jeevatkm/go-model#AddNoTraverseType)
* RemoveNoTraverseType - [usage](#addnotraversetype--removenotraversetype-methods), [godoc](https://godoc.org/github.com/jeevatkm/go-model#RemoveNoTraverseType)
* AddConversion - [usage](#addconversion--removeconversion-methods), [godoc](https://godoc.org/github.com/jeevatkm/go-model#AddConversion)
* RemoveConversion - [usage](#addconversion--removeconversion-methods), [godoc](https://godoc.org/github.com/jeevatkm/go-model#RemoveConversion)

#### Copy Method
How do I copy my struct object into another? Not to worry, go-model does deep copy.
```go
// let's say you have just decoded/unmarshalled your request body to struct object.
tempProduct, _ := myapp.ParseJSON(request.Body)

product := Product{}

// tag your Product fields with appropriate options like
// -, omitempty, notraverse to get desired result.
// Not to worry, go-model does deep copy :)
errs := model.Copy(&product, tempProduct)
fmt.Println("Errors:", errs)

fmt.Printf("\nSource: %#v\n", tempProduct)
fmt.Printf("\nDestination: %#v\n", product)
```

#### Map Method
I want to convert my struct into Map (`map[string]interface{}`). Sure, go-model does deep convert.
```go
// tag your SearchResult fields with appropriate options like
// -, name, omitempty, notraverse to get desired result.
sr, _ := myapp.GetSearchResult( /* params here */ )

// Embedded/Anonymous struct fields appear in map at same level as represented by Go
srchResMap, err := model.Map(sr)
fmt.Println("Error:", err)

fmt.Printf("\nSearch Result Map: %#v\n", srchResMap)
```

#### Clone Method
I would like to clone my struct object. That's nice, you know go-model does deep processing.
```go
input := Product { /* Product struct field values go here */ }

// have your struct fields tagged appropriately. Options like
// -, name, omitempty, notraverse to get desired result.
clonedObj := model.Clone(input)

// let's see the result
fmt.Printf("\nCloned Object: %#v\n", clonedObj)
```

#### IsZero Method
I want to check my struct object is empty or not. Of course, go-model does deep zero check.
```go
// let's say you have just decoded/unmarshalled your request body to struct object.
productInfo, _ := myapp.ParseJSON(request.Body)

// wanna check productInfo is empty or not
isEmpty := model.IsZero(productInfo)

// tag your ProductInfo fields with appropriate options like
// -, omitempty, notraverse to get desired result.
fmt.Println("Hey, I have all fields zero value:", isEmpty)
```

#### HasZero Method
I want to check my struct object has any zero/empty value. Of course, go-model does deep zero check.
```go
// let's say you have just decoded/unmarshalled your request body to struct object.
productInfo, _ := myapp.ParseJSON(request.Body)

// wanna check productInfo is empty or not
isEmpty := model.HasZero(productInfo)

// tag your ProductInfo fields with appropriate options like
// -, omitempty, notraverse to get desired result.
fmt.Println("Hey, I have zero values:", isEmpty)
```

#### IsZeroInFields Method
Is it possible to check to particular fields has zero/empty values. Of-course you can.
```go
// let's say you have just decoded/unmarshalled your request body to struct object.
product, _ := myapp.ParseJSON(request.Body)

// check particular fields has zero value or not
fieldName, isEmpty := model.IsZeroInFields(product, "SKU", "Title", "InternalIdentifier")

fmt.Println("Empty Field Name:", fieldName)
fmt.Println("Yes, I have zero value:", isEmpty)
```

#### Fields Method
You wanna all the fields from `struct`, Yes you can have it :)
```go
src := SampleStruct {
  /* struct fields go here */
}

fields, _ := model.Fields(src)
fmt.Println("Fields:", fields)
```

#### Kind Method
go-model library provides an ability to know the `reflect.Kind` in as easy way.
```go
src := SampleStruct {
  /* struct fields go here */
}

fieldKind, _ := model.Kind(src, "BookingInfoPtr")
fmt.Println("Field kind:", fieldKind)
```

#### Tag Method
I want to get Go lang supported Tag value from my `struct`. Yes, it is easy to get it.
```go
src := SampleStruct {
	BookCount      int         `json:"-"`
	BookCode       string      `json:"-"`
	ArchiveInfo    BookArchive `json:"archive_info,omitempty"`
	Region         BookLocale  `json:"region,omitempty"`
}

tag, _ := model.Tag(src, "ArchiveInfo")
fmt.Println("Tag Value:", tag.Get("json"))

// Output:
Tag Value: archive_info,omitempty
```

#### Tags Method
I would like to get all the fields Tag values from my `struct`. It's easy.
```go
src := SampleStruct {
	BookCount      int         `json:"-"`
	BookCode       string      `json:"-"`
	ArchiveInfo    BookArchive `json:"archive_info,omitempty"`
	Region         BookLocale  `json:"region,omitempty"`
}

tags, _ := model.Tags(src)
fmt.Println("Tags:", tags)
```

#### Get Method
I want to get value by field name on my `struct`. Yes, it is easy to get it.
```go
src := SampleStruct {
	BookCount: 100,
	BookCode:  "GHT67HH00",
}

value, _ := model.Get(src, "BookCode")
fmt.Println("Value:", value)

// Output:
Value: GHT67HH00
```

#### Set Method
I want to set value by field name on my `struct`. Yes, it is easy to get it.
```go
src := SampleStruct {
	BookCount: 100,
	BookCode:  "GHT67HH00",
}

err := model.Set(&src, "BookCount", 200)
fmt.Println("Error:", err)
```

#### AddNoTraverseType & RemoveNoTraverseType Methods
There are scenarios, where you want the object values but not to traverse/look inside the struct object. Use `notraverse` option in the model tag for those fields or Add it `NoTraverseTypeList`. Customize it as per your need.

Default `NoTraverseTypeList` has these types `time.Time{}`, `&time.Time{}`, `os.File{}`, `&os.File{}`, `http.Request{}`, `&http.Request{}`, `http.Response{}`, `&http.Response{}`.
```go
// If you have added your type into list then you need not mention `notraverse` option for those types.

// Adding type into NoTraverseTypeList
model.AddNoTraverseType(time.Location{}, &time.Location{})

// Removing type from NoTraverseTypeList
model.RemoveNoTraverseType(time.Location{}, &time.Location{})
```

#### AddConversion & RemoveConversion Methods

This example registers a custom conversion from the `int` to the `string` type.
```go
AddConversion((*int)(nil), (*string)(nil), func(in reflect.Value) (reflect.Value, error) {
		return reflect.ValueOf(strconv.FormatInt(in.Int(), 10)), nil
	})
```

If a an integer field on the source struct matches the name of a string field on the target struct, the provided Converter method is invoked.

Note that if you want to register a converter from `int` to `*string` you will
have to provide a pointer to a pointer as destination type ( `(**string)(nil)`
).

More examples can be found in the [AddConversion godoc](https://godoc.org/github.com/jeevatkm/go-model#AddConversion).

## Versioning
go-model releases versions according to [Semantic Versioning](http://semver.org)

`gopkg.in/jeevatkm/go-model.vX` points to appropriate tag versions; `X` denotes version number and it's a stable release. It's recommended to use version, for eg. `gopkg.in/jeevatkm/go-model.v0`. Development takes place at the master branch. Although the code in master should always compile and test successfully, it might break API's. We aim to maintain backwards compatibility, but API's and behaviour might be changed to fix a bug.

## Contributing
Welcome! If you find any improvement or issue you want to fix, feel free to send a pull request. I like pull requests that include test cases for fix/enhancement. I have done my best to bring pretty good code coverage. Feel free to write tests.

BTW, I'd like to know what you think about go-model. Kindly open an issue or send me an email; it'd mean a lot to me.

## Author
Jeevanandam M. - jeeva@myjeeva.com

## Contributors
Have a look on [Contributors](https://github.com/jeevatkm/go-model/graphs/contributors) page.

## License
go-model released under MIT license, refer [LICENSE](LICENSE) file.
