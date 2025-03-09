package server

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	// Инициализация конфигурации
	config.DaemonConfig = &config.Config{}
	config.DaemonConfig.Server.Host = "localhost"
	config.DaemonConfig.Server.Port = "0" // Случайный свободный порт

	t.Run("create server", func(t *testing.T) {
		ctx := context.Background()
		srv := NewStatsDaemonServer(ctx)
		require.NotNil(t, srv)
		require.NotNil(t, srv.grpcServer)
		require.NotNil(t, srv.metrics)
	})

	t.Run("server starts and listens", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем сервер
		srv := NewStatsDaemonServer(ctx)
		require.NotNil(t, srv)

		// Получаем свободный порт
		lis, err := net.Listen("tcp", net.JoinHostPort(config.DaemonConfig.Server.Host, "0"))
		require.NoError(t, err)
		port := lis.Addr().(*net.TCPAddr).Port
		lis.Close() // Важно закрыть листенер перед использованием порта

		// Устанавливаем порт в конфигурацию
		config.DaemonConfig.Server.Port = fmt.Sprintf("%d", port)

		// Запускаем сервер в горутине
		errCh := make(chan error, 1)
		go func() {
			errCh <- srv.Start()
		}()

		// Даем серверу время на запуск
		time.Sleep(100 * time.Millisecond)

		// Проверяем, что сервер запустился без ошибок
		select {
		case err := <-errCh:
			require.NoError(t, err)
		default:
			// Проверяем, что порт действительно прослушивается
			conn, err := net.Dial("tcp", net.JoinHostPort(config.DaemonConfig.Server.Host, config.DaemonConfig.Server.Port))
			if err == nil {
				conn.Close()
			}
			require.NoError(t, err, "Сервер должен прослушивать порт")
		}

		// Останавливаем сервер
		srv.Stop()
	})
}
