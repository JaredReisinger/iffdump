package iff

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

const (
	// CopyrightType ...
	CopyrightType TypeID = "(c) "
)

type copyrightDecoder struct{}

func (d *copyrightDecoder) Decode(typeID TypeID, r *io.SectionReader, context *Decoder, logger log.FieldLogger) (Chunk, error) {
	if err := ExpectType(CopyrightType, typeID); err != nil {
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

// Copyright is an ANNO IFF chunk
type Copyright struct {
	typeID TypeID
	length int64
	value  string
}

// TypeID ...
func (c *Copyright) TypeID() TypeID {
	return c.typeID
}

// Length ...
func (c *Copyright) Length() int64 {
	return c.length
}

func (c *Copyright) String() string {
	return fmt.Sprintf("%s (length %d): %q", c.typeID, c.length, c.value)
}
