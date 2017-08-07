package quetzal

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"

	"github.com/JaredReisinger/iffdump/iff"
)

const (
	UncompressedMemoryType iff.TypeID = "UMem"
)

type umemDecoder struct{}

func (d *umemDecoder) Decode(typeID iff.TypeID, r *io.SectionReader, context *iff.Decoder, logger log.FieldLogger) (iff.Chunk, error) {
	if err := iff.ExpectType(UncompressedMemoryType, typeID); err != nil {
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

// UMem is a FORM IFF chunk
type UMem struct {
	typeID iff.TypeID
	length int64
	raw    []byte
}

// TypeID ...
func (c *UMem) TypeID() iff.TypeID {
	return c.typeID
}

// Length ...
func (c *UMem) Length() int64 {
	return c.length
}

func (c *UMem) String() string {
	raw := iff.RenderBytes(c.raw, 16)
	return fmt.Sprintf("%s (length %d)\n%s", c.typeID, c.length, raw)
}
