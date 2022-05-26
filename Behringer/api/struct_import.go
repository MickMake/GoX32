package api

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Behringer/api/output"
	"github.com/MickMake/GoX32/Only"
	"time"
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

func (a *Aliases) Append(b Aliases) *Aliases {
	for k, v := range b {
		(*a)[k] = v
	}
	return a
}


type PointsMapFile struct {
	Refresh Refresh `json:"refresh"`
	Aliases Aliases `json:"aliases"`
	PointsMap      PointsMap      `json:"points"`
	PointsArrayMap PointsArrayMap `json:"points_array_map"`
}

type PointsArrayMap struct {
	Min       int       `json:"min"`
	Max       int       `json:"max"`
	Increment int       `json:"increment"`
	PointsMap PointsMap `json:"points"`
}

type Refresh struct {
	RefreshPointsArray RefreshPointsArray `json:"points_array_map"`
	RefreshPoints      RefreshPoints      `json:"points"`
	Delay              *int               `json:"delay"`
	BatchLimit         *int               `json:"batch_limit"`
}

type RefreshPointsArray struct {
	Delay     *int                    `json:"delay"`
	Min       int                     `json:"min"`
	Max       int                     `json:"max"`
	Increment int                     `json:"increment"`
	Points    map[string]RefreshPoint `json:"points"`
}

type RefreshPoints map[string]*RefreshPoint
// type RefreshPointsMap map[string]int

// type RefreshPointsMap []string
// type RefreshPointsMap map[string]RefreshPoint

type RefreshPoint struct {
	Delay *int `json:"delay"`
	When  time.Time `json:"-"`
}

func (rp *RefreshPoint) IsExpired() bool {
	var yes bool

	if rp.Delay == nil {
		d := 60
		rp.Delay = &d
	}
	then := rp.When.Add(time.Second * time.Duration(*rp.Delay))
	if then.Before(time.Now()) {
		yes = true
	}

	return yes
}

func (rp *RefreshPoint) Reset() {
	rp.When = time.Now()
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

			// Generate points from a min and max.
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

const DefaultBatchLimit = 200

func ImportRefresh(filenames ...string) (RefreshPoints, int, error) {
	pm := make(RefreshPoints)
	batchLimit := DefaultBatchLimit
	var err error

	for range Only.Once {
		for _, filename := range filenames {
			var pmi PointsMapFile
			err = output.FileRead(filename, &pmi)
			if err != nil {
				err = errors.New(fmt.Sprintf("Error reading points json file '%s': %s", filename, err))
				break
			}

			// if strings.HasSuffix(filename, "points_refresh.json") {
			// 	fmt.Sprintf("")
			// }

			if len(pmi.Refresh.RefreshPointsArray.Points) == 0 {
				continue
			}

			if pmi.Refresh.BatchLimit != nil {
				batchLimit = *pmi.Refresh.BatchLimit
			}

			now := time.Now()

			// Generate refresh points from an array map.
			for i := pmi.Refresh.RefreshPointsArray.Min; i <= pmi.Refresh.RefreshPointsArray.Max; i++ {
				for n, p := range pmi.Refresh.RefreshPointsArray.Points {
					if n == "" {
						continue
					}

					name := fmt.Sprintf(n, i)
					delay := 0
					switch {
						case p.Delay != nil:
							delay = *p.Delay
						case pmi.Refresh.RefreshPointsArray.Delay != nil:
							delay = *pmi.Refresh.RefreshPointsArray.Delay
						case pmi.Refresh.Delay != nil:
							delay = *pmi.Refresh.Delay
						default:
							delay = 60
					}
					pm[name] = &RefreshPoint { Delay: &delay, When: now.Add(- (time.Second * time.Duration(delay))) }

					// if p.Delay != nil {
					// 	pm[name] = RefreshPoint { Delay: p.Delay }
					// }
					// if pmi.Refresh.RefreshPointsArray.Delay != nil {
					// 	pm[name] = RefreshPoint { Delay: pmi.Refresh.RefreshPointsArray.Delay }
					// }
					// if pmi.Refresh.Delay != nil {
					// 	pm[name] = RefreshPoint { Delay: pmi.Refresh.Delay }
					// }
					// if pm[name].Delay == 0 {
					// 	pm[name] = 60
					// }
					//
					// pm[name] = RefreshPoint { Delay: &delay, When: now }
				}
			}

			for name, p := range pmi.Refresh.RefreshPoints {
				if name == "" {
					continue
				}

				delay := 0
				switch {
					case p.Delay != nil:
						delay = *p.Delay
					case pmi.Refresh.RefreshPointsArray.Delay != nil:
						delay = *pmi.Refresh.RefreshPointsArray.Delay
					case pmi.Refresh.Delay != nil:
						delay = *pmi.Refresh.Delay
					default:
						delay = 60
				}
				pm[name] = &RefreshPoint { Delay: &delay, When: now.Add(- (time.Second * time.Duration(delay))) }

				// if pmi.Refresh.Delay != nil {
				// 	pm[name] = *pmi.Refresh.Delay
				// }
				// if pmi.Refresh.RefreshPointsArray.Delay != nil {
				// 	pm[name] = *pmi.Refresh.RefreshPointsArray.Delay
				// }
				// if p.Delay != nil {
				// 	pm[name] = *p.Delay
				// }
				// if pm[name] == 0 {
				// 	pm[name] = 60
				// }
			}
		}
		if err != nil {
			break
		}
	}

	return pm, batchLimit, err
}
