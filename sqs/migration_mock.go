package sqs

import (
	"github.com/stretchr/testify/mock"
)

func NewMigrationMock() *MigrationMock {
	return &MigrationMock{}
}

type MigrationMock struct {
	mock.Mock
}

func (s *MigrationMock) Run(queue string) error {
	args := s.Called(queue)
	return args.Error(0)
}
