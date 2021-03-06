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
	Val   *T
	Valid bool
	Set   bool
}

func (s Nullable[T]) GetValue() driver.Value {
	if s.Valid {
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
		*s = Nullable[T]{Valid: false}
		return nil
	}

	var val T
	*s = Nullable[T]{Val: &val, Valid: true}
	switch src.(type) {
	case []byte:
		if v, ok := reflect.ValueOf(s.Val).Interface().(sql.Scanner); ok {
			err := v.Scan(src)
			if err != nil {
				return err
			}
			val = reflect.ValueOf(v).Elem().Interface().(T)
		}
	case T:
		val = src.(T)
	}
	return nil
}

func (s Nullable[T]) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(s.Val)
}

func (s *Nullable[T]) UnmarshalJSON(data []byte) error {
	s.Set = true
	if string(data) == "null" {
		s.Valid = false
		return nil
	}
	s.Valid = true
	return json.Unmarshal(data, &s.Val)
}
