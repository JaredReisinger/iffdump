package quetzal

import (
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/JaredReisinger/iffdump/iff"
)

const (
	StacksType iff.TypeID = "Stks"
)

type stksDecoder struct{}

func (d *stksDecoder) Decode(typeID iff.TypeID, r *io.SectionReader, context *iff.Decoder, logger log.FieldLogger) (iff.Chunk, error) {
	if err := iff.ExpectType(StacksType, typeID); err != nil {
		return nil, err
	}

	length := r.Size()

	frames := make([]*StackFrame, 0)

	off, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	for off < length {
		frame, err := readFrame(r, logger)
		if err != nil {
			return nil, err
		}

		frames = append(frames, frame)

		off, err = r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}

	c := &Stks{
		typeID: typeID,
		length: length,
		Frames: frames,
	}

	return c, nil
}

func readFrame(r io.Reader, logger log.FieldLogger) (*StackFrame, error) {
	// logger.Debug("reading stack frame...")

	returnPC, err := iff.ReadUint24(r)
	if err != nil {
		return nil, err
	}

	flags, err := iff.ReadUint8(r)
	if err != nil {
		return nil, err
	}

	result, err := iff.ReadUint8(r)
	if err != nil {
		return nil, err
	}

	arguments, err := iff.ReadUint8(r)
	if err != nil {
		return nil, err
	}

	evalWords, err := iff.ReadUint16(r)
	if err != nil {
		return nil, err
	}

	// logger.Debug("reading stack frame locals...")

	// read 'v' words for local variables... don't know how to find 'v'!
	v := 0
	locals := make([]uint16, 0, v)
	for i := 0; i < v; i++ {
		local, err := iff.ReadUint16(r)
		if err != nil {
			return nil, err
		}
		locals = append(locals, local)
	}

	// logger.Debug("reading stack frame evaluation stack...")

	evalStack := make([]uint16, 0, evalWords)

	for i := 0; i < int(evalWords); i++ {
		val, err := iff.ReadUint16(r)
		if err != nil {
			return nil, err
		}
		evalStack = append(evalStack, val)
	}

	// logger.Debug("reading stack frame... returning!")

	frame := &StackFrame{
		ReturnPC:        returnPC,
		Flags:           flags,
		Result:          result,
		Arguments:       arguments,
		EvaluationWords: evalWords,
		LocalVariables:  locals,
		EvaluationStack: evalStack,
	}

	return frame, nil
}

// Stks is a FORM IFF chunk
type Stks struct {
	typeID iff.TypeID
	length int64
	Frames []*StackFrame
}

type StackFrame struct {
	ReturnPC        uint32
	Flags           uint8
	Result          uint8
	Arguments       uint8
	EvaluationWords uint16
	LocalVariables  []uint16
	EvaluationStack []uint16
}

// TypeID ...
func (c *Stks) TypeID() iff.TypeID {
	return c.typeID
}

// Length ...
func (c *Stks) Length() int64 {
	return c.length
}

func (c *Stks) String() string {
	lines := make([]string, 0)
	for _, f := range c.Frames {
		lines = append(lines, fmt.Sprintf("%+v", f))
	}
	return fmt.Sprintf("%s (length %d) frames:\n        %s", c.typeID, c.length, strings.Join(lines, "\n        "))
}
