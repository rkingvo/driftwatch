package inspector

import (
	"context"
	"testing"
)

func sampleContainers() []*ContainerInfo {
	return []*ContainerInfo{
		{
			Name:   "web",
			Image:  "nginx:1.25",
			Env:    map[string]string{"PORT": "8080"},
			Labels: map[string]string{"app": "web"},
			Status: "running",
		},
		{
			Name:   "db",
			Image:  "postgres:15",
			Env:    map[string]string{"POSTGRES_DB": "app"},
			Labels: map[string]string{"app": "db"},
			Status: "running",
		},
	}
}

func TestMockClient_Inspect_Found(t *testing.T) {
	mock := NewMock(sampleContainers())
	info, err := mock.Inspect(context.Background(), "web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Image != "nginx:1.25" {
		t.Errorf("expected image nginx:1.25, got %q", info.Image)
	}
	if info.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", info.Env["PORT"])
	}
}

func TestMockClient_Inspect_NotFound(t *testing.T) {
	mock := NewMock(sampleContainers())
	_, err := mock.Inspect(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing container, got nil")
	}
}

func TestMockClient_InspectAll(t *testing.T) {
	mock := NewMock(sampleContainers())
	all, err := mock.InspectAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 containers, got %d", len(all))
	}
}

func TestContainerInfoFromRaw_StripLeadingSlash(t *testing.T) {
	info := &ContainerInfo{Name: "/mycontainer"}
	if len(info.Name) > 0 && info.Name[0] == '/' {
		info.Name = info.Name[1:]
	}
	if info.Name != "mycontainer" {
		t.Errorf("expected name without slash, got %q", info.Name)
	}
}

func TestMockClient_EmptyRegistry(t *testing.T) {
	mock := NewMock(nil)
	all, err := mock.InspectAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("expected 0 containers, got %d", len(all))
	}
}
