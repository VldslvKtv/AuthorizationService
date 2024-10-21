package main // точка входа
import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

// TODO: Иициализировать обьект конфига

// TODO: Иициализировать логгер (slog)

// TODO: Иициализировать само приложение (app)

// TODO: запустить gRPC-сервер

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("start app",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", cfg.GRPC.Port),
	)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	// Gracefull shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT) // ждем сигнала от ОС - и  висим на строке <-stop
	// пока отдельная горутина запущена с сервером
	check := <-stop

	log.Info("server stopped", slog.String("signal", check.String()))

	application.GRPCSrv.Stop()

	// fmt.Println(cfg)
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger { // логирование
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New( // показываем куда выводить и какой уровень логирования
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
