package prompt

import (
	"strings"
	"testing"

	"github.com/Sreenivas-Sadhu-Prabhakara/les-petits-horizons/internal/retriever"
)

func TestLevelByNameDefaultsToA1(t *testing.T) {
	if LevelByName("").Name != "a1" {
		t.Fatal("empty level should default to a1")
	}
	if LevelByName("nonsense").Name != "a1" {
		t.Fatal("unknown level should default to a1")
	}
	if LevelByName("B2").Name != "b2" {
		t.Fatal("level codes should be case-insensitive")
	}
}

func TestAllLevelsShareAudience(t *testing.T) {
	for _, code := range []string{"a1", "a2", "b1", "b2", "c1", "c2"} {
		if LevelByName(code).Audience != Audience {
			t.Fatalf("level %s should retrieve from audience %q", code, Audience)
		}
	}
}

func TestAssembleForGradesPersonaByLevel(t *testing.T) {
	chunks := []retriever.Chunk{{Title: "Salutations", Content: "bonjour, salut"}}

	a1Sys, _ := AssembleFor(LevelByName("a1"), "bonjour", chunks)
	if !strings.Contains(a1Sys, "A1") {
		t.Fatalf("a1 persona missing level guidance: %q", a1Sys)
	}
	c2Sys, _ := AssembleFor(LevelByName("c2"), "bonjour", chunks)
	if !strings.Contains(c2Sys, "C2") {
		t.Fatalf("c2 persona missing level guidance: %q", c2Sys)
	}
	if a1Sys == c2Sys {
		t.Fatal("different levels must produce different personas")
	}
	// The shared base persona ("Petit Horizon") is present at every level.
	if !strings.Contains(a1Sys, "Petit Horizon") || !strings.Contains(c2Sys, "Petit Horizon") {
		t.Fatal("base tutor persona missing")
	}
}

func TestUserPromptIncludesMaterialAndMessage(t *testing.T) {
	chunks := []retriever.Chunk{{Title: "Salutations", Content: "bonjour, salut"}}
	_, user := Assemble("comment dire hello?", chunks)
	if !strings.Contains(user, "bonjour, salut") {
		t.Fatal("user prompt missing material content")
	}
	if !strings.Contains(user, "comment dire hello?") {
		t.Fatal("user prompt missing the learner message")
	}
}

func TestUserPromptStillHelpsWithoutMaterial(t *testing.T) {
	_, user := Assemble("salut", nil)
	if !strings.Contains(user, "aucun matériel") {
		t.Fatalf("expected an explicit no-material signal: %q", user)
	}
	if !strings.Contains(strings.ToLower(user), "aide quand même") {
		t.Fatal("tutor should still help the learner when no material is found")
	}
}

func TestShortPersonasNonEmptyAndShorter(t *testing.T) {
	for _, code := range []string{"a1", "a2", "b1", "b2", "c1", "c2"} {
		l := LevelByName(code)
		if l.Short == "" {
			t.Fatalf("level %s missing Short persona", code)
		}
		if len(l.Short) >= len(l.System) {
			t.Errorf("level %s short persona should be shorter than full System", code)
		}
	}
}
