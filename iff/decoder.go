package iff

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	log "github.com/sirupsen/logrus"
)

// Decoder is the top-level IFF decoder that uses a map of ChunkDecoders to do
// its thing.
type Decoder struct {
	customDecoders map[TypeID]ChunkDecoder
	// Do we provide defaults and/or fallbacks here?
	standardDecoders map[TypeID]ChunkDecoder
	fallbackDecoder  ChunkDecoder
}

// NewDecoder creates a decoder with the provided known chunk decoders.
func NewDecoder(decoders map[TypeID]ChunkDecoder) *Decoder {
	standardDecoders := make(map[TypeID]ChunkDecoder, 0)
	standardDecoders[FormType] = &formDecoder{}
	standardDecoders[AnnotationType] = &annoDecoder{}
	standardDecoders[AuthorType] = &authDecoder{}
	standardDecoders[CopyrightType] = &copyrightDecoder{}
	return &Decoder{
		customDecoders:   decoders,
		standardDecoders: standardDecoders,
		fallbackDecoder:  &unknownChunkDecoder{},
	}
}

// Decode ...
func (d *Decoder) Decode(r ReadAtSeeker, logger log.FieldLogger) (Chunk, error) {
	typeID, err := ReadTypeID(r)
	if err != nil {
		return nil, err
	}

	length, err := ReadUint32(r)
	if err != nil {
		return nil, err
	}

	off, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	sr := io.NewSectionReader(r, off, int64(length))

	decoder, ok := d.customDecoders[typeID]
	if !ok {
		decoder, ok = d.standardDecoders[typeID]
	}
	if !ok {
		decoder = d.fallbackDecoder
	}

	logger.WithFields(log.Fields{
		"typeID":  typeID,
		"decoder": reflect.TypeOf(decoder),
	}).Debug("decoding...")

	return decoder.Decode(typeID, sr, d, logger)
}

// ExpectType is a helper for the common type ID assertion
func ExpectType(expected TypeID, actual TypeID) error {
	if expected == actual {
		return nil
	}

	return fmt.Errorf("expected type ID of %q, got %q", expected, actual)
}

// ReadTypeID reads a type ID (FourCC) from the reader.
func ReadTypeID(r io.Reader) (TypeID, error) {
	s, err := ReadString(r, 4)
	return TypeID(s), err
}

// ReadUint32 reads a 32-bit (4-byte) unsigned integer.    Note that the IFF
// documentation uses an unsigned, 32-bit integer for chunk lengths, but Go uses
// a signed, 64-bit integer for seeking and many Reader/Writer APIs. Rather than
// deal with casting all over the place, we just use Go's native int64 since it
// can represent a uint32 without problem.  (The only downside is that it *will*
// allow values that a uint32 would not... but since we're currently only
// reading the IFF, this isn't really an issue.)
// func ReadUint32(r io.Reader) (int64, error) {
// 	var val uint32
// 	err := binary.Read(r, binary.BigEndian, &val)
// 	return int64(val), err
// }
func ReadUint32(r io.Reader) (uint32, error) {
	var val uint32
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// ReadUint24 reads a 3-byte number (24-bit value)
func ReadUint24(r io.Reader) (uint32, error) {
	b, err := ReadBytes(r, 3)
	if err != nil {
		return 0, err
	}
	return (uint32(b[0]) << 16) + (uint32(b[1]) << 8) + uint32(b[2]), nil
}

// ReadUint16 reads a word (16-bit value)
func ReadUint16(r io.Reader) (uint16, error) {
	var val uint16
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// ReadUint8 reads a word (16-bit value)
func ReadUint8(r io.Reader) (uint8, error) {
	var val uint8
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// ReadString reads a string of known length
func ReadString(r io.Reader, length int64) (string, error) {
	// b := make([]byte, length)
	// _, err := io.ReadFull(r, b)
	// if err != nil {
	// 	return "", err
	// }
	b, err := ReadBytes(r, length)
	return string(b), err
}

// ReadBytes reads a byte "string" of known length
func ReadBytes(r io.Reader, length int64) ([]byte, error) {
	b := make([]byte, length)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
