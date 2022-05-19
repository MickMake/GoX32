package api

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Behringer/api/output"
	"github.com/MickMake/GoX32/Only"
	"math"
	"regexp"
	"strconv"
	"strings"
)


const (
	DefaultParentId = "virtual"

	TypeInstant = "instant"
)


type Aliases map[ConvertAlias]ConvertStruct

func (a *Aliases) Get(selector *ConvertAlias) ConvertStruct {
	if selector == nil {
		return ConvertStruct{}
	}
	if ret, ok := (*a)[*selector]; ok {
		return ret
	}
	return ConvertStruct{}
}


type PointsMapFile struct {
	Aliases   Aliases   `json:"aliases"`
	PointsMap PointsMap `json:"points"`
	PointsArrayMap PointsArrayMap `json:"points_array_map"`
}

type PointsArrayMap struct {
	Min int `json:"min"`
	Max int `json:"max"`
	Increment int `json:"increment"`
	PointsMap PointsMap `json:"points"`
}

func ImportPoints(parentId string, filenames ...string) (PointsMap, error) {
	var pm PointsMapFile
	var err error

	for range Only.Once {
		if parentId == "" {
			parentId = DefaultParentId
		}
		pm.Aliases = make(Aliases)
		pm.PointsMap = make(PointsMap)

		for _, filename := range filenames {
			var pmi PointsMapFile
			err = output.FileRead(filename, &pmi)
			if err != nil {
				err = errors.New(fmt.Sprintf("Error reading points json file '%s': %s", filename, err))
				break
			}
			pm.Aliases.Append(pmi.Aliases)
			pm.PointsMap.Append(pmi.PointsMap)

			if len(pmi.PointsArrayMap.PointsMap) == 0 {
				continue
			}

			for i := pmi.PointsArrayMap.Min; i <= pmi.PointsArrayMap.Max; i++ {
				for n, p := range pmi.PointsArrayMap.PointsMap {
					if n == "" {
						delete(pmi.PointsMap, n)
						continue
					}

					name := fmt.Sprintf(n, i)
					if p.EndPoint != "" {
						p.EndPoint = fmt.Sprintf(p.EndPoint, i)
					}
					pm.PointsMap[name] = p
				}
			}
		}
		if err != nil {
			break
		}

		for n, p := range pm.PointsMap {
			if n == "" {
				delete(pm.PointsMap, n)
				continue
			}

			p.Valid = true
			if p.EndPoint == "" {
				p.EndPoint = n
			}
			if p.Id == "" {
				p.Id = JoinStringsForId(p.EndPoint)
			}
			if p.ParentId == "" {
				p.ParentId = parentId
			}
			if p.FullId == "" {
				p.FullId = JoinStringsForId(p.ParentId, p.Id)
			}
			if p.Name == "" {
				p.Name = p.EndPoint	// JoinStrings(p.ParentId, p.EndPoint)
			}
			if p.Type == "" {
				p.Type = TypeInstant
			}

			if p.Convert.Alias != nil {
				p.Convert = pm.Aliases.Get(p.Convert.Alias)
			}

			if p.Convert.Map != nil {
				p.Unit = ""
				// p.Unit = UnitState
				// @TODO = States are binary.
			}

			if p.Convert.Range != nil {
				if (p.Convert.Range.InMin == 0) && (p.Convert.Range.InMax == 0) {
					p.Convert.Range.InMin = 0
					p.Convert.Range.InMax = 1
				}
			}

			if p.Convert.BitMap != nil {
				p.Unit = ""
			}

			if p.Convert.Asset != nil {
				p.Unit = ""
			}

			if p.Convert.Binary != nil {
				p.Unit = ""
				switch *p.Convert.Binary {
					case "":
						fallthrough
					case "normal":
						p.Convert.Map = &ConvertMap{ "0":"OFF", "1":"ON" }

					case "swap":
						fallthrough
					case "swapped":
						fallthrough
					case "invert":
						fallthrough
					case "inverted":
						p.Convert.Map = &ConvertMap{ "0":"ON", "1":"OFF" }
				}
				p.Convert.Binary = nil
			}

			if p.Convert.FloatMap != nil {
				p.Unit = ""
				p.Convert.FloatMap.Map = make(map[string]string)
				if p.Convert.FloatMap.Precision == 0 {
					p.Convert.FloatMap.Precision = 4
				}
				minFv := 1.0
				for k, v := range p.Convert.FloatMap.Values {
					var fv float64
					fv, err = strconv.ParseFloat(k, 64)
					if err != nil {
						p.Valid = false
						break
					}
					if fv < minFv {
						minFv = fv
					}
					k = strconv.FormatFloat(fv, 'f', p.Convert.FloatMap.Precision, 32)
					p.Convert.FloatMap.Map[k] = v
				}
				p.Convert.FloatMap.DefaultZero = strconv.FormatFloat(minFv, 'f', p.Convert.FloatMap.Precision, 32)
			}

			if p.Convert.Blob != nil {
				// p.Unit = ""
				// var blob []ConvertBlobData
				for _, b := range p.Convert.Blob.Order {
					switch {
						case b.Data != nil:
							if b.Data.Convert != nil {
								if b.Data.Convert.Alias != nil {
									foo := pm.Aliases.Get(b.Data.Convert.Alias)
									b.Data.Convert = &foo
								}
							}
							if b.Data.Key == "" {
								b.Data.Key = "%d"
							}
							p.Convert.Blob.Sequence = append(p.Convert.Blob.Sequence, *b.Data)
						case b.Array != nil:
							if b.Array.Data.Convert != nil {
								if b.Array.Data.Convert.Alias != nil {
									foo := pm.Aliases.Get(b.Array.Data.Convert.Alias)
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
									p.Convert.Blob.Sequence = append(p.Convert.Blob.Sequence, ConvertBlobData {
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
								p.Convert.Blob.Sequence = append(p.Convert.Blob.Sequence, ConvertBlobData {
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
				p.Convert.Blob.Order = ConvertBlobOrder{}
			}

			pm.PointsMap[n] = p
		}
	}

	return pm.PointsMap, err
}

func (a *Aliases) Append(b Aliases) *Aliases {
	for k, v := range b {
		(*a)[k] = v
	}
	return a
}

func (pm *PointsMap) Append(b PointsMap) *PointsMap {
	for k, v := range b {
		(*pm)[k] = v
	}
	return pm
}


func (p *Point) CorrectUnit(unit string) *Point {
	for range Only.Once {
		if p == nil {
			return nil
		}
		if p.Unit != "" {
			break
		}
		p.Unit = unit
	}
	return p
}

func JoinStrings(args ...string) string {
	return strings.TrimSpace(strings.Join(args, " "))
}

func JoinStringsForId(args ...string) string {
	var ret string

	for range Only.Once {
		var newargs []string
		var re = regexp.MustCompile(`(/| |:|\.)+`)
		var re2 = regexp.MustCompile(`^(-|_)+`)
		var re3 = regexp.MustCompile(`(-|_)+$`)

		for _, a := range args {
			if a == "" {
				continue
			}

			a = strings.TrimSpace(a)
			a = re.ReplaceAllString(a, `_`)
			a = re2.ReplaceAllString(a, ``)
			a = re3.ReplaceAllString(a, ``)
			// a = strings.TrimPrefix(a, `-`)
			// a = strings.TrimPrefix(a, `_`)
			// a = strings.TrimSuffix(a, `-`)
			// a = strings.TrimSuffix(a, `_`)
			newargs = append(newargs, a)
		}

		ret =  strings.Join(newargs, "-")
	}
	return ret
}


type ConvertStruct struct {
	Alias     *ConvertAlias     `json:"alias"`
	Increment *ConvertIncrement `json:"increment"`
	Range     *ConvertRange     `json:"range"`
	Map       *ConvertMap       `json:"map"`
	BitMap    *ConvertBitMap    `json:"bit_map"`
	Function  *ConvertFunction  `json:"function"`
	Binary    *ConvertBinary    `json:"binary"`
	String    *ConvertString    `json:"string"`
	Asset     *ConvertAsset     `json:"asset"`
	Array     *ConvertArray     `json:"array"`
	FloatMap  *ConvertFloatMap  `json:"float_map"`
	Integer   *ConvertInteger   `json:"integer"`
	Blob      *ConvertBlob      `json:"blob"`
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
					ret[Single] = ToRange(fmt.Sprintf("%v", value), c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
					break

				case c.Map != nil:
					val := fmt.Sprintf("%v", value)
					if v, ok := (*c.Map)[val]; ok {
						ret[Single] = v
					}
					break

				case c.BitMap != nil:
					foo := fmt.Sprintf("%v", value)
					ret[Single] = ToBitMap(foo, *c.BitMap, 0)
					break

				case c.Function != nil:
					if *c.Function == "log" {
						ret[Single] = ToLogFunc(fmt.Sprintf("%v", value), 1)
						break
					}
					break

				case c.Binary != nil:
					ret[Single] = ToBitMap(fmt.Sprintf("%v", value), *c.BitMap, 0)
					break

				case c.String != nil:
					break

				case c.Asset != nil:
					break

				case c.Array != nil:
					// 	value = strings.Join(*c.Array, ", ")
					break

				case c.FloatMap != nil:
					ret[Single] = c.FloatMap.Get(fmt.Sprintf("%v", value))
					break

				case c.Blob != nil:
					var val []byte
					for _, v := range value.([]byte) {
						val = append(val, v)
					}
					ret = c.Blob.Get(val)
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
				ret = ToRange(fmt.Sprintf("%v", value), c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
				break

			case c.Map != nil:
				val := fmt.Sprintf("%v", value)
				if v, ok := (*c.Map)[val]; ok {
					ret = v
				}
				break

			case c.BitMap != nil:
				foo := fmt.Sprintf("%v", value)
				ret = ToBitMap(foo, *c.BitMap, 0)
				break

			case c.Function != nil:
				if *c.Function == "log" {
					ret = ToLogFunc(fmt.Sprintf("%v", value), 1)
					break
				}
				break

			case c.Binary != nil:
				ret = ToBitMap(fmt.Sprintf("%v", value), *c.BitMap, 0)
				break

			case c.String != nil:
				break

			case c.Asset != nil:
				break

			case c.Array != nil:
				// 	value = strings.Join(*c.Array, ", ")
				break

			case c.FloatMap != nil:
				ret = c.FloatMap.Get(fmt.Sprintf("%v", value))
				break

			case c.Blob != nil:
				// Can't have a blob on a blob.
		}
	}

	return ret
}


func ToBitMap(value string, array []string, size uint32) string {
	for range Only.Once {
		if len(array) == 0 {
			break
		}

		if value == "" {
			value = array[0]
			break
		}

		iv, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			break
		}

		if iv == 0 {
			value = array[0]
			break
		}

		if size == 0 {
			size = uint32(len(array)) - 1
		}

		var elems []string
		for j := uint32(0); j < size; j++ {
			if iv & (1 << j) != 0 {
				elems = append(elems, array[j+1])
			}
		}

		value = strings.Join(elems, ", ")
	}

	return value
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

		value = strconv.FormatFloat(float64(d), 'f', precision, 32)
	}

	return value
}

func ToRange(value string, inMin float64, inMax float64, outMin float64, outMax float64, precision int) string {
	for range Only.Once {
		var err error
		var fv float64

		fv, err = strconv.ParseFloat(value, 64)
		if err != nil {
			break
		}

		type mapRange 	func(float64) (float64, error)

		var foo mapRange
		foo = MapRange64(RangeFloat64{inMin, inMax}, RangeFloat64{outMin, outMax})
		// Convert to linear scale.
		fv, err = foo(fv)

		if err != nil {
			break
		}

		if precision == 0 {
			precision = 1
		}

		value = strconv.FormatFloat(fv, 'f', precision, 32)
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

func (array *ConvertFloatMap) Get(value string) string {
	for range Only.Once {
		if array == nil {
			break
		}

		if len(array.Values) == 0 {
			break
		}

		if len(array.Map) == 0 {
			break
		}

		if value == "" {
			// value = array.FloatValues[0]
			break
		}

		fv, err := strconv.ParseFloat(value, 64)
		if err != nil {
			break
		}

		value = strconv.FormatFloat(fv, 'f', array.Precision, 32)

		// if fv == 0 {
		// 	// value = array.FloatValues[0]
		// 	break
		// }

		if v, ok := array.Map[value]; ok {
			value = v
			break
		}
	}

	return value
}


// func (blob *ConvertBlobData) Get(reader *bytes.Reader) string {
// 	var ret string
//
// 	for range Only.Once {
// 		var bo binary.ByteOrder
// 		if blob.BigEndian {
// 			bo = binary.BigEndian
// 		} else {
// 			bo = binary.LittleEndian
// 		}
//
// 		var err error
// 		switch blob.Type {
// 			case "int32":
// 				v := int32(0)
// 				err = binary.Read(reader, bo, &v)
// 				ret = fmt.Sprintf("int32: %d", v)
// 			case "int64":
// 				v := int64(0)
// 				err = binary.Read(reader, bo, &v)
// 				ret = fmt.Sprintf("int64: %d", v)
// 			case "float32":
// 				v := float32(0)
// 				err = binary.Read(reader, bo, &v)
// 				ret = fmt.Sprintf("float32: %f", v)
// 			case "float64":
// 				v := float64(0)
// 				err = binary.Read(reader, bo, &v)
// 				ret = fmt.Sprintf("float64: %f", v)
// 		}
//
// 		if err != nil {
// 			ret = fmt.Sprintf(`"%s": "Error: %s"`, blob.Key, err)
// 			break
// 		}
//
// 		ret = fmt.Sprintf(`"%s": "%s"`, blob.Key, ret)
// 	}
//
// 	return ret
// }

func (blob *ConvertBlobData) GetValue(reader *bytes.Reader) (string, error) {
	var ret string
	var err error

	for range Only.Once {
		var bo binary.ByteOrder
		if blob.BigEndian {
			bo = binary.BigEndian
		} else {
			bo = binary.LittleEndian
		}

		switch blob.Type {
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

		if blob.Convert == nil {
			break
		}
		ret = blob.Convert.GetValue(ret)
	}

	return ret, err
}

func (blob *ConvertBlobData) GetString(reader *bytes.Reader) string {
	var ret string

	for range Only.Once {
		var err error
		ret, err = blob.GetValue(reader)
		if err != nil {
			ret = fmt.Sprintf(`"%s": "Error with value (big_endian:%v) of type %s - %s"`,
				blob.Key, blob.BigEndian, blob.Type, err)
		}

		// ret = fmt.Sprintf(`"%s": "%s"`, blob.Key, ret)
	}

	return ret
}

func (blob *ConvertBlobDataArray) GetStrings(reader *bytes.Reader) []string {
	var ret []string

	for range Only.Once {

		// Using string keys instead of counter.
		if len(blob.Keys) > 0 {
			for _, i := range blob.Keys {
				b := ConvertBlobData {
					Type:      blob.Data.Type,
					Key:       i,
					BigEndian: blob.Data.BigEndian,
				}
				ret = append(ret, b.GetString(reader))
			}
			// ret = strings.Join(ret, ",")
			break
		}

		// Using counter instead of string keys.
		for i := 0; i < blob.Count; i++ {
			b := ConvertBlobData {
				Type:      blob.Data.Type,
				Key:       fmt.Sprintf(blob.Data.Key, i + blob.Offset),
				BigEndian: blob.Data.BigEndian,
			}
			ret = append(ret, b.GetString(reader))
		}
		// ret = strings.Join(ret, ",")
	}

	return ret
}

func (blob *ConvertBlob) Get(value []byte) map[string]string {
	ret := make(map[string]string)

	for range Only.Once {
		if blob == nil {
			break
		}

		if len(blob.Sequence) == 0 {
			break
		}

		r := bytes.NewReader(value)

		var i int
		for i = 0; (r.Len() > 0) && (i < len(blob.Sequence)); i++ {
			// fmt.Printf("Array[%d]: %d\n", i, r.Len())

			v := blob.Sequence[i].GetString(r)
			ret[blob.Sequence[i].Key] = v

			// switch {
			// 	case blob.Order[i].Data != nil:
			// 		v := blob.Order[i].Data.GetString(r)
			// 		// ret = append(ret, v)
			// 		ret[blob.Order[i].Data.Key] = v
			//
			// 	case blob.Order[i].Array != nil:
			// 		// v := (*blob)[i].Array.GetStrings(r)
			// 		// ret = append(ret, v...)
			// 		for _, v := range blob.Order[i].Array.GetStrings(r) {
			// 			ret[blob.Order[i].Array.Data.Key] = v
			// 		}
			// }
		}

		i = len(blob.Sequence) - i
		// fmt.Printf("Array[END]: %d / %d\n", r.Len(), i)
		if i > 0 {
			fmt.Printf("Expected %d more array elements.\n", i)
			// ret = append(ret, fmt.Sprintf(`"expected_count": "%d"`, i))
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
			// ret = append(ret, fmt.Sprintf(`"remaining_bytes": "%s"`, remBytes))
			ret["remaining_bytes"] = fmt.Sprintf("%s", remBytes)
		}

		// ret = "{" + strings.Join(retArray, ",") + "}"
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
