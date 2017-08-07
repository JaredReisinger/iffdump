package quetzal

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"

	"github.com/JaredReisinger/iffdump/iff"
)

const (
	TranscriptType iff.TypeID = "TxHs"
)

type txhsDecoder struct{}

func (d *txhsDecoder) Decode(typeID iff.TypeID, r *io.SectionReader, context *iff.Decoder, logger log.FieldLogger) (iff.Chunk, error) {
	if err := iff.ExpectType(TranscriptType, typeID); err != nil {
		return nil, err
	}

	length := r.Size()

	raw := make([]byte, length)
	_, err := io.ReadFull(r, raw)
	if err != nil {
		return nil, err
	}

	c := &TxHs{
		typeID: typeID,
		length: length,
		raw:    raw,
	}

	return c, nil
}

// TxHs is a FORM IFF chunk
type TxHs struct {
	typeID iff.TypeID
	length int64
	raw    []byte
}

// TypeID ...
func (c *TxHs) TypeID() iff.TypeID {
	return c.typeID
}

// Length ...
func (c *TxHs) Length() int64 {
	return c.length
}

func (c *TxHs) String() string {
	// raw := iff.RenderBytes(c.raw, 16)
	// raw := string(c.raw)
	// return fmt.Sprintf("%s (length %d)\n%s", c.typeID, c.length, raw)
	return fmt.Sprintf("%s (length %d)", c.typeID, c.length)
}
