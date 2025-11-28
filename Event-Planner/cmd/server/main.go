package main

import (
	"encoding/json"
	"log"
	"net/http"

	"event-planner/internal/auth"
	"event-planner/internal/db"
	"event-planner/internal/event"
	"event-planner/internal/invitation"
	"event-planner/internal/search"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Connect to PostgreSQL
	pool, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	//User Management
	authService := auth.NewService(pool)
	authHandler := auth.NewHandler(authService)

	//Event Management
	eventRepo := event.NewRepository(pool)
	eventService := event.NewService(eventRepo)
	eventHandler := event.NewHandler(eventService)

	//Response Management / Invitations
	invRepo := invitation.NewRepository(pool)
	invService := invitation.NewService(invRepo, eventRepo)
	invHandler := invitation.NewHandler(invService)

	// search & Filtering
	searchRepo := search.NewRepository(pool)
	searchService := search.NewService(searchRepo)
	searchHandler := search.NewHandler(searchService)

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:4200",
			"http://127.0.0.1:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Server is running"))
	})

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Events routes
	r.Route("/events", func(r chi.Router) {
		// Public endpoints (no auth required)

		// GET all events
		r.Get("/", eventHandler.GetAllEvents)

		// Advanced search 
		r.With(authHandler.AuthMiddleware).Get("/search", searchHandler.SearchEvents)

		// GET events by organizer
		r.Get("/organizer/{id}", eventHandler.GetEventsByOrganizer)

		// GET single event by ID (public)
		r.Get("/{id}", eventHandler.GetEventByID)

		// GET event attendees
		r.Get("/{id}/attendees", eventHandler.GetEventAttendees)

		// GET invitations for an event
		r.With(authHandler.AuthMiddleware).Get("/{id}/invitations", invHandler.GetEventInvitations)

		// POST create new event
		r.With(authHandler.AuthMiddleware).Post("/", eventHandler.CreateEvent)

		// PUT update event
		r.With(authHandler.AuthMiddleware).Put("/{id}", eventHandler.UpdateEvent)

		// DELETE event
		r.With(authHandler.AuthMiddleware).Delete("/{id}", eventHandler.DeleteEvent)

		// POST join event
		r.With(authHandler.AuthMiddleware).Post("/{id}/join", eventHandler.JoinEvent)

		// POST invite user to event
		r.With(authHandler.AuthMiddleware).Post("/{id}/invite", eventHandler.InviteUserToEvent)

		// PUT update attendance status
		r.With(authHandler.AuthMiddleware).Put("/{id}/attendance", eventHandler.UpdateAttendanceStatus)

		r.Route("/my", func(r chi.Router) {
			r.Use(authHandler.AuthMiddleware)

			// GET events I'm attending
			r.Get("/attending", eventHandler.GetMyAttendingEvents)

			// GET events I'm organizing
			r.Get("/organized", eventHandler.GetMyOrganizedEvents)
		})
	})

	// Invitation routes
	r.Route("/invitations", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)

		// Send invitation
		r.Post("/", invHandler.SendInvitation)

		// Get my invitations
		r.Get("/my", invHandler.GetMyInvitations)

		// Respond to invitation
		r.Put("/{id}/respond", invHandler.RespondToInvitation)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			userID, ok := auth.GetUserID(r.Context())
			if !ok {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			resp := map[string]interface{}{
				"message": "This is a protected route",
				"user_id": userID,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
	})

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

