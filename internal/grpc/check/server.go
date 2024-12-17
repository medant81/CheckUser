package checkgrpc

import (
	"CheckUser/internal/models"
	checkUserV1 "CheckUser/protos/gen/go/checkuser"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServisesCheck interface {
	TokenServises(ctx context.Context, username string, password string, appID int) (token string, err error)
	CheckUsersServises(ctx context.Context, tokenRes string, usersId []int64) (usersRes models.UsersResult, err error)
}

type serverAPI struct {
	checkUserV1.UnimplementedCheckUsersServer
	check ServisesCheck
}

func Register(gRPCServer *grpc.Server, check ServisesCheck) {
	checkUserV1.RegisterCheckUsersServer(gRPCServer, &serverAPI{check: check})
}

// Должен совпадать с proto
func (s *serverAPI) Token(ctx context.Context, req *checkUserV1.TokenRequest) (*checkUserV1.TokenResponse, error) {

	//username
	//password
	if req.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.check.TokenServises(ctx, req.GetUsername(), req.GetPassword(), 1)
	if err != nil {
		// TODO: добавить обработку ошибок
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &checkUserV1.TokenResponse{
		Token: token,
	}, nil
}

// Должен совпадать с proto
func (s *serverAPI) CheckUsers(ctx context.Context, req *checkUserV1.CheckUsersRequest) (*checkUserV1.CheckUsersResponse, error) {

	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	if len(req.GetUsers()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "users is required")
	}

	usersRes, err := s.check.CheckUsersServises(ctx, req.GetToken(), req.GetUsers())

	usersSlice := make([]*checkUserV1.TypeUsers, 0)
	//usersMap := make(map[int64]bool)
	for key, value := range usersRes.Users {
		user := checkUserV1.TypeUsers{
			Id:    key,
			Check: value,
		}
		usersSlice = append(usersSlice, &user)
		//usersMap[key] = value
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &checkUserV1.CheckUsersResponse{
		Users: usersSlice,
	}, nil
}

/*func (s *serverAPI) mustEmbedUnimplementedCheckUsersServer() {
	//TODO implement me
	panic("implement me")
}*/
