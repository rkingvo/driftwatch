package inspector

import (
	"context"
	"fmt"
)

// MockClient is an in-memory inspector used for testing.
type MockClient struct {
	containers map[string]*ContainerInfo
}

// NewMock creates a MockClient pre-populated with the provided containers.
func NewMock(containers []*ContainerInfo) *MockClient {
	m := &MockClient{containers: make(map[string]*ContainerInfo, len(containers))}
	for _, c := range containers {
		m.containers[c.Name] = c
	}
	return m
}

// Inspect returns the ContainerInfo for the given name, or an error if not found.
func (m *MockClient) Inspect(_ context.Context, nameOrID string) (*ContainerInfo, error) {
	if info, ok := m.containers[nameOrID]; ok {
		return info, nil
	}
	return nil, fmt.Errorf("mock inspector: container %q not found", nameOrID)
}

// InspectAll returns all containers registered in the mock.
func (m *MockClient) InspectAll(_ context.Context) ([]*ContainerInfo, error) {
	result := make([]*ContainerInfo, 0, len(m.containers))
	for _, c := range m.containers {
		result = append(result, c)
	}
	return result, nil
}
