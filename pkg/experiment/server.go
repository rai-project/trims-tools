package experiment

import (
	"context"

	tui "github.com/marcusolsson/tui-go"
)

type Server struct {
	ctx context.Context
}

func NewServer(ctx context.Context, ui tui.UI) *Server {
	return &Server{}
}

func (s *Server) Widget() tui.Widget {
	return nil
}
