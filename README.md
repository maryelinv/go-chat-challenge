# Go Chat Challenge

A simple real-time chat app built with Go and a lightweight React frontend.

## Features
- Real-time chat (WebSockets)
- Slash commands (e.g. `/stock=aapl.us` for stock quotes)
- Bot worker (RabbitMQ)
- SQLite message storage

## Quick Start
1. Install Go 1.22+ and Docker Desktop
2. Start RabbitMQ:
   ```sh
   docker compose up -d
   ```
3. Run servers (in separate terminals):
   ```sh
   go run ./cmd/api
   go run ./cmd/bot
   ```
4. Open [http://localhost:8080](http://localhost:8080) in your browser

## Usage
- Enter a username, send messages, and use `/stock=...` for quotes

## Project Structure
- `cmd/api` – API & WebSocket server
- `cmd/bot` – Bot worker
- `internal/` – Core logic (auth, chat, db, queue, stooq)
- `web/` – Frontend (HTML/CSS)

## License
MIT

