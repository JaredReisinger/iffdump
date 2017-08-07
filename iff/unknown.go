package iff

import (
	"errors"
	"fmt"
	"io"
	"math"

	log "github.com/sirupsen/logrus"
)

// Unknown represents an unrecognized chunk
type Unknown struct {
	typeID TypeID
	length int64
	data   []byte
}

// TypeID ...
func (c *Unknown) TypeID() TypeID {
	return c.typeID
}

// Length ...
func (c *Unknown) Length() int64 {
	return c.length
}

func (c *Unknown) String() string {
	s := RenderBytes(c.data, 16)
	return fmt.Sprintf("%s (unknown, length %d)\n%s", c.typeID, c.length, s)
}

// Data ...
func (c *Unknown) Data() []byte {
	return c.data
}

type unknownChunkDecoder struct{}

func (d *unknownChunkDecoder) Decode(typeID TypeID, r *io.SectionReader, context *Decoder, logger log.FieldLogger) (Chunk, error) {

	length := r.Size()
	if length > math.MaxUint32 {
		return nil, errors.New("chunk size exceeds uint32 maximum")
	}

	data := make([]byte, length)
	_, err := io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	c := &Unknown{
		typeID: typeID,
		length: length,
		data:   data,
	}

	return c, nil
}
