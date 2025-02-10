package socket

import "github.com/google/uuid"

func (m *Manager) FindClientByUserID(id uuid.UUID) *Client {
	if client, ok := m.clients[id]; ok {
		return client
	}
	return nil
}
