package mock

type Handler struct {
	CallCount int
	toReturn  error
}

func NewHandler(err error) *Handler {
	return &Handler{0, err}
}

func (h *Handler) Perform() error {
	h.CallCount += 1
	return h.toReturn
}
