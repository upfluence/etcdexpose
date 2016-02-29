package handler

type Handler interface {
	Perform() error
}
