package mock

import "log"

type Handler struct {
	CallCount int
}

func NewHandler() *Handler {
	return &Handler{0}
}

func (h *Handler) Run(in <-chan bool) {
	go func() {
		for {
			_, ok := <-in
			if !ok {
				log.Println("EXIT")
				return
			}

			log.Println("CALL")
			h.CallCount += 1
		}
	}()
}
