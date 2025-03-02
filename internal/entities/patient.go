package entities

type Patient struct {
	ID         int     `json:"id" gorm:"primaryKey"`
	UserID     int      `json:"user_id"`
	Name       string  `json:"name"`
	Ratings    float64 `json:"ratings"`
}
