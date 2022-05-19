package gosc

import "time"

// PackageType is used to create comparable constants.
type PackageType string

// Constants representing the different types of packages that exist.
const (
	PackageTypeMessage = PackageType("message")
	PackageTypeBundle  = PackageType("bundle")
)

// Immediately is a specific Timetag representing immediate execution of a Bundle.
const Immediately = Timetag(0x01)
const timeTo1970 = 2208988800

// Package is the generalization of both package types.
type Package interface {
	GetType() PackageType
}

// Message is the data structure for OSC message packets.
type Message struct {
	// The Address is a '/' separated string as per the specification
	Address string
	// Arguments is the array of that that is written when the package is sent.
	// Only data-types with defined writers is supported.
	Arguments []any
}

// GetType returns the package type for Messages
func (m *Message) GetType() PackageType {
	return PackageTypeMessage
}

// Timetag represents the time since 1900-01-01 00:00.
type Timetag uint64

// Bundle is the data structure for OSC bundle packets.
type Bundle struct {
	// Timetag for execution of the messages in this Bundle
	Timetag Timetag
	// List of messages to execute at Timetag. Messages are expected to be
	// handled atomically.
	Messages []*Message
	// Bundles can contain bundles, bundles in budles are not handled
	// atomically.
	Bundles []*Bundle
	// Name is the name of the packet after the '#' when encoding. If omitted
	// this is set to 'bundle'.
	Name string
}

// GetType returns the package type for Bundle
func (b *Bundle) GetType() PackageType {
	return PackageTypeBundle
}

func getPadBytes(length int) int {
	return (4 - (length % 4)) % 4
}

// Fractions will return the fractions of the Timetag with picoseconds
// resolution.
func (tt Timetag) Fractions() uint32 {
	return uint32(tt & 0x00000000FFFFFFFF)
}

// Seconds will return the seconds since the Timetag beginning.
func (tt Timetag) Seconds() uint32 {
	return uint32(tt & 0xFFFFFFFF00000000 >> 32)
}

// Time will return the Timetag as a Golang time.Time type.
func (tt Timetag) Time() time.Time {
	return time.Unix(int64(tt.Seconds()-timeTo1970), int64(tt.Fractions()))
}
