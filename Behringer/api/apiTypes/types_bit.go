package apiTypes

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"strconv"
)


// ################################################################################
type BitValue struct {
	Valid		bool
	Updated		bool
	Value		int32

	Min 		int32
	Max 		int32
}

func (me *BitValue) define(min int32, max int32) error {
	var err error

	me.Valid = false
	me.Min = min
	me.Max = max

	return err
}

func (me *BitValue) get() (int32, error) {
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

func (me *BitValue) set(v int32) error {
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

func (me *BitValue) IsValid() error {
	var err error

	for range Only.Once {
		if me.Valid {
			// err = errors.New("# Invalid bitwise value")
			break
		}

		err = me.IsInRange(me.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (me *BitValue) IsInRange(v int32) error {
	var err error

	for range Only.Once {
		if me.Value < me.Min {
			err = errors.New(fmt.Sprintf("# Value %d LT %d", v, me.Min))
			break
		}

		if me.Value > me.Max {
			err = errors.New(fmt.Sprintf("# Value %d GT %d", v, me.Max))
			break
		}
	}

	return err
}

func (me *BitValue) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = fmt.Sprintf("%X", me.Value)
	}

	return s, err
}

func (me *BitValue) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = fmt.Sprintf("%X", me.Value)
		r = strconv.FormatInt(int64(me.Value), 10)
	}

	return r, s, err
}
