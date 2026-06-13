# Les Petits Horizons — French Tutor Chatbot · Design Spec

**Date:** 2026-06-13
**Owner:** Sreenivas-Sadhu-Prabhakara
**Status:** Approved (brainstorming) — Phase 1 implementation

## 1. Summary

A self-hosted **French-language tutor** chatbot, cloned from the architecture of the
existing `llm-slm` ("Apolaki Solar Assistant") Go service but stripped of all
solar-domain logic. It speaks **immersive French, graded to the learner's CEFR level
(A1–C2)**, does conversation practice, corrects grammar, and teaches
vocabulary/translation — with answers steered by retrieved lesson material (RAG).

Audience: **everyone learning French**, beginner to advanced.

## 2. Reused vs. new infrastructure

| Concern | Decision |
|---|---|
| Generation model | `qwen-coder-32b` via LiteLLM `:4000` (shared, strong multilingual/French) |
| Embeddings | `bge-m3` via LiteLLM (shared) — `vector(1024)` |
| Vector DB | **same pgvector server `:5433`**, new **`french`** database (isolated from live `solar`) |
| Service port | `:8092` (solar bot on `:8090`/`:8091` untouched) |
| Public access | **new ephemeral ngrok tunnel → `:8092`** |
| Repo | new public GitHub repo `les-petits-horizons` |
| Website | GitHub Pages from `/docs`, public |

## 3. Architecture (mirrors the proven solar-bot pipeline)

```
client → POST /tutor/chat (SSE streaming)  [CORS-enabled]
  → (light, permissive French-learning gate)
  → retrieve lesson material (pgvector HNSW, audience="french")   internal/retriever + internal/embed (bge-m3)
  → assemble graded French-tutor persona + sources                internal/prompt + internal/prompt/levels
  → generate (LiteLLM → qwen-coder-32b, streaming)                internal/generator
  → log turn + 👍/👎 feedback                                      internal/chatlog
  → browser test page at /
```

### Key adaptations from the solar bot
- **Persona:** immersive French tutor, **graded by CEFR level** passed per request
  (`level`: a1|a2|b1|b2|c1|c2; default a1). At A1 the French is very simple with
  minimal English scaffolding only when the learner is stuck; difficulty rises with level.
- **Features in the persona:** conversation practice, gentle grammar correction with
  short explanations, vocabulary & FR↔EN translation, level-appropriate complexity.
- **Topic gate:** replaced solar keyword gate with a *permissive* gate — almost any
  learner utterance is valid French practice; only clearly unrelated task requests
  (e.g. "write me code") are softly redirected back to French practice.
- **Safety/escalation removed:** the solar bot escalated to a human when no source was
  found. A language tutor must not refuse un-sourced French (it would refuse almost
  everything). **Grounding is advisory, not gating:** retrieved lesson material steers
  content when relevant; otherwise the tutor still teaches from the model's own French.
  This is the practical reading of "pure RAG" for a tutor and is revisited once real
  curriculum material lands in Phase 2.
- **CORS:** permissive CORS + OPTIONS preflight so the GitHub Pages site can call the
  ngrok endpoint directly (the solar bot was same-origin only).

## 4. Data model
Identical schema to the solar bot (`knowledge_documents`, `knowledge_chunks` with
`vector(1024)` HNSW cosine index, `conversations`, `messages`, `feedback`). All seed
documents use `audience = "french"`, `language = "fr"`.

## 5. Phase 1 (today) vs. Phase 2 (tomorrow)

**Phase 1 — today (this build):**
1. Clone + adapt the Go service into `les-petits-horizons`.
2. Generate a small **synthetic French seed corpus** (graded grammar notes, vocab,
   example dialogues) so RAG works end-to-end and there is a live URL.
3. Provision the `french` DB, migrate, ingest.
4. Run server on `:8092` + ngrok tunnel; smoke test through the tunnel.
5. Build the GitHub Pages website (How-it-was-built · ELI5 · live Chat tab).
6. Create the public GitHub repo, push, enable Pages.

**Phase 2 — tomorrow (after real materials land):**
1. Replace the seed corpus with real French learning materials → re-ingest.
2. Run the LoRA training pipeline (carried over in `training/`) on that material.
3. Swap the tuned model in as the LiteLLM primary; fresh URL.

## 6. The website (GitHub Pages, public)
Static site in `/docs`:
- **How it was built** — architecture diagram, RAG pipeline, model/infra choices, training plan.
- **If I were 5 (ELI5)** — the same system explained in plain, playful language.
- **Chat** — a live chat widget that streams from `POST /tutor/chat` (SSE) at the
  ngrok URL. The backend URL is read from `docs/config.js` and overridable via a
  field saved in `localStorage` (the free-tier ngrok URL rotates on restart).

Consequence: the Chat tab only works when the Mac + tunnel are up; the other tabs always work.

## 7. Out of scope (Phase 1)
- Real curriculum content (lands Phase 2).
- LoRA fine-tuning (Phase 2).
- Persistent/reserved ngrok domain (free ephemeral tunnel for now).
- Auth beyond the inherited permissive dev stub.
