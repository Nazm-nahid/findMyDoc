package controllers

import (
	"encoding/json"
	"findMyDoc/internal/utils"
	"findMyDoc/users/requests"
	"findMyDoc/users/usecases"
	"net/http"

	doctorsRepositories "findMyDoc/doctors/repositories"

	patientsRepositories "findMyDoc/patients/repositories"

	usersRepositories "findMyDoc/users/repositories"
)

type UserController struct {
	usecase              usecases.UserUsecase
	doctorsRepositories  doctorsRepositories.DoctorRepository
	patientsRepositories patientsRepositories.PatientRepository
	usersRepositories    usersRepositories.UserRepository
}

func NewUserController(uc usecases.UserUsecase, dc doctorsRepositories.DoctorRepository, pc patientsRepositories.PatientRepository, ur usersRepositories.UserRepository) *UserController {
	return &UserController{usecase: uc, doctorsRepositories: dc, patientsRepositories: pc,usersRepositories: ur}
}

func (c *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := c.usecase.LoginUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (c *UserController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := c.usecase.RegisterUser(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully!! Please verify your email!!",
	})
}

func (c *UserController) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Verification code is required", http.StatusBadRequest)
		return
	}

	// Step 1: Find user by verification code
	user, err := c.usersRepositories.GetByVerificationCode(code)
	if err != nil || user == nil || user.ID == 0 {
		http.Error(w, "Invalid or expired verification code", http.StatusBadRequest)
		return
	}

	// Step 2: Update user verification status
	user.IsVerified = true
	user.VerificationCode = ""

	if err := c.usersRepositories.Update(user); err != nil {
		http.Error(w, "Failed to verify user", http.StatusInternalServerError)
		return
	}

	// Step 3: Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email verified successfully",
	})
}

func (c *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	userID := utils.ExtractUserIDFromToken(authHeader)
	role := utils.ExtractRoleFromToken(authHeader)

	if role == "doctor" {
		doctors, err := c.doctorsRepositories.GetDoctorById(userID)
		if err != nil {
			http.Error(w, "Error retrieving doctors", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(doctors)
	} else {
		patients, err := c.patientsRepositories.GetPatientById(userID)
		if err != nil {
			http.Error(w, "Error retrieving doctors", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(patients)
	}

}
