package apiTypes

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"strconv"
)


type EnumValue struct {
	Valid		bool
	Updated		bool
	Value		int32

	Options		map[string]int32
}

func (v *EnumValue) Define(sa []string) error {
	var err error

	v.Valid = false

	if v.Options == nil {
		v.Options = make(map[string]int32)
	}

	var i int
	var s string
	for i, s = range sa {
		// s = strings.ToLower(s)
		v.Options[s] = int32(i)
	}

	return err
}

func (v *EnumValue) get() (int32, error) {
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}

		v.Updated = false
	}

	return v.Value, err
}

func (v *EnumValue) set(value int32) error {
	var err error

	for range Only.Once {
		// Check value is within range.
		value, err = v.IsInRange(value)
		if err != nil {
			break
		}

		// If not currently valid, update structure.
		if !v.Valid {
			v.Valid = true
			v.Updated = true
			v.Value = value
			break
		}

		// If there's no change, exit.
		if v.Value == value {
			break
		}

		// for s, i := range me.Options {
		// 	if i == v {
		// 		me.Value = me.Options[s]
		// 		// me.StringValue = s
		// 		break
		// 	}
		// }

		v.Updated = true
		v.Value = value
	}

	return err
}

func (v *EnumValue) IsValid() error {
	var err error

	for range Only.Once {
		if v.Valid {
			// err = errors.New("Invalid value")
			break
		}

		_, err = v.IsInRange(v.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (v *EnumValue) IsInRange(value int32) (int32, error) {
	var err error

	for range Only.Once {
		err = errors.New(fmt.Sprintf("# Value %d invalid", value))
		for _, c := range v.Options {
			if c == value {
				err = nil
				break
			}
		}
	}

	return value, err
}

func (v *EnumValue) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		var val int32
		for s, val = range v.Options {
			if val == v.Value {
				break
			}
		}
	}

	return s, err
}

func (v *EnumValue) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		var val int32
		for r, val = range v.Options {
			if val == v.Value {
				break
			}
		}

		s = strconv.FormatInt(int64(v.Value), 10)
	}

	return r, s, err
}
