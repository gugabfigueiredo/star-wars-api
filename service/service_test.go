package service

import (
	"errors"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/test"
	"github.com/gugabfigueiredo/swapi"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAPIHandler_SchedulePlanetRefsUpdate(t *testing.T) {

	tests := []struct{
		name               	string
		stub               	*test.Stub
		expectedUpdates		int
	}{
		{
			name: "update planet refs",
			stub: &test.Stub{},
			expectedUpdates: 3,
		},
		{
			name: "fail to get updated refs",
			stub: &test.Stub{Error: errors.New("failed to get updated refs")},
		},
	}

	logger := log.New(&log.Config{
		Context:               "sw-api-test",
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			swapiStub := &SwapiStub{Error: tt.stub.Error}

			s := &APIService{
				IRepo: tt.stub,
				SwapiClient: swapiStub,
				Logger:   logger,
			}

			schedule := s.SchedulePlanetUpdate(time.Second)

			time.Sleep(3 * time.Second)
			schedule <- false
			close(schedule)
			assert.Equal(t, tt.expectedUpdates, swapiStub.SwapiUpdates)
		})
	}
}

type SwapiStub struct {
	Error	error
	SwapiUpdates int
}


func (s *SwapiStub) AllPlanets() ([]swapi.Planet, error) {
	if s.Error == nil {
		s.SwapiUpdates += 1
	}
	return []swapi.Planet{}, s.Error
}