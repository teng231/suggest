package compression

import (
	"io"

	"github.com/teng231/suggest/pkg/store"
)

// VBEncoder returns new instance of vbEnc that encodes posting list using
// delta encoding
// variable length byte string compression
func VBEncoder() Encoder {
	return &vbEnc{}
}

// VBDecoder decodes given bytes array to posting list which was encoded by VBEncoder
func VBDecoder() Decoder {
	return &vbEnc{}
}

// vbEnc implements VBEncoder and VBDecoder
type vbEnc struct{}

// Encode encodes the given positing list into the buf array
// Returns a number of written bytes
func (b *vbEnc) Encode(list []uint32, out store.Output) (int, error) {
	return varIntEncode(list, out, 0)
}

// Decode decodes the given byte array to the buf list
// Returns a number of elements encoded
func (b *vbEnc) Decode(in store.Input, buf []uint32) (int, error) {
	return varIntDecode(in, buf, 0)
}

func varIntEncode(list []uint32, out store.Output, prev uint32) (int, error) {
	total := 0

	for _, v := range list {
		delta := v - prev
		prev = v

		n, err := out.WriteVUInt32(delta)
		total += n

		if err != nil {
			if err == io.EOF {
				return total, nil
			}

			return total, err
		}
	}

	return total, nil
}

func varIntDecode(in store.Input, buf []uint32, prev uint32) (int, error) {
	total := 0

	for total < len(buf) {
		v, err := in.ReadVUInt32()

		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return total, nil
			}

			return total, err
		}

		prev += v
		buf[total] = prev
		total++
	}

	return total, nil
}
