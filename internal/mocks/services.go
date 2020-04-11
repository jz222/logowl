package mocks

import "github.com/jz222/loggy/internal/models"

type LoggingService struct{}

func (l *LoggingService) SaveError(e models.Error) {}
