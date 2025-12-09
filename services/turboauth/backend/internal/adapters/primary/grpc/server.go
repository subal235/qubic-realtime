package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "turboauth/api/proto/api/proto"
	"turboauth/internal/domain/auth"
	"turboauth/pkg/metrics"
)

// Server implements the gRPC AuthService
type Server struct {
	pb.UnimplementedAuthServiceServer
	authService *auth.Service
}

// NewServer creates a new gRPC server
func NewServer(authService *auth.Service) *Server {
	return &Server{
		authService: authService,
	}
}

// GetStatus retrieves authentication status for a wallet
func (s *Server) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	start := time.Now()

	walletAuth, err := s.authService.GetStatus(ctx, req.WalletAddress)
	if err != nil {
		metrics.GRPCRequestsTotal.WithLabelValues("GetStatus", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to get status: %v", err)
	}

	metrics.GRPCRequestsTotal.WithLabelValues("GetStatus", "success").Inc()
	metrics.GRPCRequestDuration.WithLabelValues("GetStatus").Observe(time.Since(start).Seconds())

	return &pb.GetStatusResponse{
		Status:          string(walletAuth.Status),
		TrustScore:      int32(walletAuth.TrustScore),
		UpdatedAt:       walletAuth.UpdatedAt.Unix(),
		ContractAddress: walletAuth.ContractAddress,
	}, nil
}

// SetStatus updates authentication status (admin only)
func (s *Server) SetStatus(ctx context.Context, req *pb.SetStatusRequest) (*pb.SetStatusResponse, error) {
	start := time.Now()

	setReq := &auth.SetStatusRequest{
		WalletAddress:  req.WalletAddress,
		Status:         auth.AuthStatus(req.Status),
		TrustScore:     int(req.TrustScore),
		AdminSignature: req.AdminSignature,
	}

	txHash, err := s.authService.SetStatus(ctx, setReq)
	if err != nil {
		metrics.GRPCRequestsTotal.WithLabelValues("SetStatus", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to set status: %v", err)
	}

	metrics.GRPCRequestsTotal.WithLabelValues("SetStatus", "success").Inc()
	metrics.GRPCRequestDuration.WithLabelValues("SetStatus").Observe(time.Since(start).Seconds())

	return &pb.SetStatusResponse{
		Success: true,
		TxHash:  txHash,
		Message: "Status updated successfully",
	}, nil
}

// VerifyWallet verifies a wallet signature
func (s *Server) VerifyWallet(ctx context.Context, req *pb.VerifyWalletRequest) (*pb.VerifyWalletResponse, error) {
	start := time.Now()

	verifyReq := &auth.VerifyRequest{
		WalletAddress: req.WalletAddress,
		Signature:     req.Signature,
		Message:       req.Message,
	}

	walletAuth, verified, err := s.authService.VerifyWallet(ctx, verifyReq)
	if err != nil {
		metrics.GRPCRequestsTotal.WithLabelValues("VerifyWallet", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to verify wallet: %v", err)
	}

	metrics.GRPCRequestsTotal.WithLabelValues("VerifyWallet", "success").Inc()
	metrics.GRPCRequestDuration.WithLabelValues("VerifyWallet").Observe(time.Since(start).Seconds())

	return &pb.VerifyWalletResponse{
		Verified:   verified,
		Status:     string(walletAuth.Status),
		TrustScore: int32(walletAuth.TrustScore),
	}, nil
}

// BatchGetStatus retrieves status for multiple wallets (high-performance)
func (s *Server) BatchGetStatus(ctx context.Context, req *pb.BatchGetStatusRequest) (*pb.BatchGetStatusResponse, error) {
	start := time.Now()

	statuses, err := s.authService.BatchGetStatus(ctx, req.WalletAddresses)
	if err != nil {
		metrics.GRPCRequestsTotal.WithLabelValues("BatchGetStatus", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to batch get status: %v", err)
	}

	metrics.GRPCRequestsTotal.WithLabelValues("BatchGetStatus", "success").Inc()
	metrics.GRPCRequestDuration.WithLabelValues("BatchGetStatus").Observe(time.Since(start).Seconds())

	// Convert to protobuf responses
	pbStatuses := make([]*pb.GetStatusResponse, len(statuses))
	for i, s := range statuses {
		pbStatuses[i] = &pb.GetStatusResponse{
			Status:          string(s.Status),
			TrustScore:      int32(s.TrustScore),
			UpdatedAt:       s.UpdatedAt.Unix(),
			ContractAddress: s.ContractAddress,
		}
	}

	return &pb.BatchGetStatusResponse{
		Statuses: pbStatuses,
	}, nil
}
