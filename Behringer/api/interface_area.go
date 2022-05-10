package api

import "github.com/MickMake/GoX32/Behringer/api/output"

type Area interface {
	Init(*Web) AreaStruct
	GetAreaName() AreaName
	GetEndPoints() TypeEndPoints
	Call(name EndPointName) output.Json
	SetRequest(name EndPointName, ref interface{}) error
	GetRequest(name EndPointName) output.Json
	GetResponse(name EndPointName) output.Json
	GetData(name EndPointName) output.Json
	IsValid(name EndPointName) error
	GetError(name EndPointName) error
}

