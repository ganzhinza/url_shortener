package db

type Interface interface {
	MakeShort(url string) (string, error)
	GetOriginal(url string) (string, error)
	URLSize() uint
}
