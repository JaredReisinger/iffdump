package quetzal

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"

	"github.com/JaredReisinger/iffdump/iff"
)

const (
	CompressedMemoryType iff.TypeID = "CMem"
)

type cmemDecoder struct{}

func (d *cmemDecoder) Decode(typeID iff.TypeID, r *io.SectionReader, context *iff.Decoder, logger log.FieldLogger) (iff.Chunk, error) {
	if err := iff.ExpectType(CompressedMemoryType, typeID); err != nil {
		return nil, err
	}

	length := r.Size()

	raw := make([]byte, length)
	_, err := io.ReadFull(r, raw)
	if err != nil {
		return nil, err
	}

	c := &CMem{
		typeID: typeID,
		length: length,
		raw:    raw,
	}

	return c, nil
}

// CMem is a FORM IFF chunk
type CMem struct {
	typeID iff.TypeID
	length int64
	raw    []byte
}

// TypeID ...
func (c *CMem) TypeID() iff.TypeID {
	return c.typeID
}

// Length ...
func (c *CMem) Length() int64 {
	return c.length
}

func (c *CMem) String() string {
	raw := iff.RenderBytes(c.raw, 16)
	return fmt.Sprintf("%s (length %d)\n%s", c.typeID, c.length, raw)
}
