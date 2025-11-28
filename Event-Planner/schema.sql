-- ==========================
-- USERS TABLE (AUTH)
-- ==========================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Index on email for faster lookups
CREATE INDEX idx_users_email ON users(email);


-- ==========================
-- EVENTS TABLE
-- ==========================
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    time TIME NOT NULL,
    location TEXT NOT NULL,
    organizer_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for faster lookups
CREATE INDEX idx_events_organizer ON events(organizer_id);
CREATE INDEX idx_events_date_time ON events(date, time);


-- ==========================
-- EVENT_ATTENDEES TABLE
-- ==========================
CREATE TABLE event_attendees (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id INT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('organizer', 'attendee', 'collaborator')),
    status TEXT NOT NULL DEFAULT 'going' CHECK (status IN ('going', 'maybe', 'not_going')),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, event_id)
);

-- Indexes for faster lookups (search & filters)
CREATE INDEX idx_event_attendees_user ON event_attendees(user_id);
CREATE INDEX idx_event_attendees_event ON event_attendees(event_id);
CREATE INDEX idx_event_attendees_status ON event_attendees(status);
CREATE INDEX idx_event_attendees_role ON event_attendees(role);


-- ==========================
-- INVITATIONS TABLE
-- ==========================
-- matches internal/invitation models & repo:
--  ID, EventID, InviterID, InviteeEmail, InviteeID (nullable),
--  Role ('attendee','collaborator','organizer'),
--  Status ('pending','accepted','declined'),
--  Message, CreatedAt, RespondedAt
CREATE TABLE invitations (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    inviter_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    invitee_email TEXT NOT NULL,
    invitee_id INT REFERENCES users(id) ON DELETE SET NULL,

    role TEXT NOT NULL CHECK (role IN ('attendee', 'collaborator', 'organizer')),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined')),

    message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    responded_at TIMESTAMP NULL
);

-- Indexes to speed up:
-- - "my invitations" lookup by email
-- - event invitations listing
-- - filtering by inviter / status
CREATE INDEX idx_invitations_invitee_email ON invitations(invitee_email);
CREATE INDEX idx_invitations_event ON invitations(event_id);
CREATE INDEX idx_invitations_inviter ON invitations(inviter_id);
CREATE INDEX idx_invitations_status ON invitations(status);
CREATE INDEX idx_invitations_created_at ON invitations(created_at);
