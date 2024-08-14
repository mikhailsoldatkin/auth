package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	userAPI "github.com/mikhailsoldatkin/auth/internal/api/user"
	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service"
	serviceMocks "github.com/mikhailsoldatkin/auth/internal/service/mocks"
	"github.com/mikhailsoldatkin/auth/internal/service/user/converter"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.CreateRequest
	}

	var (
		ctx           = context.Background()
		mc            = minimock.NewController(t)
		id            = gofakeit.Int64()
		name          = gofakeit.Name()
		email         = gofakeit.Email()
		role          = gofakeit.RandomString([]string{"USER", "ADMIN"})
		password      = "12345678"
		wrongPassword = "123456789"
		wrongEmail    = "some@com"
		req           = &pb.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            pb.Role(pb.Role_value[role]),
		}
		invalidPasswordReq = &pb.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: wrongPassword,
			Role:            pb.Role(pb.Role_value[role]),
		}
		invalidEmailReq = &pb.CreateRequest{
			Name:            name,
			Email:           wrongEmail,
			Password:        password,
			PasswordConfirm: password,
			Role:            pb.Role(pb.Role_value[role]),
		}

		wantResp = &pb.CreateResponse{
			Id: id,
		}
		//wantUser = &model.User{
		//	Name:  name,
		//	Email: email,
		//	Role:  role,
		//}
		wantErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *pb.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: wantResp,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.FromProtobufToServiceCreate(req)).Return(id, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  customerrors.ConvertError(wantErr),
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, converter.FromProtobufToServiceCreate(req)).Return(0, wantErr)
				return mock
			},
		},
		{
			name: "invalid password case",
			args: args{
				ctx: ctx,
				req: invalidPasswordReq,
			},
			want: nil,
			err:  status.Errorf(codes.InvalidArgument, "password validation failed: passwords don't match"),
			userServiceMock: func(_ *minimock.Controller) service.UserService {
				return nil
			},
		},
		{
			name: "invalid email case",
			args: args{
				ctx: ctx,
				req: invalidEmailReq,
			},
			want: nil,
			err:  status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid email format: %s", wrongEmail)),
			userServiceMock: func(_ *minimock.Controller) service.UserService {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(userServiceMock)

			resp, grpcErr := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, grpcErr)
			require.Equal(t, tt.want, resp)
		})
	}
}
