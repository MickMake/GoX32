package gosc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func readPackage(r *bufio.Reader) (pack Package, err error) {
	firstByte, err := r.Peek(1)
	if err != nil {
		return nil, err
	}
	if firstByte[0] == '/' {
		pack, err = readMessage(r)
	} else if firstByte[0] == '#' {
		pack, err = readBundle(r)
	}

	return
}

func readBundle(r *bufio.Reader) (bundle *Bundle, err error) {
	name, err := readPaddedString(r)
	if err != nil {
		return nil, err
	}
	tt, err := timeTagReader(r)
	if err != nil {
		return nil, err
	}
	bundle = &Bundle{
		Timetag:  tt.(Timetag),
		Messages: make([]*Message, 0, 10),
		Bundles:  make([]*Bundle, 0, 10),
		Name:     name,
	}

	for {
		length := int32(0)
		if err := binary.Read(r, binary.BigEndian, &length); err == io.EOF {
			return bundle, nil
		} else if err != nil {
			return nil, err
		}
		buf := make([]byte, length)
		bufR := bufio.NewReader(bytes.NewReader(buf))
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return nil, err
		}
		pack, err := readPackage(bufR)
		if err != nil {
			return nil, err
		}
		if pack.GetType() == PackageTypeMessage {
			bundle.Messages = append(bundle.Messages, pack.(*Message))
		} else if pack.GetType() == PackageTypeBundle {
			bundle.Bundles = append(bundle.Bundles, pack.(*Bundle))
		} else {
			return nil, errors.New("bundle contained unknown type")
		}
	}
}

func readMessage(r *bufio.Reader) (msg *Message, err error) {
	address, err := readPaddedString(r)
	if err != nil {
		return nil, err
	}
	args, err := readArguments(r)
	if err != nil {
		return nil, err
	}
	return &Message{
		Address:   address,
		Arguments: args,
	}, nil
}

func readPaddedString(r *bufio.Reader) (str string, err error) {
	str, err = r.ReadString(0)
	if err != nil {
		return "", err
	}

	_, err = r.Discard(getPadBytes(len(str)))
	if err != nil {
		return "", err
	}
	str = str[:len(str)-1]
	return str, nil
}

func readArguments(r *bufio.Reader) ([]any, error) {
	typeTags, err := readPaddedString(r)
	if err != nil {
		return nil, err
	}
	if typeTags[0] != ',' {
		return nil, errors.New("typetag format error")
	}
	typeTags = typeTags[1:]
	res := make([]any, 0, len(typeTags))

	for _, tt := range typeTags {
		decoder, ok := readerMap[typeTag(tt)]
		if !ok {
			return nil, fmt.Errorf("unknown tag type '%s'", string(tt))
		}
		val, err := decoder(r)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	return res, nil
}
