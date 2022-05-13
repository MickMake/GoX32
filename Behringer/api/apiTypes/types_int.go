package apiTypes

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"strconv"
)


// ################################################################################
type Int32Value struct {
	Valid		bool
	Updated		bool
	Value		int32

	Min 		int32
	Max 		int32
}

func (me *Int32Value) define(min int32, max int32) error {
	var err error

	me.Valid = false
	me.Min = min
	me.Max = max

	return err
}

func (me *Int32Value) get() (int32, error) {
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

func (me *Int32Value) set(v int32) error {
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

func (me *Int32Value) IsValid() error {
	var err error

	for range Only.Once {
		if me.Valid {
			// err = errors.New("Invalid value")
			break
		}

		err = me.IsInRange(me.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (me *Int32Value) IsInRange(v int32) error {
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

func (me *Int32Value) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = fmt.Sprintf("%d", me.Value)
	}

	return s, err
}

func (me *Int32Value) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = strconv.FormatInt(int64(me.Value), 10)
		// s = fmt.Sprintf("%d", me.Value)
		r = s
	}

	return r, s, err
}


// ################################################################################
type Int64Value struct {
	Valid	bool
	Updated	bool
	Value	int64

	Min 	int64
	Max 	int64
}

func (me *Int64Value) define(min int64, max int64) error {
	var err error

	me.Valid = false
	me.Min = min
	me.Max = max

	return err
}

func (me *Int64Value) get() (int64, error) {
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

func (me *Int64Value) set(v int64) error {
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

func (me *Int64Value) IsValid() error {
	var err error

	for range Only.Once {
		if me.Valid {
			// err = errors.New("Invalid value")
			break
		}

		err = me.IsInRange(me.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (me *Int64Value) IsInRange(v int64) error {
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

func (me *Int64Value) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = fmt.Sprintf("%d", me.Value)
	}

	return s, err
}

func (me *Int64Value) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = strconv.FormatInt(me.Value, 10)
		// s = fmt.Sprintf("%d", me.Value)
		r = s
	}

	return r, s, err
}
