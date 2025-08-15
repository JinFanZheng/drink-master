package models

import (
	"database/sql/driver"
	"fmt"
)

// BitBool is a custom type to handle MySQL bit(1) fields
// MySQL bit(1) fields are returned as []byte, which cannot be directly scanned into Go's bool or int8
type BitBool int8

// Scan implements the sql.Scanner interface
func (b *BitBool) Scan(value interface{}) error {
	if value == nil {
		*b = 0
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// MySQL bit(1) returns []byte{0x00} for false, []byte{0x01} for true
		if len(v) > 0 && v[0] == 1 {
			*b = 1
		} else {
			*b = 0
		}
	case int64:
		if v == 1 {
			*b = 1
		} else {
			*b = 0
		}
	case bool:
		if v {
			*b = 1
		} else {
			*b = 0
		}
	default:
		return fmt.Errorf("cannot scan %T into BitBool", value)
	}
	return nil
}

// Value implements the driver.Valuer interface
func (b BitBool) Value() (driver.Value, error) {
	if b == 1 {
		return int64(1), nil
	}
	return int64(0), nil
}

// Bool returns the boolean representation
func (b BitBool) Bool() bool {
	return b == 1
}

// Int8 returns the int8 representation
func (b BitBool) Int8() int8 {
	return int8(b)
}

// String returns the string representation
func (b BitBool) String() string {
	if b == 1 {
		return "true"
	}
	return "false"
}
