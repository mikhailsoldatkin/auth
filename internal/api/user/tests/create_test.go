package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	userAPI "github.com/mikhailsoldatkin/auth/internal/api/user"
	"github.com/mikhailsoldatkin/auth/internal/service"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	"github.com/stretchr/testify/require"

	serviceMocks "github.com/mikhailsoldatkin/auth/internal/service/mocks"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id          = gofakeit.Int64()
		name        = gofakeit.Name()
		email       = gofakeit.Email()
		role        = gofakeit.RandomString([]string{"USER", "ADMIN"})
		currentTime = time.Now()

		//serviceErr = fmt.Errorf("service error")

		req = &pb.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        "12345678",
			PasswordConfirm: "12345678",
			Role:            pb.Role(pb.Role_value[role]),
		}

		res = &pb.CreateResponse{
			Id: id,
		}

		user = &model.User{
			//ID:    1,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
	)
	defer t.Cleanup(mc.Finish)

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
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(id, nil)
				return mock
			},
		},
		//{
		//	name: "service error case",
		//	args: args{
		//		ctx: ctx,
		//		req: req,
		//	},
		//	want: nil,
		//	err:  serviceErr,
		//	userServiceMock: func(mc *minimock.Controller) service.UserService {
		//		mock := serviceMocks.NewUserServiceMock(mc)
		//		mock.CreateMock.Expect(ctx, user).Return(0, serviceErr)
		//		return mock
		//	},
		//},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(noteServiceMock)

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
