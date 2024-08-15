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
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.DeleteRequest
	}

	var (
		ctx      = context.Background()
		mc       = minimock.NewController(t)
		id       = gofakeit.Int64()
		req      = &pb.DeleteRequest{Id: id}
		wantResp = &emptypb.Empty{}
		wantErr  = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
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
				mock.DeleteMock.Expect(ctx, id).Return(nil)
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
				mock.DeleteMock.Expect(ctx, id).Return(wantErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(userServiceMock)

			resp, grpcErr := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, grpcErr)
			require.Equal(t, tt.want, resp)
		})
	}
}
