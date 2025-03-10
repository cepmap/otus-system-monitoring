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

	config.DaemonConfig = &config.Config{}
	config.DaemonConfig.Server.Host = "localhost"
	config.DaemonConfig.Server.Port = "0"

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

		srv := NewStatsDaemonServer(ctx)
		require.NotNil(t, srv)

		lis, err := net.Listen("tcp", net.JoinHostPort(config.DaemonConfig.Server.Host, "0"))
		require.NoError(t, err)
		port := lis.Addr().(*net.TCPAddr).Port
		lis.Close()

		config.DaemonConfig.Server.Port = fmt.Sprintf("%d", port)

		errCh := make(chan error, 1)
		go func() {
			errCh <- srv.Start()
		}()

		time.Sleep(100 * time.Millisecond)

		select {
		case err := <-errCh:
			require.NoError(t, err)
		default:

			conn, err := net.Dial("tcp", net.JoinHostPort(config.DaemonConfig.Server.Host, config.DaemonConfig.Server.Port))
			if err == nil {
				conn.Close()
			}
			require.NoError(t, err, "Сервер должен прослушивать порт")
		}

		srv.Stop()
	})
}
