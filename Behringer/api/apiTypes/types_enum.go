package apiTypes

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"strconv"
)


// ################################################################################
type EnumValue struct {
	Valid		bool
	Updated		bool
	Value		int32

	Options		map[string]int32
}

func (me *EnumValue) define(sa []string) error {
	var err error

	me.Valid = false

	if me.Options == nil {
		me.Options = make(map[string]int32)
	}

	var i int
	var s string
	for i, s = range sa {
		// s = strings.ToLower(s)
		me.Options[s] = int32(i)
	}

	return err
}

func (me *EnumValue) get() (int32, error) {
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

func (me *EnumValue) set(v int32) error {
	var err error

	for range Only.Once {
		// Check value is within range.
		v, err = me.IsInRange(v)
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

		// for s, i := range me.Options {
		// 	if i == v {
		// 		me.Value = me.Options[s]
		// 		// me.StringValue = s
		// 		break
		// 	}
		// }

		me.Updated = true
		me.Value = v
	}

	return err
}

func (me *EnumValue) IsValid() error {
	var err error

	for range Only.Once {
		if me.Valid {
			// err = errors.New("Invalid value")
			break
		}

		_, err = me.IsInRange(me.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (me *EnumValue) IsInRange(v int32) (int32, error) {
	var err error

	for range Only.Once {
		err = errors.New(fmt.Sprintf("# Value %d invalid", v))
		for _, c := range me.Options {
			if c == v {
				err = nil
				break
			}
		}
	}

	return v, err
}

func (me *EnumValue) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		var v int32
		for s, v = range me.Options {
			if v == me.Value {
				break
			}
		}
	}

	return s, err
}

func (me *EnumValue) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = me.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		var v int32
		for r, v = range me.Options {
			if v == me.Value {
				break
			}
		}

		s = strconv.FormatInt(int64(me.Value), 10)
	}

	return r, s, err
}
