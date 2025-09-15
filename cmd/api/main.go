package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/maryelinv/go-chat-challenge/internal/auth"
	"github.com/maryelinv/go-chat-challenge/internal/chat"
	"github.com/maryelinv/go-chat-challenge/internal/db"
	"github.com/maryelinv/go-chat-challenge/internal/queue"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

type postMessageReq struct{ Room, Username, Text string }

type wsEnvelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func main() {
	g, _ := db.Open("chat.db")
	_ = db.Migrate(g)
	g.Exec("UPDATE rooms SET is_private = ? WHERE is_private IS NULL", false)

	hub := chat.NewHub()
	q := queue.MustConnectFromEnv()

	r := chi.NewRouter()

	r.Post("/register", register(g))
	r.Post("/login", login(g))
	r.Post("/logout", logout())
	r.Get("/me", me(g))

	r.Post("/messages", func(w http.ResponseWriter, r *http.Request) {
		var in postMessageReq
		json.NewDecoder(r.Body).Decode(&in)

		if isStockCmd(in.Text) {
			code := parseStockCode(in.Text)
			_ = q.PublishRequest(code, in.Room)
			w.WriteHeader(202)
			return
		}

		m := db.Message{Room: in.Room, Username: in.Username, Text: in.Text, CreatedAt: time.Now()}
		if err := db.SaveMessage(g, &m); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		b, _ := json.Marshal(wsEnvelope{Type: "message", Payload: m})
		hub.Broadcast(in.Room, b)
		w.WriteHeader(201)
	})

	r.Get("/ws/{room}", func(w http.ResponseWriter, r *http.Request) {
		room := chi.URLParam(r, "room")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c := &chat.Client{Conn: conn, Send: make(chan []byte, 8), Room: room}
		hub.Join(c)
		defer func() { hub.Leave(c); conn.Close() }()
		go func() {
			for msg := range c.Send {
				conn.WriteMessage(websocket.TextMessage, msg)
			}
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	})

	r.Get("/messages/{room}", func(w http.ResponseWriter, r *http.Request) {
		room := chi.URLParam(r, "room")
		user := r.URL.Query().Get("user")
		limit := 50

		var msgs []db.Message
		var err error
		if user != "" {
			msgs, err = db.LastNMessagesByUser(g, room, user, limit)
		} else {
			msgs, err = db.LastNMessages(g, room, limit)
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
			msgs[i], msgs[j] = msgs[j], msgs[i]
		}
		json.NewEncoder(w).Encode(msgs)
	})

	r.With(requireAuth).Post("/rooms", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Name      string `json:"name"`
			IsPrivate *bool  `json:"is_private"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		uid, ok := auth.CurrentUserID(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var u db.User
		if err := g.First(&u, uid).Error; err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		name := strings.ToLower(strings.TrimSpace(in.Name))
		if name == "" {
			http.Error(w, "name required", http.StatusBadRequest)
			return
		}

		isPrivate := false
		if in.IsPrivate != nil {
			isPrivate = *in.IsPrivate
		}

		rm := db.Room{Name: name, Creator: u.Username, IsPrivate: isPrivate}
		if err := g.Create(&rm).Error; err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "unique") {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/rooms", func(w http.ResponseWriter, r *http.Request) {
		type nameOnly struct{ Name string }

		var public []nameOnly
		if err := g.Model(&db.Room{}).
			Select("name").
			Where("is_private = ? OR is_private IS NULL", false).
			Find(&public).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := strings.TrimSpace(r.URL.Query().Get("user"))
		var privateOwned []nameOnly
		if user != "" {
			_ = g.Model(&db.Room{}).
				Select("name").
				Where("creator = ? AND is_private = ?", user, true).
				Find(&privateOwned).Error
		}

		seen := map[string]struct{}{}
		add := func(list []nameOnly) {
			for _, r := range list {
				if r.Name != "" {
					seen[r.Name] = struct{}{}
				}
			}
		}
		add(public)
		add(privateOwned)
		seen["general"] = struct{}{}

		out := make([]string, 0, len(seen))
		for k := range seen {
			out = append(out, k)
		}
		_ = json.NewEncoder(w).Encode(out)
	})

	r.With(requireAuth).Post("/messages", func(w http.ResponseWriter, r *http.Request) {
		var in postMessageReq
		json.NewDecoder(r.Body).Decode(&in)

		if uid, ok := auth.CurrentUserID(r); ok {
			var u db.User
			if err := g.First(&u, uid).Error; err == nil {
				in.Username = u.Username
			}
		}

		if isStockCmd(in.Text) {
			code := parseStockCode(in.Text)
			_ = q.PublishRequest(code, in.Room)
			w.WriteHeader(202)
			return
		}
		m := db.Message{Room: in.Room, Username: in.Username, Text: in.Text, CreatedAt: time.Now()}
		if err := db.SaveMessage(g, &m); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		b, _ := json.Marshal(wsEnvelope{Type: "message", Payload: m})
		hub.Broadcast(in.Room, b)
		w.WriteHeader(201)
	})

	go func() {
		q.ConsumeReplies(func(room, text string) {
			m := db.Message{Room: room, Username: "stock-bot", Text: text, CreatedAt: time.Now()}
			if err := db.SaveMessage(g, &m); err != nil {
				log.Println("save:", err)
				return
			}
			b, _ := json.Marshal(wsEnvelope{Type: "message", Payload: m})
			hub.Broadcast(room, b)
		})
	}()

	fs := http.FileServer(http.Dir("web"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
	r.Handle("/*", fs)

	log.Println("API on :8080")
	http.ListenAndServe(":8080", r)
}
