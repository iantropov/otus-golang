package grpc

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		timeStart := time.Now()
		res, err := handler(ctx, req)
		elapsed := time.Since(timeStart)

		remoteAddr := "unknown"

		p, ok := peer.FromContext(ctx)
		if ok {
			remoteAddr = strings.SplitN(p.Addr.String(), ":", 2)[0]
		}

		intCode := int(status.Code(err))

		// 66.249.65.3 [25/Feb/2020:19:11:24 +0600] GET /hello?q=1 HTTP/1.1 200 30 "Mozilla/5.0"
		logger.Infof(
			"%s [%s %s] %s GRPC %d %dms\n",
			remoteAddr,
			timeStart.Format("01/Jan/2006"),
			timeStart.Format("15:04:05 -0700"),
			info.FullMethod,
			intCode,
			elapsed.Milliseconds(),
		)

		return res, err
	}
}
