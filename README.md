# Les Petits Horizons 🇫🇷

A self-hosted **French-language tutor** chatbot. It speaks **immersive French, graded to your
CEFR level (A1–C2)**, does conversation practice, corrects grammar, and teaches
vocabulary/translation — with answers steered by retrieved lesson material (RAG). Runs entirely
on a Mac Studio at **≈ $0 per query**.

It reuses the architecture of an existing self-hosted RAG assistant, re-pointed at a new domain:
teaching French.

- **Website (how it was built · ELI5 · live chat):** https://Sreenivas-Sadhu-Prabhakara.github.io/les-petits-horizons/
- **Tutor backend:** self-hosted on `:8092`, exposed via a Cloudflare quick tunnel (ephemeral URL — see `docs/config.js`).

## Architecture

```
client → POST /tutor/chat (SSE, CORS-enabled)
  → retrieve lesson material (pgvector HNSW, audience="french")   internal/retriever + internal/embed (bge-m3)
  → assemble graded French-tutor persona (A1–C2) + sources        internal/prompt
  → generate (LiteLLM → qwen-coder-32b, streaming)                internal/generator
  → log turn + 👍/👎 feedback                                      internal/chatlog
  → browser test page at /
```

- **Generation:** Qwen 32B (`qwen-coder-32b`) via LiteLLM `:4000`.
- **Embeddings:** BGE-M3 (`bge-m3`), 1024-dim, via LiteLLM.
- **Vector DB:** shared pgvector server on `:5433`, isolated `french` database.
- **One persona, six grades:** a single tutor (“Petit Horizon”) graded by CEFR level — simple
  French at A1, near-native at C2. Lesson material steers content; the model supplies the language.

## Run locally

```bash
cp .env.example .env                 # defaults reuse pgvector :5433 + LiteLLM :4000
make db-create                       # create the `french` database in the shared pgvector container
make migrate                         # apply DB migrations
make seed && make ingest             # generate + ingest the synthetic French seed corpus
make run                             # serve the tutor + chat UI at http://localhost:8092/
make test                            # full Go test suite

# CLI smoke test:
go run ./cmd/ask -level a1 "Comment dire 'I am hungry' en français ?"
```

### Expose publicly (own URL, separate from anything else on the box)

```bash
cloudflared tunnel --url http://localhost:8092
# copy the printed https://*.trycloudflare.com URL into docs/config.js (or paste it
# into the "Backend URL" field on the website's Chat tab — it's saved in localStorage)
```

## Layout

```
cmd/         ask · ingest · server
internal/    config db embed ingest retriever prompt generator chatlog httpapi personalizer
data/        synthetic French seed corpus generator
docs/        the GitHub Pages site (How it was built · ELI5 · live Chat) + the design spec
migrations/  pgvector schema, knowledge tables, conversations/feedback
```

## Phases

- **Phase 1 (this build):** Go service + pgvector/HNSW RAG + BGE-M3 + Qwen via LiteLLM, a
  synthetic graded French seed corpus, the public tunnel, and the showcase website.
- **Phase 2 (next):** replace the seed with real French learning materials → re-ingest →
  **LoRA fine-tune** on that material → serve the tuned model as the LiteLLM primary.

See `docs/superpowers/specs/` for the full design spec.
