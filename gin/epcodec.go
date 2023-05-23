package gin

type JsonCodecEp struct{}

func (d JsonCodecEp) Codec() Codec {
	return JSONCodec{}
}
