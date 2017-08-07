package iff

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

const (
	// AnnotationType ...
	AnnotationType TypeID = "ANNO"
)

type annoDecoder struct{}

func (d *annoDecoder) Decode(typeID TypeID, r *io.SectionReader, context *Decoder, logger log.FieldLogger) (Chunk, error) {
	if err := ExpectType(AnnotationType, typeID); err != nil {
		return nil, err
	}

	length := r.Size()
	value, err := ReadString(r, length)
	if err != nil {
		return nil, err
	}

	c := &Annotation{
		typeID: typeID,
		length: length,
		value:  value,
	}

	return c, nil
}

// Annotation is an ANNO IFF chunk
type Annotation struct {
	typeID TypeID
	length int64
	value  string
}

// TypeID ...
func (c *Annotation) TypeID() TypeID {
	return c.typeID
}

// Length ...
func (c *Annotation) Length() int64 {
	return c.length
}

func (c *Annotation) String() string {
	return fmt.Sprintf("%s (length %d): %q", c.typeID, c.length, c.value)
}
