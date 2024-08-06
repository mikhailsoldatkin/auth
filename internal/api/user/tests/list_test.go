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
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestList(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.ListRequest
	}

	var (
		ctx       = context.Background()
		mc        = minimock.NewController(t)
		id1       = gofakeit.Int64()
		id2       = gofakeit.Int64()
		name1     = gofakeit.Name()
		name2     = gofakeit.Name()
		email1    = gofakeit.Email()
		email2    = gofakeit.Email()
		role1     = gofakeit.RandomString([]string{"USER", "ADMIN"})
		role2     = gofakeit.RandomString([]string{"USER", "ADMIN"})
		now       = gofakeit.Date()
		limit     = 2
		wantUsers = []*model.User{
			{
				ID:        id1,
				Name:      name1,
				Email:     email1,
				Role:      role1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        id2,
				Name:      name2,
				Email:     email2,
				Role:      role2,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}
		req = &pb.ListRequest{
			Limit: int64(limit),
		}
		wantResp = &pb.ListResponse{
			Users: []*pb.User{
				{
					Id:        id1,
					Name:      name1,
					Email:     email1,
					Role:      pb.Role(pb.Role_value[role1]),
					CreatedAt: timestamppb.New(now),
					UpdatedAt: timestamppb.New(now),
				},
				{
					Id:        id2,
					Name:      name2,
					Email:     email2,
					Role:      pb.Role(pb.Role_value[role2]),
					CreatedAt: timestamppb.New(now),
					UpdatedAt: timestamppb.New(now),
				},
			},
		}
		wantErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *pb.ListResponse
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
				mock.ListMock.Expect(ctx, req).Return(wantUsers, nil)
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
				mock.ListMock.Expect(ctx, req).Return(nil, wantErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(userServiceMock)

			resp, grpcErr := api.List(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, grpcErr)
			require.Equal(t, tt.want, resp)
		})
	}
}
