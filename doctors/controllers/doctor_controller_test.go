package controllers

import (
	"encoding/json"
	"findMyDoc/internal/entities"
	"findMyDoc/doctors/usecases/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchDoctors(t *testing.T) {
	mockUsecase := new(mocks.DoctorUsecase)
	controller := NewDoctorController(mockUsecase)

	expectedDoctors := []entities.Doctor{
		{ID: 1, Name: "Dr. Smith", Speciality: "Cardiologist", Latitude: 40.7128, Longitude: -74.0060},
		{ID: 2, Name: "Dr. Jane", Speciality: "Cardiologist", Latitude: 40.7130, Longitude: -74.0070},
	}

	// Set up mock expectation
	mockUsecase.On("SearchDoctors", "Cardiologist", 40.7128, -74.0060).Return(expectedDoctors, nil)

	// Prepare the request
	req := httptest.NewRequest("GET", "/doctors/search?speciality=Cardiologist&latitude=40.7128&longitude=-74.0060", nil)
	rec := httptest.NewRecorder()

	// Call the handler
	controller.SearchDoctors(rec, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rec.Code)

	var actualDoctors []entities.Doctor
	err := json.Unmarshal(rec.Body.Bytes(), &actualDoctors)
	assert.NoError(t, err)
	assert.Equal(t, expectedDoctors, actualDoctors)

	// Assert the mock was called
	mockUsecase.AssertExpectations(t)
}
