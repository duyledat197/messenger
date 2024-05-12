package service

import (
	"context"
	"database/sql"
	"errors"
	"os"

	"github.com/reddit/jwt-go"
	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openmyth/messgener/internal/user/entity"
	"openmyth/messgener/internal/user/repository"
	pb "openmyth/messgener/pb/user"
	"openmyth/messgener/util"
	"openmyth/messgener/util/database"
)

type authService struct {
	userRepo     repository.UserRepository
	db           database.Executor
	tknGenerator util.JWTAuthenticator

	pb.UnimplementedAuthServiceServer
}

// NewAuthService creates a new instance of the authService struct.
func NewAuthService(db database.Executor, tknGenerator util.JWTAuthenticator, userRepo repository.UserRepository) pb.AuthServiceServer {
	return &authService{
		db:           db,
		userRepo:     userRepo,
		tknGenerator: tknGenerator,
	}
}

// userToPb converts a User entity to a User protobuf message.
// It takes a pointer to a User entity as a parameter and returns a pointer to a User protobuf message.
func userToPb(user *entity.User) *pb.User {
	return &pb.User{
		Username: user.Username.String,
		Email:    user.Email.String,
		Facebook: user.Facebook.String,
		Discord:  user.Discord.String,
		Github:   user.Github.String,
		Gmail:    user.Google.String,
	}
}

// Login handles the login process for a user.
func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.userRepo.RetrieveByUserName(ctx, s.db, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err)
	}

	if err := util.CheckPassword(req.GetPassword(), user.HashedPassword.String); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	tkn, err := s.tknGenerator.Generate(&jwt.StandardClaims{
		Id:       user.ID.String,
		Audience: os.Getenv("PLATFORM"), // TODO: make this configurable
	}, cast.ToDuration(os.Getenv("TOKEN_TTL")), // TODO: make this configurable
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to generate token: %v", err)
	}

	return &pb.LoginResponse{
		User:  userToPb(user),
		Token: tkn,
	}, nil
}

func (s *authService) ForgotPassword(_ context.Context, _ *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) Disable2FA(_ context.Context, _ *pb.Disable2FARequest) (*pb.Disable2FAResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) GenerateOTP(_ context.Context, _ *pb.GenerateOTPRequest) (*pb.GenerateOTPResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) VerifyOTP(_ context.Context, _ *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) ResetPassword(_ context.Context, _ *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) ChangePassword(_ context.Context, _ *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) Enable2FA(_ context.Context, _ *pb.Enable2FARequest) (*pb.Enable2FAResponse, error) {
	panic("not implemented") // TODO: Implement
}
