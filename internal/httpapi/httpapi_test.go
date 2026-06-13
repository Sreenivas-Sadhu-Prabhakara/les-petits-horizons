package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sreenivas-Sadhu-Prabhakara/les-petits-horizons/internal/prompt"
	"github.com/Sreenivas-Sadhu-Prabhakara/les-petits-horizons/internal/retriever"
)

type fakeRetriever struct {
	chunks      []retriever.Chunk
	gotTenant   *string
	gotAudience string
}

func (f *fakeRetriever) Search(ctx context.Context, q string, tenant *string, aud string, k int) ([]retriever.Chunk, error) {
	f.gotTenant = tenant
	f.gotAudience = aud
	return f.chunks, nil
}

type fakeGenerator struct {
	called bool
	toks   []string
	gotSys string
}

func (f *fakeGenerator) Stream(ctx context.Context, sys, user string, onToken func(string)) error {
	f.called = true
	f.gotSys = sys
	for _, t := range f.toks {
		onToken(t)
	}
	return nil
}

type fakeLogger struct {
	feedbackCalled bool
	loggedAnswer   string
	gotMode        string
	gotTenant      *string
}

func (f *fakeLogger) StartConversation(ctx context.Context, tenant, user *string, mode, channel string) (string, error) {
	f.gotMode = mode
	return "conv-1", nil
}
func (f *fakeLogger) LogTurn(ctx context.Context, convID string, tenant *string, q, a, model string, ids []string, ms int) (string, string, error) {
	f.loggedAnswer = a
	f.gotTenant = tenant
	return "user-1", "asst-1", nil
}
func (f *fakeLogger) RecordFeedback(ctx context.Context, msgID string, tenant, user *string, rating string, solved *bool, note string) error {
	f.feedbackCalled = true
	return nil
}

func newTestServer(r Retriever, g *fakeGenerator, l Logger) http.Handler {
	return Handler(Deps{Retriever: r, Generator: g, Logger: l, GenModel: "qwen-coder-32b"})
}

func post(h http.Handler, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestChatHappyPathStreamsAndLogs(t *testing.T) {
	gen := &fakeGenerator{toks: []string{"Bonjour", " !"}}
	log := &fakeLogger{}
	h := newTestServer(&fakeRetriever{chunks: []retriever.Chunk{
		{ChunkID: "11111111-1111-1111-1111-111111111111", Title: "Salutations", Content: "bonjour", Score: 0.2},
	}}, gen, log)

	rec := post(h, "/tutor/chat", `{"message":"comment dire hello en français?"}`, nil)
	body := rec.Body.String()

	if !gen.called {
		t.Fatal("generator should have been called")
	}
	if !strings.Contains(body, `"token":"Bonjour"`) {
		t.Fatalf("streamed tokens missing: %q", body)
	}
	if !strings.Contains(body, "event: done") || !strings.Contains(body, `"message_id":"asst-1"`) {
		t.Fatalf("done event / message id missing: %q", body)
	}
	if !strings.Contains(body, "Salutations") {
		t.Fatalf("source title missing in done event: %q", body)
	}
	if log.loggedAnswer != "Bonjour !" {
		t.Fatalf("logged answer wrong: %q", log.loggedAnswer)
	}
}

func TestChatStillAnswersWithoutMaterial(t *testing.T) {
	// A language tutor must keep helping even when retrieval finds nothing —
	// unlike the closed-domain solar bot which escalated to a human.
	gen := &fakeGenerator{toks: []string{"Salut"}}
	h := newTestServer(&fakeRetriever{chunks: nil}, gen, &fakeLogger{})
	rec := post(h, "/tutor/chat", `{"message":"salut"}`, nil)
	if !gen.called {
		t.Fatal("generator should still run with no retrieved material")
	}
	if !strings.Contains(rec.Body.String(), `"token":"Salut"`) {
		t.Fatalf("expected streamed answer: %q", rec.Body.String())
	}
}

func TestChatUsesGradedPersonaForLevel(t *testing.T) {
	gen := &fakeGenerator{toks: []string{"ok"}}
	h := newTestServer(&fakeRetriever{chunks: []retriever.Chunk{{Title: "x", Content: "y"}}}, gen, &fakeLogger{})
	post(h, "/tutor/chat", `{"message":"bonjour","level":"c1"}`, nil)
	if gen.gotSys != prompt.LevelByName("c1").System {
		t.Fatalf("expected the c1 graded persona, got %q", gen.gotSys)
	}
}

func TestChatDefaultsToA1Level(t *testing.T) {
	gen := &fakeGenerator{toks: []string{"ok"}}
	log := &fakeLogger{}
	h := newTestServer(&fakeRetriever{chunks: []retriever.Chunk{{Title: "x", Content: "y"}}}, gen, log)
	post(h, "/tutor/chat", `{"message":"bonjour"}`, nil)
	if gen.gotSys != prompt.LevelByName("a1").System {
		t.Fatalf("expected a1 persona by default, got %q", gen.gotSys)
	}
	if log.gotMode != "a1" {
		t.Fatalf("expected conversation logged with level a1, got %q", log.gotMode)
	}
}

func TestChatRetrievesFromFrenchAudience(t *testing.T) {
	ret := &fakeRetriever{chunks: []retriever.Chunk{{Title: "x", Content: "y"}}}
	h := newTestServer(ret, &fakeGenerator{toks: []string{"ok"}}, &fakeLogger{})
	post(h, "/tutor/chat", `{"message":"bonjour"}`, nil)
	if ret.gotAudience != prompt.Audience {
		t.Fatalf("expected retrieval audience %q, got %q", prompt.Audience, ret.gotAudience)
	}
}

func TestChatUsesShortPromptWhenEnabled(t *testing.T) {
	gen := &fakeGenerator{toks: []string{"ok"}}
	h := Handler(Deps{
		Retriever: &fakeRetriever{chunks: []retriever.Chunk{{Title: "x", Content: "y"}}},
		Generator: gen, Logger: &fakeLogger{}, GenModel: "tuned", ShortPrompt: true,
	})
	post(h, "/tutor/chat", `{"message":"bonjour","level":"b1"}`, nil)
	if gen.gotSys != prompt.LevelByName("b1").Short {
		t.Fatalf("expected short b1 persona, got %q", gen.gotSys)
	}
}

func TestChatRejectsEmptyMessage(t *testing.T) {
	h := newTestServer(&fakeRetriever{}, &fakeGenerator{}, &fakeLogger{})
	rec := post(h, "/tutor/chat", `{"message":"   "}`, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty message, got %d", rec.Code)
	}
}

func TestCORSPreflight(t *testing.T) {
	h := newTestServer(&fakeRetriever{}, &fakeGenerator{}, &fakeLogger{})
	req := httptest.NewRequest(http.MethodOptions, "/tutor/chat", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for preflight, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing permissive CORS header")
	}
}

func TestFeedbackValidatesAndRecords(t *testing.T) {
	log := &fakeLogger{}
	h := newTestServer(&fakeRetriever{}, &fakeGenerator{}, log)

	bad := post(h, "/tutor/feedback", `{"message_id":"asst-1"}`, nil)
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing rating, got %d", bad.Code)
	}
	if log.feedbackCalled {
		t.Fatal("feedback should not be recorded on invalid input")
	}

	ok := post(h, "/tutor/feedback", `{"message_id":"asst-1","rating":"down"}`, nil)
	if ok.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", ok.Code)
	}
	if !log.feedbackCalled {
		t.Fatal("valid feedback should be recorded")
	}
}

func TestPermissiveAuthPassesTenantHeader(t *testing.T) {
	ret := &fakeRetriever{chunks: []retriever.Chunk{{Title: "x", Content: "y", Score: 0.2}}}
	h := newTestServer(ret, &fakeGenerator{toks: []string{"ok"}}, &fakeLogger{})

	post(h, "/tutor/chat", `{"message":"bonjour"}`,
		map[string]string{"X-Tenant-Id": "tenant-xyz"})
	if ret.gotTenant == nil || *ret.gotTenant != "tenant-xyz" {
		t.Fatalf("expected tenant from header to reach retriever, got %v", ret.gotTenant)
	}
}
