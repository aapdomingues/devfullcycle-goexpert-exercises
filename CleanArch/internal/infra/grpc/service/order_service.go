package service

import (
	"context"

	"CleanArch/internal/infra/grpc/pb"
	"CleanArch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrdersUseCase usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(context.Context, *pb.Blank) (*pb.OrderList, error) {
	listOrdersDto, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	orders := []*pb.Order{}
	for _, orderDto := range listOrdersDto {
		orders = append(orders, &pb.Order{
			Id:         orderDto.ID,
			Price:      float32(orderDto.Price),
			Tax:        float32(orderDto.Tax),
			FinalPrice: float32(orderDto.FinalPrice),
		})
	}

	return &pb.OrderList{Orders: orders}, nil
}
