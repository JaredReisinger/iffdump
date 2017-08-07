package quetzal

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"

	"github.com/JaredReisinger/iffdump/iff"
)

const (
	InteractiveFictionHeaderType iff.TypeID = "IFhd"
)

type ifhdDecoder struct{}

func (d *ifhdDecoder) Decode(typeID iff.TypeID, r *io.SectionReader, context *iff.Decoder, logger log.FieldLogger) (iff.Chunk, error) {
	logger.Debug("decoding IFhd...")
	if typeID != InteractiveFictionHeaderType {
		return nil, fmt.Errorf("expected type ID of %q, got %q", InteractiveFictionHeaderType, typeID)
	}

	releaseNum, err := iff.ReadUint16(r)
	if err != nil {
		return nil, err
	}

	serialNum, err := iff.ReadBytes(r, 6)
	if err != nil {
		return nil, err
	}

	checksum, err := iff.ReadUint16(r)
	if err != nil {
		return nil, err
	}

	pc, err := iff.ReadUint24(r)
	if err != nil {
		return nil, err
	}

	c := &InteractiveFictionHeader{
		typeID:         typeID,
		length:         r.Size(),
		ReleaseNumber:  releaseNum,
		SerialNumber:   serialNum,
		Checksum:       checksum,
		ProgramCounter: pc,
	}

	logger.WithField("IFhd", c).Debug("decoded IFhd...")

	return c, nil
}

// InteractiveFictionHeader is a FORM IFF chunk
type InteractiveFictionHeader struct {
	typeID         iff.TypeID
	length         int64
	ReleaseNumber  int
	SerialNumber   []byte
	Checksum       int
	ProgramCounter int
}

// TypeID ...
func (c *InteractiveFictionHeader) TypeID() iff.TypeID {
	return c.typeID
}

// Length ...
func (c *InteractiveFictionHeader) Length() int64 {
	return c.length
}

func (c *InteractiveFictionHeader) String() string {
	return fmt.Sprintf("%s (length %d): release %d, serial %v (%s), checksum %d, pc %d", c.typeID, c.length, c.ReleaseNumber, c.SerialNumber, string(c.SerialNumber), c.Checksum, c.ProgramCounter)
}
