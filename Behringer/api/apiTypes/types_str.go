package apiTypes

import (
	"github.com/MickMake/GoX32/Only"
)


type StringValue struct {
	Valid		bool
	Updated		bool
	Value		string

	Parse		string
}

func (v *StringValue) Define(p string) error {
	var err error

	v.Valid = false
	v.Parse = p

	return err
}

func (v *StringValue) get() (string, error) {
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

func (v *StringValue) set(value string) error {
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

func (v *StringValue) IsValid() error {
	var err error

	for range Only.Once {
		if v.Valid {
			// err = errors.New("Invalid string value")
			break
		}

		err = v.IsInRange(v.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (v *StringValue) IsInRange(value string) error {
	var err error

	for range Only.Once {
		// err = errors.New(fmt.Sprintf("# Value %d LT %d", v, me.Min))
		// break
		// @TODO - Add in me.Parse checking.
	}

	return err
}

func (v *StringValue) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = v.Value
	}

	return s, err
}

func (v *StringValue) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = v.Value
		r = s
	}

	return r, s, err
}
