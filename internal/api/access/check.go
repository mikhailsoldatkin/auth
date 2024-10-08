package access

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	pb "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Check verifies access to the specified endpoint.
func (i *Implementation) Check(ctx context.Context, req *pb.CheckRequest) (*emptypb.Empty, error) {
	err := i.accessService.Check(ctx, req.GetEndpoint())
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}
	return &emptypb.Empty{}, nil
}
