package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

type NullableValuer interface {
	IsSet() bool
	GetValue() driver.Value
}

type Nullable[T driver.Value] struct {
	Val *T
	Set bool
}

func (s Nullable[T]) GetValue() driver.Value {
	if s.Val != nil {
		return s.Val
	} else {
		return nil
	}
}

func (s Nullable[T]) IsSet() bool {
	return s.Set
}

func (s Nullable[T]) Value() (driver.Value, error) {
	return s.GetValue(), nil
}

func (s *Nullable[T]) Scan(src interface{}) error {
	if src == nil {
		*s = Nullable[T]{}
		return nil
	}

	var val T
	*s = Nullable[T]{Val: &val}
	switch src.(type) {
	case T:
		val = src.(T)
	default:
		if v, ok := reflect.ValueOf(s.Val).Interface().(sql.Scanner); ok {
			err := v.Scan(src)
			if err != nil {
				return err
			}
			val = reflect.ValueOf(v).Elem().Interface().(T)
		}
	}
	return nil
}

func (s Nullable[T]) MarshalJSON() ([]byte, error) {
	if s.Val == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(s.Val)
}

func (s *Nullable[T]) UnmarshalJSON(data []byte) error {
	s.Set = true
	if string(data) == "null" {
		return nil
	}
	return json.Unmarshal(data, &s.Val)
}
