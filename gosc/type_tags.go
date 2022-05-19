package gosc

import (
	"bufio"
	"encoding/binary"
	"io"
	"reflect"
)

type valueReader func(r *bufio.Reader) (any, error)
type valueWriter func(w *bufio.Writer, data any) (typeTag, error)

type typeTag byte

// Constants for all TypeTags as defined in the specification.
const (
	TypeTagInt32           = typeTag('i')
	TypeTagFloat32         = typeTag('f')
	TypeTagString          = typeTag('s')
	TypeTagBlob            = typeTag('b') // TODO: Implement write for TypeTagBlob
	TypeTagBigInt          = typeTag('h') // TODO: Implement read/write for TypeTagBigInt
	TypeTagTimetag         = typeTag('t')
	TypeTagDouble          = typeTag('d') // TODO: Implement read/write for TypeTagDouble
	TypeTagStringAlternate = typeTag('S') // TODO: Implement read/write for TypeTagStringAlternate
	TypeTagChar            = typeTag('c') // TODO: Implement read/write for TypeTagChar
	TypeTagRGBA            = typeTag('r') // TODO: Implement read/write for TypeTagRGBA
	TypeTagMIDI            = typeTag('m') // TODO: Implement read/write for TypeTagMIDI
	TypeTagTrue            = typeTag('T') // TODO: Implement read/write for TypeTagTrue
	TypeTagFalse           = typeTag('F') // TODO: Implement read/write for TypeTagFalse
	TypeTagNil             = typeTag('N') // TODO: Implement read/write for TypeTagNil
	TypeTagInfinite        = typeTag('I') // TODO: Implement read/write for TypeTagInfinite
	TypeTagArrayStart      = typeTag('[') // TODO: Implement read/write for TypeTagArrayStart
	TypeTagArrayStop       = typeTag(']') // TODO: Implement read/write for TypeTagArrayStop
)

// readerMap is used to map typeTag to the correct reader function.
var readerMap = map[typeTag]valueReader{
	TypeTagInt32:   int32Reader,
	TypeTagFloat32: float32Reader,
	TypeTagString:  stringReader,
	TypeTagBlob:    blobReader,
	TypeTagTimetag: timeTagReader,
}

// writeMap is used to map the golang types to the correct writer.
var writerMap = map[reflect.Type]valueWriter{
	reflect.TypeOf(int32(0)):   int32Writer,
	reflect.TypeOf(float32(0)): float32Writer,
	reflect.TypeOf(""):         stringWriter,
	reflect.TypeOf(Timetag(0)): timeTagWriter,
}

func stringWriter(w *bufio.Writer, data any) (typeTag, error) {
	return TypeTagString, writePaddedString(w, data.(string))
}

func float32Writer(w *bufio.Writer, data any) (typeTag, error) {
	return TypeTagFloat32, binary.Write(w, binary.BigEndian, data)
}

func int32Writer(w *bufio.Writer, data any) (typeTag, error) {
	return TypeTagInt32, binary.Write(w, binary.BigEndian, data)
}

func timeTagWriter(w *bufio.Writer, data any) (typeTag, error) {
	return TypeTagTimetag, binary.Write(w, binary.BigEndian, data)
}

func blobReader(r *bufio.Reader) (any, error) {
	n := int32(0)
	if err := binary.Read(r, binary.BigEndian, &n); err != nil {
		return nil, err
	}
	res := make([]byte, n)
	if _, err := io.ReadFull(r, res); err != nil {
		return nil, err
	}
	_, err := r.Discard(getPadBytes(len(res) + 4))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func stringReader(r *bufio.Reader) (any, error) {
	return readPaddedString(r)
}

func float32Reader(r *bufio.Reader) (any, error) {
	res := float32(0)
	if err := binary.Read(r, binary.BigEndian, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func int32Reader(r *bufio.Reader) (any, error) {
	res := int32(0)
	if err := binary.Read(r, binary.BigEndian, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func timeTagReader(r *bufio.Reader) (any, error) {
	tt := Timetag(0)
	if err := binary.Read(r, binary.BigEndian, &tt); err != nil {
		return nil, err
	}
	return tt, nil
}
