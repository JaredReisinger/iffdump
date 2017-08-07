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
	logger.Debug("decoding CMem...")
	if typeID != CompressedMemoryType {
		return nil, fmt.Errorf("expected type ID of %q, got %q", CompressedMemoryType, typeID)
	}

	c := &CMem{
		typeID: typeID,
		length: r.Size(),
	}

	logger.WithField("IFhd", c).Debug("decoded IFhd...")

	return c, nil
}

// CMem is a FORM IFF chunk
type CMem struct {
	typeID iff.TypeID
	length int64
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
	return fmt.Sprintf("%s (length %d)", c.typeID, c.length)
}
