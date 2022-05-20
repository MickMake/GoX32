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


func (u UnitValue) String() string {
	unit := u.Unit
	if unit == "-" {
		unit = ""
	}
	return fmt.Sprintf("%s%s", u.Value, unit)
}

func (u UnitValueMap) String() string {
	var ret string
	for range Only.Once {
		if len(u) == 1 {
			for _, n := range u {
				ret += fmt.Sprintf("%s", n)
			}
			break
		}
		for _, n := range u.Sort() {
			ret += fmt.Sprintf("%s:%s\t", n, u[n])
		}
	}
	return ret
}

func (u *UnitValueMap) Sort() []string {
	var ret []string
	for n := range *u {
		ret = append(ret, n)
	}
	sort.Strings(ret)
	return ret
}


func (u *UnitValueMap) GetFirst() *UnitValue {
	var ret UnitValue
	for range Only.Once {
		if len(*u) == 1 {
			for _, ret = range *u {
				break
			}
			break
		}
		for _, n := range u.Sort() {
			if n == "0" {
				ret = (*u)[n]
				break
			}
		}
	}
	return &ret
}

func (u *UnitValue) UnitValueFix() UnitValue {
	if u.Unit == "W" {
		fvs, err := DivideByThousand(u.Value)
		// fv, err := strconv.ParseFloat(p.Value, 64)
		// fv = fv / 1000
		if err == nil {
			// p.Value = fmt.Sprintf("%.3f", fv)
			u.Value = fvs
			u.Unit = "kW"
		}
	}

	if u.Unit == "Wh" {
		fvs, err := DivideByThousand(u.Value)
		// fv, err := strconv.ParseFloat(p.Value, 64)
		// fv = fv / 1000
		if err == nil {
			// p.Value = fmt.Sprintf("%.3f", fv)
			u.Value = fvs
			u.Unit = "kWh"
		}
	}

	u.ValueFloat, _ = strconv.ParseFloat(u.Value, 64)

	return *u
}

func (u *UnitValue) UnitValueToPoint(psId string, point string, name string) *Point {
	uv := u.UnitValueFix()

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
