package api

import (
	"database/sql"
	"net/http"

	"sustainwear/internal/api/handlers"
	"sustainwear/internal/api/middleware"
	"sustainwear/internal/config"
	"sustainwear/internal/domain/analytics"
	"sustainwear/internal/domain/donation"
	"sustainwear/internal/domain/inventory"
	"sustainwear/internal/domain/organisation"
	"sustainwear/internal/domain/user"

	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Config, db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// APPLY GLOBAL MIDDLEWARE
	router.Use(middleware.LoggerMiddleware)
	router.Use(middleware.CORSMiddleware(cfg))

	// INITIALIZE REPOSITORIES
	userRepo := user.NewRepository(db)
	donationRepo := donation.NewRepository(db)
	inventoryRepo := inventory.NewRepository(db)
	organisationRepo := organisation.NewRepository(db)
	analyticsRepo := analytics.NewRepository(db)

	// INITIALIZE SERVICES
	userService := user.NewService(userRepo)
	donationService := donation.NewService(donationRepo)
	inventoryService := inventory.NewService(inventoryRepo)
	organisationService := organisation.NewService(organisationRepo)
	analyticsService := analytics.NewService(analyticsRepo)

	// INITIALIZE HANDLERS
	authHandler := handlers.NewAuthHandler(userService, cfg)
	userHandler := handlers.NewUserHandler(userService, cfg)
	donationHandler := handlers.NewDonationHandler(donationService, cfg)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService, cfg)
	organisationHandler := handlers.NewOrganisationHandler(organisationService, cfg)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, cfg)
	uploadHandler := handlers.NewUploadHandler(cfg)

	// FOR PROTECTED ROUTES
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg))

	// AUTH ROUTES
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")   // UNPROTECTED
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")         // UNPROTECTED
	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")        // PROTECTED
	protected.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST") // PROTECTED

	// USER ROUTES
	protected.HandleFunc("/users/profile", userHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/users/profile", userHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/users", userHandler.List).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandler.GetByID).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandler.Delete).Methods("DELETE")

	// DONATION ROUTES
	protected.HandleFunc("/donations", donationHandler.Create).Methods("POST")
	protected.HandleFunc("/donations", donationHandler.List).Methods("GET")
	protected.HandleFunc("/donations/my", donationHandler.GetMyDonations).Methods("GET")
	protected.HandleFunc("/donations/{id}", donationHandler.GetByID).Methods("GET")
	protected.HandleFunc("/donations/{id}/status", donationHandler.UpdateStatus).Methods("PUT")
	protected.HandleFunc("/donations/{id}/approve", donationHandler.Approve).Methods("POST")
	protected.HandleFunc("/donations/{id}/reject", donationHandler.Reject).Methods("POST")
	protected.HandleFunc("/donations/{id}", donationHandler.Delete).Methods("DELETE")

	// INVENTORY ROUTES
	protected.HandleFunc("/inventory", inventoryHandler.List).Methods("GET")
	protected.HandleFunc("/inventory/{id}", inventoryHandler.GetByID).Methods("GET")
	protected.HandleFunc("/inventory/{id}", inventoryHandler.Update).Methods("PUT")
	protected.HandleFunc("/inventory/{id}/allocate", inventoryHandler.Allocate).Methods("POST")
	protected.HandleFunc("/inventory/{id}/distribute", inventoryHandler.Distribute).Methods("POST")
	protected.HandleFunc("/inventory/{id}/deallocate", inventoryHandler.Deallocate).Methods("POST")
	protected.HandleFunc("/inventory/{id}", inventoryHandler.Delete).Methods("DELETE")
	protected.HandleFunc("/inventory/stats", inventoryHandler.GetStats).Methods("GET")

	// ORGANISATION ROUTES
	protected.HandleFunc("/organisations", organisationHandler.Create).Methods("POST")
	protected.HandleFunc("/organisations", organisationHandler.List).Methods("GET")
	protected.HandleFunc("/organisations/{id}", organisationHandler.GetByID).Methods("GET")
	protected.HandleFunc("/organisations/email", organisationHandler.GetByEmail).Methods("GET")
	protected.HandleFunc("/organisations/{id}", organisationHandler.Update).Methods("PUT")
	protected.HandleFunc("/organisations/{id}", organisationHandler.Delete).Methods("DELETE")
	protected.HandleFunc("/organisations/{id}/stats", organisationHandler.GetStats).Methods("GET")

	// ANALYTICS ROUTES
	protected.HandleFunc("/analytics/trends", analyticsHandler.GetDonationTrends).Methods("GET")
	protected.HandleFunc("/analytics/categories", analyticsHandler.GetCategoryBreakdown).Methods("GET")
	protected.HandleFunc("/analytics/sustainability", analyticsHandler.GetSustainabilityMetrics).Methods("GET")
	protected.HandleFunc("/analytics/donor-impact", analyticsHandler.GetDonorImpact).Methods("GET")
	protected.HandleFunc("/analytics/org-performance", analyticsHandler.GetOrgPerformance).Methods("GET")
	protected.HandleFunc("/analytics/system-overview", analyticsHandler.GetSystemOverview).Methods("GET")

	// UPLOAD ROUTES
	protected.HandleFunc("/uploads/images", uploadHandler.UploadImages).Methods("POST")                                                       // PROTECTED
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.FileUpload.UploadDir)))).Methods("GET") // UNPROTECTED

	// HEALTH CHECK
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}
