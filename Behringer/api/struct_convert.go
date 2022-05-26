package api

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Behringer/api/apiReflect"
	"github.com/MickMake/GoX32/Only"
	"math"
	"strconv"
	"strings"
)


const (
	Off = "OFF"
	On = "ON"
	DefaultPrecision = 3
	Single = "0"
)


type ConvertStruct struct {
	Alias     *ConvertAlias     `json:"alias,omitempty"`
	Increment *ConvertIncrement `json:"increment,omitempty"`
	Range     *ConvertRange     `json:"range,omitempty"`
	Map       *ConvertMap       `json:"map,omitempty"`
	BitMap    *ConvertBitMap    `json:"bit_map,omitempty"`
	Function  *ConvertFunction  `json:"function,omitempty"`
	Binary    *ConvertBinary    `json:"binary,omitempty"`
	String    *ConvertString    `json:"string,omitempty"`
	Asset     *ConvertAsset     `json:"asset,omitempty"`
	Array     *ConvertArray     `json:"array,omitempty"`
	FloatMap  *ConvertFloatMap  `json:"float_map,omitempty"`
	Integer   *ConvertInteger   `json:"integer,omitempty"`
	Blob      *ConvertBlob      `json:"blob,omitempty"`
	Index     *ConvertIndex     `json:"index,omitempty"`
}

// type ValuesMap UnitValueMap		// map[string]UnitValue	// string
//
// func (v *ValuesMap) Add(key string, value string) {
// 	for range Only.Once {
// 		(*v)[key] = value
// 	}
// }
//
// func (v *ValuesMap) Append(values ValuesMap) {
// 	for range Only.Once {
// 		for key, value := range values {
// 			(*v)[key] = value
// 		}
// 	}
// }

func (c *ConvertStruct) GetValues(values ...any) UnitValueMap {
	ret := make(UnitValueMap)

	for range Only.Once {
		// if c.Array != nil {
		// 	ret = c.Array.Convert(values...)
		// 	break
		// }

		for index, value := range values {
			switch {
				case c.Alias != nil:
					break

				case c.Increment != nil:
					ret.Add(Single, c.Increment.Convert(value), "")
					// value = ToLinearDb(value, c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
					break

				case c.Range != nil:
					ret.Add(Single, c.Range.Convert(value), "")
					break

				case c.Map != nil:
					ret.Add(Single, c.Map.Convert(value), "")
					break

				case c.Index != nil:
					ret.Add(Single, c.Index.Convert(value), "")
					break

				case c.BitMap != nil:
					ret.Add(Single, c.BitMap.Convert(value, 0), "")
					break

				case c.Function != nil:
					ret.Add(Single, c.Function.Convert(value), "")
					break

				case c.Binary != nil:
					ret.Add(Single, c.Binary.Convert(value), "")
					break

				case c.String != nil:
					ret.Add(Single, c.String.Convert(value), "")
					break

				case c.Asset != nil:
					ret.Add(Single, c.Asset.Convert(value), "")
					break

				case c.Array != nil:
					ret.Append(c.Array.Convert(index, value))
					break

				case c.FloatMap != nil:
					ret.Add(Single, c.FloatMap.Convert(value), "")
					break

				case c.Blob != nil:
					ret = c.Blob.Convert(value)
					break
			}
		}
	}

	return ret
}

func (c *ConvertStruct) GetValue(value any) string {
	var ret string

	for range Only.Once {
		switch {
			case c.Alias != nil:
				break

			case c.Increment != nil:
				ret = c.Increment.Convert(value)
				// value = ToLinearDb(value, c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
				break

			case c.Range != nil:
				ret = c.Range.Convert(value)
				break

			case c.Map != nil:
				ret = c.Map.Convert(value)
				break

			case c.Index != nil:
				ret = c.Index.Convert(value)
				break

			case c.BitMap != nil:
				ret = c.BitMap.Convert(value, 0)
				break

			case c.Function != nil:
				ret = c.Function.Convert(value)
				break

			case c.Binary != nil:
				ret = c.Binary.Convert(value)
				break

			case c.String != nil:
				ret = c.String.Convert(value)
				break

			case c.Asset != nil:
				ret = c.Asset.Convert(value)
				break

			case c.FloatMap != nil:
				ret = c.FloatMap.Convert(value)
				break

			case c.Array != nil:
			// Can't have an array within a blob.

			case c.Blob != nil:
				// Can't have a blob within a blob.
		}
	}

	return ret
}

func (c *ConvertStruct) SetValue(value any) any {
	var ret any

	for range Only.Once {
		switch {
			case c.Alias != nil:
				// Can't set an Alias.
				break

			case c.Increment != nil:
				ret = c.Increment.Set(value)
				// value = ToLinearDb(value, c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
				break

			case c.Range != nil:
				ret = c.Range.Set(value)
				break

			case c.Map != nil:
				ret = c.Map.Set(value)
				break

			case c.Index != nil:
				ret = c.Index.Set(value)
				break

			case c.BitMap != nil:
				ret = c.BitMap.Set(value)
				break

			case c.Function != nil:
				ret = c.Function.Set(value)
				break

			case c.Binary != nil:
				ret = c.Binary.Set(value)
				break

			case c.String != nil:
				ret = c.String.Set(value)
				break

			case c.Asset != nil:
				ret = c.Asset.Set(value)
				break

			case c.FloatMap != nil:
				ret = c.FloatMap.Set(value)
				break

			case c.Array != nil:
				// Can't set an Array.

			case c.Blob != nil:
				// Can't set a Blob.
		}
	}

	return ret
}

func (c *ConvertStruct) GetConvertType() string {
	var ret string

	for range Only.Once {
		ret = apiReflect.GetJsonTagIfNotNil(*c)
		if c.Map == nil {
			break
		}

		// Special cases.
		if len(*c.Map) == 2 {
			if (*c.Map)[Single] == Off {
				ret = "binary"
				break
			}
			if (*c.Map)[Single] == On {
				ret = "binary"
				break
			}
		}
	}

	return ret
}


type ConvertAlias string

func (c *ConvertAlias) Import() ConvertStruct {
	var ret ConvertStruct
	// @TODO - Not yet tested.

	for range Only.Once {
		// ret = c.Get(p.Convert.Alias)
	}

	return ret
}


type ConvertIncrement struct {
	Min 		float64 `json:"min"`
	Max 		float64 `json:"max"`
	Increment	float64 `json:"increment"`
	Precision   int     `json:"precision"`
}
func (c *ConvertIncrement) Convert(value any) string {
	var ret string
	// @TODO - Not yet tested.

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if (c.Min == 0) && (c.Max == 0) {
			c.Min = 0
			c.Max = 1
		}

		if ret == "" {
			break
		}
	}

	return ret
}

func (c *ConvertIncrement) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertIncrement) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}

		if (c.Min == 0) && (c.Max == 0) {
			c.Min = 0
			c.Max = 1
		}

		if c.Precision == 0 {
			c.Precision = DefaultPrecision
		}

		if c.Increment == 0 {
			c.Increment = 0.1
		}
	}

	return err
}


type ConvertRange struct {
	InMin 		float64 `json:"in_min"`
	InMax 		float64 `json:"in_max"`
	OutMin 		float64 `json:"out_min"`
	OutMax 		float64 `json:"out_max"`
	Precision   int     `json:"precision"`
}
func (c *ConvertRange) Convert(value any) string {
	var ret string

	for range Only.Once {
		var err error
		var fv float64

		ret = fmt.Sprintf("%v", value)

		fv, err = strconv.ParseFloat(ret, 64)
		if err != nil {
			break
		}

		if (c.InMin == 0) && (c.InMax == 0) {
			c.InMin = 0
			c.InMax = 1
		}

		type mapRange 	func(float64) (float64, error)

		var foo mapRange
		foo = MapRange64(RangeFloat64{c.InMin, c.InMax}, RangeFloat64{c.OutMin, c.OutMax})
		// Convert to linear scale.
		fv, err = foo(fv)

		if err != nil {
			break
		}

		if c.Precision == 0 {
			c.Precision = DefaultPrecision
		}

		ret = strconv.FormatFloat(fv, 'f', c.Precision, 32)
	}

	return ret
}

func (c *ConvertRange) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertRange) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}

		if (c.InMin == 0) && (c.InMax == 0) {
			c.InMin = 0
			c.InMax = 1
		}
	}

	return err
}


type ConvertMap map[string]string
func (c *ConvertMap) Convert(value any) string {
	var ret string

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if len(*c) == 0 {
			break
		}

		if v, ok := (*c)[ret]; ok {
			ret = v
		}
	}

	return ret
}

func (c *ConvertMap) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if len(*c) == 0 {
			break
		}

		for k, v := range *c {
			if ret == v {
				ret = k
				break
			}
		}
	}

	return ret
}

func (c *ConvertMap) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}

		if len(*c) == 0 {
			err = errors.New("empty map")
			break
		}
	}

	return err
}

func (c *ConvertMap) GetOptions() []string {
	var ret []string

	for range Only.Once {
		if c == nil {
			break
		}

		ret = make([]string, len(*c))
		for key, value := range *c {
			i, err := strconv.ParseUint(key, 10, 64)
			if err != nil {
				continue
			}
			if i >= uint64(len(*c)) {
				continue
			}
			ret[i] = value
		}
	}

	return ret
}


type ConvertIndex []string
func (c *ConvertIndex) Convert(value any) string {
	var ret string

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if len(*c) == 0 {
			break
		}

		iv, err := strconv.ParseUint(ret, 10, 32)
		if err != nil {
			break
		}

		ret = (*c)[iv]
	}

	return ret
}

func (c *ConvertIndex) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if len(*c) == 0 {
			break
		}

		for k, v := range *c {
			if ret == v {
				ret = int32(k)
				break
			}
		}
	}

	return ret
}

func (c *ConvertIndex) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}

		if len(*c) == 0 {
			err = errors.New("empty map")
			break
		}
	}

	return err
}

func (c *ConvertIndex) GetOptions() []string {
	var ret []string

	for range Only.Once {
		if c == nil {
			break
		}

		ret = *c
	}

	return ret
}


type ConvertBitMap []string
func (c *ConvertBitMap) Convert(value any, size uint32) string {
	var ret string

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if len(*c) == 0 {
			break
		}

		if ret == "" {
			ret = (*c)[0]
			break
		}

		iv, err := strconv.ParseUint(ret, 10, 64)
		if err != nil {
			break
		}

		if iv == 0 {
			ret = (*c)[0]
			break
		}

		if size == 0 {
			size = uint32(len(*c)) - 1
		}

		var elems []string
		for j := uint32(0); j < size; j++ {
			if iv & (1 << j) != 0 {
				elems = append(elems, (*c)[j+1])
			}
		}

		ret = strings.Join(elems, ", ")
	}

	return ret
}

func (c *ConvertBitMap) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertBitMap) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}

		if len(*c) == 0 {
			err = errors.New("empty map")
			break
		}
	}

	return err
}

func (c *ConvertBitMap) GetOptions() []string {
	var ret []string

	for range Only.Once {
		if c == nil {
			break
		}

		for _, o := range *c {
			ret = append(ret, o)
		}
	}

	return ret
}


type ConvertFunction string
func (c *ConvertFunction) Convert(value any) string {
	var ret string
	// @TODO - Not yet tested.

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if len(*c) == 0 {
			break
		}

		if *c == "log" {
			ret = ToLogFunc(fmt.Sprintf("%v", value), 1) // DefaultPrecision)
			break
		}

		if ret == "" {
			break
		}
	}

	return ret
}
func ToLogFunc(value any, precision int) string {
	var ret string

	for range Only.Once {
		// s = strconv.FormatFloat(float64(fv), 'f', -1, 32)
		ret = fmt.Sprintf("%v", value)
		fv, err := strconv.ParseFloat(ret, 64)
		if err != nil {
			ret = "-inf"
			break
		}

		fv = toFixed(fv, 4)

		if fv == 0 {
			ret = "-inf"
			break
		}

		// Convert to log scale.
		var d float64
		if fv >= 0.5 {
			d = float64(fv * 40.0 - 30.0)

		} else if fv >= 0.25 {
			d = float64(fv * 80.0 - 50.0)

		} else if fv >= 0.0625 {
			d = float64(fv * 160.0 - 70.0)

		} else if fv >= 0.0 {
			d = float64(fv * 480.0 - 90.0)
		}

		// def float_to_db(f):
		//    if (f >= 0.5):
		//        db = f * 40. - 30. # max dB value: +10
		//    elif (f >= 0.25):
		//        db = f * 80. - 50.
		//    elif (f >= 0.0625):
		//        db = f * 160. - 70.
		//    elif (f >= 0.0):
		//        db = f * 480. - 90. # min db value: -90
		//    return db

		if precision == 0 {
			precision = DefaultPrecision
		}

		ret = strconv.FormatFloat(float64(d), 'f', precision, 32)
	}

	return ret
}

func (c *ConvertFunction) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertFunction) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}

		if *c == "" {
			err = errors.New("empty function")
			break
		}
	}

	return err
}


type ConvertBinary struct {
	On          string `json:"on"`
	Off         string `json:"off"`
	Type        string `json:"type"`
	IsSwitch    bool   `json:"-"`
	IsMomentary bool   `json:"-"`
	NameOn      string `json:"name_on"`
	NameOff     string `json:"name_off"`
}
func (c *ConvertBinary) Convert(value any) string {
	var ret string

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if ret == "" {
			break
		}

		if ret == "0" {
			ret = c.Off
			break
		}

		if ret == "OFF" {
			ret = c.Off
			break
		}

		if ret == "1" {
			ret = c.On
			break
		}

		if ret == "ON" {
			ret = c.On
			break
		}
	}

	return ret
}

func (c *ConvertBinary) Set(value any) any {
	var ret any

	for range Only.Once {
		if c == nil {
			break
		}

		switch strings.ToUpper(fmt.Sprintf("%v", value)) {
			case c.On:
				ret = int32(1)

			case c.Off:
				ret = int32(0)
		}
	}

	return ret
}

func (c *ConvertBinary) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}

		normal := Labels{"", "normal"}
		inverted := Labels{"swap", "swapped", "invert", "inverted"}
		oneshot := Labels{"oneshot", "momentary"}

		switch {
			case (c.On == "") && (c.Off != ""):
				err = errors.New("missing On binary value")

			case (c.On != "") && (c.Off == ""):
				err = errors.New("missing Off binary value")

			case (c.On != "") && (c.Off == "") && (c.Type == ""):
				fallthrough
			case normal.ValueExists(c.Type):
				// case c.Type == "":
				// 	fallthrough
				// case c.Type == "normal":
				// 	fallthrough
				c.On = On
				c.Off = Off
				// ret = &ConvertMap{ "0":Off, "1":On }
				c.IsSwitch = true

			case inverted.ValueExists(c.Type):
				// case c.Type == "swap":
				// 	fallthrough
				// case c.Type == "swapped":
				// 	fallthrough
				// case c.Type == "invert":
				// 	fallthrough
				// case c.Type == "inverted":
				c.On = Off
				c.Off = On
				// ret = &ConvertMap{ "0":On, "1":Off }
				c.IsSwitch = true
				c.IsMomentary = false

			case oneshot.ValueExists(c.Type):
			// case c.Type == "oneshot":
			// 	fallthrough
			// case c.Type == "momentary":
				c.On = On
				c.Off = Off
				c.IsSwitch = false
				c.IsMomentary = true

			default:
				c.IsSwitch = true
				c.IsMomentary = false
		}

		// if (strings.ToUpper(c.On) == On) && (strings.ToUpper(c.Off) == Off) {
		// 	c.IsSwitch = true
		// }
		// if (strings.ToUpper(c.On) == Off) && (strings.ToUpper(c.Off) == On) {
		// 	c.IsSwitch = true
		// }
	}

	return err
}


type ConvertString struct {
	Size int `json:"size"`
}
func (c *ConvertString) Convert(value any) string {
	var ret string
	// @TODO - Not yet tested.

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if c.Size == 0 {
			// break
		}

		if ret == "" {
			break
		}
	}

	return ret
}

func (c *ConvertString) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertString) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}
	}

	return err
}


type ConvertAsset struct {
	Url    bool `json:"url"`
	Icon   bool `json:"icon"`
	String bool `json:"string"`
}
func (c *ConvertAsset) Convert(value any) string {
	var ret string
	// @TODO - Not yet tested.

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if ret == "" {
			break
		}
	}

	return ret
}

func (c *ConvertAsset) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertAsset) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}
	}

	return err
}


type ConvertArray struct {
	Expected int      `json:"expected"`
	Names    []string `json:"names"`
}
func (c *ConvertArray) Convert(index int, value any) UnitValueMap {
	ret := make(UnitValueMap)

	for range Only.Once {
		if c == nil {
			break
		}

		name := fmt.Sprintf("%d", index)
		if c.Names != nil {
			if index < len(c.Names) {
				name = c.Names[index]
			}
		}

		// if c.Expected == 0 {
		// 	break
		// }

		ret.Add(name, value, "")
	}

	return ret
}

func (c *ConvertArray) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertArray) Import() error {
	var err error

	for range Only.Once {
		if c == nil {
			break
		}

		// if c.Expected == 0 {
		// 	err = errors.New("empty array")
		// 	break
		// }

		if len(c.Names) == 0 {
			err = errors.New("empty array")
			break
		}
	}

	return err
}


type ConvertFloatMap struct {
	Values      map[string]string `json:"values"`
	Precision   int               `json:"precision"`
	Map         map[string]string `json:"-"`
	DefaultZero string            `json:"-"`
}
func (c *ConvertFloatMap) Convert(value any) string {
	var ret string

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if len(c.Values) == 0 {
			break
		}

		// if len(c.Map) == 0 {
		// 	break
		// }

		if ret == "" {
			// value = array.FloatValues[0]
			break
		}

		fv, err := strconv.ParseFloat(ret, 64)
		if err != nil {
			break
		}

		if c.Precision == 0 {
			c.Precision = DefaultPrecision
		}

		ret = strconv.FormatFloat(fv, 'f', c.Precision, 32)
		if v, ok := c.Values[ret]; ok {
			ret = v
			break
		}
	}

	return ret
}

func (c *ConvertFloatMap) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertFloatMap) Import() error {
	var err error
	// @TODO - Not yet tested.

	for range Only.Once {
		if c == nil {
			err = errors.New("nil structure")
			break
		}

		c.Map = make(map[string]string)
		if c.Precision == 0 {
			c.Precision = 4
		}
		minFv := 1.0
		for k, v := range c.Values {
			var fv float64
			fv, err = strconv.ParseFloat(k, 64)
			if err != nil {
				break
			}
			if fv < minFv {
				minFv = fv
			}
			k = strconv.FormatFloat(fv, 'f', c.Precision, 32)
			c.Map[k] = v
		}
		c.DefaultZero = strconv.FormatFloat(minFv, 'f', c.Precision, 32)
	}

	return err
}


type ConvertInteger struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
func (c *ConvertInteger) Convert(value any) string {
	var ret string
	// @TODO - Not yet tested.

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}

		if (c.Min == 0) && (c.Max == 0) {
			c.Min = 0
			c.Max = 1
		}

		if ret == "" {
			break
		}
	}

	return ret
}

func (c *ConvertInteger) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertInteger) Import() error {
	var err error
	// @TODO - Not yet tested.

	for range Only.Once {
		if c == nil {
			err = errors.New("nil structure")
			break
		}

		if err != nil {
			break
		}
	}

	return err
}


type ConvertBlob struct {
	Order []struct {
		Data  *ConvertBlobData      `json:"data"`
		Array *ConvertBlobDataArray `json:"array"`
	} `json:"order"`
	Sequence []ConvertBlobData `json:"-"`
}
func (c *ConvertBlob) Convert(value any) UnitValueMap {		// map[string]UnitValue {	// string {
	ret := make(UnitValueMap)	// map[string]UnitValue)	// string)

	for range Only.Once {
		if c == nil {
			break
		}

		if len(c.Sequence) == 0 {
			break
		}

		var val []byte
		for _, v := range value.([]byte) {
			val = append(val, v)
		}
		r := bytes.NewReader(val)

		var i int
		for i = 0; (r.Len() > 0) && (i < len(c.Sequence)); i++ {
			// unitValue := ret[c.Sequence[i].Key]
			// v := c.Sequence[i].ReaderString(r)
			// unitValue.Set(v)

			v := c.Sequence[i].ReaderString(r)
			ret.Add(c.Sequence[i].Key, v, "")
		}

		i = len(c.Sequence) - i
		if i > 0 {
			fmt.Printf("Expected %d more array elements.\n", i)
			// ret["expected_count"] = fmt.Sprintf("%d", i)
			ret.Add("expected_count", i, "")
		}

		if r.Len() > 0 {
			var remBytes string
			for r.Len() > 0 {
				v := byte(0)
				err := binary.Read(r, binary.BigEndian, &v)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					// ret = append(ret, fmt.Sprintf(`"error": "%s"`, err))
					// ret["error"].Set() = fmt.Sprintf("%s", err)
					ret.Add("error", err, "")
					break
				}
				remBytes += fmt.Sprintf("%.2X,", v)
			}
			fmt.Printf("Remaining bytes[%d]: %s\n", len(remBytes) / 3, remBytes)
			// ret["remaining_bytes"] = fmt.Sprintf("%s", remBytes)
			ret.Add("remaining_bytes", remBytes, "")
		}
	}

	return ret
}

func (c *ConvertBlob) Set(value any) any {
	var ret any

	for range Only.Once {
		ret = fmt.Sprintf("%v", value)

		if c == nil {
			break
		}
	}

	return ret
}

func (c *ConvertBlob) Import(aliases Aliases) error {
	var err error

	for range Only.Once {
		if c == nil {
			err = errors.New("nil structure")
			break
		}

		for _, b := range c.Order {
			switch {
				case b.Data != nil:
					if b.Data.Convert != nil {
						if b.Data.Convert.Alias != nil {
							foo := aliases.Get(b.Data.Convert.Alias)
							b.Data.Convert = &foo
						}
					}
					if b.Data.Key == "" {
						b.Data.Key = "%d"
					}
					c.Sequence = append(c.Sequence, *b.Data)

				case b.Array != nil:
					if b.Array.Data.Convert != nil {
						if b.Array.Data.Convert.Alias != nil {
							foo := aliases.Get(b.Array.Data.Convert.Alias)
							b.Array.Data.Convert = &foo
						}
					}

					if b.Array.Keys != nil {
						for _, v := range b.Array.Keys {
							if b.Array.Data.Key == "" {
								if len(b.Array.Keys) > 0 {
									b.Array.Data.Key = "%s"
								} else {
									b.Array.Data.Key = "%d"
								}
							}
							c.Sequence = append(c.Sequence, ConvertBlobData {
								Convert:   b.Array.Data.Convert,
								Unit:      b.Array.Data.Unit,
								Key:       fmt.Sprintf("%s", v),
								Type:      b.Array.Data.Type,
								BigEndian: b.Array.Data.BigEndian,
							})
						}
						continue
					}

					for i := 0; i < b.Array.Count; i++ {
						if b.Array.Data.Key == "" {
							if len(b.Array.Keys) > 0 {
								b.Array.Data.Key = "%s"
							} else {
								b.Array.Data.Key = "%d"
							}
						}
						c.Sequence = append(c.Sequence, ConvertBlobData {
							Convert:   b.Array.Data.Convert,
							Unit:      b.Array.Data.Unit,
							Key:       fmt.Sprintf(b.Array.Data.Key, i + b.Array.Offset),	// , b.Array.Data.Unit, v),
							Type:      b.Array.Data.Type,
							BigEndian: b.Array.Data.BigEndian,
						})
					}
					continue
			}
		}
		c.Order = ConvertBlobOrder{}
	}

	return err
}


type ConvertBlobOrder []struct {
	Data  *ConvertBlobData      `json:"data"`
	Array *ConvertBlobDataArray `json:"array"`
}


type ConvertBlobData struct {
	Convert   *ConvertStruct `json:"convert"`
	Unit      string         `json:"unit"`
	Key       string         `json:"key"`
	Type      string         `json:"type"`
	BigEndian bool           `json:"big_endian"`
}
func (c *ConvertBlobData) ReaderValue(reader *bytes.Reader) (string, error) {
	var ret string
	var err error

	for range Only.Once {
		var bo binary.ByteOrder
		if c.BigEndian {
			bo = binary.BigEndian
		} else {
			bo = binary.LittleEndian
		}

		switch c.Type {
			case "int32":
				v := int32(0)
				err = binary.Read(reader, bo, &v)
				ret = fmt.Sprintf("%d", v)
			case "int64":
				v := int64(0)
				err = binary.Read(reader, bo, &v)
				ret = fmt.Sprintf("%d", v)
			case "float32":
				v := float32(0)
				err = binary.Read(reader, bo, &v)
				ret = fmt.Sprintf("%f", v)
			case "float64":
				v := float64(0)
				err = binary.Read(reader, bo, &v)
				ret = fmt.Sprintf("%f", v)
		}

		if err != nil {
			break
		}

		if c.Convert == nil {
			break
		}
		ret = c.Convert.GetValue(ret)
	}

	return ret, err
}

func (c *ConvertBlobData) ReaderString(reader *bytes.Reader) string {
	var ret string

	for range Only.Once {
		var err error
		ret, err = c.ReaderValue(reader)
		if err != nil {
			ret = fmt.Sprintf(`"%s": "Error with value (big_endian:%t) of type %s - %s"`,
				c.Key, c.BigEndian, c.Type, err)
		}
	}

	return ret
}


type ConvertBlobDataArray struct {
	Data   ConvertBlobData `json:"data"`
	Count  int             `json:"count"`
	Offset int             `json:"offset"`
	Keys   []string        `json:"keys"`
	// NameFormat string `json:"name_format"`
}
func (c *ConvertBlobDataArray) ReaderStrings(reader *bytes.Reader) []string {
	var ret []string

	for range Only.Once {

		// Using string keys instead of counter.
		if len(c.Keys) > 0 {
			for _, i := range c.Keys {
				b := ConvertBlobData {
					Type:      c.Data.Type,
					Key:       i,
					BigEndian: c.Data.BigEndian,
				}
				ret = append(ret, b.ReaderString(reader))
			}
			break
		}

		// Using counter instead of string keys.
		for i := 0; i < c.Count; i++ {
			b := ConvertBlobData {
				Type:      c.Data.Type,
				Key:       fmt.Sprintf(c.Data.Key, i + c.Offset),
				BigEndian: c.Data.BigEndian,
			}
			ret = append(ret, b.ReaderString(reader))
		}
	}

	return ret
}


type RangeFloat32 struct {
	Min float32
	Max float32
}

func MapRange32(xr RangeFloat32, yr RangeFloat32) func(float32) (float32, error) {
	// normalize direction of ranges so that out-of-range test works
	if xr.Min > xr.Max {
		xr.Min, xr.Max = xr.Max, xr.Min
		yr.Min, yr.Max = yr.Max, yr.Min
	}
	// compute slope, intercept
	m := (yr.Max - yr.Min) / (xr.Max - xr.Min)
	b := yr.Min - m*xr.Min

	// return function literal
	return func(x float32) (y float32, ok error) {
		if x < xr.Min || x > xr.Max {
			return 0, errors.New("Value out of range") // out of range
		}
		f2 := m*x + b
		return f2, nil
	}
}


type RangeFloat64 struct {
	Min float64
	Max float64
}

func MapRange64(xr RangeFloat64, yr RangeFloat64) func(float64) (float64, error) {
	// normalize direction of ranges so that out-of-range test works
	if xr.Min > xr.Max {
		xr.Min, xr.Max = xr.Max, xr.Min
		yr.Min, yr.Max = yr.Max, yr.Min
	}
	// compute slope, intercept
	m := (yr.Max - yr.Min) / (xr.Max - xr.Min)
	b := yr.Min - m*xr.Min

	// return function literal
	return func(x float64) (y float64, ok error) {
		if x < xr.Min || x > xr.Max {
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


type Labels []string
func (l *Labels) ValueExists(value string) bool {
	var ok bool
	for _, l := range *l {
		if l == value {
			ok = true
		}
	}
	return ok
}