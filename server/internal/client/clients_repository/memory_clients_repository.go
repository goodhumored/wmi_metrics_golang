package clients_repository

import (
	"fmt"
	"goodhumored/wmi-metrics-server/internal/client"
)

type MemoryClientsRepository struct {
	clients map[string]*client.Client
}

func New() *MemoryClientsRepository {
	return &MemoryClientsRepository{
		clients: make(map[string]*client.Client),
	}
}

func (r *MemoryClientsRepository) AddClient(c *client.Client) error {
	r.clients[c.ID] = c
	return nil
}

func (r *MemoryClientsRepository) UpdateClient(c *client.Client) error {
	if _, exists := r.clients[c.ID]; !exists {
		return fmt.Errorf("client with ID %s does not exist", c.ID)
	}
	r.clients[c.ID] = c
	return nil
}

func (r *MemoryClientsRepository) GetClient(id string) (*client.Client, bool) {
	c, exists := r.clients[id]
	return c, exists
}

func (r *MemoryClientsRepository) RemoveClient(id string) {
	delete(r.clients, id)
}

func (r *MemoryClientsRepository) GetAllClients() []*client.Client {
	clients := make([]*client.Client, 0, len(r.clients))
	for _, c := range r.clients {
		clients = append(clients, c)
	}
	return clients
}
