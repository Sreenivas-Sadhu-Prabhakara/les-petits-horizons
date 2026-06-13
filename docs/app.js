// Tab switching
const tabs = document.querySelectorAll(".tab");
const panels = document.querySelectorAll(".panel");
tabs.forEach((t) => {
  t.addEventListener("click", () => {
    tabs.forEach((x) => x.classList.remove("active"));
    panels.forEach((p) => p.classList.remove("active"));
    t.classList.add("active");
    document.getElementById(t.dataset.tab).classList.add("active");
    if (t.dataset.tab === "chat") checkHealth();
  });
});

// ---- Backend URL: localStorage override > config.js default ----
const cfgDefault = (window.LPH_CONFIG && window.LPH_CONFIG.BACKEND_URL) || "";
function backendURL() {
  return (localStorage.getItem("lph_backend") || cfgDefault).replace(/\/+$/, "");
}
const backendInput = document.getElementById("backend");
backendInput.value = backendURL();
document.getElementById("saveBackend").addEventListener("click", () => {
  localStorage.setItem("lph_backend", backendInput.value.trim().replace(/\/+$/, ""));
  checkHealth();
});

// ---- Health indicator ----
const dot = document.getElementById("statusDot");
const statusText = document.getElementById("statusText");
async function checkHealth() {
  const url = backendURL();
  if (!url) { setStatus("bad", "no backend URL set"); return; }
  setStatus("", "checking…");
  try {
    const r = await fetch(url + "/tutor/health", { method: "GET" });
    const j = await r.json();
    if (j.status === "ok") setStatus("ok", "tutor online");
    else setStatus("bad", "backend: " + j.status);
  } catch (e) {
    setStatus("bad", "offline (the Mac or tunnel is down)");
  }
}
function setStatus(cls, text) {
  dot.className = "dot " + cls;
  statusText.textContent = text;
}

// ---- Chat ----
let conversationId = "";
const log = document.getElementById("log");
const input = document.getElementById("msg");
const send = document.getElementById("send");
const level = document.getElementById("level");

function add(cls, text) {
  const d = document.createElement("div");
  d.className = "msg " + cls;
  d.textContent = text || "";
  log.appendChild(d);
  log.scrollTop = log.scrollHeight;
  return d;
}

async function ask() {
  const message = input.value.trim();
  const url = backendURL();
  if (!message) return;
  if (!url) { add("bot", "⚠️ Set the Backend URL first."); return; }
  input.value = "";
  send.disabled = true;
  add("user", message);
  const bot = add("bot", "");
  const textNode = document.createTextNode("");
  bot.appendChild(textNode);
  try {
    const res = await fetch(url + "/tutor/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message, level: level.value, conversation_id: conversationId }),
    });
    if (!res.ok || !res.body) throw new Error("HTTP " + res.status);
    const reader = res.body.getReader();
    const dec = new TextDecoder();
    let buf = "";
    while (true) {
      const { value, done } = await reader.read();
      if (done) break;
      buf += dec.decode(value, { stream: true });
      const frames = buf.split("\n\n");
      buf = frames.pop();
      for (const f of frames) {
        let event = "", data = "";
        for (const line of f.split("\n")) {
          if (line.startsWith("event:")) event = line.slice(6).trim();
          if (line.startsWith("data:")) data += line.slice(5).trim();
        }
        if (!data) continue;
        let obj;
        try { obj = JSON.parse(data); } catch { continue; }
        if (event === "done") {
          conversationId = obj.conversation_id || conversationId;
          if (obj.sources && obj.sources.length) {
            const s = document.createElement("div");
            s.className = "sources";
            s.textContent = "📚 Matériel : " + obj.sources.join(" · ");
            bot.appendChild(s);
          }
        } else if (obj.token) {
          textNode.textContent += obj.token;
          log.scrollTop = log.scrollHeight;
        } else if (obj.error) {
          textNode.textContent = "⚠️ " + obj.error;
        }
      }
    }
  } catch (e) {
    textNode.textContent = "⚠️ Could not reach the tutor (" + e.message + "). It only runs while the Mac + tunnel are up.";
    setStatus("bad", "offline");
  } finally {
    send.disabled = false;
    input.focus();
  }
}

send.addEventListener("click", ask);
input.addEventListener("keydown", (e) => { if (e.key === "Enter") ask(); });
