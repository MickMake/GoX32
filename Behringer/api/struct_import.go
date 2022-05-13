package api

import (
	"errors"
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

	UnitArray = "array"
	UnitState = "state"
	UnitToggle = "toggle"
	UnitToggleInvert = "toggle-invert"
	UnitSourceSelect = "source-select"
	UnitOutputSelect = "output-select"
	UnitFilterTypeSelect = "filter-type-select"
	UnitColourSelect = "colour-select"
	UnitIconSelect = "icon-select"
	UnitRatioSelect = "ratio-select"
	UnitEqModeSelect = "eq-type-select"
	UnitRecPosSelect = "rec-pos-select"
	UnitMonitorSourceSelect = "monitor-source-select"
	UnitString = "string"
)


// type StringMap map[string]string
type Aliases map[ConvertAlias]ConvertStruct

func (s *Aliases) Get(selector *ConvertAlias) ConvertStruct {
	if selector == nil {
		return ConvertStruct{}
	}
	if ret, ok := (*s)[*selector]; ok {
		return ret
	}
	return ConvertStruct{}
}


type PointsMapFile struct {
	Aliases   Aliases   `json:"aliases"`
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
				break
			}
			pm.Aliases.Append(pmi.Aliases)
			pm.PointsMap.Append(pmi.PointsMap)
		}

		for n, p := range pm.PointsMap {
			if n == "" {
				delete(pm.PointsMap, n)
				continue
			}
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
			if p.Convert.Map != nil {
				// p.Unit = UnitState
				// @TODO = States are binary.
			}

			if p.Convert.Alias != nil {
				p.Convert = pm.Aliases.Get(p.Convert.Alias)
			}

			if p.Convert.Range != nil {
				if (p.Convert.Range.InMin == 0) && (p.Convert.Range.InMax == 0) {
					p.Convert.Range.InMin = 0
					p.Convert.Range.InMax = 1
				}
			}

			// switch {
			// 	// case sm.Map != nil:
			// 	// 	p.Convert.Map = sm.Map
			// 	// 	p.Unit = ""
			// 	//
			// 	// case p.Unit == UnitToggle:
			// 	// 	p.Convert.Map = &ConvertMap {
			// 	// 		"0": "off",
			// 	// 		"1": "on",
			// 	// 	}
			// 	// 	p.Unit = "state"
			// 	//
			// 	// case p.Unit == UnitToggleInvert:
			// 	// 	p.Convert.Map = &ConvertMap {
			// 	// 		"0": "on",
			// 	// 		"1": "off",
			// 	// 	}
			// 	// 	p.Unit = "state"
			// 	//
			// 	// case p.Unit == UnitSourceSelect:
			// 	// 	p.States = map[string]string{
			// 	// 		"0":"Self",
			// 	// 		"1":"Ch 01", "2":"Ch 02", "3":"Ch 03", "4":"Ch 04", "5":"Ch 05", "6":"Ch 06", "7":"Ch 07", "8":"Ch 08", "9":"Ch 09", "10":"Ch 10", "11":"Ch 11", "12":"Ch 12", "13":"Ch 13", "14":"Ch 14", "15":"Ch 15", "16":"Ch 16", "17":"Ch 17", "18":"Ch 18", "19":"Ch 19", "20":"Ch 20", "21":"Ch 21", "22":"Ch 22", "23":"Ch 23", "24":"Ch 24", "25":"Ch 25", "26":"Ch 26", "27":"Ch 27", "28":"Ch 28", "29":"Ch 29", "30":"Ch 30", "31":"Ch 31", "32":"Ch 32",
			// 	// 		"33":"Aux 01", "34":"Aux 02", "35":"Aux 03", "36":"Aux 04", "37":"Aux 05", "38":"Aux 06", "39":"Aux 07", "40":"Aux 08",
			// 	// 		"41":"Fx 1L", "42":"Fx 1R", "43":"Fx 2L", "44":"Fx 2R", "45":"Fx 3L", "46":"Fx 3R", "47":"Fx 4L", "48":"Fx 4R",
			// 	// 		"49":"Bus 01", "50":"Bus 02", "51":"Bus 03", "52":"Bus 04", "53":"Bus 05", "54":"Bus 06", "55":"Bus 07", "56":"Bus 08", "57":"Bus 09", "58":"Bus 10", "59":"Bus 11", "60":"Bus 12", "61":"Bus 13", "62":"Bus 14", "63":"Bus 15", "64":"Bus 16",
			// 	// 	}
			// 	// 	p.Unit = ""
			// 	//
			// 	// case p.Unit == UnitOutputSelect:
			// 	// 	p.States = map[string]string{
			// 	// 		"0":"OFF",
			// 	// 		"1":"Main L","2":"Main R","3":"Main C/M",
			// 	// 		"4":"MixBus 1","5":"MixBus 2","6":"MixBus 3","7":"MixBus 4",
			// 	// 		"8":"MixBus 5","9":"MixBus 6","10":"MixBus 7","11":"MixBus 8",
			// 	// 		"12":"MixBus 9","13":"MixBus 10","14":"MixBus 11","15":"MixBus 12",
			// 	// 		"16":"MixBus 13","17":"MixBus 14","18":"MixBus 15","19":"MixBus 16",
			// 	// 		"20":"Matrix 1","21":"Matrix 2","22":"Matrix 3","23":"Matrix 4","24":"Matrix 5","25":"Matrix 6",
			// 	// 		"26":"DirOut Ch 1","27":"DirOut Ch 2","28":"DirOut Ch 3","29":"DirOut Ch 4",
			// 	// 		"30":"DirOut Ch 5","31":"DirOut Ch 6","32":"DirOut Ch 7","33":"DirOut Ch 8",
			// 	// 		"34":"DirOut Ch 9","35":"DirOut Ch 10","36":"DirOut Ch 11","37":"DirOut Ch 12",
			// 	// 		"38":"DirOut Ch 13","39":"DirOut Ch 14","40":"DirOut Ch 15","41":"DirOut Ch 16",
			// 	// 		"42":"DirOut Ch 17","43":"DirOut Ch 18","44":"DirOut Ch 19","45":"DirOut Ch 20",
			// 	// 		"46":"DirOut Ch 21","47":"DirOut Ch 22","48":"DirOut Ch 23","49":"DirOut Ch 24",
			// 	// 		"50":"DirOut Ch 25","51":"DirOut Ch 26","52":"DirOut Ch 27","53":"DirOut Ch 28",
			// 	// 		"54":"DirOut Ch 29","55":"DirOut Ch 30","56":"DirOut Ch 31","57":"DirOut Ch 32",
			// 	// 		"58":"DirOut Aux 1","59":"DirOut Aux 2","60":"DirOut Aux 3","61":"DirOut Aux 4",
			// 	// 		"62":"DirOut Aux 5","63":"DirOut Aux 6","64":"DirOut Aux 7","65":"DirOut Aux 8",
			// 	// 		"66":"DirOut FX 1L","67":"DirOut FX 1R","68":"DirOut FX 2L","69":"DirOut FX 2R",
			// 	// 		"70":"DirOut FX 3L","71":"DirOut FX 3R","72":"DirOut FX 4L","73":"DirOut FX 4R",
			// 	// 		"74":"Monitor L","75":"Monitor R","76":"Talkback",
			// 	// 	}
			// 	// 	p.Unit = ""
			// 	//
			// 	// case p.Unit == UnitFilterTypeSelect:
			// 	// 	p.States = map[string]string{
			// 	// 		"0": "LC6", "1": "LC12", "2": "HC6", "3": "HC12", "4": "1","5": "2", "6": "3","7": "5", "8": "10",
			// 	// 	}
			// 	// 	p.Unit = ""
			// 	//
			// 	// case p.Unit == UnitColourSelect:
			// 	// 	p.States = map[string]string{
			// 	// 		"0": "BLACK", "1": "RED", "2": "GREEN", "3": "YELLOW", "4": "BLUE","5": "PINK", "6": "CYAN","7": "WHITE",
			// 	// 	}
			// 	// 	p.Unit = ""
			// 	//
			// 	// case p.Unit == UnitIconSelect:
			// 	// 	// @TODO - Fetch icons.
			// 	// 	// p.States = map[string]string{
			// 	// 	// 	"0": "BLACK", "1": "RED", "2": "GREEN", "3": "YELLOW", "4": "BLUE","5": "PINK", "6": "CYAN","7": "WHITE",
			// 	// 	// }
			// 	// 	p.Unit = ""
			// 	//
			// 	// case p.Unit == UnitRatioSelect:
			// 	// 	p.States = map[string]string{
			// 	// 		"0":"1.1:1","1":"1.3:1","2":"1.5:1","3":"2:1","4":"2.5:1","5":"3:1","6":"4:1","7":"5:1","8":"7:1","9":"10:1","10":"20:1","11":"100:1",
			// 	// 	}
			// 	// 	p.Unit = ""
			// 	//
			// 	// case p.Unit == UnitEqModeSelect:
			// 	// 	p.States = map[string]string{
			// 	// 		"0":"LCut","1":"LShv","2":"PEQ","3":"VEQ","4":"HShv","5":"HCut",
			// 	// 	}
			// 	// 	p.Unit = ""
			// 	//
			// 	// case p.Unit == UnitRecPosSelect:
			// 	// 	p.States = map[string]string{
			// 	// 		"0":"IN/Low Cut",
			// 	// 		"1":"IN/Low Cut +M",
			// 	// 		"2":"Pre EQ",
			// 	// 		"3":"Pre EQ +M",
			// 	// 		"4":"Post EQ",
			// 	// 		"5":"Post EQ +M",
			// 	// 		"6":"Pre Fader",
			// 	// 		"7":"Pre Fader +M",
			// 	// 		"8":"Post Fader",
			// 	// 	}
			// 	// 	p.Unit = ""
			// 	//
			// 	// case p.Unit == UnitMonitorSourceSelect:
			// 	// 	p.States = map[string]string{
			// 	// 		"0":"Off",
			// 	// 		"1":"LR Bus",
			// 	// 		"2":"LR+M/C",
			// 	// 		"3":"LR PFL",
			// 	// 		"4":"LR AFL",
			// 	// 		"5":"Aux 5/6",
			// 	// 		"6":"Aux 7/8",
			// 	// 	}
			// 	// 	p.Unit = ""
			//
			// 	case p.Unit == UnitString:
			// 		p.Unit = ""
			// }

			p.Valid = true
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

func (a *PointsMap) Append(b PointsMap) *PointsMap {
	for k, v := range b {
		(*a)[k] = v
	}
	return a
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
	var newargs []string
	var re = regexp.MustCompile(`(/| |:|\.)+`)
	var re2 = regexp.MustCompile(`^(-|_)+`)
	for _, a := range args {
		if a == "" {
			continue
		}
		a = strings.TrimSpace(a)
		a = re.ReplaceAllString(a, `_`)
		a = re2.ReplaceAllString(a, ``)
		a = strings.TrimPrefix(a, `-`)
		a = strings.TrimPrefix(a, `_`)
		a = strings.TrimSuffix(a, `-`)
		a = strings.TrimSuffix(a, `_`)
		newargs = append(newargs, a)
	}
	// return strings.ReplaceAll(strings.TrimSpace(strings.Join(args, ".")), ".", "_")
	return strings.Join(newargs, "-")
}


// type ConvFunc    func(ConvertStruct) float64

type ConvertStruct struct {
	// Min 		float64
	// Max 		float64
	// Linear		bool
	// Increment	float64
	// MapRange 	func(float64) (float64, error)

	Increment  *ConvertIncrement `json:"increment"`
	Range         *ConvertRange `json:"range"`
	Map   *ConvertMap   `json:"map"`
	Alias *ConvertAlias `json:"alias"`
	Function *ConvertFunction `json:"function"`
	BitMap   *ConvertBitMap   `json:"bit_map"`
}

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

type ConvertAlias string

type ConvertFunction string

func (c *ConvertStruct) Get(value string) string {
	for range Only.Once {
		switch {
			case c.Increment != nil:
				break

			case c.Range != nil:
				value = ToLinDb(value, c.Range.InMin, c.Range.InMax, c.Range.OutMin, c.Range.OutMax, c.Range.Precision)
				break

			case c.Map != nil:
				if v, ok := (*c.Map)[value]; ok {
					value = v
				}
				break

			case c.Alias != nil:
				break

			case c.BitMap != nil:
				value = ToBitMap(value, *c.BitMap)
				break

			case c.Function != nil:
				if *c.Function == "log" {
					value = ToLogDb(value, 1)
					break
				}
				break
		}
	}
	return value
}


func ToBitMap(value string, array []string) string {
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

		var elems []string
		for j := 0; j < 8; j++ {
			if iv & (1 << byte(j)) != 0 {
				elems = append(elems, array[j+1])
			}
		}

		value = strings.Join(elems, ", ")
	}

	return value
}

func ToLogDb(value string, precision int) string {
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
		var d float32
		if fv >= 0.5 {
			d = float32(fv * 40.0 - 30.0)

		} else if fv >= 0.25 {
			d = float32(fv * 80.0 - 50.0)

		} else if fv >= 0.0625 {
			d = float32(fv * 160.0 - 70.0)

		} else if fv >= 0.0 {
			d = float32(fv * 480.0 - 90.0)
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

func ToLinDb(value string, inMin float64, inMax float64, outMin float64, outMax float64, precision int) string {
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

func ToLinDbString(value string, inMin string, inMax string, outMin string, outMax string) string {
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

		outMinFloat, err = strconv.ParseFloat(inMin, 64)
		if err != nil {
			break
		}

		outMaxFloat, err = strconv.ParseFloat(inMax, 64)
		if err != nil {
			break
		}

		outMinFloat, err = strconv.ParseFloat(outMin, 64)
		if err != nil {
			break
		}

		outMaxFloat, err = strconv.ParseFloat(outMax, 64)
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
