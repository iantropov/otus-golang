package internalhttp

import (
	"context"
)

type helloRequest struct {
	ID string `json:"id"`
}

type helloResponse struct {
	Success bool `json:"success"`
}

func (s *Server) getHello(ctx context.Context, req helloRequest) (helloResponse, error) {
	s.logger.Info("got /hello request: " + req.ID)

	return helloResponse{
		Success: true,
	}, nil
}
