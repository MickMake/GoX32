package apiTypes

import (
	"github.com/MickMake/GoX32/Only"
	"strconv"
)

const (
	BoolOn = "ON"
	BoolOff = "OFF"
)


// ################################################################################
type BoolValue struct {
	Valid		bool
	Updated		bool
	Value		bool

	OffString 	string
	OnString	string
}

func (v *BoolValue) Define(off string, on string) error {
	var err error

	v.Valid = false
	v.OffString = off
	v.OnString = on

	return err
}

func (v *BoolValue) get() (bool, error) {
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

func (v *BoolValue) set(value bool) error {
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

func (v *BoolValue) IsValid() error {
	var err error

	for range Only.Once {
		if v.Valid {
			// err = errors.New("Invalid boolean value")
			break
		}

		// @TODO - Not really needed
		// err = me.IsInRange(me.Value)
		// if err != nil {
		// 	break
		// }
	}

	return err
}

func (v *BoolValue) IsInRange(value bool) error {
	var err error

	for range Only.Once {
		// err = errors.New(fmt.Sprintf("# Value %d LT %d", v, me.Min))
		// break
		// @TODO - Not really needed
	}

	return err
}

func (v *BoolValue) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		if v.Value {
			s = v.OnString
		} else {
			s = v.OffString
		}
	}

	return s, err
}

func (v *BoolValue) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		if v.Value {
			r = v.OnString
		} else {
			r = v.OffString
		}

		s = strconv.FormatBool(v.Value)
	}

	return r, s, err
}
