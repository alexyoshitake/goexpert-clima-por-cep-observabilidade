package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func Run(ctx context.Context, server *http.Server, shutdownTelemetry func(context.Context) error) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	var serverErr error
	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			serverErr = err
		}
	case <-ctx.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		serverErr = server.Shutdown(shutdownContext)
	}

	if shutdownTelemetry == nil {
		return serverErr
	}

	telemetryContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return errors.Join(serverErr, shutdownTelemetry(telemetryContext))
}
