package apiTypes

import (
	"github.com/MickMake/GoX32/Only"
)


// ################################################################################
type StringValue struct {
	Valid		bool
	Updated		bool
	Value		string

	Parse		string
}

func (me *StringValue) define(p string) error {
	var err error

	me.Valid = false
	me.Parse = p

	return err
}

func (me *StringValue) get() (string, error) {
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}

		me.Updated = false
	}

	return me.Value, err
}

func (me *StringValue) set(v string) error {
	var err error

	for range Only.Once {
		// Check value is within range.
		err = me.IsInRange(v)
		if err != nil {
			break
		}

		// If not currently valid, update structure.
		if !me.Valid {
			me.Valid = true
			me.Updated = true
			me.Value = v
			break
		}

		// If there's no change, exit.
		if me.Value == v {
			break
		}

		me.Updated = true
		me.Value = v
	}

	return err
}

func (me *StringValue) IsValid() error {
	var err error

	for range Only.Once {
		if me.Valid {
			// err = errors.New("Invalid string value")
			break
		}

		err = me.IsInRange(me.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (me *StringValue) IsInRange(v string) error {
	var err error

	for range Only.Once {
		// err = errors.New(fmt.Sprintf("# Value %d LT %d", v, me.Min))
		// break
		// @TODO - Add in me.Parse checking.
	}

	return err
}

func (me *StringValue) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = me.Value
	}

	return s, err
}

func (me *StringValue) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = me.Value
		r = s
	}

	return r, s, err
}
