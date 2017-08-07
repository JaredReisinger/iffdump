package quetzal

import "github.com/JaredReisinger/iffdump/iff"

// RegisterDecoders registers the defined decoders for use by iff.Decode.
func RegisterDecoders(decoders map[iff.TypeID]iff.ChunkDecoder) {
	decoders[InteractiveFictionHeaderType] = &ifhdDecoder{}
	decoders[CompressedMemoryType] = &cmemDecoder{}
}
