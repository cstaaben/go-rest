package client

import "net/http"

type Client struct {
	Client *http.Client
}
