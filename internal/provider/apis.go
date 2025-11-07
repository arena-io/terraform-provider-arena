// Copyright (c) ArenaML Labs Pvt Ltd.

package provider

import (
	"net/http"

	"github.com/arena-ml/terraform-provider-arenaml/generator/client"
)

func NewClient(addr string) (*client.ClientWithResponses, error) {
	if addr == "" {
		addr = "http://localhost:18080"
	}

	hc := http.Client{}
	return client.NewClientWithResponses(addr, client.WithHTTPClient(&hc))
}
