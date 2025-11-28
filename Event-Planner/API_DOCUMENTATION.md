# Event Planner API Documentation

Base URL: `http://localhost:8080`

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer YOUR_JWT_TOKEN
```

---

## 🔐 Authentication Endpoints

### Register User
**POST** `/auth/register`

Create a new user account.

**Request:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Email and password are required"
}
```

---

### Login User
**POST** `/auth/login`

Authenticate an existing user.

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
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "invalid email or password"
}
```

---

## 🎉 Event Management Endpoints

### Get All Events
**GET** `/events/`

Retrieve all events in the system (public endpoint).

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
    },
    {
      "id": 2,
      "title": "Team Meeting",
      "description": "Weekly team sync",
      "date": "2025-11-30",
      "time": "14:00:00",
      "location": "Conference Room A",
      "organizer_id": 2,
      "created_at": "2025-11-26T11:00:00Z"
    }
  ]
}
```

---

### Get Single Event
**GET** `/events/{id}`

Retrieve a specific event by ID (public endpoint).

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

**Error Response (404 Not Found):**
```json
{
  "error": "event not found"
}
```

---

### Get Events by Organizer
**GET** `/events/organizer/{id}`

Retrieve all events created by a specific organizer (public endpoint).

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

Create a new event (requires authentication).

**Headers:**
```
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

**Error Response (400 Bad Request):**
```json
{
  "error": "event title is required"
}
```

---

### Update Event
**PUT** `/events/{id}` 🔒

Update an existing event (requires authentication and ownership).

**Headers:**
```
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

**Error Response (403 Forbidden):**
```json
{
  "error": "you are not authorized to update this event"
}
```

---

### Delete Event
**DELETE** `/events/{id}` 🔒

Delete an event (requires authentication and ownership).

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Response (200 OK):**
```json
{
  "message": "event deleted successfully"
}
```

**Error Response (403 Forbidden):**
```json
{
  "error": "you are not authorized to delete this event"
}
```

---

## 👥 Event Attendance Endpoints

### Get Event Attendees
**GET** `/events/{id}/attendees`

Get all attendees for a specific event (public endpoint).

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
    },
    {
      "id": 3,
      "user_id": 3,
      "event_id": 1,
      "role": "collaborator",
      "status": "going",
      "created_at": "2025-11-26T11:30:00Z"
    }
  ]
}
```

---

### Join Event
**POST** `/events/{id}/join` 🔒

Join an event as an attendee (requires authentication).

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Response (200 OK):**
```json
{
  "message": "successfully joined event"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "failed to join event: duplicate key value..."
}
```

---

### Add Attendee to Event
**POST** `/events/{id}/attendees` 🔒

Add a user to an event with a specific role (requires authentication).

**Headers:**
```
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

**Response (200 OK):**
```json
{
  "message": "attendee added successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "invalid role: must be 'attendee' or 'collaborator'"
}
```

---

### Update Attendance Status
**PUT** `/events/{id}/attendance` 🔒

Update your attendance status for an event (requires authentication).

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json
```

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

**Error Response (400 Bad Request):**
```json
{
  "error": "invalid status: must be 'going', 'maybe', or 'not_going'"
}
```

---

## 📋 Personal Event Management

### Get My Attending Events
**GET** `/events/my/attending` 🔒

Get all events I'm attending (as organizer, attendee, or collaborator).

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
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
      "organizer_id": 2,
      "created_at": "2025-11-26T10:30:00Z",
      "role": "attendee",
      "status": "going"
    },
    {
      "id": 3,
      "title": "Workshop",
      "description": "Training workshop",
      "date": "2025-12-20",
      "time": "13:00:00",
      "location": "Training Room",
      "organizer_id": 1,
      "created_at": "2025-11-26T12:00:00Z",
      "role": "organizer",
      "status": "going"
    }
  ]
}
```

---

### Get My Organized Events
**GET** `/events/my/organized` 🔒

Get all events I'm organizing.

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

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

## 🛡️ Protected Route Example

### Get Profile
**GET** `/api/profile` 🔒

Example protected route showing user information.

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Response (200 OK):**
```json
{
  "message": "This is a protected route",
  "user_id": 1
}
```

---

## 🔍 Health Check

### Health Check
**GET** `/health`

Check if the server is running (public endpoint).

**Response (200 OK):**
```
Server is running
```

---

## 📊 Status Codes

- **200 OK** - Request successful
- **201 Created** - Resource created successfully
- **400 Bad Request** - Invalid request data
- **401 Unauthorized** - Authentication required or invalid
- **403 Forbidden** - Access denied (insufficient permissions)
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server error

---

## 🎯 Attendance Status Values

- **`going`** - User will attend the event
- **`maybe`** - User might attend the event
- **`not_going`** - User will not attend the event

---

## 👑 Role Values

- **`organizer`** - Event creator (automatically assigned)
- **`attendee`** - Regular participant
- **`collaborator`** - Helper with event management privileges

---

## 🚀 Example Workflow

1. **Register/Login** to get JWT token
2. **Create Event** using the token
3. **Add Attendees** to your event
4. **Update Attendance Status** as needed
5. **View Events** you're attending or organizing

---

## 📝 Notes

- All timestamps are in ISO 8601 format (UTC)
- Date format: `YYYY-MM-DD`
- Time format: `HH:MM:SS` (24-hour)
- JWT tokens expire after 7 days
- Event dates/times must be in the future when creating/updating
- Organizers are automatically added as attendees with "going" status