package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	repoMocks "github.com/mikhailsoldatkin/auth/internal/repository/mocks"
	"github.com/mikhailsoldatkin/auth/internal/service/user"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()
	type userRepoMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *pb.ListRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id1    = gofakeit.Int64()
		name1  = gofakeit.Name()
		email1 = gofakeit.Email()
		role1  = gofakeit.RandomString([]string{"USER", "ADMIN"})
		now1   = time.Now()

		id2    = gofakeit.Int64()
		name2  = gofakeit.Name()
		email2 = gofakeit.Email()
		role2  = gofakeit.RandomString([]string{"USER", "ADMIN"})
		now2   = time.Now()

		wantResp = []*model.User{
			{
				ID:        id1,
				Name:      name1,
				Email:     email1,
				Role:      role1,
				CreatedAt: now1,
				UpdatedAt: now1,
			},
			{
				ID:        id2,
				Name:      name2,
				Email:     email2,
				Role:      role2,
				CreatedAt: now2,
				UpdatedAt: now2,
			},
		}
		req = &pb.ListRequest{
			Limit: 2,
		}
		wantErr = fmt.Errorf("repository error")
	)

	tests := []struct {
		name         string
		args         args
		want         []*model.User
		err          error
		userRepoMock userRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: wantResp,
			err:  nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.ListMock.Expect(ctx, req).Return(wantResp, nil)
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
			err:  wantErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.ListMock.Expect(ctx, req).Return(nil, wantErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepoMock(mc)
			service := user.NewMockService(userRepoMock)

			resp, repoErr := service.List(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, repoErr)
			require.Equal(t, tt.want, resp)
		})
	}
}
