package apiTypes

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"reflect"
	"strconv"
)


const (
	TypeString = 0
	TypeByte = 1
	TypeInt32 = 2
	TypeInt64 = 3
	TypeEnum = 4
	TypeBool = 5
	TypeFloat32 = 6
	TypeFloat64 = 7
	TypeBit = 8
	TypeTime = 9

	StrString = "s"
	StrByte = "b"
	StrInt32 = "i"
	StrInt64 = "h"
	StrEnum = "i"
	StrBool = "i"           // Integer type with only two values - keep it simple.
	StrFloat32 = "f"
	StrFloat64 = "d"
	StrBit = "i"
	StrTime = "t"
)

type CmdType struct {
	Name 		string
	// Description string

	Type  		int8
	TypeString	string
	// Value		string
	Units		string

	BitValue
	BoolValue
	EnumValue
	Float32Value
	Float64Value
	Int32Value
	Int64Value
	StringValue
}

type CmdRef map[string]*CmdType

func (ct *CmdType) get() (string, error) {
	var v string
	var err error

	for range Only.Once {
		switch ct.Type {
			case TypeBit:
				v, err = ct.BitValue.getString()

			case TypeBool:
				v, err = ct.BoolValue.getString()

			case TypeEnum:
				v, err = ct.EnumValue.getString()

			case TypeFloat32:
				v, err = ct.Float32Value.getString()

			case TypeFloat64:
				v, err = ct.Float64Value.getString()

			case TypeInt32:
				v, err = ct.Int32Value.getString()

			case TypeInt64:
				v, err = ct.Int64Value.getString()

			case TypeString:
				v, err = ct.StringValue.getString()

			default:
				err = errors.New("Mismatched value type: " + string(ct.Type))
				break
		}
	}

	return v, err
}

func (ct *CmdType) getReal() (string, string, string, error) {
	var r string	// Real value.
	var u string	// Units.
	var s string	// Stored value.
	var err error

	for range Only.Once {
		switch ct.Type {
			case TypeBit:
				r, s, err = ct.BitValue.getReal()

			case TypeBool:
				r, s, err = ct.BoolValue.getReal()

			case TypeEnum:
				r, s, err = ct.EnumValue.getReal()

			case TypeFloat32:
				r, s, err = ct.Float32Value.getReal()

			case TypeFloat64:
				r, s, err = ct.Float64Value.getReal()

			case TypeInt32:
				r, s, err = ct.Int32Value.getReal()

			case TypeInt64:
				r, s, err = ct.Int64Value.getReal()

			case TypeString:
				r, s, err = ct.StringValue.getReal()

			default:
				err = errors.New("Mismatched value type: " + string(ct.Type))
				break
		}

		if r == "" {
			break
		}

		if ct.Units != "" {
			u = ct.Units
		}
	}

	return r, u, s, err
}

func (ct *CmdType) isUpdated() (bool, error) {
	var ok bool
	var err error

	for range Only.Once {
		switch ct.Type {
			case TypeBit:
				ok = ct.BitValue.Updated

			case TypeBool:
				ok = ct.BoolValue.Updated

			case TypeEnum:
				ok = ct.EnumValue.Updated

			case TypeFloat32:
				ok = ct.Float32Value.Updated

			case TypeFloat64:
				ok = ct.Float64Value.Updated

			case TypeInt32:
				ok = ct.Int32Value.Updated

			case TypeInt64:
				ok = ct.Int64Value.Updated

			case TypeString:
				ok = ct.StringValue.Updated

			default:
				err = errors.New("Mismatched value type: " + string(ct.Type))
				break
		}
	}

	return ok, err
}

func (ct *CmdType) set(v interface{}) error {
	var err error

	for range Only.Once {
		// var vType int8
		var value string

		value, err = to_string(v)
		if err != nil {
			break
		}

		switch ct.Type {
			case TypeBit:
				t, _ := to_int32(value)
				err = ct.BitValue.set(t)
				if err != nil {
					break
				}

			case TypeBool:
				t, _ := to_bool(value)
				err = ct.BoolValue.set(t)
				if err != nil {
					break
				}

			case TypeEnum:
				t, _ := to_int32(value)
				err = ct.EnumValue.set(t)
				if err != nil {
					break
				}

			case TypeFloat32:
				t, _ := to_float32(value)
				err = ct.Float32Value.set(t)
				if err != nil {
					break
				}

			case TypeFloat64:
				t, _ := to_float64(value)
				err = ct.Float64Value.set(t)
				if err != nil {
					break
				}

			case TypeInt32:
				t, _ := to_int32(value)
				err = ct.Int32Value.set(t)
				if err != nil {
					break
				}

			case TypeInt64:
				t, _ := to_int64(value)
				err = ct.Int64Value.set(t)
				if err != nil {
					break
				}

			case TypeString:
				err = ct.StringValue.set(value)
				if err != nil {
					break
				}

			default:
				err = errors.New(fmt.Sprintf("Mismatched value type: %v", ct.Type))
				break
		}
	}

	return err
}


func to_bool(i interface{}) (bool, error) {
	var o bool
	var err error

	for range Only.Once {
		switch i.(type) {
			case bool:
				o = i.(bool)
				break

			case int:
				if i.(int) == 0 {
					o = false
				} else {
					o = true
				}
				break

			case int32:
				if i.(int32) == 0 {
					o = false
				} else {
					o = true
				}
				break

			case int64:
				if i.(int64) == 0 {
					o = false
				} else {
					o = true
				}
				break

			case float32:
				if i.(float32) == 0 {
					o = false
				} else {
					o = true
				}
				break

			case float64:
				if i.(float64) == 0 {
					o = false
				} else {
					o = true
				}
				break

			case string:
				o, err = strconv.ParseBool(i.(string))

			default:
				err = errors.New("Can't convert to bool")
				break
		}
	}

	return o, err
}

func to_int(i interface{}) (int, error) {
	var o int
	var err error

	for range Only.Once {
		switch i.(type) {
		case bool:
			if i.(bool) {
				o = 1
			} else {
				o = 0
			}
			break

		case int:
			o = i.(int)
			break

		case int32:
			o = int(i.(int32))
			break

		case int64:
			o = int(i.(int64))
			break

		case float32:
			o = int(i.(float32))
			break

		case float64:
			o = int(i.(float64))
			break

		case string:
			var t int64
			t, err = strconv.ParseInt(i.(string), 10, 32)
			o = int(t)

		default:
			err = errors.New("Can't convert to int")
			break
		}
	}

	return o, err
}

func to_int32(i interface{}) (int32, error) {
	var o int32
	var err error

	for range Only.Once {
		switch i.(type) {
		case bool:
			if i.(bool) {
				o = 1
			} else {
				o = 0
			}
			break

		case int:
			o = int32(i.(int))
			break

		case int32:
			o = i.(int32)
			break

		case int64:
			o = int32(i.(int64))
			break

		case float32:
			o = int32(i.(float32))
			break

		case float64:
			o = int32(i.(float64))
			break

		case string:
			var t int64
			t, err = strconv.ParseInt(i.(string), 10, 32)
			o = int32(t)

		default:
			err = errors.New("Can't convert to int32")
			break
		}
	}

	return o, err
}

func to_int64(i interface{}) (int64, error) {
	var o int64
	var err error

	for range Only.Once {
		switch i.(type) {
			case bool:
				if i.(bool) {
					o = 1
				} else {
					o = 0
				}
				break

			case int:
				o = int64(i.(int))
				break

			case int32:
				o = int64(i.(int32))
				break

			case int64:
				o = i.(int64)
				break

			case float32:
				o = int64(i.(float32))
				break

			case float64:
				o = int64(i.(float64))
				break

			case string:
				o, err = strconv.ParseInt(i.(string), 10, 64)

			default:
				err = errors.New("Can't convert to int64")
				break
		}
	}

	return o, err
}

func to_float32(i interface{}) (float32, error) {
	var o float32
	var err error

	for range Only.Once {
		switch i.(type) {
			case bool:
				if i.(bool) {
					o = 1
				} else {
					o = 0
				}
				break

			case int:
				o = float32(i.(int))
				break

			case int32:
				o = float32(i.(int32))
				break

			case int64:
				o = float32(i.(int64))
				break

			case float32:
				o = i.(float32)
				break

			case float64:
				o = float32(i.(float64))
				break

			case string:
				var t float64
				t, err = strconv.ParseFloat(i.(string),32)
				o = float32(t)

			default:
				err = errors.New("Can't convert to float32")
				break
		}
	}

	return o, err
}

func to_float64(i interface{}) (float64, error) {
	var o float64
	var err error

	for range Only.Once {
		switch i.(type) {
			case bool:
				if i.(bool) {
					o = 1
				} else {
					o = 0
				}
				break

			case int:
				o = float64(i.(int))
				break

			case int32:
				o = float64(i.(int32))
				break

			case int64:
				o = float64(i.(int64))
				break

			case float32:
				o = float64(i.(float32))
				break

			case float64:
				o = i.(float64)
				break

			case string:
				o, err = strconv.ParseFloat(i.(string),64)

			default:
				err = errors.New("Can't convert to float64")
				break
		}
	}

	return o, err
}

func to_string(i interface{}) (string, error) {
	var o string
	var err error

	for range Only.Once {
		switch i.(type) {
			case bool:
				if i.(bool) {
					o = "true"
				} else {
					o = "false"
				}
				break

			case int:
				// o = strconv.Itoa(i.(int))
				o = strconv.FormatInt(int64(i.(int)), 10)
				break

			case int32:
				o = strconv.FormatInt(int64(i.(int32)), 10)
				break

			case int64:
				o = strconv.FormatInt(i.(int64), 10)
				break

			case float32:
				o = strconv.FormatFloat(float64(i.(float32)), 'f', -1, 32)
				break

			case float64:
				o = strconv.FormatFloat(i.(float64), 'f', -1, 64)
				break

			case string:
				o = i.(string)

			default:
				err = errors.New("Can't convert to string")
				break
		}
	}

	return o, err
}


func (ct *CmdType) GetInterface() *CmdType {
	return ct
}

func (ct *CmdType) GetName() string {
	return ct.Name
}

func (ct *CmdType) GetTypeString() string {
	return ct.TypeString
}

func (ct *CmdType) GetType() int8 {
	return ct.Type
}

func (ct *CmdType) GetUnits() string {
	return ct.Units
}

func (ct *CmdType) DefineBit(n string, d string, units string, min int32, max int32) error {
	var err error

	ct.Name = n
	ct.Units = units
	// me.Description = d
	ct.Type = TypeBit
	ct.TypeString = StrBit

	err = ct.BitValue.define(min, max)

	return err
}
func (ct *CmdType) GetBit() (bool, error) {
	return ct.BoolValue.get()
}
func (ct *CmdType) SetBit(v int32) error {
	return ct.BitValue.set(v)
}


func (ct *CmdType) DefineBool(n string, d string, off string, on string) error {
	var err error

	ct.Name = n
	// me.Description = d
	ct.Type = TypeBool
	ct.TypeString = StrBool

	err = ct.BoolValue.define(off, on)

	return err
}
func (ct *CmdType) GetBool() (bool, error) {
	return ct.BoolValue.get()
}
func (ct *CmdType) SetBool(v bool) error {
	return ct.BoolValue.set(v)
}


func (ct *CmdType) DefineEnum(n string, d string, units string, sa []string) error {
	var err error

	ct.Name = n
	ct.Units = units
	// me.Description = d
	ct.Type = TypeEnum
	ct.TypeString = StrEnum

	err = ct.EnumValue.define(sa)

	return err
}
func (ct *CmdType) GetEnum() (int32, error) {
	return ct.EnumValue.get()
}
func (ct *CmdType) SetEnum(v int32) error {
	return ct.EnumValue.set(v)
}


func (ct *CmdType) DefineFloat32(n string, d string, units string, min float32, max float32, linear bool, inc float32) error {
	var err error

	ct.Name = n
	ct.Units = units
	// me.Description = d
	ct.Type = TypeFloat32
	ct.TypeString = StrFloat32

	err = ct.Float32Value.define(min, max, linear, inc)

	return err
}
func (ct *CmdType) GetFloat32() (float32, error) {
	return ct.Float32Value.get()
}
func (ct *CmdType) SetFloat32(v float32) error {
	return ct.Float32Value.set(v)
}


func (ct *CmdType) DefineFloat64(n string, d string, units string, min float64, max float64, linear bool, inc float64) error {
	var err error

	ct.Name = n
	ct.Units = units
	// me.Description = d
	ct.Type = TypeFloat64
	ct.TypeString = StrFloat64

	err = ct.Float64Value.define(min, max, linear, inc)

	return err
}
func (ct *CmdType) GetFloat64() (float64, error) {
	return ct.Float64Value.get()
}
func (ct *CmdType) SetFloat64(v float64) error {
	return ct.Float64Value.set(v)
}


func (ct *CmdType) DefineInt32(n string, d string, units string, min int32, max int32) error {
	var err error

	ct.Name = n
	ct.Units = units
	// me.Description = d
	ct.Type = TypeInt32
	ct.TypeString = StrInt32

	err = ct.Int32Value.define(min, max)

	return err
}
func (ct *CmdType) GetInteger32() (int32, error) {
	return ct.Int32Value.get()
}
func (ct *CmdType) SetInt32(v int32) error {
	return ct.Int32Value.set(v)
}


func (ct *CmdType) DefineInt64(n string, d string, units string, min int64, max int64) error {
	var err error

	ct.Name = n
	ct.Units = units
	// me.Description = d
	ct.Type = TypeInt64
	ct.TypeString = StrInt64

	err = ct.Int64Value.define(min, max)

	return err
}
func (ct *CmdType) GetInteger() (int64, error) {
	return ct.Int64Value.get()
}
func (ct *CmdType) SetInteger(v int64) error {
	return ct.Int64Value.set(v)
}


func (ct *CmdType) DefineString(n string, d string, p string) error {
	var err error

	ct.Name = n
	// me.Description = d
	ct.Type = TypeString
	ct.TypeString = StrString

	err = ct.StringValue.define(p)

	return err
}
func (ct *CmdType) GetString() (string, error) {
	return ct.StringValue.get()
}
func (ct *CmdType) SetString(v string) error {
	return ct.StringValue.set(v)
}


func (ct *CmdType) PrintName() {
	fmt.Printf("%s\n", ct.Name)
}


// Define Use reflection to traverse the OSC command structure
func (r *CmdRef) Define(x interface{}) { // , r CmdRef

	for range Only.Once {
		// if len(*r) == 0 {
		// 	break
		// }

		valueOf := reflect.ValueOf(x)
		reflectKind := valueOf.Kind()

		if reflectKind == reflect.Ptr {
			break
		}

		if reflectKind == reflect.Slice {
			for i := 0; i < valueOf.Len(); i++ {
				ref := valueOf.Index(i)
				// names = append(names, GetNames(ref.Interface())...)
				r.GetCmds(ref.Interface())
			}
			break
		}

		if reflectKind == reflect.Struct {
			for i := 0; i < valueOf.NumField(); i++ {
				// @TODO - Handle auto-creation of structure variables.
				// typeField := valueOf.Type().Field(i)
				// tagJson := typeField.Tag.Get("json")
				// tagType := typeField.Tag.Get("type")
				// if tagType == "enum" {
				// 	tagArray := typeField.Tag.Get("array")
				// 	// foo := reflect.ValueOf(AutoMixGroup)
				// 	fmt.Printf("(%s):'%s'\tjson:'%s'\tarray:'%s'\n", tagType , typeField.Name, tagJson, tagArray)
				// } else if tagType == "float32" {
				// 	fmt.Printf("(%s):'%s'\tjson:'%s'\n", tagType , typeField.Name, tagJson)
				// } else if tagType == "" {
				// 	fmt.Printf("(EMPTY):'%s'\tjson:'%s'\n", tagType , typeField.Name, tagJson)
				// } else {
				// 	fmt.Printf("(%s):'%s'\tjson:'%s'\n", tagType , typeField.Name, tagJson)
				// }

				// valueField := valueOf.Field(i)
				// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))

				ref := valueOf.Field(i).Interface()
				typeName := valueOf.Field(i).Type().Name()

				if typeName == "CmdType" {
					ct := ref.(CmdType)
					if ct.Name != "" {
						// names = append(names, ct.Name)
						fmt.Printf("Define>> %s\n", ct.Name)
						(*r)[ct.Name] = &ct
					}
				} else {
					// names = append(names, GetNames(ref)...)
					r.GetCmds(ref)
				}
			}

			break
		}
	}

}

// GetCmds Use reflection to traverse the OSC command structure
func (r *CmdRef) GetCmds(x interface{}) { // , r CmdRef

	for range Only.Once {
		// if len(*r) == 0 {
		// 	break
		// }

		valueOf := reflect.ValueOf(x)
		reflectKind := valueOf.Kind()

		if reflectKind == reflect.Ptr {
			break
		}

		if reflectKind == reflect.Slice {
			for i := 0; i < valueOf.Len(); i++ {
				ref := valueOf.Index(i)
				// names = append(names, GetNames(ref.Interface())...)
				r.GetCmds(ref.Interface())
			}
			break
		}

		if reflectKind == reflect.Struct {
			for i := 0; i < valueOf.NumField(); i++ {
				// @TODO - Handle auto-creation of structure variables.
				// typeField := valueOf.Type().Field(i)
				// tagJson := typeField.Tag.Get("json")
				// tagType := typeField.Tag.Get("type")
				// if tagType == "enum" {
				// 	tagArray := typeField.Tag.Get("array")
				// 	fmt.Printf("(%s):'%s'\tjson:'%s'\tarray:'%s'\n", tagType , typeField.Name, tagJson, tagArray)
				// } else if tagType == "float32" {
				// 	fmt.Printf("(%s):'%s'\tjson:'%s'\n", tagType , typeField.Name, tagJson)
				// } else if tagType == "" {
				// 	fmt.Printf("(EMPTY):'%s'\tjson:'%s'\n", typeField.Name, tagJson)
				// } else {
				// 	fmt.Printf("(%s):'%s'\tjson:'%s'\n", tagType , typeField.Name, tagJson)
				// }

				// valueField := valueOf.Field(i)
				// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))

				ref := valueOf.Field(i).Interface()
				typeName := valueOf.Field(i).Type().Name()

				if typeName == "CmdType" {
					ct := ref.(CmdType)
					if ct.Name != "" {
						// names = append(names, ct.Name)
						// fmt.Printf("GetCmds>> %s\n", ct.Name)
						(*r)[ct.Name] = &ct
					}
				} else {
					// names = append(names, GetNames(ref)...)
					r.GetCmds(ref)
				}
			}

			break
		}
	}

}

// GetNames Use reflection to traverse the OSC command structure
func GetNames(x interface{}) []string {
	var names []string

	for range Only.Once {
		valueOf := reflect.ValueOf(x)

		// reflectName := valueOf.Type().Name()
		// reflectName := reflect.TypeOf(x).Name()
		// fmt.Printf("reflectName(ValueOf) == %s\n", reflectName)
		reflectKind := valueOf.Kind()
		// reflectKind := reflect.TypeOf(x).Kind()
		// fmt.Printf("reflectKind(ValueOf) == %s\n", reflectKind)

		if reflectKind == reflect.Ptr {
			break
		}

		if reflectKind == reflect.Slice {
			// ref1 := v.Interface()
			// ref2 := reflect.TypeOf(x).Elem()
			// fmt.Printf("ref1:%v\n", ref1)
			// fmt.Printf("ref2:%v\n", ref2)

			for i := 0; i < valueOf.Len(); i++ {
				ref := valueOf.Index(i)
				names = append(names, GetNames(ref.Interface())...)

				// r1 := ref.Type()
				// fmt.Printf("reflectType == %s\n", r1)
				// r2 := ref.Kind()
				// fmt.Printf("reflectKind == %s\n", r2)
			}
			break
		}

		if reflectKind == reflect.Struct {
			// v := reflect.ValueOf(x)
			// valueOf = valueOf.Elem()
			// valueOf.Len()
			for i := 0; i < valueOf.NumField(); i++ {
				// typeField := valueOf.Type().Field(i)
				// tag := typeField.Tag
				// fmt.Printf("Field Name: %s,\tTag Value: %s\n", typeField.Name, tag.Get("json"))

				// valueField := valueOf.Field(i)
				// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))

				ref := valueOf.Field(i).Interface()
				typeName := valueOf.Field(i).Type().Name()

				if typeName == "CmdType" {
					ct := ref.(CmdType)
					// fmt.Printf("(%s).Name == %s\t", typeName, ct.Name)
					// fmt.Printf("%s\n", ct.Name)
					// fmt.Printf("GetNames>> %s\n", ct.Name)
					if ct.Name != "" {
						names = append(names, ct.Name)
					}
				} else {
					// fmt.Printf("(%s).Name\n", typeName)
					names = append(names, GetNames(ref)...)
				}
			}

			break
		}
	}

	// fmt.Println(names)
	return names
}
