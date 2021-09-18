package handler

type ConsoleHandler struct {
	HelloHandler HelloHandler
}

func NewConsoleHandler(helloHandler HelloHandler) *ConsoleHandler {
	return &ConsoleHandler{
		HelloHandler: helloHandler,
	}
}
