package mock

type Handler struct {
	CallCount int
}

func NewHandler() *Handler {
	return &Handler{0}
}

func (h *Handler) Run(in <-chan bool) {
	go func() {
		_, ok := <-in
		if !ok {
			return
		}
		h.CallCount += 1
	}()
}
