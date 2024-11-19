package metrics

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	vm "github.com/VictoriaMetrics/metrics"
	"go.xsfx.dev/glucose_exporter/api"
	"go.xsfx.dev/glucose_exporter/internal/cache"
	"go.xsfx.dev/glucose_exporter/internal/config"
)

var ErrOddDataLength = errors.New("odd data length")

func Handler(w http.ResponseWriter, r *http.Request) {
	if err := glucose(r.Context(), w); err != nil {
		slog.Error("getting glucose", "err", err)
		vm.GetOrCreateCounter("errors_total").Inc()
	}

	vm.WritePrometheus(w, false)
}

func glucose(ctx context.Context, w io.Writer) error {
	c, err := cache.Load()
	if err != nil {
		return fmt.Errorf("loading cache: %w", err)
	}

	// No token exists or token is expired
	if c.JWT == "" || time.Now().After(time.Time(c.Expires)) {
		slog.Debug("needs a fresh token")

		if err := token(ctx); err != nil {
			return fmt.Errorf("refreshing token: %w", err)
		}
	}

	// Loading cache again with the new token inside.
	c, err = cache.Load()
	if err != nil {
		return fmt.Errorf("loading cache: %w", err)
	}

	// Getting connections, which includes the glucose data.
	resp, err := api.Connections(ctx, api.BaseURL, c.JWT, c.AccountID)
	if err != nil {
		return fmt.Errorf("connections: %w", err)
	}

	// Storing new token to cache.
	if err := cache.Save(cache.Cache{
		JWT:     resp.Ticket.Token,
		Expires: resp.Ticket.Expires,
	}); err != nil {
		return fmt.Errorf("saving cache: %w", err)
	}

	if len(resp.Data) != 1 {
		return ErrOddDataLength
	}

	labels := fmt.Sprintf(
		`patient_id="%s",sensor_serialnumber="%s"`,
		resp.Data[0].PatientID,
		resp.Data[0].Sensor.Serialnumber,
	)

	// Writing metrics.
	vm.WriteGaugeUint64(
		w,
		fmt.Sprintf("value_in_mg_per_dl{%s}", labels),
		uint64(resp.Data[0].GlucoseMeasurement.ValueInMgPerDl),
	)

	vm.WriteGaugeUint64(
		w,
		fmt.Sprintf("trend_arrow{%s}", labels),
		uint64(resp.Data[0].GlucoseMeasurement.TrendArrow),
	)

	return nil
}

func token(ctx context.Context) error {
	pass, err := config.Cfg.GetPassword()
	if err != nil {
		return fmt.Errorf("getting password from config: %w", err)
	}

	token, err := api.Login(
		ctx,
		api.BaseURL,
		config.Cfg.Email,
		pass,
	)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	slog.Debug("got token", "token", token)

	if err := cache.Save(
		cache.Cache{JWT: token.Data.AuthTicket.Token, Expires: token.Data.AuthTicket.Expires, AccountID: token.Data.User.ID},
	); err != nil {
		return fmt.Errorf("saving cache: %w", err)
	}

	return nil
}
