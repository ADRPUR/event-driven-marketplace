package grpc

import (
	"context"
	"github.com/ADRPUR/event-driven-marketplace/internal/product/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ADRPUR/event-driven-marketplace/api/proto/product/v1"
	"github.com/ADRPUR/event-driven-marketplace/internal/product/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	pb.UnimplementedProductServiceServer
	svc *service.ProductService
}

func NewGRPCServer(svc *service.ProductService) pb.ProductServiceServer {
	return &grpcServer{svc: svc}
}

func (s *grpcServer) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	p, err := s.svc.Create(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{Product: toProto(p)}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.svc.Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{Product: toProto(p)}, nil
}

func (s *grpcServer) ListProducts(ctx context.Context, in *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	list, err := s.svc.List(ctx, int(in.Page), int(in.PageSize))
	if err != nil {
		return nil, err
	}
	out := make([]*pb.Product, len(list))
	for i, p := range list {
		out[i] = toProto(&p)
	}
	return &pb.ListProductsResponse{Products: out, Total: int32(len(out))}, nil
}

func (s *grpcServer) UpdateProduct(ctx context.Context, in *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	p, err := s.svc.Update(ctx, in.Id, &in.Name, &in.Description, &in.Price)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductResponse{Product: toProto(p)}, nil
}

func (s *grpcServer) DeleteProduct(ctx context.Context, in *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.svc.Delete(ctx, in.Id)
}

// helper
func toProto(m *model.Product) *pb.Product {
	return &pb.Product{
		Id:          m.ID.String(),
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		CreatedAt:   timestamppb.New(m.CreatedAt),
	}
}
