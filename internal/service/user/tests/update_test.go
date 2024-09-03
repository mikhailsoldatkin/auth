package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/mikhailsoldatkin/auth/internal/repository"
	repoMocks "github.com/mikhailsoldatkin/auth/internal/repository/mocks"
	"github.com/mikhailsoldatkin/auth/internal/service/user"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userRepoMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *model.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = gofakeit.Int64()
		name = gofakeit.Name()
		req  = &model.User{
			ID:       id,
			Username: name,
		}
		wantErr = fmt.Errorf("repository error")
	)

	tests := []struct {
		name         string
		args         args
		err          error
		userRepoMock userRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: wantErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(wantErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepoMock(mc)
			service := user.NewMockUserService(userRepoMock)

			repoErr := service.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, repoErr)
		})
	}
}
