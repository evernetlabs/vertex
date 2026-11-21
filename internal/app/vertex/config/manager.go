package config

type Manager struct {
}

func (m *Manager) GetVertexEndpoint() string {
	return ""
}

func (m *Manager) GetJwtSigningKey() string {
	return ""
}
