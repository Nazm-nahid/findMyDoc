package repositories

import (
	"findMyDoc/internal/entities"
	"findMyDoc/internal/utils"
	"gorm.io/gorm"
)

type DoctorRepository interface {
	SearchDoctors(speciality string, latitude, longitude float64) ([]entities.Doctor, error)
}

type doctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{
		db: db,
	}
}

func (r *doctorRepository) SearchDoctors(speciality string, latitude, longitude float64) ([]entities.Doctor, error) {
	var doctors []entities.Doctor
	if err := r.db.Where("speciality = ?", speciality).Find(&doctors).Error; err != nil {
		return nil, err
	}

	// Filter doctors by location (within 1km)
	var filteredDoctors []entities.Doctor
	for _, doc := range doctors {
		distance := utils.CalculateDistance(latitude, longitude, doc.Latitude, doc.Longitude)
		if distance <= 1.0 {
			filteredDoctors = append(filteredDoctors, doc)
		}
	}

	return filteredDoctors, nil
}
