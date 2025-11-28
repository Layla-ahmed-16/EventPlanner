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
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	// Connect to PostgreSQL
	pool, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// ==== Auth setup (Requirement 1: User Management) ====
	authService := auth.NewService(pool)
	authHandler := auth.NewHandler(authService)

	// ==== Event setup (Requirement 2: Event Management) ====
	eventRepo := event.NewRepository(pool)
	eventService := event.NewService(eventRepo)
	eventHandler := event.NewHandler(eventService)

	// ==== Invitation setup (Requirement 3: Response Management / Invitations) ====
	invRepo := invitation.NewRepository(pool)
	invService := invitation.NewService(invRepo)
	invHandler := invitation.NewHandler(invService)

	// ==== Search setup (Requirement 4: Search & Filtering) ====
	searchRepo := search.NewRepository(pool)
	searchService := search.NewService(searchRepo)
	searchHandler := search.NewHandler(searchService)

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	// ===== Auth routes =====
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// ===== Events routes (Req 2) =====
	r.Route("/events", func(r chi.Router) {
		// Public endpoints (no auth required)

		// GET all events
		r.Get("/", eventHandler.GetAllEvents)

		// Advanced search (Req 4 – uses search package, needs auth)
		r.With(authHandler.AuthMiddleware).Get("/search", searchHandler.SearchEvents)

		// GET events by organizer (public)
		r.Get("/organizer/{id}", eventHandler.GetEventsByOrganizer)

		// GET single event by ID (public)
		r.Get("/{id}", eventHandler.GetEventByID)

		// GET event attendees (public أو تخليها protected لو حابب)
		r.Get("/{id}/attendees", eventHandler.GetEventAttendees)

		// GET invitations for an event (logical to protect)
		r.With(authHandler.AuthMiddleware).Get("/{id}/invitations", invHandler.GetEventInvitations)

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

	// ===== Invitation routes (Req 3) =====
	r.Route("/invitations", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)

		// Send invitation
		r.Post("/", invHandler.SendInvitation)

		// Get my invitations
		r.Get("/my", invHandler.GetMyInvitations)

		// Respond to invitation
		r.Put("/{id}/respond", invHandler.RespondToInvitation)
	})

	// ===== Example protected API group =====
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
			json.NewEncoder(w).Encode(resp)
		})
	})

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

