package liveview

type tokenizer interface {
	Encode(any) (string, error)
	Decode(string, any) error
}
