package suite // простой тестовый клиент

import (
	"context"
	"net"
	"sso/internal/config"
	"strconv"
	"testing"

	ssov1 "github.com/VldslvKtv/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T // для вызоворв методов testing.T внутри
	Cfg        *config.Config
	AuthClient ssov1.AuthClient // клиент для вхаимодействия с grpc-сервером
}

const grpcHost = "localhost"

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath("../config/local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper() // указание что функция вспомагательня для тестов и нужно лучше показывать ошибки
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(), // grpc-клиент
		grpcAddres(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // использование небезопасного соединения для простоты
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddres(cfd *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfd.GRPC.Port))
}
