package interceptors

import (
	"context"
	"errors"
	"github.com/go-xorm/xorm"
	"github.com/tv2169145/golang-grpc/repos"
	"github.com/tv2169145/golang-grpc/utils"
	"google.golang.org/grpc"
	"reflect"
)

func GlobalRepoInjector(db *xorm.Engine) grpc.UnaryServerInterceptor {
	// 前面的 grpc.UnaryServerInterceptor() 作為型別轉換, 將後面定義的function 轉為 grpc.UnaryServerInterceptor類型
	return grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		globalRepo := repos.GlobalRepo(db)
		newContext := utils.SetGlobalRepoOnContext(ctx, globalRepo)

		// before the request
		v := reflect.Indirect(reflect.ValueOf(req))
		vField := reflect.Indirect(v.FieldByName("JWT"))

		// if there is no JWT field on the request
		if !vField.IsValid() {
			return handler(newContext, req)
		}
		jwtToken := vField.String()

		if len(jwtToken) == 0 {
			return nil, errors.New("unauthorized")
		}

		user, err := globalRepo.Auth().GetDataFromToken(jwtToken)
		if err != nil {
			return nil, errors.New("unauthorized")
		}
		newContext = utils.SetUserOnContext(newContext, user)

		// make the actual request
		res, err := handler(newContext, req)

		// after the request - the database has presumably already been called!!

		return res, err
	})
}
