package prompt

import (
	"fmt"
	"strings"

	"github.com/Sreenivas-Sadhu-Prabhakara/les-petits-horizons/internal/retriever"
)

// Assemble builds (systemPrompt, userPrompt) for the default (A1) level.
func Assemble(message string, chunks []retriever.Chunk) (string, string) {
	return AssembleFor(Default, message, chunks)
}

// AssembleFor builds (systemPrompt, userPrompt) using the given level's full persona.
func AssembleFor(l Level, message string, chunks []retriever.Chunk) (string, string) {
	return l.System, buildUserPrompt(message, chunks)
}

// AssembleForShort is AssembleFor with the level's distilled short persona — used
// when serving the tuned (prompt-distilled) model. The user prompt is identical.
func AssembleForShort(l Level, message string, chunks []retriever.Chunk) (string, string) {
	return l.Short, buildUserPrompt(message, chunks)
}

// buildUserPrompt renders the MATÉRIEL (sources) + learner message + grounding
// instruction. Grounding is advisory: the tutor uses lesson material when
// relevant but still helps the learner practice French when none is found.
func buildUserPrompt(message string, chunks []retriever.Chunk) string {
	var b strings.Builder
	if len(chunks) == 0 {
		b.WriteString("MATÉRIEL : (aucun matériel pédagogique trouvé)\n\n")
	} else {
		b.WriteString("MATÉRIEL :\n")
		for i, c := range chunks {
			fmt.Fprintf(&b, "[%d] %s : %s\n", i+1, c.Title, c.Content)
		}
		b.WriteString("\n")
	}
	fmt.Fprintf(&b, "MESSAGE DE L'APPRENANT : %s\n", message)
	b.WriteString("\nRéponds en français adapté au niveau. Utilise le MATÉRIEL ci-dessus " +
		"quand c'est pertinent ; sinon, aide quand même l'apprenant. Corrige ses erreurs " +
		"avec bienveillance et termine par une question pour le faire continuer.")
	return b.String()
}
