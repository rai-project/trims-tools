package micro

import tui "github.com/marcusolsson/tui-go"

type Server struct {
}

func NewServer(ui tui.UI) *Server {
	return &Server{}
}

func (s *Server) Widget() tui.Widget {
	return nil
}
