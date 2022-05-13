package api

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"sort"
	"strconv"
	"strings"
)


type UnitValue struct {
	Unit  string `json:"unit"`
	Value string `json:"value"`
	ValueFloat float64 `json:"value_float,omitempty"`
	ValueInt int64 `json:"value_int,omitempty"`
}
type UnitValues []UnitValue
type UnitValueMap map[string]UnitValue


func (u *UnitValueMap) Sort() []string {
	var ret []string
	for n := range *u {
		ret = append(ret, n)
	}
	sort.Strings(ret)
	return ret
}


func (ref *UnitValue) UnitValueFix() UnitValue {
	if ref.Unit == "W" {
		fvs, err := DivideByThousand(ref.Value)
		// fv, err := strconv.ParseFloat(p.Value, 64)
		// fv = fv / 1000
		if err == nil {
			// p.Value = fmt.Sprintf("%.3f", fv)
			ref.Value = fvs
			ref.Unit = "kW"
		}
	}

	if ref.Unit == "Wh" {
		fvs, err := DivideByThousand(ref.Value)
		// fv, err := strconv.ParseFloat(p.Value, 64)
		// fv = fv / 1000
		if err == nil {
			// p.Value = fmt.Sprintf("%.3f", fv)
			ref.Value = fvs
			ref.Unit = "kWh"
		}
	}

	ref.ValueFloat, _ = strconv.ParseFloat(ref.Value, 64)

	return *ref
}

func (ref *UnitValue) UnitValueToPoint(psId string, point string, name string) *Point {
	uv := ref.UnitValueFix()

	// u := ref.Unit
	//
	// if ref.Unit == "W" {
	// 	fvs, err := DivideByThousand(ref.Value)
	// 	// fv, err := strconv.ParseFloat(p.Value, 64)
	// 	// fv = fv / 1000
	// 	if err == nil {
	// 		// p.Value = fmt.Sprintf("%.3f", fv)
	// 		ref.Value = fvs
	// 		ref.Unit = "kW"
	// 	}
	// }
	//
	// if ref.Unit == "Wh" {
	// 	fvs, err := DivideByThousand(ref.Value)
	// 	// fv, err := strconv.ParseFloat(p.Value, 64)
	// 	// fv = fv / 1000
	// 	if err == nil {
	// 		// p.Value = fmt.Sprintf("%.3f", fv)
	// 		ref.Value = fvs
	// 		ref.Unit = "kWh"
	// 	}
	// }

	if name == "" {
		name = PointToName(point)
	}

	vt := GetPoint(psId, point)
	if !vt.Valid {
		vt = &Point {
			ParentId: psId,
			Id:       point,
			Name:     name,
			Unit:     uv.Unit,
			Type:     "PointTypeInstant",
			Valid:    true,
		}
	}

	return vt
}


func JsonToUnitValue(j string) UnitValue {
	var ret UnitValue

	for range Only.Once {
		err := json.Unmarshal([]byte(j), &ret)
		if err != nil {
			break
		}
	}

	return ret
}

func Float32ToString(num float64) string {
	s := fmt.Sprintf("%.6f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

func Float64ToString(num float64) string {
	s := fmt.Sprintf("%.6f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

func DivideByThousand(num string) (string, error) {
	fv, err := strconv.ParseFloat(num, 64)
	if err == nil {
		num = Float64ToString(fv / 1000)
	}
	return num, err
}
