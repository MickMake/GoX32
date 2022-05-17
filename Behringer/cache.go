package Behringer

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Behringer/api/output"
	"github.com/MickMake/GoX32/Only"
	"github.com/loffa/gosc"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)



func (x *X32) SetCacheDir(basedir string) error {
	for range Only.Once {
		x.cacheDir = filepath.Join(basedir)
		_, x.Error = os.Stat(x.cacheDir)
		if os.IsExist(x.Error) {
			x.Error = nil
			break
		}

		x.Error = os.MkdirAll(x.cacheDir, 0700)
		if x.Error != nil {
			break
		}
	}

	return x.Error
}

func (x *X32) GetCacheDir() string {
	return x.cacheDir
}

func (x *X32) SetCacheTimeout(duration time.Duration) {
	if duration == 0 {
		duration = time.Minute
	}
	x.cacheTimeout = duration
}

func (x *X32) GetCacheTimeout() time.Duration {
	return x.cacheTimeout
}

// CheckCache Retrieves cache data from a local file.
func (x *X32) CheckCache() bool {
	var ok bool
	for range Only.Once {
		fn := filepath.Join(x.cacheDir, "states.json")

		var f os.FileInfo
		f, x.Error = os.Stat(fn)
		if x.Error != nil {
			if os.IsNotExist(x.Error) {
				x.Error = nil
			}
			break
		}

		if f.IsDir() {
			x.Error = errors.New("file is a directory")
			break
		}

		duration := x.GetCacheTimeout()
		then := f.ModTime()
		then = then.Add(duration)
		now := time.Now()
		if then.Before(now) {
			break
		}

		ok = true
	}

	return ok
}

// CacheRead Retrieves cache data from a local file.
func (x *X32) CacheRead() error {
	for range Only.Once {
		fn := filepath.Join(x.cacheDir, "states.json")
		x.Error = output.FileRead(fn, &x.cache)
		if x.Error != nil {
			if x.Error.Error() == "EOF" {
				x.Error = nil
				break
			}
		}
		for n := range x.cache {
			x.cache[n].SeenBefore = false	// Force refresh on re-read.
			x.cache[n].LastSeen = time.Now().Add(- x.cacheTimeout)
		}
	}
	return x.Error
}

// CacheRemove Removes a cache file.
func (x *X32) CacheRemove() error {
	fn := filepath.Join(x.cacheDir, "states.json")
	return output.FileRemove(fn)
}

// CacheWrite Saves cache data to a file path.
func (x *X32) CacheWrite() error {
	for range Only.Once {
		if x.CheckCache() {
			break
		}

		fn := filepath.Join(x.cacheDir, "states.json")
		x.Error = output.FileWrite(fn, x.cache, output.DefaultFileMode)
	}
	return x.Error
}


func (x *X32) UpdateCache(msg *gosc.Message) *Message {
	for range Only.Once {
		if x.cache == nil {
			x.cache = make(MessageMap)
		}
		if !x.MessageExists(msg.Address) {
			x.cache[msg.Address] = &Message{
				Message:    msg,
				SeenBefore: false,
				Counter:    1,
				LastSeen:   time.Now(),
			}
			break
		}

		then := x.cache[msg.Address].LastSeen
		then = then.Add(x.cacheTimeout)
		now := time.Now()
		if then.Before(now) {
			x.cache[msg.Address].Counter++
			x.cache[msg.Address].SeenBefore = false
			x.cache[msg.Address].Message = msg
			x.cache[msg.Address].LastSeen = time.Now()
			break
		}

		x.cache[msg.Address].Counter++
		x.cache[msg.Address].SeenBefore = true
		x.cache[msg.Address].Message = msg
		x.cache[msg.Address].LastSeen = time.Now()
	}

	return x.cache[msg.Address]
}


func (x *X32) MessageExists(address string) bool {
	return x.cache.Exists(address)
}

func (x *X32) GetMessage(address string) *Message {
	return x.cache.Get(address)
}

func (x *X32) MessageSeenBefore(address string) bool {
	return x.cache.SeenBefore(address)
}


type MessageMap map[string]*Message

func (m *MessageMap) Exists(address string) bool {
	_, ok := (*m)[address]
	return ok
}

func (m *MessageMap) Get(address string) *Message {
	if ret, ok := (*m)[address]; ok {
		return ret
	}
	return nil
}

func (m *MessageMap) SeenBefore(address string) bool {
	ret := m.Get(address)
	if ret == nil {
		return false
	}
	return ret.SeenBefore
}


type Message struct {
	*gosc.Message `json:"Message"`
	SeenBefore   bool      `json:"seen_before"`
	LastSeen     time.Time `json:"last_seen"`
	Counter      int       `json:"counter"`
	Type         string    `json:"type"`
	Error        error     `json:"-"`
}

func (m *Message) DetermineType() string {
	var ret string
	for range Only.Once {
		for _, v := range m.Arguments {
			ret = reflect.TypeOf(v).Name()
		}
	}
	return ret
}

func (m *Message) CacheStale() bool {
	var ok bool
	for range Only.Once {
		if !m.SeenBefore {
			ok = true
		}

		// then := m.LastSeen
		// then = then.Add()
		// now := time.Now()
		// if then.Before(now) {
		// 	break
		// }
		//
		// ok = true
	}
	return ok
}

func (m *Message) GetType() string {
	var ret string
	for range Only.Once {
		if m.Type != "" {
			ret = m.Type
			break
		}

		// if len(m.Arguments) > 1 {
		// 	ret = api.UnitArray
		// 	break
		// }

		ret = reflect.TypeOf(m.Arguments[0]).Name()
	}
	return ret
}

func (m *Message) GetStringValue() string {
	var ret string
	for range Only.Once {
		for _, a := range m.Arguments {
			ret += fmt.Sprintf("%v ", a)
		}
		ret = strings.TrimSpace(ret)
	}
	return ret
}

func (m *Message) GetBoolValue() bool {
	var ok bool
	for range Only.Once {
		for _, a := range m.Arguments {
			if a == "ON" {
				ok = true
			} else {
				ok = false
			}
			break
		}
	}
	return ok
}

func (m *Message) GetFloatValue() float64 {
	var ret float64
	for range Only.Once {
		var err error
		for _, a := range m.Arguments {
			ret, err = strconv.ParseFloat(fmt.Sprintf("%v", a), 64)
			if err == nil {
				break
			}
		}
	}
	return ret
}
