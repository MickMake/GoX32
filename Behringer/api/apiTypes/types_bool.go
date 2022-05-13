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

func (me *BoolValue) define(off string, on string) error {
	var err error

	me.Valid = false
	me.OffString = off
	me.OnString = on

	return err
}

func (me *BoolValue) get() (bool, error) {
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

func (me *BoolValue) set(v bool) error {
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

func (me *BoolValue) IsValid() error {
	var err error

	for range Only.Once {
		if me.Valid {
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

func (me *BoolValue) IsInRange(v bool) error {
	var err error

	for range Only.Once {
		// err = errors.New(fmt.Sprintf("# Value %d LT %d", v, me.Min))
		// break
		// @TODO - Not really needed
	}

	return err
}

func (me *BoolValue) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		if me.Value {
			s = me.OnString
		} else {
			s = me.OffString
		}
	}

	return s, err
}

func (me *BoolValue) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		if me.Value {
			r = me.OnString
		} else {
			r = me.OffString
		}

		s = strconv.FormatBool(me.Value)
	}

	return r, s, err
}
