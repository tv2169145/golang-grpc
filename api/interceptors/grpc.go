package interceptors

import (
	"context"
	"github.com/go-xorm/xorm"
	"github.com/tv2169145/golang-grpc/repos"
	"github.com/tv2169145/golang-grpc/utils"
	"google.golang.org/grpc"
)

func globalRepoInjector(db *xorm.Engine) grpc.UnaryServerInterceptor {
	// 前面的 grpc.UnaryServerInterceptor() 作為型別轉換, 將後面定義的function 轉為 grpc.UnaryServerInterceptor類型
	return grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		globalRepo := repos.GlobalRepo(db)
		newContext := utils.SetGlobalRepoOnContext(ctx, globalRepo)

		// before the request

		// make the actual request
		res, err := handler(newContext, req)

		// after the request - the database has presumably already been called!!

		return res, err
	})
}
