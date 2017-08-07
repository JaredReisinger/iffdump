package iff

import (
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Well-known IFF type IDs
const (
	FormType TypeID = "FORM"
)

type formDecoder struct{}

func (d *formDecoder) Decode(typeID TypeID, r *io.SectionReader, context *Decoder, logger log.FieldLogger) (Chunk, error) {
	if err := ExpectType(FormType, typeID); err != nil {
		return nil, err
	}

	formID, err := ReadTypeID(r)
	if err != nil {
		return nil, err
	}

	c := &Form{
		typeID: typeID,
		length: r.Size(),
		FormID: formID,
	}

	logger.WithField("form", c).Debugf("decoding %s...", typeID)

	// Now loop through all of the child chunks? or do we delegate to a
	// subtype-specific decoder?  The definition of FORM says that the record
	// type is "followed by nested chunks specifying the record fields"
	var off int64 = 4 // we've only read the subtype...
	limit := r.Size()

	for off < limit {
		logger.WithFields(log.Fields{
			"offset": off,
			"limit":  limit,
		}).Debugf("decoding next %s chunk...", typeID)
		child, err := context.Decode(r, logger)
		if err != nil {
			return nil, err
		}

		c.Chunks = append(c.Chunks, child)

		off += 8 + child.Length()
		if off%2 == 1 {
			off++ // pad to even address!
		}
		off, err = r.Seek(off, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Form is a FORM IFF chunk
type Form struct {
	typeID TypeID
	length int64
	FormID TypeID
	Chunks []Chunk
}

// TypeID ...
func (c *Form) TypeID() TypeID {
	return c.typeID
}

// Length ...
func (c *Form) Length() int64 {
	return c.length
}

func (c *Form) String() string {
	lines := make([]string, 0, 1+len(c.Chunks))

	lines = append(lines, fmt.Sprintf("%s %s (length %d, %d children):", c.typeID, c.FormID, c.length, len(c.Chunks)))

	for _, child := range c.Chunks {
		lines = append(lines, fmt.Sprintf("  - %s", child.String()))
	}

	return strings.Join(lines, "\n")
}
