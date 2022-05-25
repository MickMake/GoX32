package gosc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
)

func writePackage(pack Package, w *bufio.Writer) error {
	switch v := pack.(type) {
	case *Message:
		err := writePaddedString(w, v.Address)
		if err != nil {
			return err
		}
		err = writeArguments(w, v.Arguments)
		if err != nil {
			return err
		}
	case *Bundle:
		if v.Name == "" {
			v.Name = "#bundle"
		}
		err := writePaddedString(w, v.Name)
		if err != nil {
			return err
		}
		_, err = timeTagWriter(w, v.Timetag)
		if err != nil {
			return err
		}
		for _, msg := range v.Messages {
			err = writePayload(w, msg)
			if err != nil {
				return err
			}
		}
		for _, bundle := range v.Bundles {
			err = writePayload(w, bundle)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown package type (%v)", pack)
	}
	return nil
}

func writePayload(w *bufio.Writer, pack Package) error {
	buf := bytes.Buffer{}
	bufW := bufio.NewWriter(&buf)
	err := writePackage(pack, bufW)
	if err != nil {
		return err
	}
	_ = bufW.Flush()
	data := buf.Bytes()
	if err := binary.Write(w, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func writeArguments(w *bufio.Writer, arguments []any) error {
	buf := bytes.Buffer{}
	bufWriter := bufio.NewWriter(&buf)
	ttString := strings.Builder{}
	_, err := ttString.WriteRune(',')
	if err != nil {
		return err
	}

	for _, a := range arguments {
		typ := reflect.TypeOf(a)
		// fmt.Printf("TYPE: %s\n", typ.Name())
		if writer, ok := writerMap[typ]; ok {
			tt, err := writer(bufWriter, a)
			if err != nil {
				return err
			}
			if err := ttString.WriteByte(byte(tt)); err != nil {
				return err
			}
		}
	}

	err = writePaddedString(w, ttString.String())
	if err != nil {
		return err
	}
	_ = bufWriter.Flush()
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func writePaddedString(w *bufio.Writer, str string) error {
	if !strings.HasSuffix(str, "\x00") {
		str += "\x00"
	}
	_, err := w.WriteString(str)
	if err != nil {
		return err
	}
	padding := make([]byte, getPadBytes(len(str)))
	_, err = w.Write(padding)
	if err != nil {
		return err
	}
	return nil
}
