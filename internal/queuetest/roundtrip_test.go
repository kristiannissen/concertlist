// Package queuetest exercises the publish/receive/acknowledge round trip
// against a fake HTTP server that reproduces the Vercel Queues API contract
// (https://{region}.vercel-queue.com/api/v3/...) in memory.
//
// This is a unit test: it has no network dependency and no OIDC token
// requirement, so it always runs, in CI and everywhere else. It exists to
// pin down the request/response shapes this codebase relies on (URL
// patterns, status codes, the base64-encoded ndjson receive format, the
// lease-based acknowledge) so a change that breaks that contract fails
// immediately, without needing a live Vercel project or credentials.
//
// It does NOT verify that the real Vercel Queues API still behaves this
// way, nor that push-mode delivery reaches EventScrapeConsumer/
// EventExtractConsumer in production. Verifying against the live API is a
// manual/integration concern, not something a unit test should depend on.
package queuetest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	testTopic    = "local-roundtrip-test"
	testConsumer = "local-roundtrip-consumer"
	fakeToken    = "fake-oidc-token-for-tests"
)

// fakeQueue is a minimal in-memory stand-in for a single Vercel Queues
// topic/consumer pair: messages are published, then handed out one at a
// time on receive (with a receipt handle), then removed on acknowledge.
type fakeQueue struct {
	mu      sync.Mutex
	pending [][]byte          // published, not yet received
	leased  map[string][]byte // receiptHandle -> message body, received but not yet acked
	nextID  int
}

func newFakeQueueServer(t *testing.T) *httptest.Server {
	t.Helper()
	q := &fakeQueue{leased: map[string][]byte{}}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v3/topic/{topic}", func(w http.ResponseWriter, r *http.Request) {
		if !authorized(w, r) {
			return
		}
		if r.PathValue("topic") != testTopic {
			http.NotFound(w, r)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		q.mu.Lock()
		q.pending = append(q.pending, body)
		q.mu.Unlock()

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{}`))
	})

	mux.HandleFunc("POST /api/v3/topic/{topic}/consumer/{consumer}", func(w http.ResponseWriter, r *http.Request) {
		if !authorized(w, r) {
			return
		}
		if r.PathValue("topic") != testTopic || r.PathValue("consumer") != testConsumer {
			http.NotFound(w, r)
			return
		}

		q.mu.Lock()
		defer q.mu.Unlock()
		if len(q.pending) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		msg := q.pending[0]
		q.pending = q.pending[1:]
		q.nextID++
		msgID := fmt.Sprintf("fake-msg-%d", q.nextID)
		handle := fmt.Sprintf("fake-lease-%d", q.nextID)
		q.leased[handle] = msg

		line, err := json.Marshal(struct {
			MessageID     string `json:"messageId"`
			ReceiptHandle string `json:"receiptHandle"`
			Body          string `json:"body"`
		}{
			MessageID:     msgID,
			ReceiptHandle: handle,
			Body:          base64.StdEncoding.EncodeToString(msg),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusOK)
		w.Write(line)
		w.Write([]byte("\n"))
	})

	mux.HandleFunc("DELETE /api/v3/topic/{topic}/consumer/{consumer}/lease/{handle}", func(w http.ResponseWriter, r *http.Request) {
		if !authorized(w, r) {
			return
		}
		if r.PathValue("topic") != testTopic || r.PathValue("consumer") != testConsumer {
			http.NotFound(w, r)
			return
		}
		handle, err := url.QueryUnescape(r.PathValue("handle"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		q.mu.Lock()
		defer q.mu.Unlock()
		if _, ok := q.leased[handle]; !ok {
			http.NotFound(w, r)
			return
		}
		delete(q.leased, handle)
		w.WriteHeader(http.StatusNoContent)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func authorized(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Authorization") != "Bearer "+fakeToken {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Authentication required - provide a valid Vercel OIDC token"}`))
		return false
	}
	return true
}

func TestQueuePublishReceiveAcknowledgeRoundTrip(t *testing.T) {
	srv := newFakeQueueServer(t)

	client := resty.New().SetAuthToken(fakeToken)
	base := srv.URL + "/api/v3"

	want := fmt.Sprintf("local-test-%d", time.Now().UnixNano())
	body, err := json.Marshal(map[string]string{"check": want})
	if err != nil {
		t.Fatalf("marshal test payload: %v", err)
	}

	// 1. SendMessage
	sendResp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(base + "/topic/" + testTopic)
	if err != nil {
		t.Fatalf("publish request failed: %v", err)
	}
	if sendResp.StatusCode() != 201 && sendResp.StatusCode() != 202 {
		t.Fatalf("publish returned %d: %s", sendResp.StatusCode(), sendResp.String())
	}
	t.Logf("published: %s", sendResp.String())

	// 2. ReceiveMessages — poll briefly since the fake mimics async delivery.
	var receiptHandle, msgID, got string
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		recvResp, err := client.R().
			SetHeader("Accept", "application/x-ndjson").
			Post(base + "/topic/" + testTopic + "/consumer/" + testConsumer)
		if err != nil {
			t.Fatalf("receive request failed: %v", err)
		}
		if recvResp.StatusCode() == 204 {
			time.Sleep(20 * time.Millisecond)
			continue
		}
		if recvResp.StatusCode() != 200 {
			t.Fatalf("receive returned %d: %s", recvResp.StatusCode(), recvResp.String())
		}

		firstLine := strings.SplitN(strings.TrimSpace(recvResp.String()), "\n", 2)[0]
		var line struct {
			MessageID     string `json:"messageId"`
			ReceiptHandle string `json:"receiptHandle"`
			Body          string `json:"body"`
		}
		if err := json.Unmarshal([]byte(firstLine), &line); err != nil {
			t.Fatalf("decode receive response %q: %v", firstLine, err)
		}

		decoded, err := base64.StdEncoding.DecodeString(line.Body)
		if err != nil {
			t.Fatalf("base64-decode message body: %v", err)
		}
		var payload struct {
			Check string `json:"check"`
		}
		if err := json.Unmarshal(decoded, &payload); err != nil {
			t.Fatalf("unmarshal message payload %q: %v", decoded, err)
		}

		msgID, receiptHandle, got = line.MessageID, line.ReceiptHandle, payload.Check
		break
	}

	if receiptHandle == "" {
		t.Fatal("publish succeeded but the message was never received within the deadline")
	}
	if got != want {
		t.Fatalf("received payload %q, want %q", got, want)
	}
	t.Logf("received message %s matching what was published", msgID)

	// 3. AcknowledgeMessage, so the fake topic doesn't accumulate messages.
	ackResp, err := client.R().
		Delete(base + "/topic/" + testTopic + "/consumer/" + testConsumer + "/lease/" + url.QueryEscape(receiptHandle))
	if err != nil {
		t.Fatalf("acknowledge request failed: %v", err)
	}
	if ackResp.StatusCode() != 204 {
		t.Fatalf("acknowledge returned %d: %s", ackResp.StatusCode(), ackResp.String())
	}
}
