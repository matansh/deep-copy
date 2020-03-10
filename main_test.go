package main

import (
	"bytes"
	"reflect"
	"regexp"
	"testing"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name    string
		types   typesVal
		path    string
		pointer bool
		skips   skipsVal
		want    []byte
	}{
		{name: "foo", types: typesVal{"Foo"}, path: "./testdata", want: []byte(FooFile)},
		{name: "foo - pointer", types: typesVal{"Foo"}, pointer: true, path: "./testdata", want: []byte(FooPointerFile)},
		{name: "foo - pointer, skip slice", types: typesVal{"Foo"}, pointer: true, skips: skipsVal{{"Slice": struct{}{}}}, path: "./testdata", want: []byte(FooPointerSkipSliceFile)},
		{name: "foo, skip map member", types: typesVal{"Foo"}, skips: skipsVal{{"Map[k]": struct{}{}}}, path: "./testdata", want: []byte(FooSkipMapFile)},
		{name: "alpha - with DeepCopy method", types: typesVal{"Alpha"}, path: "./testdata", want: []byte(AlphaPointer)},
		{name: "slicepointer, skip slice member", types: typesVal{"SlicePointer"}, skips: skipsVal{{"[i]": struct{}{}}}, path: "./testdata", want: []byte(SlicePointer)},
		{name: "foo, alpha, skips", types: typesVal{"Foo", "Alpha"}, skips: skipsVal{{"Map[k]": struct{}{}, "ch": struct{}{}}, {"D": struct{}{}, "E": struct{}{}}}, path: "./testdata", want: []byte(FooAlphaSkips)},
		{name: "issue 3, struct with slice of simple structs", types: typesVal{"I3WithSlice"}, pointer: true, path: "./testdata", want: []byte(Issue3SliceSimpleStruct)},
		{name: "issue 3, struct with map of simple struct keys", types: typesVal{"I3WithMap"}, pointer: true, path: "./testdata", want: []byte(Issue3MapSimpleStructKey)},
		{name: "issue 3, struct with map of simple struct values", types: typesVal{"I3WithMapVal"}, path: "./testdata", want: []byte(Issue3MapSimpleStructVal)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := run(tt.path, tt.types, tt.skips, tt.pointer)
			if err != nil {
				t.Fatal(err)
			}
			got = normalizeComment(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateFile() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

var re = regexp.MustCompile(`generated by .*deep-copy.*; DO NOT EDIT.`)

func normalizeComment(in []byte) []byte {
	return re.ReplaceAll(bytes.TrimSpace(in), []byte("generated by deep-copy; DO NOT EDIT."))
}

const (
	FooFile = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of Foo
func (o Foo) DeepCopy() Foo {
	var cp Foo = o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k, v := range o.Map {
			var cpv *Bar
			if v != nil {
				cpv = new(Bar)
				*cpv = *v
				if v.Slice != nil {
					cpv.Slice = make([]string, len(v.Slice))
					copy(cpv.Slice, v.Slice)
				}
			}
			cp.Map[k] = cpv
		}
	}
	if o.ch != nil {
		cp.ch = make(chan float32, cap(o.ch))
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return cp
}`
	FooPointerFile = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *Foo
func (o *Foo) DeepCopy() *Foo {
	var cp Foo = *o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k, v := range o.Map {
			var cpv *Bar
			if v != nil {
				cpv = new(Bar)
				*cpv = *v
				if v.Slice != nil {
					cpv.Slice = make([]string, len(v.Slice))
					copy(cpv.Slice, v.Slice)
				}
			}
			cp.Map[k] = cpv
		}
	}
	if o.ch != nil {
		cp.ch = make(chan float32, cap(o.ch))
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return &cp
}`
	FooPointerSkipSliceFile = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *Foo
func (o *Foo) DeepCopy() *Foo {
	var cp Foo = *o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k, v := range o.Map {
			var cpv *Bar
			if v != nil {
				cpv = new(Bar)
				*cpv = *v
			}
			cp.Map[k] = cpv
		}
	}
	if o.ch != nil {
		cp.ch = make(chan float32, cap(o.ch))
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return &cp
}`
	FooSkipMapFile = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of Foo
func (o Foo) DeepCopy() Foo {
	var cp Foo = o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k, v := range o.Map {
			cp.Map[k] = v
		}
	}
	if o.ch != nil {
		cp.ch = make(chan float32, cap(o.ch))
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return cp
}`
	AlphaPointer = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of Alpha
func (o Alpha) DeepCopy() Alpha {
	var cp Alpha = o
	if o.B != nil {
		cp.B = o.B.DeepCopy()
	}
	cp.G = o.G.DeepCopy()
	if o.D != nil {
		retV := o.D.DeepCopy()
		cp.D = &retV
	}
	{
		retV := o.E.DeepCopy()
		cp.E = *retV
	}
	return cp
}`
	SlicePointer = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of SlicePointer
func (o SlicePointer) DeepCopy() SlicePointer {
	var cp SlicePointer = o
	if o != nil {
		cp = make([]*int, len(o))
		copy(cp, o)
	}
	return cp
}`
	FooAlphaSkips = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of Foo
func (o Foo) DeepCopy() Foo {
	var cp Foo = o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k, v := range o.Map {
			cp.Map[k] = v
		}
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return cp
}

// DeepCopy generates a deep copy of Alpha
func (o Alpha) DeepCopy() Alpha {
	var cp Alpha = o
	if o.B != nil {
		cp.B = o.B.DeepCopy()
	}
	cp.G = o.G.DeepCopy()
	return cp
}`

	Issue3SliceSimpleStruct = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *I3WithSlice
func (o *I3WithSlice) DeepCopy() *I3WithSlice {
	var cp I3WithSlice = *o
	if o.a != nil {
		cp.a = make([]I3SimpleStruct, len(o.a))
		copy(cp.a, o.a)
	}
	return &cp
}`
	Issue3MapSimpleStructKey = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *I3WithMap
func (o *I3WithMap) DeepCopy() *I3WithMap {
	var cp I3WithMap = *o
	if o.a != nil {
		cp.a = make(map[I3SimpleStruct]string, len(o.a))
		for k, v := range o.a {
			cp.a[k] = v
		}
	}
	return &cp
}`
	Issue3MapSimpleStructVal = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of I3WithMapVal
func (o I3WithMapVal) DeepCopy() I3WithMapVal {
	var cp I3WithMapVal = o
	if o.a != nil {
		cp.a = make(map[string]I3SimpleStruct, len(o.a))
		for k, v := range o.a {
			cp.a[k] = v
		}
	}
	return cp
}`
)
