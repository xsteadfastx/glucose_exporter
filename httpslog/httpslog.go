package httpslog

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"go.xsfx.dev/glucose_exporter/httpslog/internal/mutil"
)

func Handler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := mutil.WrapWriter(w)
			defer func() {
				// Logging.
				slog.Info(
					"handled request",
					"status",
					lw.Status(),
					"bytes",
					lw.BytesWritten(),
					"duration",
					time.Since(start),
					"method",
					r.Method,
					"url",
					r.URL,
					"user-agent",
					r.UserAgent(),
					"proto",
					r.Proto,
				)

				// Metrics.
				labels := fmt.Sprintf(
					`method="%s",url="%s",status="%d"`,
					r.Method,
					r.URL,
					lw.Status(),
				)

				metrics.GetOrCreateCounter(fmt.Sprintf(`requests_total{%s}`, labels)).Inc()

				metrics.GetOrCreateSummary(
					fmt.Sprintf(`requests_duration_seconds{%s}`, labels),
				).UpdateDuration(start)

				metrics.GetOrCreateHistogram(
					fmt.Sprintf("response_size{%s}", labels),
				).Update(float64(lw.BytesWritten()))
			}()
			next.ServeHTTP(lw, r)
		})
	}
}
