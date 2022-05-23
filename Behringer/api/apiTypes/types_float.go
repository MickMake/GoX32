package apiTypes

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"math"
	"strconv"
)


const (
	// InfinityFloat = -math.MaxFloat32
	InfinityFloat = -90
)


type Float32Value struct {
	Valid		bool
	Updated		bool
	Value		float32

	Min 		float32
	Max 		float32
	Linear		bool
	Increment	float32
	Precision   int

	mapRange 	func(float32) (float32, error)
}

func (v *Float32Value) Define(min float32, max float32, linear bool, inc float32) error {
	var err error

	v.Valid = false
	v.Min = min
	v.Max = max
	v.Linear = linear
	v.Increment = inc
	v.mapRange = mapRange32(rangeBounds32{0, 1}, rangeBounds32{v.Min, v.Max})

	return err
}

func (v *Float32Value) get() (float32, error) {
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

func (v *Float32Value) set(value float32) error {
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

func (v *Float32Value) IsValid() error {
	var err error

	for range Only.Once {
		if v.Valid {
			// err = errors.New("Invalid value")
			break
		}

		err = v.IsInRange(v.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (v *Float32Value) IsInRange(value float32) error {
	var err error

	for range Only.Once {
		// @TODO - Need to figure out how to do mapRange64 properly.
		// if me.Value < me.Min {
		if value < 0 {
			err = errors.New(fmt.Sprintf("# Value %f LT %f", value, v.Min))
			break
		}

		// if me.Value > me.Max {
		if value > 1 {
			err = errors.New(fmt.Sprintf("# Value %f GT %f", value, v.Max))
			break
		}
	}

	return err
}

func (v *Float32Value) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = strconv.FormatFloat(float64(v.Value), 'f', -1, 32)
		// s = fmt.Sprintf("%f", me.Value)
	}

	return s, err
}

func (v *Float32Value) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = strconv.FormatFloat(float64(v.Value), 'f', -1, 32)

		if v.Linear {
			// Convert to linear scale.
			var val float32
			val, err = v.mapRange(v.Value)
			if err != nil {
				break
			}

			r = strconv.FormatFloat(toFixed(float64(val), 2), 'f', -1, 32)
			break
		}

		// Convert to log scale.
		var d float32
		if v.Value >= 0.5 {
			d = v.Value * 40. - 30.

		} else if v.Value >= 0.25 {
			d = v.Value * 80. -50.

		} else if v.Value >= 0.0625 {
			d = v.Value * 160. - 70.

		} else if v.Value >= 0.0 {
			d = v.Value * 480. - 90.
		}

		r = strconv.FormatFloat(toFixed(float64(d), 2), 'f', -1, 32)
	}

	return r, s, err
}


type Float64Value struct {
	Valid		bool
	Updated		bool
	Value		float64

	Min 		float64
	Max 		float64
	Linear		bool
	Increment	float64
	Precision   int

	mapRange func(float64) (float64, error)
}

func (v *Float64Value) Define(min float64, max float64, linear bool, inc float64) error {
	var err error

	v.Valid = false
	v.Min = min
	v.Max = max
	v.Linear = linear
	v.Increment = inc
	v.mapRange = mapRange64(rangeBounds64{0, 1}, rangeBounds64{v.Min, v.Max})

	return err
}

func (v *Float64Value) get() (float64, error) {
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

func (v *Float64Value) set(value float64) error {
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

func (v *Float64Value) IsValid() error {
	var err error

	for range Only.Once {
		if v.Valid {
			// err = errors.New("# Invalid value")
			break
		}

		err = v.IsInRange(v.Value)
		if err != nil {
			break
		}
	}

	return err
}

func (v *Float64Value) IsInRange(value float64) error {
	var err error

	for range Only.Once {
		// @TODO - Need to figure out how to do mapRange64 properly.
		// if me.Value < me.Min {
		if value < 0 {
			err = errors.New(fmt.Sprintf("# Value %f LT %f", value, v.Min))
			break
		}

		// if me.Value > me.Max {
		if value > 1 {
			err = errors.New(fmt.Sprintf("# Value %f GT %f", value, v.Max))
			break
		}
	}

	return err
}

func (v *Float64Value) getString() (string, error) {
	var s string
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}

		s = strconv.FormatFloat(v.Value, 'f', -1, 64)
		// s = fmt.Sprintf("%f", me.Value)
	}

	return s, err
}

func (v *Float64Value) getReal() (string, string, error) {
	var r string	// Real value.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		err = v.IsValid()
		if err != nil {
			break
		}
		// me.Updated = false

		s = strconv.FormatFloat(v.Value, 'f', -1, 64)

		if v.Linear {
			// Convert to linear scale.
			var val float64
			val, err = v.mapRange(v.Value)
			if err != nil {
				break
			}

			r = strconv.FormatFloat(toFixed(val, 2), 'f', -1, 32)
			break
		}

		// Convert to log scale.
		var d float64
		if v.Value >= 0.5 {
			d = v.Value * 40. - 30.

		} else if v.Value >= 0.25 {
			d = v.Value * 80. -50.

		} else if v.Value >= 0.0625 {
			d = v.Value * 160. - 70.

		} else if v.Value >= 0.0 {
			d = v.Value * 480. - 90.
		}

		r = strconv.FormatFloat(toFixed(d, 2), 'f', -1, 32)
	}

	return r, s, err
}


type rangeBounds32 struct {
	b1, b2 float32
}
func mapRange32(xr, yr rangeBounds32) func(float32) (float32, error) {
	// normalize direction of ranges so that out-of-range test works
	if xr.b1 > xr.b2 {
		xr.b1, xr.b2 = xr.b2, xr.b1
		yr.b1, yr.b2 = yr.b2, yr.b1
	}
	// compute slope, intercept
	m := (yr.b2 - yr.b1) / (xr.b2 - xr.b1)
	b := yr.b1 - m*xr.b1

	// return function literal
	return func(x float32) (y float32, ok error) {
		if x < xr.b1 || x > xr.b2 {
			return 0, errors.New("Value out of range") // out of range
		}
		f2 := m*x + b
		return f2, nil
	}
}

type rangeBounds64 struct {
	b1, b2 float64
}
func mapRange64(xr, yr rangeBounds64) func(float64) (float64, error) {
	// normalize direction of ranges so that out-of-range test works
	if xr.b1 > xr.b2 {
		xr.b1, xr.b2 = xr.b2, xr.b1
		yr.b1, yr.b2 = yr.b2, yr.b1
	}
	// compute slope, intercept
	m := (yr.b2 - yr.b1) / (xr.b2 - xr.b1)
	b := yr.b1 - m*xr.b1

	// return function literal
	return func(x float64) (y float64, ok error) {
		if x < xr.b1 || x > xr.b2 {
			return 0, errors.New("Value out of range") // out of range
		}
		f2 := m*x + b
		return f2, nil
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num * output)) / output
}
