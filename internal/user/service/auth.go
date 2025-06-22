package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/reddit/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openmyth/messgener/config"
	"openmyth/messgener/internal/user/entity"
	"openmyth/messgener/internal/user/repository"
	pb "openmyth/messgener/pb/user"
	"openmyth/messgener/pkg/metadata"
	"openmyth/messgener/util"
	"openmyth/messgener/util/database"
)

type authService struct {
	userRepo repository.UserRepository
	db       database.Executor

	pb.UnimplementedAuthServiceServer
}

// NewAuthService creates a new instance of the authService struct.
func NewAuthService(db database.Executor, userRepo repository.UserRepository) pb.AuthServiceServer {
	return &authService{
		db:       db,
		userRepo: userRepo,
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

	tkn, err := util.GenerateToken(&jwt.StandardClaims{
		Id:       user.ID.String,
		Audience: config.GetGlobalConfig().Platform,
	}, config.GetGlobalConfig().TokenTTL,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to generate token: %v", err)
	}

	return &pb.LoginResponse{
		User:  userToPb(user),
		Token: tkn,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	user, err := s.userRepo.RetrieveByPhone(ctx, s.db, req.GetPhone())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err)
	}

	tkn, err := util.GenerateToken(&jwt.StandardClaims{
		Id:       user.ID.String,
		Audience: config.GetGlobalConfig().Platform,
	}, config.GetGlobalConfig().TokenTTL,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to generate token: %v", err)
	}

	return &pb.ForgotPasswordResponse{
		Token: tkn,
	}, nil
}

func (s *authService) Disable2FA(ctx context.Context, _ *pb.Disable2FARequest) (*pb.Disable2FAResponse, error) {
	userCtx, ok := metadata.ExtractUserInfoFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user is invalid")
	}

	user, err := s.userRepo.RetrieveByID(ctx, s.db, userCtx.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err)
	}

	user.Enable2FA.Scan(false)

	if err := s.userRepo.UpdateInfoByID(ctx, s.db, userCtx.UserID, user); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to update user: %v", err)
	}

	return &pb.Disable2FAResponse{}, nil
}

func (s *authService) GenerateOTP(ctx context.Context, req *pb.GenerateOTPRequest) (*pb.GenerateOTPResponse, error) {
	userCtx, ok := metadata.ExtractUserInfoFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user is invalid")
	}

	user, err := s.userRepo.RetrieveByID(ctx, s.db, userCtx.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "openmyth",
		AccountName: user.Email.String,
		Secret:      []byte(user.OTPSecret.String),
		SecretSize:  16,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate OTP: %v", err)
	}

	return &pb.GenerateOTPResponse{
		OtpAuthUrl: key.URL(),
		Base32:     key.Secret(),
	}, nil
}

func (s *authService) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	userCtx, ok := metadata.ExtractUserInfoFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user is invalid")
	}

	user, err := s.userRepo.RetrieveByID(ctx, s.db, userCtx.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err)
	}

	valid := totp.Validate(req.GetOtp(), user.OTPSecret.String)
	if !valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid OTP")
	}

	return &pb.VerifyOTPResponse{}, nil
}

func (s *authService) ResetPassword(_ context.Context, _ *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) ChangePassword(_ context.Context, _ *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) Enable2FA(ctx context.Context, _ *pb.Enable2FARequest) (*pb.Enable2FAResponse, error) {
	userCtx, ok := metadata.ExtractUserInfoFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user is invalid")
	}

	user, err := s.userRepo.RetrieveByID(ctx, s.db, userCtx.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err)
	}

	user.Enable2FA.Scan(true)

	if err := s.userRepo.UpdateInfoByID(ctx, s.db, userCtx.UserID, user); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to update user: %v", err)
	}

	return &pb.Enable2FAResponse{}, nil
}
