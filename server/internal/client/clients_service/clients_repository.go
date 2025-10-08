package clients_service

import "goodhumored/wmi-metrics-server/internal/client"

type ClientsRepository interface {
	AddClient(c *client.Client) error
	GetClient(id string) (*client.Client, bool)
	UpdateClient(c *client.Client) error
	RemoveClient(id string)
}
