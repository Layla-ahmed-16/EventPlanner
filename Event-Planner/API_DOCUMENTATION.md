# Event Planner API Documentation

**Base URL:** `http://localhost:8080`  

All protected endpoints require a JWT token:

```http
Authorization: Bearer YOUR_JWT_TOKEN
````
## Requirement 1 – User Management (`/auth`)

### Register User

**POST** `/auth/register`

Create a new user.

**Request:**

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
}
```

**Error (400 Bad Request):**

```json
{
  "error": "Email and password are required"
}
```

---

### Login User

**POST** `/auth/login`

Authenticate and get JWT token.

**Request:**

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
}
```

**Error (401 Unauthorized):**

```json
{
  "error": "invalid email or password"
}
```

---

## Requirement 2 – Event Management (`/events`)

### Get All Events

**GET** `/events/`

Public – retrieve all events.

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 1,
      "title": "Tech Conference 2025",
      "description": "Annual technology conference",
      "date": "2025-12-15",
      "time": "09:00:00",
      "location": "Convention Center",
      "organizer_id": 1,
      "created_at": "2025-11-26T10:30:00Z"
    }
  ]
}
```

---

### Get Single Event

**GET** `/events/{id}`

Public – retrieve a specific event by ID.

**Response (200 OK):**

```json
{
  "data": {
    "id": 1,
    "title": "Tech Conference 2025",
    "description": "Annual technology conference",
    "date": "2025-12-15",
    "time": "09:00:00",
    "location": "Convention Center",
    "organizer_id": 1,
    "created_at": "2025-11-26T10:30:00Z"
  }
}
```

**Error (404 Not Found):**

```json
{
  "error": "failed to get event: sql: no rows in result set"
}
```

---

### Get Events by Organizer

**GET** `/events/organizer/{id}`

Public – get all events created by a specific organizer.

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 1,
      "title": "Tech Conference 2025",
      "description": "Annual technology conference",
      "date": "2025-12-15",
      "time": "09:00:00",
      "location": "Convention Center",
      "organizer_id": 1,
      "created_at": "2025-11-26T10:30:00Z"
    }
  ]
}
```

---

### Create Event

**POST** `/events/` 🔒

Create a new event (current user is the organizer).

**Headers:**

```http
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json
```

**Request:**

```json
{
  "title": "Tech Conference 2025",
  "description": "Annual technology conference",
  "date": "2025-12-15",
  "time": "09:00:00",
  "location": "Convention Center"
}
```

**Response (201 Created):**

```json
{
  "message": "event created successfully",
  "data": {
    "id": 1,
    "title": "Tech Conference 2025",
    "description": "Annual technology conference",
    "date": "2025-12-15",
    "time": "09:00:00",
    "location": "Convention Center",
    "organizer_id": 1,
    "created_at": "2025-11-26T10:30:00Z"
  }
}
```

**Validation Error (400 Bad Request):**

```json
{
  "error": "event date is required"
}
```

---

### Update Event

**PUT** `/events/{id}` 🔒

Requires authentication and ownership (must be organizer).

**Headers:**

```http
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json
```

**Request:**

```json
{
  "title": "Updated Tech Conference 2025",
  "description": "Updated description",
  "date": "2025-12-16",
  "time": "10:00:00",
  "location": "Updated Location"
}
```

**Response (200 OK):**

```json
{
  "message": "event updated successfully",
  "data": {
    "id": 1,
    "title": "Updated Tech Conference 2025",
    "description": "Updated description",
    "date": "2025-12-16",
    "time": "10:00:00",
    "location": "Updated Location",
    "organizer_id": 1,
    "created_at": "2025-11-26T10:30:00Z"
  }
}
```

**Error (403 Forbidden):**

```json
{
  "error": "you are not authorized to update this event"
}
```

---

### Delete Event

**DELETE** `/events/{id}` 🔒

Requires authentication and organizer ownership.

**Headers:**

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

**Response (200 OK):**

```json
{
  "message": "event deleted successfully"
}
```

**Error (403 Forbidden):**

```json
{
  "error": "you are not authorized to delete this event"
}
```

---

##  Requirement 3 – Response Management

### A Attendance Management (`event_attendees`)

#### Get Event Attendees

**GET** `/events/{id}/attendees`

Public – list all attendees for an event, with roles and statuses.

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "event_id": 1,
      "role": "organizer",
      "status": "going",
      "created_at": "2025-11-26T10:30:00Z"
    },
    {
      "id": 2,
      "user_id": 2,
      "event_id": 1,
      "role": "attendee",
      "status": "maybe",
      "created_at": "2025-11-26T11:00:00Z"
    }
  ]
}
```

---

#### Join Event

**POST** `/events/{id}/join` 🔒

Current user joins an event as an attendee.

**Headers:**

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

**Response (200 OK):**

```json
{
  "message": "successfully joined event"
}
```

**Error (400 Bad Request):**

```json
{
  "error": "failed to join event: duplicate key value..."
}
```

---

#### Update Attendance Status

**PUT** `/events/{id}/attendance` 🔒

Update your attendance status.

**Headers:**

```http
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json
```

**Allowed values:** `"going"`, `"maybe"`, `"not_going"`

**Request:**

```json
{
  "status": "maybe"
}
```

**Response (200 OK):**

```json
{
  "message": "attendance status updated successfully"
}
```

**Error (400 Bad Request):**

```json
{
  "error": "invalid status: must be 'going', 'maybe', or 'not_going'"
}
```

---

### B Invitations Management

#### Invite User to Event (by user_id)

**POST** `/events/{id}/invite` 🔒

Organizer invites an existing user to an event.

**Headers:**

```http
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json
```

**Request:**

```json
{
  "user_id": 2,
  "role": "attendee"
}
```

**Role values:** `"attendee"`, `"collaborator"`, `"organizer"`

**Response (200 OK):**

```json
{
  "message": "user invited to event successfully"
}
```

**Error (403 Forbidden):**

```json
{
  "error": "only the event creator can invite users to this event"
}
```

---

#### Send Invitation by Email

**POST** `/invitations` 🔒

Send an invitation by email.

**Request:**

```json
{
  "event_id": 1,
  "invitee_email": "guest@example.com",
  "role": "attendee",
  "message": "Please join our tech conference."
}
```

**Response (201 Created):**

```json
{
  "message": "invitation sent successfully",
  "data": {
    "id": 10,
    "event_id": 1,
    "inviter_id": 1,
    "invitee_email": "guest@example.com",
    "invitee_id": 2,
    "role": "attendee",
    "status": "pending",
    "message": "Please join our tech conference.",
    "created_at": "2025-11-26T12:00:00Z"
  }
}
```

---

#### Get My Invitations

**GET** `/invitations/my?email=user@example.com` 🔒

Retrieve invitations sent to the given email.

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 10,
      "event_id": 1,
      "inviter_id": 1,
      "invitee_email": "user@example.com",
      "invitee_id": 2,
      "role": "attendee",
      "status": "pending",
      "message": "Please join our tech conference.",
      "created_at": "2025-11-26T12:00:00Z",
      "responded_at": null,
      "event_title": "Tech Conference 2025",
      "event_date": "2025-12-15",
      "event_time": "09:00:00",
      "event_location": "Convention Center",
      "inviter_email": "organizer@example.com"
    }
  ]
}
```

---

#### Get Invitations for an Event

**GET** `/events/{id}/invitations` 🔒

Retrieve all invitations for a specific event.

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 10,
      "event_id": 1,
      "inviter_id": 1,
      "invitee_email": "user@example.com",
      "invitee_id": 2,
      "role": "attendee",
      "status": "pending",
      "message": "Please join our tech conference.",
      "created_at": "2025-11-26T12:00:00Z",
      "responded_at": null,
      "event_title": "Tech Conference 2025",
      "event_date": "2025-12-15",
      "event_time": "09:00:00",
      "event_location": "Convention Center",
      "inviter_email": "organizer@example.com"
    }
  ]
}
```

---

#### Respond to Invitation

**PUT** `/invitations/{id}/respond?email=user@example.com` 🔒

Accept or decline an invitation.

**Request:**

```json
{
  "status": "accepted"
}
```

**Response (200 OK):**

```json
{
  "message": "invitation response recorded successfully"
}
```

**Error (400 Bad Request):**

```json
{
  "error": "invalid status: must be 'accepted' or 'declined'"
}
```

---

##  Requirement 4 – Search & Filtering

### Advanced Event Search

**GET** `/events/search` 🔒

Search events the current user is involved in.

**Headers:**

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

**Query Parameters:**

* `q` – keyword in title/description (optional)
* `date_from` – `YYYY-MM-DD` (optional)
* `date_to` – `YYYY-MM-DD` (optional)
* `role` – `organizer` | `attendee` | `collaborator` (optional)
* `status` – `going` | `maybe` | `not_going` (optional)

**Example:**

```http
GET /events/search?q=conference&date_from=2025-12-01&date_to=2025-12-31&role=organizer&status=going
```

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 1,
      "title": "Tech Conference 2025",
      "description": "Annual technology conference",
      "date": "2025-12-15",
      "time": "09:00:00",
      "location": "Convention Center",
      "organizer_id": 1,
      "created_at": "2025-11-26T10:30:00Z",
      "role": "organizer",
      "status": "going"
    }
  ]
}
```

---

##  Personal Event Views

### Get My Attending Events

**GET** `/events/my/attending` 🔒

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 1,
      "title": "Tech Conference 2025",
      "description": "Annual technology conference",
      "date": "2025-12-15",
      "time": "09:00:00",
      "location": "Convention Center",
      "organizer_id": 2,
      "created_at": "2025-11-26T10:30:00Z",
      "role": "attendee",
      "status": "going"
    }
  ]
}
```

---

### Get My Organized Events

**GET** `/events/my/organized` 🔒

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 3,
      "title": "Workshop",
      "description": "Training workshop",
      "date": "2025-12-20",
      "time": "13:00:00",
      "location": "Training Room",
      "organizer_id": 1,
      "created_at": "2025-11-26T12:00:00Z"
    }
  ]
}
```

---

##  Protected Route Example

### Get Profile

**GET** `/api/profile` 🔒

**Response (200 OK):**

```json
{
  "message": "This is a protected route",
  "user_id": 1
}
```

---

##  Health Check

### Health Check

**GET** `/health`

**Response (200 OK):**

```text
Server is running
```

---

##  Status Codes

* `200 OK` – success
* `201 Created` – resource created
* `400 Bad Request` – invalid input
* `401 Unauthorized` – missing/invalid token
* `403 Forbidden` – not enough permissions
* `404 Not Found` – resource not found
* `500 Internal Server Error` – unexpected server error