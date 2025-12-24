# Membership Tracker

A full-stack web application for managing member registrations and visits. Built with a Go backend and a React frontend.

## Features

- **Member Management**: Register new members with name and phone number.
- **Visit Tracking**: Track and update the number of visits for each member.
- **Search**: Quickly find members by phone number.
- **CRUD Operations**: Full Create, Read, Update, and Delete capabilities.

## Tech Stack

- **Frontend**: React, Vite, React Router
- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **Deployment**: Render (Backend), Vercel (Frontend)

## Getting Started

### Prerequisites

- Go 1.23+
- Node.js & npm
- PostgreSQL database

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/membership-tracker.git
   cd membership-tracker
   ```

2. **Backend Setup**
   ```bash
   cd backend
   
   # Install dependencies
   go mod tidy
   
   # Set up environment variables (Windows PowerShell)
   $env:DATABASE_URL="postgres://user:password@host:5432/dbname"
   $env:PORT="8080"
   
   # Run the server
   go run .
   ```

3. **Frontend Setup**
   ```bash
   cd frontend
   
   # Install dependencies
   npm install
   
   # Create a .env file pointing to your backend
   # Create a file named .env in the frontend folder with:
   # VITE_API_URL=http://localhost:8080
   
   # Start the development server
   npm run dev
   ```

## API Endpoints

- `GET /get-member?phone_number=...` - Retrieve member details
- `POST /add-member` - Register a new member
- `PUT /update-member` - Update member details or visit count
- `DELETE /delete-member?phone_number=...` - Remove a member
