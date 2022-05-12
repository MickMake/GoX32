package api

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Behringer/api/apiReflect"
	"github.com/MickMake/GoX32/Behringer/api/output"
	"github.com/MickMake/GoX32/Only"
	"os"
	"path/filepath"
)


func (ep *EndPointStruct) SetFilenamePrefix(format string, args ...interface{}) string {
	if format != "" {
		ep.FileNamePrefix = fmt.Sprintf(format, args...)
		// ep.FileNamePrefix = fmt.Sprintf("%s_%s-%s", ep.Area, ep.Name, ep.FileNamePrefix)
		ep.FileNamePrefix = string(ep.Area) + "_" + string(ep.Name) + "-" + ep.FileNamePrefix
	} else {
		ep.FileNamePrefix = string(ep.Area) + "_" + string(ep.Name)
	}
	return ep.FileNamePrefix
}

func (ep *EndPointStruct) GetCsvFilename() string {
	if ep.FileNamePrefix == "" {
		ep.SetFilenamePrefix("")
	}
	return ep.FileNamePrefix + ".csv"
}

func (ep *EndPointStruct) GetJsonFilename() string {
	if ep.FileNamePrefix == "" {
		ep.SetFilenamePrefix("")
	}
	return ep.FileNamePrefix + ".json"
}

func (ep *EndPointStruct) GetImageFilename() string {
	if ep.FileNamePrefix == "" {
		ep.SetFilenamePrefix("")
	}
	return ep.FileNamePrefix + ".png"
}

func (ep *EndPointStruct) GetFilePath() string {
	var ret string

	for range Only.Once {
		var hd string
		hd, ep.Error = os.UserHomeDir()
		if ep.Error != nil {
			break
		}

		ret = filepath.Join(hd, ".GoSungrow", ep.GetJsonFilename())
	}

	return ret
}

// FileExists Checks for existance of a local file.
func (ep *EndPointStruct) FileExists(fn string) bool {
	var ok bool

	for range Only.Once {
		if fn == "" {
			fn = ep.GetJsonFilename()
			if ep.Error != nil {
				break
			}
		}

		var f os.FileInfo
		f, ep.Error = os.Stat(fn)
		if ep.Error != nil {
			if os.IsNotExist(ep.Error) {
				ep.Error = nil
			}
			break
		}
		if f.IsDir() {
			ep.Error = errors.New("file is a directory")
			break
		}

		ok = true
	}

	return ok
}

// FileRead Retrieves data from a local file.
func (ep *EndPointStruct) ApiReadDataFile(ref interface{}) error {
	return output.FileRead(ep.GetJsonFilename(), ref)

	// // var ret []byte
	//
	// for range Only.Once {
	// 	if fn == "" {
	// 		fn = ep.GetJsonFilename()
	// 		if ep.Error != nil {
	// 			break
	// 		}
	// 	}
	//
	// 	var f *os.File
	// 	f, ep.Error = os.Open(fn)
	// 	if ep.Error != nil {
	// 		if os.IsNotExist(ep.Error) {
	// 			ep.Error = nil
	// 		}
	// 		break
	// 	}
	//
	// 	//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
	// 	defer f.Close()
	//
	// 	ep.Error = json.NewDecoder(f).Decode(&ref)
	// }
	//
	// // for range Only.Once {
	// //	fn := ep.GetFilename()
	// //	if ep.Error != nil {
	// //		break
	// //	}
	// //
	// //	ret, ep.Error = os.FileRead(fn)
	// //	if ep.Error != nil {
	// //		break
	// //	}
	// // }
	//
	// return ep.Error
}

// FileWrite Saves data to a file path.
func (ep *EndPointStruct) ApiWriteDataFile(ref interface{}) error {
	return output.FileWrite(ep.GetJsonFilename(), ref, output.DefaultFileMode)
}

func (ep *EndPointStruct) ApiCacheFilename(request interface{}) string {
	postfix := apiReflect.GetFingerprint(request)
	return fmt.Sprintf("%s_%s-%s.json", ep.Area, ep.Name, postfix)
}

func (ep *EndPointStruct) ApiFingerprint(request interface{}) string {
	return apiReflect.GetFingerprint(request)
}
