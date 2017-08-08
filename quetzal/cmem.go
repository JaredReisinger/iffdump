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
		typeID:   typeID,
		length:   length,
		raw:      raw,
		expanded: expand(raw),
	}

	return c, nil
}

func expand(raw []byte) []byte {
	length := len(raw)
	expanded := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		// Quetzal doc 3.2: The data is compressed by exclusive-oring the
		// current contents of dynamic memory with the original (from the
		// original story file). The result is then compressed with a simple
		// run-length scheme: a non-zero byte in the output represents the byte
		// itself, but a zero byte is followed by a length byte, and the pair
		// represent a block of n+1 zero bytes, where n is the value of the
		// length byte.
		//
		// Note: this means that 0x00 0x00 represents *one* 0x00 byte, hence the
		// greater-than-*or-equal* in the inner loop.
		if raw[i] == 0 {
			i++
			if i < length {
				c := int(raw[i])
				for j := 0; j <= c; j++ {
					expanded = append(expanded, 0)
				}
			}
		} else {
			expanded = append(expanded, raw[i])
		}
	}

	return expanded
}

// CMem is a FORM IFF chunk
type CMem struct {
	typeID   iff.TypeID
	length   int64
	raw      []byte
	expanded []byte
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
	expanded := iff.RenderBytes(c.expanded, 16)
	return fmt.Sprintf("%s (length %d)\n      RAW:\n%s\n      EXPANDED:\n%s", c.typeID, c.length, raw, expanded)
}
