package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	playerpb "study/api"
	"study/internal/usecase"
)

type HandlerGRPC struct {
	playerpb.UnimplementedPlayerServiceServer
	uc *usecase.Usecase
}

func NewGRPC(uc *usecase.Usecase) *HandlerGRPC {
	return &HandlerGRPC{uc: uc}
}

func (h *HandlerGRPC) CreatePlayer(
	ctx context.Context,
	req *playerpb.CreatePlayerRequest,
) (*playerpb.CreatePlayerResponse, error) {

	if err := h.uc.CreatePlayer(ctx, req.GetName(), req.GetBalance()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create: %v", err)
	}

	return &playerpb.CreatePlayerResponse{
		Player: &playerpb.Player{
			Name:    req.GetName(),
			Balance: req.GetBalance(),
		},
	}, nil
}

func (h *HandlerGRPC) GetPlayer(
	ctx context.Context,
	req *playerpb.GetPlayerRequest,
) (*playerpb.GetPlayerResponse, error) {

	p, err := h.uc.GetPlayer(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get: %v", err)
	}

	return &playerpb.GetPlayerResponse{
		Player: &playerpb.Player{
			Id:      p.Id,
			Name:    p.Name,
			Balance: p.Balance,
		},
	}, nil
}
