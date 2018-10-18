// Copyright 2013 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on src/pkg/reflect/deepequal.go from official
// golang's stdlib.

package mergo

import (
	"reflect"
)

func hasExportedField(dst reflect.Value) (exported bool) {
	for i, n := 0, dst.NumField(); i < n; i++ {
		field := dst.Type().Field(i)
		if field.Anonymous && dst.Field(i).Kind() == reflect.Struct {
			exported = exported || hasExportedField(dst.Field(i))
		} else {
			exported = exported || len(field.PkgPath) == 0
		}
	}
	return
}

type MergeStrategy int

const (
	Overwrite       MergeStrategy = iota
	AppendAdditive  MergeStrategy = iota
	UniqueFirstSeen MergeStrategy = iota
	UniqueLastSeen  MergeStrategy = iota
)

func (mt MergeStrategy) String() string {
	names := [...]string{"overwrite", "additive", "unique-first-seen", "unique-last-seen"}
	if mt < AppendAdditive || mt > UniqueLastSeen {
		return names[0]
	}
	return names[mt]
}

type Config struct {
	SkipEmptyFields bool
	Override        bool
	AppendSlice     bool
	MergeStrategy   MergeStrategy
	Transformers    Transformers
}

type Transformers interface {
	Transformer(reflect.Type) func(dst, src reflect.Value) error
}

func mergeSlice(dst, src reflect.Value, config *Config) {
	new := src

	switch config.MergeStrategy {
	case AppendAdditive:
		new.Set(reflect.AppendSlice(dst, src))

	case UniqueFirstSeen:
		unique := &NamedSet{}
		// Start with our target (aka dst) first so that those values are "seen"
		// first, and subsequent duplicate values in src are ignored
		for i := 0; i < dst.Len(); i++ {
			unique.Add(dst.Index(i))
		}
		for i := 0; i < src.Len(); i++ {
			unique.Add(src.Index(i))
		}

		new = reflect.MakeSlice(dst.Type(), 0, unique.Cardinality())
		for _, obj := range unique.ToSortedSlice() {
			new = reflect.Append(new, obj)
		}
		break

	case UniqueLastSeen:
		unique := &NamedSet{}
		// Start with our target (aka dst) first so that those values are "seen"
		// first, and subsequent duplicate values in src override the ones in target
		for i := 0; i < dst.Len(); i++ {
			unique.Upsert(dst.Index(i))
		}
		for i := 0; i < src.Len(); i++ {
			unique.Upsert(src.Index(i))
		}

		new = reflect.MakeSlice(dst.Type(), 0, unique.Cardinality())
		for _, obj := range unique.ToSortedSlice() {
			new = reflect.Append(new, obj)
		}
		break
	}

	dst.Set(new)
}

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deepMerge(dst, src reflect.Value, visited map[uintptr]*visit, depth int, config *Config) (err error) {
	// fmt.Printf("deep merge called with:\n\t dst: %v\n\tsrc: %v\n\tdepth: %d\n\tconfig: %v\n", dst, src, depth, config)
	override := config.Override

	if !src.IsValid() {
		return
	}

	if isEmptyValue(dst) && config.SkipEmptyFields {
		return
	}

	if dst.CanAddr() {
		addr := dst.UnsafeAddr()
		h := 17 * addr
		seen := visited[h]
		typ := dst.Type()
		for p := seen; p != nil; p = p.next {
			if p.ptr == addr && p.typ == typ {
				return nil
			}
		}
		// Remember, remember...
		visited[h] = &visit{addr, typ, seen}
	}

	if config.Transformers != nil && !isEmptyValue(dst) {
		if fn := config.Transformers.Transformer(dst.Type()); fn != nil {
			err = fn(dst, src)
			return
		}
	}

	switch dst.Kind() {
	case reflect.Struct:
		if hasExportedField(dst) {
			for i, n := 0, dst.NumField(); i < n; i++ {
				if err = deepMerge(dst.Field(i), src.Field(i), visited, depth+1, config); err != nil {
					return
				}
			}
		} else {
			if dst.CanSet() && !isEmptyValue(src) && (override || isEmptyValue(dst)) {
				dst.Set(src)
			}
		}
	case reflect.Map:
		if dst.IsNil() && !src.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}
		for _, key := range src.MapKeys() {
			srcElement := src.MapIndex(key)
			if !srcElement.IsValid() {
				continue
			}
			dstElement := dst.MapIndex(key)
			switch srcElement.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Interface, reflect.Slice:
				if srcElement.IsNil() {
					continue
				}
				fallthrough
			default:
				if !srcElement.CanInterface() {
					continue
				}
				switch reflect.TypeOf(srcElement.Interface()).Kind() {
				case reflect.Struct:
					fallthrough
				case reflect.Ptr:
					fallthrough
				case reflect.Map:
					if err = deepMerge(dstElement, srcElement, visited, depth+1, config); err != nil {
						return
					}
				case reflect.Slice:
					srcSlice := reflect.ValueOf(srcElement.Interface())

					var dstSlice reflect.Value
					if !dstElement.IsValid() || dstElement.IsNil() {
						dstSlice = reflect.MakeSlice(srcSlice.Type(), 0, srcSlice.Len())
					} else {
						dstSlice = reflect.ValueOf(dstElement.Interface())
					}

					if config.MergeStrategy == Overwrite {
						dstSlice = reflect.MakeSlice(srcSlice.Type(), 0, srcSlice.Len())
						dstSlice = reflect.AppendSlice(dstSlice, srcSlice)
					} else {
						mergeSlice(dstSlice, srcSlice, config)
					}

					dst.SetMapIndex(key, dstSlice)
				}
			}
			if dstElement.IsValid() && reflect.TypeOf(srcElement.Interface()).Kind() == reflect.Map {
				continue
			}

			if srcElement.IsValid() && (override || (!dstElement.IsValid() || isEmptyValue(dst))) {
				if dst.IsNil() {
					dst.Set(reflect.MakeMap(dst.Type()))
				}
				dst.SetMapIndex(key, srcElement)
			}
		}
	case reflect.Slice:
		if !dst.CanSet() {
			break
		} else {
			if config.MergeStrategy == Overwrite {
				dst.Set(src)
			} else {
				mergeSlice(dst, src, config)
			}
		}

	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if src.IsNil() {
			break
		}
		if src.Kind() != reflect.Interface {
			if dst.IsNil() || override {
				if dst.CanSet() && (override || isEmptyValue(dst)) {
					dst.Set(src)
				}
			} else if src.Kind() == reflect.Ptr {
				if err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
					return
				}
			} else if dst.Elem().Type() == src.Type() {
				if err = deepMerge(dst.Elem(), src, visited, depth+1, config); err != nil {
					return
				}
			} else {
				return ErrDifferentArgumentsTypes
			}
			break
		}
		if dst.IsNil() || override {
			if dst.CanSet() && (override || isEmptyValue(dst)) {
				dst.Set(src)
			}
		} else if err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
			return
		}
	default:
		if dst.CanSet() && !isEmptyValue(src) && (override || isEmptyValue(dst)) {
			dst.Set(src)
		}
	}
	return
}

// Merge will fill any empty for value type attributes on the dst struct using corresponding
// src attributes if they themselves are not empty. dst and src must be valid same-type structs
// and dst must be a pointer to struct.
// It won't merge unexported (private) fields and will do recursively any exported field.
func Merge(dst, src interface{}, opts ...func(*Config)) error {
	return merge(dst, src, opts...)
}

// MergeWithOverwrite will do the same as Merge except that non-empty dst attributes will be overriden by
// non-empty src attribute values.
// Deprecated: use Merge(…) with WithOverride
func MergeWithOverride(dst, src interface{}, opts ...func(*Config)) error {
	return merge(dst, src, append(opts, WithOverride)...)
}

// WithTransformers adds transformers to merge, allowing to customize the merging of some types.
func WithTransformers(transformers Transformers) func(*Config) {
	return func(config *Config) {
		config.Transformers = transformers
	}
}

// WithOverride will make merge override non-empty dst attributes with non-empty src attributes values.
func WithOverride(config *Config) {
	config.Override = true
}

// WithAppendSlice will make merge append slices instead of overwriting it
func WithAppendSlice(config *Config) {
	config.AppendSlice = true
	config.MergeStrategy = AppendAdditive
}

// WithAppendSlice will make merge append slices instead of overwriting it
func WithAppend(config *Config) {
	config.AppendSlice = true
	config.MergeStrategy = AppendAdditive
}

// WithUniqueFirstSeen  will make merge add the first unique slice elements it find
// instead of blindly appending
func WithUniqueFirstSeen(config *Config) {
	config.AppendSlice = true
	config.MergeStrategy = UniqueFirstSeen
}

// WithUniqueLastSeen will make merge add the last unique slice elements it finds
// instead of blindly appending
func WithUniqueLastSeen(config *Config) {
	config.AppendSlice = true
	config.MergeStrategy = UniqueLastSeen
}

func merge(dst, src interface{}, opts ...func(*Config)) error {
	var (
		vDst, vSrc reflect.Value
		err        error
	)

	config := &Config{}

	for _, opt := range opts {
		opt(config)
	}

	if vDst, vSrc, err = resolveValues(dst, src); err != nil {
		return err
	}
	if vDst.Type() != vSrc.Type() {
		return ErrDifferentArgumentsTypes
	}
	return deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0, config)
}
