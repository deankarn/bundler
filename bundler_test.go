package bundler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"runtime"
	"testing"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//

func IsEqual(t *testing.T, val1, val2 interface{}) bool {
	v1 := reflect.ValueOf(val1)
	v2 := reflect.ValueOf(val2)

	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}

	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}

	if !v1.IsValid() && !v2.IsValid() {
		return true
	}

	v1Underlying := reflect.Zero(reflect.TypeOf(v1)).Interface()
	v2Underlying := reflect.Zero(reflect.TypeOf(v2)).Interface()

	if v1 == v1Underlying {
		if v2 == v2Underlying {
			goto CASE4
		} else {
			goto CASE3
		}
	} else {
		if v2 == v2Underlying {
			goto CASE2
		} else {
			goto CASE1
		}
	}

CASE1:
	return reflect.DeepEqual(v1.Interface(), v2.Interface())

CASE2:
	return reflect.DeepEqual(v1.Interface(), v2)
CASE3:
	return reflect.DeepEqual(v1, v2.Interface())
CASE4:
	return reflect.DeepEqual(v1, v2)
}

func Equal(t *testing.T, val1, val2 interface{}) {
	EqualSkip(t, 2, val1, val2)
}

func EqualSkip(t *testing.T, skip int, val1, val2 interface{}) {

	if !IsEqual(t, val1, val2) {

		_, file, line, _ := runtime.Caller(skip)
		fmt.Printf("%s:%d %v does not equal %v\n", path.Base(file), line, val1, val2)
		t.FailNow()
	}
}

func NotEqual(t *testing.T, val1, val2 interface{}) {
	NotEqualSkip(t, 2, val1, val2)
}

func NotEqualSkip(t *testing.T, skip int, val1, val2 interface{}) {

	if IsEqual(t, val1, val2) {
		_, file, line, _ := runtime.Caller(skip)
		fmt.Printf("%s:%d %v should not be equal %v\n", path.Base(file), line, val1, val2)
		t.FailNow()
	}
}

func PanicMatches(t *testing.T, fn func(), matches string) {
	PanicMatchesSkip(t, 2, fn, matches)
}

func PanicMatchesSkip(t *testing.T, skip int, fn func(), matches string) {

	_, file, line, _ := runtime.Caller(skip)

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)

			if err != matches {
				fmt.Printf("%s:%d Panic...  expected [%s] received [%s]", path.Base(file), line, matches, err)
				t.FailNow()
			}
		}
	}()

	fn()
}

func TestByFile(t *testing.T) {
	results := bytes.NewBuffer([]byte{})

	err := ByFile("files/test1.css", results)
	Equal(t, err, nil)
}

func TestByIo(t *testing.T) {
	input, err := ioutil.ReadFile("files/test1.css")
	if err != nil {
		panic(err)
	}

	results := bytes.NewBuffer([]byte{})

	err = ByIo(bytes.NewReader(input), "files", results)
	Equal(t, err, nil)
}

func BenchmarkByFile(b *testing.B) {

	for n := 0; n < b.N; n++ {
		results := bytes.NewBuffer([]byte{})

		ByFile("files/test1.css", results)
	}
}

func BenchmarkByIo(b *testing.B) {

	input, err := ioutil.ReadFile("files/test1.css")
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {

		results := bytes.NewBuffer([]byte{})

		ByIo(bytes.NewReader(input), "files", results)
	}
}
