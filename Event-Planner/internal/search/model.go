package search

// EventsFilter holds filters for searching events
type EventsFilter struct {
	Query    string // keyword (title / description)
	DateFrom string // YYYY-MM-DD (optional)
	DateTo   string // YYYY-MM-DD (optional)
	Role     string // 'organizer', 'attendee', 'collaborator' (optional)
	Status   string // 'going', 'maybe', 'not_going' (optional)
	UserID   int    // current user ID (required)
}
