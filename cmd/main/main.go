package main

import (
	"log"
	"net/http"

	"findMyDoc/pkg/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"findMyDoc/pkg/aisearch"
	"findMyDoc/pkg/email"

	"findMyDoc/middlewares"

	doctorsControllers "findMyDoc/doctors/controllers"

	doctorsUsecases "findMyDoc/doctors/usecases"

	doctorsRepositories "findMyDoc/doctors/repositories"

	patientsRepositories "findMyDoc/patients/repositories"

	appointmentsControllers "findMyDoc/appoinments/controllers"

	appointmentsUsecases "findMyDoc/appoinments/usecases"

	appointmentsRepositories "findMyDoc/appoinments/repositories"

	usersControllers "findMyDoc/users/controllers"

	usersUsecases "findMyDoc/users/usecases"

	usersRepositories "findMyDoc/users/repositories"
)

func main() {
	// Setup database connection
	connStr := "user=asif password=mySecret dbname=find_my_doc sslmode=disable"
	database, err := db.NewPostgresDB(connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Setup doctor-related components
	doctorRepo := doctorsRepositories.NewDoctorRepository(database)
	doctorUsecase := doctorsUsecases.NewDoctorUsecase(doctorRepo)
	doctorController := doctorsControllers.NewDoctorController(doctorUsecase)

	patientRepo := patientsRepositories.NewPatientRepository(database)

	// Setup appointment-related components
	appointmentRepo := appointmentsRepositories.NewAppointmentRepository(database)
	appointmentUsecase := appointmentsUsecases.NewAppointmentUsecase(appointmentRepo)
	appointmentController := appointmentsControllers.NewAppointmentController(appointmentUsecase)

	// Load SMTP credentials from environment variables
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpEmail := "heyno143@gmail.com"
	smtpPass := "vmig eyey hqyc mpje"
	appHost := "http://localhost:3001"

	// User authentication setup
	emailService := email.NewSMTPService(smtpHost, smtpPort, smtpEmail, smtpPass)
	userRepo := usersRepositories.NewUserRepository(database)
	userUsecase := usersUsecases.NewUserUsecase(userRepo, emailService, appHost)
	userController := usersControllers.NewUserController(userUsecase, doctorRepo, patientRepo, userRepo)

	deepSeek := aisearch.NewDeepSeekService("#")
	diagnosisHandler := aisearch.NewDiagnosisHandler(deepSeek)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/register", userController.RegisterHandler)
	r.Post("/login", userController.LoginHandler)
	r.Get("/verify", userController.VerifyEmail)
	r.Post("/suggest-doctor", diagnosisHandler.SuggestHandler)

	r.Route("/api", func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Get("/doctors", doctorController.SearchDoctors)                                             // search doctor
		r.Post("/appointments", appointmentController.BookAppointmentHandler)                         // book an appoinment
		r.Get("/doctors/appointments/pending", appointmentController.GetPendingAppointmentsHandler)   // pending appoinment list
		r.Patch("/appointments/{id}/accept", appointmentController.AcceptAppointmentHandler)          // accept appoinment
		r.Get("/doctors/appointments/accepted", appointmentController.GetAcceptedAppointmentsHandler) // accepted appoinment list
		r.Get("/profile", userController.GetProfile)                                                  // get profile
	})

	// Start server
	log.Println("Server running on port 3000")
	http.ListenAndServe(":3000", r)
}
