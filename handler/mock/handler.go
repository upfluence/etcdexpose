package mock

type Handler struct {
	CallCount int
}

func NewHandler() *Handler {
	return &Handler{0}
}

func (h *Handler) Perform() {
	h.CallCount += 1
}
