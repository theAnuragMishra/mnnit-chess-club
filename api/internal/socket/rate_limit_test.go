package socket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/ratelimit"
)

func TestMessageLimitAcrossConnections(t *testing.T) {
	manager := NewManager(func(event Event, client *Client) error {
		client.Send(event)
		return nil
	}, func(*Client) {})
	manager.messageLimit = ratelimit.New[int32](0.01, 2, 8)
	username := "player"
	session := database.GetSessionRow{UserID: 1, Username: &username}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), auth.MiddlewareSentSession, session)
		manager.ServeWS(w, r.WithContext(ctx))
	}))
	defer server.Close()

	var connections []*websocket.Conn
	for range 2 {
		conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		connections = append(connections, conn)
		if err := conn.WriteJSON(Event{Type: "move"}); err != nil {
			t.Fatal(err)
		}
		var echo Event
		if err := conn.ReadJSON(&echo); err != nil || echo.Type != "move" {
			t.Fatalf("expected an echoed move, got %v (%v)", echo, err)
		}
	}
	conn := connections[1]
	if err := conn.WriteJSON(Event{Type: "move"}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := conn.ReadMessage(); !websocket.IsCloseError(err, websocket.ClosePolicyViolation) {
		t.Fatalf("expected a rate-limit close, got %v", err)
	}
}
