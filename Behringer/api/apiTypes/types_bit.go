package apiTypes

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"strconv"
)


type BitValue struct {
	Valid		bool
	Updated		bool
	Value		int32

	Min 		int32
	Max 		int32
}

func (v *BitValue) Define(min int32, max int32) error {
	var err error

	v.Valid = false
	v.Min = min
	v.Max = max

	return err
}

func (v *BitValue) get() (int32, error) {
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

func (v *BitValue) set(value int32) error {
	var err error

	for range Only.Once {
		// Check value is within range.
		err = v.IsInRange(value)
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

		v.Updated = true
		v.Value = value
	}

	return err
}

func (v *BitValue) IsValid() error {
	var err error

	for range Only.Once {
		if v.Valid {
			// err = errors.New("# Invalid bitwise value")
			break
		}

		err = v.IsInRange(v.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (v *BitValue) IsInRange(value int32) error {
	var err error

	for range Only.Once {
		if v.Value < v.Min {
			err = errors.New(fmt.Sprintf("# Value %d LT %d", value, v.Min))
			break
		}

		if v.Value > v.Max {
			err = errors.New(fmt.Sprintf("# Value %d GT %d", value, v.Max))
			break
		}
	}

	return err
}

func (v *BitValue) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = fmt.Sprintf("%X", v.Value)
	}

	return s, err
}

func (v *BitValue) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = fmt.Sprintf("%X", v.Value)
		r = strconv.FormatInt(int64(v.Value), 10)
	}

	return r, s, err
}
