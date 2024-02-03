package utils

import (
	"errors"
	"reflect"
)

func SliceToMap[K string | int | int64 | int32, V any](s []V, filedName string) (map[K]V, error) {

	//get filed value by filedName
	m := make(map[K]V)

	for i := range s {

		t := s[i]
		v := reflect.ValueOf(&t).Elem()
		f := v.FieldByName(filedName)
		if f.IsValid() {
			m[f.Interface().(K)] = s[i]
			continue
		}

		return nil, errors.New("FiledName not found")

	}

	return m, nil
}

func SliceGroupBy[K string | int | int64 | int32, V any](s []V, filedName string) (map[K][]V, error) {

	m := make(map[K][]V)

	for i := range s {

		t := s[i]
		v := reflect.ValueOf(&t).Elem()
		f := v.FieldByName(filedName)
		if f.IsValid() {
			m[f.Interface().(K)] = append(m[f.Interface().(K)], s[i])
			continue
		}

		return nil, errors.New("FiledName not found")

	}

	return m, nil
}
