package socket

func (m *Manager) FindClientByUserID(id int32) *Client {
	if client, ok := m.clients[id]; ok {
		return client
	}
	return nil
}
