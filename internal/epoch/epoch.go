package epoch

import (
	"fmt"
	"strconv"
	"time"
)

type Epoch time.Time

func (e Epoch) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(e).Unix(), 10)), nil
}

func (e *Epoch) UnmarshalJSON(b []byte) error {
	q, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return fmt.Errorf("parse int: %w", err)
	}

	*(*time.Time)(e) = time.Unix(q, 0)

	return nil
}

func (e Epoch) String() string { return time.Time(e).String() }
