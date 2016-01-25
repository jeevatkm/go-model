# go-model [![Build Status](https://travis-ci.org/jeevatkm/go-model.svg?branch=master)](https://travis-ci.org/jeevatkm/go-model) [![GoCover](http://gocover.io/_badge/github.com/jeevatkm/go-model)](http://gocover.io/github.com/jeevatkm/go-model) [![GoDoc](https://godoc.org/github.com/jeevatkm/go-model?status.svg)](https://godoc.org/github.com/jeevatkm/go-model) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Robust & Easy to use model mapper and model utility methods for Go. Typical methods increase productivity and make Go developement more fun :smile:

**v0.1 released and tagged on Jan 22, 2016**

go-model tested with Go `v1.2` and above.

## Features
go-model provides handy methods (`Copy`, `Map`, `Clone`, `IsZero`, etc.) to process struct with below highlighted features. It's born from typical need while developing Go application or utility. I hope it's helpful!
* Embedded/Anonymous struct
* Multi-level nested struct/map/slice
* Pointer and non-pointer within struct/map/slice
* Struct within map and slice
* Embedded/Anonymous fields appear in map at same level as represented by Go
* Interface within struct/map/slice
* Add global no traverse type to the list or use `notraverse` option in the struct field
* Options to name map key, omit empty fields, and instruct not to traverse with struct/map/slice

## Installation

#### Stable - Version
Please refer section [Versioning](#versioning) for detailed info.

```sh
# install the library
go get gopkg.in/jeevatkm/go-model.v0
```

#### Latest
```sh
# install the latest & greatest library
go get github.com/jeevatkm/go-model
```

## Usage
Import go-model into your code and refer it as `model`. Have a look on [model test cases](model_test.go) to know more possibilities.
```go
import (
  "gopkg.in/jeevatkm/go-model.v0"
)
```

### Methods
* Copy - [usage](#copy-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Copy)
* Map - [usage](#map-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Map)
* Clone - [usage](#clone-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#Clone)
* IsZero - [usage](#iszero-method), [godoc](https://godoc.org/github.com/jeevatkm/go-model#IsZero)
* AddNoTraverseType - [usage](#addnotraversetype--removenotraversetype-methods), [godoc](https://godoc.org/github.com/jeevatkm/go-model#AddNoTraverseType)
* RemoveNoTraverseType - [usage](#addnotraversetype--removenotraversetype-methods), [godoc](https://godoc.org/github.com/jeevatkm/go-model#RemoveNoTraverseType)

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
I want to convert my struct into Map (`map[string]interface{}`). Not to worry, go-model does deep convert.
```go
// tag your SearchResult fields with appropriate options like 
// -, name, omitempty, notraverse to get desired result.
sr, _ := myapp.GetSearchResult( /* params here */ )

srchResMap, err := model.Map(sr)
fmt.Println("Error:", err)

fmt.Printf("\nSearch Result Map: %#v\n", srchResMap)
```

#### Clone Method
I would like to clone my struct object. Not to worry, go-model does deep processing.
```go
input := Product { /* Product struct field values go here */ }

// have your struct fields tagged appropriately. Options like 
// -, name, omitempty, notraverse to get desired result.
clonedObj := model.Clone(input)

// let's see the result
fmt.Printf("\nCloned Object: %#v\n", clonedObj)
```

#### IsZero Method
I want to check my struct object is empty or not. Not to worry, go-model does deep zero check.
```go
// let's say you have just decoded/unmarshalled your request body to struct object.
productInfo, _ := myapp.ParseJSON(request.Body)

// wanna check productInfo is empty or not
isEmpty := model.IsZero(productInfo)

// tag your ProductInfo fields with appropriate options like 
// -, omitempty, notraverse to get desired result.
fmt.Println("Hey, I have zero values:", isEmpty)
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

## Versioning
go-model releases versions according to [Semantic Versioning](http://semver.org)

`gopkg.in/jeevatkm/go-model.vX` points to appropriate tag versions; `X` denotes version number and it's a stable release. It's recommended to use version, for eg. `gopkg.in/jeevatkm/go-model.v0`. Development takes place at the master branch. Although the code in master should always compile and test successfully, it might break API's. We aim to maintain backwards compatibility, but API's and behaviour might be changed to fix a bug.
 
## Contributing
Welcome! If you find any improvement or issue you want to fix, feel free to send a pull request. I like pull requests that include test cases for fix/enhancement. I have done my best to bring pretty good code coverage. Feel free to write tests.

BTW, I'd like to know what you think about go-model. Kindly open an issue or send me an email; it'd mean a lot to me.

## Author
Jeevanandam M. - jeeva@myjeeva.com

## License
go-model released under MIT license, refer [LICENSE](LICENSE) file.
