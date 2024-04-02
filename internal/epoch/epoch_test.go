package epoch_test

import (
	"errors"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.xsfx.dev/glucose_exporter/internal/epoch"
)

func TestMarshal(t *testing.T) {
	tables := []struct {
		name     string
		input    epoch.Epoch
		expected []byte
		err      error
	}{
		{
			name:     "00",
			input:    epoch.Epoch(time.Date(2024, 9, 14, 8, 21, 23, 0, time.UTC)),
			expected: []byte("1726302083"),
			err:      nil,
		},
		{
			name:     "01",
			input:    epoch.Epoch(time.Time{}),
			expected: []byte("-62135596800"),
			err:      nil,
		},
	}

	is := is.New(t)

	for _, tt := range tables {
		t.Run(tt.name, func(_ *testing.T) {
			b, err := tt.input.MarshalJSON()
			if tt.err == nil {
				is.Equal(b, tt.expected)
			} else {
				is.True(errors.Is(err, tt.err))
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	tables := []struct {
		name     string
		input    []byte
		expected epoch.Epoch
	}{
		{
			name:     "00",
			input:    []byte("1726302083"),
			expected: epoch.Epoch(time.Date(2024, 9, 14, 8, 21, 23, 0, time.UTC)),
		},
	}

	is := is.New(t)

	for _, tt := range tables {
		t.Run(tt.name, func(_ *testing.T) {
			e := epoch.Epoch{}

			err := e.UnmarshalJSON(tt.input)
			is.NoErr(err)

			is.Equal(time.Time(e).UTC(), time.Time(tt.expected).UTC())
		})
	}
}
