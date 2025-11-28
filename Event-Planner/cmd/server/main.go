package main

import (
	"encoding/json"
	"event-planner/internal/auth"
	"event-planner/internal/db"
	"event-planner/internal/event"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	pool, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Initialize services and handlers
	authService := auth.NewService(pool)
	authHandler := auth.NewHandler(authService)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Public routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	eventRepo := event.NewRepository(pool)
	eventService := event.NewService(eventRepo)
	eventHandler := event.NewHandler(eventService)

	r.Route("/events", func(r chi.Router) {
		// Public endpoints (no auth required)
		// GET all events
		r.Get("/", eventHandler.GetAllEvents)

		// GET events by organizer
		r.Get("/organizer/{id}", eventHandler.GetEventsByOrganizer)

		// GET single event by ID
		r.Get("/{id}", eventHandler.GetEventByID)

		// GET event attendees
		r.Get("/{id}/attendees", eventHandler.GetEventAttendees)

		// Protected endpoints (auth required)
		// POST create new event (requires auth)
		r.With(authHandler.AuthMiddleware).Post("/", eventHandler.CreateEvent)

		// PUT update event (requires auth + ownership)
		r.With(authHandler.AuthMiddleware).Put("/{id}", eventHandler.UpdateEvent)

		// DELETE event (requires auth + ownership)
		r.With(authHandler.AuthMiddleware).Delete("/{id}", eventHandler.DeleteEvent)

		// POST join event (requires auth)
		r.With(authHandler.AuthMiddleware).Post("/{id}/join", eventHandler.JoinEvent)

		// POST invite user to event (requires auth, creator only)
		r.With(authHandler.AuthMiddleware).Post("/{id}/invite", eventHandler.InviteUserToEvent)

		// PUT update attendance status (requires auth)
		r.With(authHandler.AuthMiddleware).Put("/{id}/attendance", eventHandler.UpdateAttendanceStatus)

		// Protected routes for user's own events
		r.Route("/my", func(r chi.Router) {
			r.Use(authHandler.AuthMiddleware)

			// GET events I'm attending (as organizer or attendee)
			r.Get("/attending", eventHandler.GetMyAttendingEvents)

			// GET events I'm organizing
			r.Get("/organized", eventHandler.GetMyOrganizedEvents)
		})
	})

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)

		// Example protected route
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			userID, ok := auth.GetUserID(r.Context())
			if !ok {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			response := map[string]interface{}{
				"message": "This is a protected route",
				"user_id": userID,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
	})
	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
