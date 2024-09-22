package liveview

type HttpError interface {
	Code() int
	Error() string
}
