package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/mikhailsoldatkin/auth/internal/repository"
	repoMocks "github.com/mikhailsoldatkin/auth/internal/repository/mocks"
	"github.com/mikhailsoldatkin/auth/internal/service/user"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userRepoMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = gofakeit.RandomString([]string{"USER", "ADMIN"})
		now   = time.Now()
		req   = id

		wantUser = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: now,
			UpdatedAt: now,
		}
		wantErr = fmt.Errorf("repository error")
	)

	tests := []struct {
		name         string
		args         args
		want         *model.User
		err          error
		userRepoMock userRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: wantUser,
			err:  nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, req).Return(wantUser, nil)
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
				mock.GetMock.Expect(ctx, req).Return(nil, wantErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepoMock(mc)
			service := user.NewMockService(userRepoMock)

			resp, repoErr := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, repoErr)
			require.Equal(t, tt.want, resp)
		})
	}
}
