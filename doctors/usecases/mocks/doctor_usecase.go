package mocks

import (
	"findMyDoc/internal/entities"
	"github.com/stretchr/testify/mock"
)

type DoctorUsecase struct {
	mock.Mock
}

func (m *DoctorUsecase) SearchDoctors(speciality string, lat, long float64) ([]entities.Doctor, error) {
	args := m.Called(speciality, lat, long)
	return args.Get(0).([]entities.Doctor), args.Error(1)
}
