// Ejemplo del Capítulo 37: logging estructurado con log/slog — niveles,
// atributos y loggers "hijo" con contexto fijo.
package main

import (
	"log/slog"
	"os"
	"time"
)

func procesarPedido(logger *slog.Logger, pedidoID int, monto float64) error {
	inicio := time.Now()
	logger.Info("procesando pedido", "pedido_id", pedidoID, "monto", monto)

	if monto <= 0 {
		logger.Warn("monto inválido, se rechaza el pedido", "pedido_id", pedidoID, "monto", monto)
		return nil
	}

	logger.Debug("pedido procesado exitosamente",
		"pedido_id", pedidoID,
		"duracion_ms", time.Since(inicio).Milliseconds(),
	)
	return nil
}

func main() {
	// Logger en JSON, típico para producción (fácil de parsear por
	// sistemas de observabilidad).
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// With() crea un logger "hijo" que siempre incluye estos atributos:
	// útil para adjuntar un request_id o un nombre de servicio.
	logerServicio := logger.With("servicio", "pedidos", "version", "1.0")

	logerServicio.Info("servicio iniciado")
	procesarPedido(logerServicio, 1001, 49.90)
	procesarPedido(logerServicio, 1002, -10)

	// slog.Group agrupa atributos relacionados.
	logerServicio.Info("resumen de sesión",
		slog.Group("metricas",
			"pedidos_procesados", 2,
			"errores", 0,
		),
	)
}
