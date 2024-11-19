package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.xsfx.dev/glucose_exporter/internal/epoch"
)

const BaseURL = "https://api.libreview.io/llu"

// client is our default HTTP client to use.
//
//nolint:gochecknoglobals
var client = http.Client{
	Timeout: 2 * time.Second,
}

func request(
	ctx context.Context,
	method string,
	url string,
	token, accountID string,
	data []byte,
) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return &http.Response{}, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Product", "llu.android")
	req.Header.Add("Version", "4.12.0")

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	if accountID != "" {
		h := sha256.New()

		h.Write([]byte(accountID))

		req.Header.Add("Account-Id", hex.EncodeToString(h.Sum(nil)))
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return &http.Response{}, fmt.Errorf("doing the request: %w", err)
	}

	return resp, nil
}

type Ticket struct {
	Token   string      `json:"token"`
	Expires epoch.Epoch `json:"expires"`
}

type User struct {
	ID string `json:"id"`
}

type LoginResponse struct {
	Data struct {
		AuthTicket Ticket `json:"authTicket"`
		User       User   `json:"user"`
	} `json:"data"`
}

func Login(ctx context.Context, baseURL, username, password string) (LoginResponse, error) {
	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{username, password}

	d, err := json.Marshal(data)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("marshalling data: %w", err)
	}

	url, err := url.JoinPath(baseURL, "/auth/login")
	if err != nil {
		return LoginResponse{}, fmt.Errorf("joining url: %w", err)
	}

	resp, err := request(ctx, http.MethodPost, url, "", "", d)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("reading body: %w", err)
	}

	var authTicket LoginResponse
	if err := json.Unmarshal(body, &authTicket); err != nil {
		return LoginResponse{}, fmt.Errorf("unmarshal login response: %w", err)
	}

	return authTicket, nil
}

type ConnectionsResponse struct {
	Data []struct {
		PatientID string `json:"patientID"`

		Sensor struct {
			DeviceID     string `json:"deviceId"`
			Serialnumber string `json:"sn"`
		} `json:"sensor"`

		GlucoseMeasurement struct {
			ValueInMgPerDl int `json:"ValueInMgPerDl"`
			TrendArrow     int `json:"TrendArrow"`
		} `json:"glucoseMeasurement"`
	} `json:"data"`

	Ticket Ticket `json:"ticket"`
}

func Connections(
	ctx context.Context,
	baseURL, token, accountID string,
) (ConnectionsResponse, error) {
	url, err := url.JoinPath(baseURL, "/connections")
	if err != nil {
		return ConnectionsResponse{}, fmt.Errorf("joining url: %w", err)
	}

	resp, err := request(ctx, http.MethodGet, url, token, accountID, []byte{})
	if err != nil {
		return ConnectionsResponse{}, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ConnectionsResponse{}, fmt.Errorf("reading body: %w", err)
	}

	var connResp ConnectionsResponse

	if err := json.Unmarshal(body, &connResp); err != nil {
		return ConnectionsResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return connResp, nil
}
