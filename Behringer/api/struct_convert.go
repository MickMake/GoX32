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
}

type ConvertAlias string

type ConvertIncrement struct {
	Min 		float64 `json:"min"`
	Max 		float64 `json:"max"`
	Increment	float64 `json:"increment"`
	Precision   int     `json:"precision"`
}

type ConvertRange struct {
	InMin 		float64 `json:"in_min"`
	InMax 		float64 `json:"in_max"`
	OutMin 		float64 `json:"out_min"`
	OutMax 		float64 `json:"out_max"`
	Precision   int     `json:"precision"`
}

type ConvertMap map[string]string

type ConvertBitMap []string

type ConvertFunction string

type ConvertBinary string

type ConvertString struct {
	Size int `json:"size"`
}

type ConvertAsset struct {
	Url    bool `json:"url"`
	Icon   bool `json:"icon"`
	String bool `json:"string"`
}

type ConvertArray struct {
	Expected int      `json:"expected"`
	Names    []string `json:"names"`
}

type ConvertFloatMap struct {
	Values      map[string]string `json:"values"`
	Precision   int               `json:"precision"`
	Map         map[string]string `json:"-"`
	DefaultZero string            `json:"-"`
}

type ConvertInteger struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type ConvertBlob struct {
	Order []struct {
		Data  *ConvertBlobData      `json:"data"`
		Array *ConvertBlobDataArray `json:"array"`
	} `json:"order"`
	Sequence []ConvertBlobData `json:"-"`
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
type ConvertBlobDataArray struct {
	Data   ConvertBlobData `json:"data"`
	Count  int             `json:"count"`
	Offset int             `json:"offset"`
	Keys   []string        `json:"keys"`
	// NameFormat string `json:"name_format"`
}
type BlobReturn struct {
}

const Single = "0"
func (c *ConvertStruct) GetValues(values ...any) map[string]string {
	ret := make(map[string]string)

	for range Only.Once {
		for _, value := range values {
			switch {
				case c.Alias != nil:
					break

				case c.Increment != nil:
					// value = ToLinearDb(value, c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
					break

				case c.Range != nil:
					// ret[Single] = ToRange(fmt.Sprintf("%v", value), c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
					ret[Single] = c.Range.Convert(fmt.Sprintf("%v", value))
					break

				case c.Map != nil:
					val := fmt.Sprintf("%v", value)
					if v, ok := (*c.Map)[val]; ok {
						ret[Single] = v
					}
					break

				case c.BitMap != nil:
					ret[Single] = c.BitMap.Convert(value, 0)
					break

				case c.Function != nil:
					if *c.Function == "log" {
						ret[Single] = ToLogFunc(fmt.Sprintf("%v", value), 1)
						break
					}
					break

				case c.Binary != nil:
					ret[Single] = c.BitMap.Convert(value, 0)
					break

				case c.String != nil:
					break

				case c.Asset != nil:
					break

				case c.Array != nil:
					// 	value = strings.Join(*c.Array, ", ")
					break

				case c.FloatMap != nil:
					ret[Single] = c.FloatMap.Convert(value)
					break

				case c.Blob != nil:
					var val []byte
					for _, v := range value.([]byte) {
						val = append(val, v)
					}
					ret = c.Blob.Convert(val)
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
				// value = ToLinearDb(value, c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
				break

			case c.Range != nil:
				// ret = ToRange(fmt.Sprintf("%v", value), c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
				ret = c.Range.Convert(fmt.Sprintf("%v", value))
				break

			case c.Map != nil:
				val := fmt.Sprintf("%v", value)
				if v, ok := (*c.Map)[val]; ok {
					ret = v
				}
				break

			case c.BitMap != nil:
				ret = c.BitMap.Convert(value, 0)
				break

			case c.Function != nil:
				if *c.Function == "log" {
					ret = ToLogFunc(fmt.Sprintf("%v", value), 1)
					break
				}
				break

			case c.Binary != nil:
				ret = c.BitMap.Convert(value, 0)
				break

			case c.String != nil:
				break

			case c.Asset != nil:
				break

			case c.Array != nil:
				// 	value = strings.Join(*c.Array, ", ")
				break

			case c.FloatMap != nil:
				ret = c.FloatMap.Convert(value)
				break

			case c.Blob != nil:
				// Can't have a blob on a blob.
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


func (c *ConvertBitMap) Convert(value any, size uint32) string {
	var ret string

	for range Only.Once {
		if len(*c) == 0 {
			break
		}

		ret = fmt.Sprintf("%v", value)
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

func (c *ConvertRange) Convert(value string) string {
	// var ret string

	for range Only.Once {
		var err error
		var fv float64

		fv, err = strconv.ParseFloat(value, 64)
		if err != nil {
			break
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
			c.Precision = 1
		}

		value = strconv.FormatFloat(fv, 'f', c.Precision, 32)
	}

	return value
}

func (c *ConvertFloatMap) Convert(value any) string {
	var val string

	for range Only.Once {
		if c == nil {
			break
		}

		if len(c.Values) == 0 {
			break
		}

		if len(c.Map) == 0 {
			break
		}

		val = fmt.Sprintf("%v", value)
		if val == "" {
			// value = array.FloatValues[0]
			break
		}

		fv, err := strconv.ParseFloat(val, 64)
		if err != nil {
			break
		}

		val = strconv.FormatFloat(fv, 'f', c.Precision, 32)
		if v, ok := c.Map[val]; ok {
			value = v
			break
		}
	}

	return val
}

func (c *ConvertBlob) Convert(value []byte) map[string]string {
	ret := make(map[string]string)

	for range Only.Once {
		if c == nil {
			break
		}

		if len(c.Sequence) == 0 {
			break
		}

		r := bytes.NewReader(value)

		var i int
		for i = 0; (r.Len() > 0) && (i < len(c.Sequence)); i++ {
			v := c.Sequence[i].ReaderString(r)
			ret[c.Sequence[i].Key] = v
		}

		i = len(c.Sequence) - i
		if i > 0 {
			fmt.Printf("Expected %d more array elements.\n", i)
			ret["expected_count"] = fmt.Sprintf("%d", i)
		}

		if r.Len() > 0 {
			var remBytes string
			for r.Len() > 0 {
				v := byte(0)
				err := binary.Read(r, binary.BigEndian, &v)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					// ret = append(ret, fmt.Sprintf(`"error": "%s"`, err))
					ret["error"] = fmt.Sprintf("%s", err)
					break
				}
				remBytes += fmt.Sprintf("%.2X,", v)
			}
			fmt.Printf("Remaining bytes[%d]: %s\n", len(remBytes) / 3, remBytes)
			ret["remaining_bytes"] = fmt.Sprintf("%s", remBytes)
		}
	}

	return ret
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
			ret = fmt.Sprintf(`"%s": "Error with value (big_endian:%v) of type %s - %s"`,
				c.Key, c.BigEndian, c.Type, err)
		}
	}

	return ret
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

func ToLogFunc(value string, precision int) string {
	for range Only.Once {
		// s = strconv.FormatFloat(float64(fv), 'f', -1, 32)
		fv, err := strconv.ParseFloat(value, 64)
		if err != nil {
			value = "-inf"
			break
		}

		if fv == 0 {
			value = "-inf"
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
			precision = 1
		}

		value = strconv.FormatFloat(d, 'f', precision, 32)
	}

	return value
}

func ToLinear(value string, inMin string, inMax string, outMin string, outMax string) string {
	for range Only.Once {
		var err error
		var fv float64
		var inMinFloat float64
		var inMaxFloat float64
		var outMinFloat float64
		var outMaxFloat float64

		fv, err = strconv.ParseFloat(value, 64)
		if err != nil {
			break
		}

		type mapRange 	func(float64) (float64, error)

		var foo mapRange
		foo = MapRange64(RangeFloat64{inMinFloat, inMaxFloat}, RangeFloat64{outMinFloat, outMaxFloat})
		// Convert to linear scale.
		fv, err = foo(fv)

		if err != nil {
			break
		}

		value = strconv.FormatFloat(fv, 'f', -1, 32)
	}

	return value
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
