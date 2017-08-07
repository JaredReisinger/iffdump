package iff

import (
	"io"

	log "github.com/sirupsen/logrus"
)

// The raw data in an IFF chunk is a 4-byte ASCII (0x20-0x7e) string/FourCC type
// ID, followed by a uint32 length, and then that many bytes of data.  The next
// chunk (if any) follows the data, possibly padded by a single zero byte to
// ensure that it starts on an even address boundary.

// TypeID is a type-specific value for the four-byte ASCII type ID.
type TypeID string

// func (t TypeID) String() string {
// 	return string(t)
// }

// // Well-known IFF type IDs
// const (
// 	List      TypeID = "LIST"
// 	Cat       TypeID = "CAT "
// 	Property  TypeID = "PROP"
// )

// ReadAtSeeker encapsulates the same functionality as io.SectionReader.
// Conveniently, it is *also* implemented by os.File!
type ReadAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

// A Chunk is the bare minimum interface that all chuncks must provide.  Note
// that the IFF documentation uses an unsigned, 32-bit integer for chunk
// lengths, but Go uses a signed, 64-bit integer for seeking and many
// Reader/Writer APIs. Rather than deal with casting all over the place, we just
// use Go's native int64 since it can represent a uint32 without problem.  (The
// only downside is that it *will* allow values that a uint32 would not... but
// since we're currently only reading the IFF, this isn't really an issue.)
type Chunk interface {
	TypeID() TypeID
	Length() int64
	String() string
}

// ChunkDecoder represents the functionality required for decoding a chunk.
// When Decode() is called, the type ID and length have *already* been
// read/parsed, and the SectionReader is just the data for the chunk.
type ChunkDecoder interface {
	Decode(typeID TypeID, r *io.SectionReader, context *Decoder, logger log.FieldLogger) (Chunk, error)
}
