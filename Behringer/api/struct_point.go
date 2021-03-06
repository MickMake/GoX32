package api

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"time"
)


type Point struct {
	EndPoint    string  `json:"endpoint"`
	FullId      string  `json:"full_id"`
	ParentId    string  `json:"parent_id"`
	Id          string  `json:"id"`
	GroupName 	string  `json:"group_name"`
	Name 		string	`json:"name"`
	Unit 		string	`json:"unit"`
	Type        string	`json:"type"`
	Valid       bool	`json:"valid"`
	Info        string	`json:"info"`

	Convert     ConvertStruct `json:"convert"`
}


func (p *Point) WhenReset() string {
	var ret string

	for range Only.Once {
		var err error
		now := time.Now()

		switch {
		case p.IsInstant():
			ret = ""

		case p.IsDaily():
			now, err = time.Parse("2006-01-02T15:04:05", now.Format("2006-01-02") + "T00:00:00")
			// ret = fmt.Sprintf("%d", now.Unix())
			ret = now.Format("2006-01-02T15:04:05") + ""

		case p.IsMonthly():
			now, err = time.Parse("2006-01-02T15:04:05", now.Format("2006-01") + "-01T00:00:00")
			ret = fmt.Sprintf("%d", now.Unix())
			ret = now.Format("2006-01-02T15:04:05") + ""

		case p.IsYearly():
			now, err = time.Parse("2006-01-02T15:04:05", now.Format("2006") + "-01-01T00:00:00")
			ret = fmt.Sprintf("%d", now.Unix())
			ret = now.Format("2006-01-02T15:04:05") + ""

		case p.IsTotal():
			// ret = "0"
			ret = "1970-01-01T00:00:00"

		default:
			// ret = "0"
			ret = "1970-01-01T00:00:00"
		}
		if err != nil {
			now := time.Now()
			ret = fmt.Sprintf("%d", now.Unix())
		}
	}

	return ret
}

func (p Point) String() string {
	return fmt.Sprintf("endpoint:%s\tid:%s\tname:%s\ttype:%s", p.EndPoint, p.FullId, p.Name, p.Convert.GetConvertType())
}

func (p Point) Json() string {
	j, _ := json.Marshal(p)
	return string(j)
}

func (p Point) IsInstant() bool {
	if p.Type == PointTypeInstant {
		return true
	}
	return false
}

func (p Point) IsDaily() bool {
	if p.Type == PointTypeDaily {
		return true
	}
	return false
}

func (p Point) IsMonthly() bool {
	if p.Type == PointTypeMonthly {
		return true
	}
	return false
}

func (p Point) IsYearly() bool {
	if p.Type == PointTypeYearly {
		return true
	}
	return false
}

func (p Point) IsTotal() bool {
	if p.Type == PointTypeTotal {
		return true
	}
	return false
}

func (p *Point) IsBinaryMap() bool {
	var ok bool
	for range Only.Once {
		if p.Convert.Binary == nil {
			break
		}
		if (p.Convert.Binary.Off != Off) && (p.Convert.Binary.Off != On) {
			ok = true
			break
		}
		if (p.Convert.Binary.On != Off) && (p.Convert.Binary.On != On) {
			ok = true
			break
		}
	}
	return ok
}

func (p *Point) IsBinary() bool {
	var ok bool
	for range Only.Once {
		if p.Convert.Binary == nil {
			break
		}
		ok = true
	}
	return ok
}

func (p *Point) IsSwitch() bool {
	var ok bool
	for range Only.Once {
		if p.Convert.Binary == nil {
			break
		}
		if p.Convert.Binary.IsSwitch {
			ok = true
			break
		}
	}
	return ok
}

func (p *Point) IsMomentary() bool {
	var ok bool
	for range Only.Once {
		if p.Convert.Binary == nil {
			break
		}
		if p.Convert.Binary.IsMomentary {
			ok = true
			break
		}
	}
	return ok
}

func (p *Point) IsIndex() bool {
	var ok bool
	for range Only.Once {
		if p.Convert.Index != nil {
			ok = true
			break
		}
	}
	return ok
}

func (p *Point) GetIndexOptions() []string {
	var ret []string
	for range Only.Once {
		if p.Convert.Index != nil {
			ret = p.Convert.Index.GetOptions()
			break
		}
		if p.Convert.Map != nil {
			ret = p.Convert.Map.GetOptions()
			break
		}
		if p.Convert.BitMap != nil {
			ret = p.Convert.BitMap.GetOptions()
			break
		}
		if p.Convert.Binary != nil {
			ret = []string{ p.Convert.Binary.Off, p.Convert.Binary.On}
			break
		}
	}
	return ret
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
