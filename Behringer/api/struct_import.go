package api

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Behringer/api/output"
	"github.com/MickMake/GoX32/Only"
	"regexp"
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

			switch {
				case p.Convert.Alias != nil:
					p.Convert = pm.Aliases.Get(p.Convert.Alias)

				case p.Convert.Map != nil:
					err = p.Convert.Map.Import()
					if err != nil {
						break
					}
					if p.Unit == "" {
						p.Unit = "-"
					}

				case p.Convert.Range != nil:
					err = p.Convert.Range.Import()
					if err != nil {
						break
					}

				case p.Convert.BitMap != nil:
					err = p.Convert.BitMap.Import()
					if err != nil {
						break
					}
					if p.Unit == "" {
						p.Unit = "-"
					}

				case p.Convert.Asset != nil:
					err = p.Convert.Asset.Import()
					if err != nil {
						break
					}
					if p.Unit == "" {
						p.Unit = "-"
					}

				case p.Convert.Binary != nil:
					err = p.Convert.Binary.Import()
					if err != nil {
						break
					}
					if p.Unit == "" {
						p.Unit = "-"
					}

				case p.Convert.FloatMap != nil:
					err = p.Convert.FloatMap.Import()
					if err != nil {
						break
					}
					// p.Convert.FloatMap.Map = make(map[string]string)
					// if p.Convert.FloatMap.Precision == 0 {
					// 	p.Convert.FloatMap.Precision = 4
					// }
					// minFv := 1.0
					// for k, v := range p.Convert.FloatMap.Values {
					// 	var fv float64
					// 	fv, err = strconv.ParseFloat(k, 64)
					// 	if err != nil {
					// 		p.Valid = false
					// 		break
					// 	}
					// 	if fv < minFv {
					// 		minFv = fv
					// 	}
					// 	k = strconv.FormatFloat(fv, 'f', p.Convert.FloatMap.Precision, 32)
					// 	p.Convert.FloatMap.Map[k] = v
					// }
					// p.Convert.FloatMap.DefaultZero = strconv.FormatFloat(minFv, 'f', p.Convert.FloatMap.Precision, 32)

				case p.Convert.Blob != nil:
					err = p.Convert.Blob.Import(pm.Aliases)
					if err != nil {
						break
					}
					// for _, b := range p.Convert.Blob.Order {
					// 	switch {
					// 	case b.Data != nil:
					// 		if b.Data.Convert != nil {
					// 			if b.Data.Convert.Alias != nil {
					// 				foo := pm.Aliases.Get(b.Data.Convert.Alias)
					// 				b.Data.Convert = &foo
					// 			}
					// 		}
					// 		if b.Data.Key == "" {
					// 			b.Data.Key = "%d"
					// 		}
					// 		p.Convert.Blob.Sequence = append(p.Convert.Blob.Sequence, *b.Data)
					// 	case b.Array != nil:
					// 		if b.Array.Data.Convert != nil {
					// 			if b.Array.Data.Convert.Alias != nil {
					// 				foo := pm.Aliases.Get(b.Array.Data.Convert.Alias)
					// 				b.Array.Data.Convert = &foo
					// 			}
					// 		}
					//
					// 		if b.Array.Keys != nil {
					// 			for _, v := range b.Array.Keys {
					// 				if b.Array.Data.Key == "" {
					// 					if len(b.Array.Keys) > 0 {
					// 						b.Array.Data.Key = "%s"
					// 					} else {
					// 						b.Array.Data.Key = "%d"
					// 					}
					// 				}
					// 				p.Convert.Blob.Sequence = append(p.Convert.Blob.Sequence, ConvertBlobData {
					// 					Convert:   b.Array.Data.Convert,
					// 					Unit:      b.Array.Data.Unit,
					// 					Key:       fmt.Sprintf("%s", v),
					// 					Type:      b.Array.Data.Type,
					// 					BigEndian: b.Array.Data.BigEndian,
					// 				})
					// 			}
					// 			continue
					// 		}
					//
					// 		for i := 0; i < b.Array.Count; i++ {
					// 			if b.Array.Data.Key == "" {
					// 				if len(b.Array.Keys) > 0 {
					// 					b.Array.Data.Key = "%s"
					// 				} else {
					// 					b.Array.Data.Key = "%d"
					// 				}
					// 			}
					// 			p.Convert.Blob.Sequence = append(p.Convert.Blob.Sequence, ConvertBlobData {
					// 				Convert:   b.Array.Data.Convert,
					// 				Unit:      b.Array.Data.Unit,
					// 				Key:       fmt.Sprintf(b.Array.Data.Key, i + b.Array.Offset),	// , b.Array.Data.Unit, v),
					// 				Type:      b.Array.Data.Type,
					// 				BigEndian: b.Array.Data.BigEndian,
					// 			})
					// 		}
					// 		continue
					// 	}
					// }
					// p.Convert.Blob.Order = ConvertBlobOrder{}

				case p.Convert.String != nil:
					err = p.Convert.String.Import()
					if err != nil {
						break
					}
					if p.Unit == "" {
						p.Unit = "-"
					}

				case p.Convert.Increment != nil:
					err = p.Convert.Increment.Import()
					if err != nil {
						break
					}

				case p.Convert.Function != nil:
					err = p.Convert.Function.Import()
					if err != nil {
						break
					}

				case p.Convert.Array != nil:
					err = p.Convert.Array.Import()
					if err != nil {
						break
					}
					if p.Unit == "" {
						p.Unit = "-"
					}

				case p.Convert.Integer != nil:
					err = p.Convert.Integer.Import()
					if err != nil {
						break
					}
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
			a = strings.ToLower(a)
			newargs = append(newargs, a)
		}

		ret =  strings.Join(newargs, "-")
	}

	return ret
}
