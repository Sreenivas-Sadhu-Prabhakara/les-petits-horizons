// Backend URL for the live Chat tab.
//
// The tutor runs self-hosted on a Mac Studio behind a Cloudflare quick tunnel.
// That URL is EPHEMERAL — it changes every time the tunnel restarts. When it
// changes, either update the value below and redeploy, or just paste the fresh
// URL into the "Backend URL" field on the Chat tab (it's saved in your browser).
window.LPH_CONFIG = {
  // Current public tunnel for the tutor backend:
  BACKEND_URL: "https://revision-precipitation-system-informational.trycloudflare.com",
};
