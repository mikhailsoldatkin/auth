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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.UpdateRequest
	}

	var (
		ctx          = context.Background()
		mc           = minimock.NewController(t)
		id           = gofakeit.Int64()
		name         = gofakeit.Name()
		validEmail   = gofakeit.Email()
		invalidEmail = "email@com"

		validReq = &pb.UpdateRequest{
			Id:    id,
			Name:  wrapperspb.String(name),
			Email: wrapperspb.String(validEmail),
		}
		invalidReq = &pb.UpdateRequest{
			Id:    id,
			Name:  wrapperspb.String(name),
			Email: wrapperspb.String(invalidEmail),
		}
		wantResp = &emptypb.Empty{}

		wantErr      = fmt.Errorf("service error")
		wantEmailErr = status.Errorf(codes.InvalidArgument, "invalid email format: %v", invalidEmail)
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
				req: validReq,
			},
			want: wantResp,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, validReq).Return(nil)
				return mock
			},
		},
		{
			name: "invalid email format",
			args: args{
				ctx: ctx,
				req: invalidReq,
			},
			want: nil,
			err:  wantEmailErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
		{
			name: "service error",
			args: args{
				ctx: ctx,
				req: validReq,
			},
			want: nil,
			err:  customerrors.ConvertError(wantErr),
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, validReq).Return(wantErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(userServiceMock)

			resp, grpcErr := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, grpcErr)
			require.Equal(t, tt.want, resp)
		})
	}
}
