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
	donationHandler := handlers.NewDonationHandler(donationService, inventoryService, userService, cfg)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService, userService, cfg)
	organisationHandler := handlers.NewOrganisationHandler(organisationService, userService, cfg)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, userService, cfg)
	uploadHandler := handlers.NewUploadHandler(cfg)

	// ROLE BASED MIDDLEWARE
	requireAdmin := middleware.RequireRole("admin")
	requireStaffOrAdmin := middleware.RequireRole("charity_staff", "admin")

	// FOR PROTECTED ROUTES
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg))

	// AUTH ROUTES
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")   // REGISTER - UNPROTECTED
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")         // LOGIN - UNPROTECTED
	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")        // LOGOUT - PROTECTED
	protected.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST") // REFRESH AUTH TOKEN - PROTECTED

	// USER ROUTES
	protected.HandleFunc("/users/profile", userHandler.GetProfile).Methods("GET")                         // GET USER PROFILE
	protected.HandleFunc("/users/profile", userHandler.UpdateProfile).Methods("PUT")                      // UPDATE USER PROFILE
	protected.Handle("/users", requireAdmin(http.HandlerFunc(userHandler.List))).Methods("GET")           // LIST USERS - REQUIRES ADMIN
	protected.Handle("/users/{id}", requireAdmin(http.HandlerFunc(userHandler.GetByID))).Methods("GET")   // GET USER BY ID - REQUIRES ADMIN
	protected.Handle("/users/{id}", requireAdmin(http.HandlerFunc(userHandler.Delete))).Methods("DELETE") // DELETE USER BY ID - REQUIRES ADMIN

	// DONATION ROUTES
	protected.HandleFunc("/donations", donationHandler.Create).Methods("POST")                                                     // CREATE DONATION
	protected.HandleFunc("/donations/my", donationHandler.GetMyDonations).Methods("GET")                                           // GET MY DONATIONS
	protected.Handle("/donations", requireStaffOrAdmin(http.HandlerFunc(donationHandler.List))).Methods("GET")                     // LIST DONATIONS - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/donations/{id}", requireStaffOrAdmin(http.HandlerFunc(donationHandler.GetByID))).Methods("GET")             // GET DONATION BY ID - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/donations/{id}/status", requireStaffOrAdmin(http.HandlerFunc(donationHandler.UpdateStatus))).Methods("PUT") // UPDATE DONATION STATUS BY ID - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/donations/{id}/approve", requireStaffOrAdmin(http.HandlerFunc(donationHandler.Approve))).Methods("POST")    // APPROVE DONATION BY ID - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/donations/{id}/reject", requireStaffOrAdmin(http.HandlerFunc(donationHandler.Reject))).Methods("POST")      // REJECT DONATION BY ID - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/donations/{id}", requireAdmin(http.HandlerFunc(donationHandler.Delete))).Methods("DELETE")                  // DELETE DONATION BY ID - REQUIRES ADMIN

	// INVENTORY ROUTES
	protected.HandleFunc("/inventory", inventoryHandler.List).Methods("GET")                                                           // LIST INVENTORY ITEMS
	protected.Handle("/inventory/stats", requireStaffOrAdmin(http.HandlerFunc(inventoryHandler.GetStats))).Methods("GET")              // INVENTORY STATS - REQUIRES CHARITY STAFF OR ADMIN
	protected.HandleFunc("/inventory/{id}", inventoryHandler.GetByID).Methods("GET")                                                   // GET INVENTORY ITEM BY ID
	protected.Handle("/inventory", requireStaffOrAdmin(http.HandlerFunc(inventoryHandler.Create))).Methods("POST")                     // MANUAL CREATE INVENTORY ITEM - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/inventory/{id}", requireStaffOrAdmin(http.HandlerFunc(inventoryHandler.Update))).Methods("PUT")                 // UPDATE INVENTORY ITEM BY ID - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/inventory/{id}/allocate", requireStaffOrAdmin(http.HandlerFunc(inventoryHandler.Allocate))).Methods("POST")     // ALLOCATE INVENTORY BY ID - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/inventory/{id}/distribute", requireStaffOrAdmin(http.HandlerFunc(inventoryHandler.Distribute))).Methods("POST") // DISTRIBUTE INVENTORY BY ID - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/inventory/{id}/deallocate", requireStaffOrAdmin(http.HandlerFunc(inventoryHandler.Deallocate))).Methods("POST") // DEALLOCATE INVENTORY BY ID - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/inventory/{id}", requireAdmin(http.HandlerFunc(inventoryHandler.Delete))).Methods("DELETE")                     // DELETE INVENTORY ITEM - REQUIRES ADMIN

	// ORGANISATION ROUTES
	protected.HandleFunc("/organisations", organisationHandler.List).Methods("GET")                                                   // LIST ORGANISATIONS
	protected.HandleFunc("/organisations/{id}", organisationHandler.GetByID).Methods("GET")                                           // GET ORGANISATION BY ID
	protected.HandleFunc("/organisations/email/{email}", organisationHandler.GetByEmail).Methods("GET")                               // GET ORGANISATION BY EMAIL
	protected.Handle("/organisations/{id}/stats", requireStaffOrAdmin(http.HandlerFunc(organisationHandler.GetStats))).Methods("GET") // GET ORGANISATION STATS BY ID - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/organisations", requireAdmin(http.HandlerFunc(organisationHandler.Create))).Methods("POST")                    // CREATE ORGANISATION - REQUIRES ADMIN
	protected.Handle("/organisations/{id}", requireAdmin(http.HandlerFunc(organisationHandler.Update))).Methods("PUT")                // UPDATE ORGANISATION BY ID - REQUIRES ADMIN
	protected.Handle("/organisations/{id}", requireAdmin(http.HandlerFunc(organisationHandler.Delete))).Methods("DELETE")             // DELETE ORGANISATION BY ID - REQUIRES ADMIN

	// ANALYTICS ROUTES
	protected.HandleFunc("/analytics/donor-impact", analyticsHandler.GetDonorImpact).Methods("GET")                                                // GET DONOR SUSTAINABILITY IMPACT
	protected.Handle("/analytics/trends", requireStaffOrAdmin(http.HandlerFunc(analyticsHandler.GetDonationTrends))).Methods("GET")                // GET DONATION TRENDS - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/analytics/categories", requireStaffOrAdmin(http.HandlerFunc(analyticsHandler.GetCategoryBreakdown))).Methods("GET")         // GET CATEGORY BREAKDOWN - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/analytics/sustainability", requireStaffOrAdmin(http.HandlerFunc(analyticsHandler.GetSustainabilityMetrics))).Methods("GET") // GET SUSTAINABILITY METRICS - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/analytics/org-performance", requireStaffOrAdmin(http.HandlerFunc(analyticsHandler.GetOrgPerformance))).Methods("GET")       // GET ORGANISATION PERFORMANCE - REQUIRES CHARITY STAFF OR ADMIN
	protected.Handle("/analytics/system-overview", requireAdmin(http.HandlerFunc(analyticsHandler.GetSystemOverview))).Methods("GET")              // GET SYSTEM OVERVIEW - REQUIRES ADMIN

	// UPLOAD ROUTES
	protected.HandleFunc("/uploads/images", uploadHandler.UploadImages).Methods("POST")                                                       // MULTIPART EP TO UPLOAD IMAGES - PROTECTED
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.FileUpload.UploadDir)))).Methods("GET") // FILESERVER TO VIEW UPLOADED FILES - UNPROTECTED

	// HEALTH CHECK
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}
