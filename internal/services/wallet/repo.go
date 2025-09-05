package walletservice

type Reader interface {
}

type Writer interface {
}

type ReadWriter interface {
	Reader
	Writer
}
