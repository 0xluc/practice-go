package main

type MockStore struct {
}

func (m *MockStore) CreateTask(*Task) (*Task, error) {
	return &Task{}, nil
}
func (m *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}
