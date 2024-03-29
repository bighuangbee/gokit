package userAccess

import (
	"context"
	"errors"
	pb "github.com/bighuangbee/gokit/api/common/v1"
	kitKratos "github.com/bighuangbee/gokit/kratos"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"strconv"
)

// CheckToken Check Token middleware
func CheckToken(access IUserAccess) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if token := tr.RequestHeader().Get("jwtToken"); token != "" {

					userClaimsm, err := access.Validate(token)
					if err != nil {
						return nil, kitKratos.ResponseErr(ctx, pb.ErrorUnauthenticated)
					}

					ctx = context.WithValue(ctx, UserClaims{}, userClaimsm)

					return handler(ctx, req)
				}
			}
			return nil, pb.ErrorUnauthenticated("CheckToken failed.")
		}
	}
}

func GetUserToken(ctx context.Context) (*UserClaims, error) {
	jwtToken := ctx.Value(UserClaims{})
	if jwtToken == nil {
		return nil, errors.New("token data empty.")
	}

	val, ok := jwtToken.(*UserClaims)
	if !ok{
		return nil, errors.New("token data error.")
	}
	return val, nil
}

func GetUserId(ctx context.Context) (int64) {
	user, err := GetUserToken(ctx)
	if err != nil{
		return 0
	}

	id, err := strconv.ParseInt(user.JwtClaims.Id, 10, 64)
	if err != nil {
		return 0
	}

	return id
}
