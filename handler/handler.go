package handler

type Handler interface {
	Run(<-chan bool)
}
