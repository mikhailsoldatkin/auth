package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	userAPI "github.com/mikhailsoldatkin/auth/internal/api/user"
	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service"
	serviceMocks "github.com/mikhailsoldatkin/auth/internal/service/mocks"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.GetRequest
	}

	var (
		ctx   = context.Background()
		mc    = minimock.NewController(t)
		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = gofakeit.RandomString([]string{"USER", "ADMIN"})
		now   = time.Now()
		req   = &pb.GetRequest{
			Id: id,
		}
		wantUser = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: now,
			UpdatedAt: now,
		}
		wantResp = &pb.GetResponse{
			User: &pb.User{
				Id:        id,
				Name:      name,
				Email:     email,
				Role:      pb.Role(pb.Role_value[role]),
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
			},
		}
		wantErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *pb.GetResponse
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
				mock.GetMock.Expect(ctx, id).Return(wantUser, nil)
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
				mock.GetMock.Expect(ctx, id).Return(nil, wantErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(userServiceMock)

			resp, grpcErr := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, grpcErr)
			require.Equal(t, tt.want, resp)
		})
	}
}
