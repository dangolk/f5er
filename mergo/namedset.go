package mergo

import (
	"fmt"
	"reflect"
	"sort"
)

type NamedSet map[string]reflect.Value

func newNamedSet() NamedSet {
	return make(NamedSet)
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

func getName(value reflect.Value) (string, error) {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		return value.Interface().(string), nil
	case reflect.Struct:
		value = value.FieldByName("Name")
		if value.IsValid() {
			return value.Interface().(string), nil
		}
	}
	return "", fmt.Errorf("Object of type %[1]T (%[1]v) is not a string, and not a struct with a Name property", value)
}

func (set *NamedSet) Cardinality() int {
	return len(*set)
}

func (set *NamedSet) Add(i reflect.Value) bool {
	key, err := getName(i)
	if key == "" || err != nil {
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return false
	}

	if _, found := (*set)[key]; found {
		return false //False if it existed already
	}

	// fmt.Printf(" -- Adding: %s -> %v\n", key, i)
	(*set)[key] = i
	return true
}

func (set *NamedSet) Upsert(i reflect.Value) bool {
	key, err := getName(i)
	if key == "" || err != nil {
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return false
	}

	// fmt.Printf(" -- Upserting: %s -> %v\n", key, i)
	(*set)[key] = i
	return true
}

func (set *NamedSet) ToSortedSlice() []reflect.Value {
	structs := make([]reflect.Value, 0, set.Cardinality())
	for _, elem := range *set {
		structs = append(structs, elem)
	}

	sort.Slice(structs, func(i, j int) bool {
		iname, _ := getName(structs[i])
		jname, _ := getName(structs[j])
		return iname < jname
	})

	return structs
}
