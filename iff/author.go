package iff

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

const (
	// AuthorType ...
	AuthorType TypeID = "AUTH"
)

type authDecoder struct{}

func (d *authDecoder) Decode(typeID TypeID, r *io.SectionReader, context *Decoder, logger log.FieldLogger) (Chunk, error) {
	if err := ExpectType(AuthorType, typeID); err != nil {
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

// Author is an ANNO IFF chunk
type Author struct {
	typeID TypeID
	length int64
	value  string
}

// TypeID ...
func (c *Author) TypeID() TypeID {
	return c.typeID
}

// Length ...
func (c *Author) Length() int64 {
	return c.length
}

func (c *Author) String() string {
	return fmt.Sprintf("%s (length %d): %q", c.typeID, c.length, c.value)
}
