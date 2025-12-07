package grpcserver

import (
	"context"
	"fmt"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// UnaryLoggingInterceptor возвращает grpc.UnaryServerInterceptor, который логирует
// входящие unary-вызовы: клиентский адрес, метод, latency, ошибку и входящие metadata.
func UnaryLoggingInterceptor(logg *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		// получаем IP клиента (если есть).
		peerInfo, _ := peer.FromContext(ctx)
		clientAddr := ""
		if peerInfo != nil && peerInfo.Addr != nil {
			clientAddr = peerInfo.Addr.String()
		}

		// метаданные (user-agent и пр.).
		md, _ := metadata.FromIncomingContext(ctx)

		// вызываем реальный handler.
		resp, err = handler(ctx, req)
		latency := time.Since(start)

		msg := fmt.Sprintf("%s [%s] %s error=%v latency=%s md=%v",
			clientAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			info.FullMethod,
			err,
			latency.String(),
			md,
		)
		if err != nil {
			logg.Error(msg)
		} else {
			logg.Info(msg)
		}

		return resp, err
	}
}
