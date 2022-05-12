package api

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"strconv"
	"strings"
)


func CleanString(s string) string {
	// var ret string
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			(b == '-') ||
			(b == '_') ||
			(b == '.') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}

	// dupes := regexp.MustCompile(`\s+`)
	// ret = dupes.ReplaceAllString(result.String(), )
	return result.String()
}


func ResolvePoint(point string) *Point {
	return Points.Resolve(point)
}

func GetPoint(device string, point string) *Point {
	return Points.Get(device, point)
}

func GetPointInt(device string, point int64) *Point {
	return Points.Get(device, strconv.FormatInt(point, 10))
}

func GetDevicePoint(devicePoint string) *Point {
	return Points.GetDevicePoint(devicePoint)
}

// func GetPointName(device string, point int64) string {
// 	return fmt.Sprintf("%s.%d", device, point)
// }

func NameDevicePointInt(device string, point int64) string {
	return fmt.Sprintf("%s.%d", device, point)
}

func NameDevicePoint(device string, point string) string {
	if device == "" {
		return point
	}
	return fmt.Sprintf("%s.%s", device, point)
}

func SetPoint(point string) string {
	for range Only.Once {
		p := strings.TrimPrefix(point, "p")
		_, err := strconv.ParseInt(p, 10, 64)
		if err == nil {
			point = "p" + p
			break
		}
	}
	return point
}

func PointToName(s string) string {
	s = CleanString(s)
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.Title(s)
	return s
}


// CheckString RequestCommon checks
func CheckString(name string, rc string) error {
	var err error
	for range Only.Once {
		if rc == "" {
			err = errors.New(name + ": empty string")
			break
		}
		if strings.TrimSpace(rc) == "" {
			err = errors.New(name + ": empty string with spaces")
			break
		}
	}
	return err
}
