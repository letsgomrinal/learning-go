/* Learning Go: progress tracking and small reading aids for lesson pages.
 *
 * Progress lives in localStorage, so it stays in this browser only.
 *   go-roadmap-v1        { done: {c00: true, ...}, start }  shared with roadmap.html
 *   go-course-checks-v1  { c00: [true, false, ...] }         ticks in each lesson's checklist
 *   go-course-last       "c05"                               last lesson opened, for "continue"
 */
(function () {
  "use strict";

  var COURSE = window.COURSE || { total: 0, lessons: [] };
  var KEY = "go-roadmap-v1", CHECKS = "go-course-checks-v1", LAST = "go-course-last", THEME = "go-course-theme";
  var body = document.body;
  var lessonId = body.dataset.lesson || null;

  /* ── storage (never throws; private windows can block it) ── */
  function read(key, fallback) {
    try { var v = JSON.parse(localStorage.getItem(key) || "null"); return v == null ? fallback : v; }
    catch (e) { return fallback; }
  }
  function write(key, value) { try { localStorage.setItem(key, JSON.stringify(value)); } catch (e) {} }
  function roadmap() {
    var s = read(KEY, null);
    if (!s || typeof s !== "object") s = { done: { install: true }, start: null };
    if (!s.done || typeof s.done !== "object") s.done = {};
    return s;
  }
  function isDone(id) { return !!roadmap().done[id]; }
  function setDone(id, on) {
    var s = roadmap();
    if (on) s.done[id] = true; else delete s.done[id];
    write(KEY, s);
  }
  function checks() { var c = read(CHECKS, {}); return c && typeof c === "object" ? c : {}; }

  var $ = function (sel, root) { return (root || document).querySelector(sel); };
  var $$ = function (sel, root) { return Array.prototype.slice.call((root || document).querySelectorAll(sel)); };
  var byId = {};
  COURSE.lessons.forEach(function (l, i) { l.index = i; byId[l.id] = l; });

  function normPath(p) { return p.replace(/index\.html$/, "").replace(/README\.md$/, ""); }

  /* ── render progress everywhere on the page ── */
  function refresh() {
    var done = roadmap().done;
    var n = COURSE.lessons.filter(function (l) { return done[l.id]; }).length;
    var pct = COURSE.total ? Math.round((n / COURSE.total) * 100) : 0;

    $$("[data-count-done]").forEach(function (el) { el.textContent = n; });
    var ring = $(".ring-fg");
    if (ring) ring.style.strokeDasharray = pct + " 100";
    var pill = $(".progress-pill");
    if (pill) pill.setAttribute("aria-label", n + " of " + COURSE.total + " lessons complete");

    $$("[data-nav-lesson]").forEach(function (a) {
      var d = !!done[a.dataset.navLesson];
      a.toggleAttribute("data-done", d);
      var sr = $("[data-sr-status]", a);
      if (sr) sr.textContent = d ? " (completed)" : "";
    });
    $$("[data-part-count]").forEach(function (el) {
      var links = $$('[data-part="' + el.dataset.partCount + '"] [data-nav-lesson]');
      var k = links.filter(function (a) { return done[a.dataset.navLesson]; }).length;
      el.textContent = k + "/" + links.length;
    });

    $$("a[data-lesson-link]").forEach(function (a) { a.toggleAttribute("data-done", !!done[a.dataset.lessonLink]); });

    if (body.hasAttribute("data-home")) refreshHome(done, n, pct);
    if (lessonId) refreshLesson(!!done[lessonId]);
  }

  function nextLesson(done) {
    var last = byId[read(LAST, "")];
    if (last && !done[last.id]) return { lesson: last, resume: true };
    for (var i = 0; i < COURSE.lessons.length; i++) if (!done[COURSE.lessons[i].id]) return { lesson: COURSE.lessons[i], resume: false };
    return null;
  }

  function refreshHome(done, n, pct) {
    var head = $("[data-hp-headline]"), btn = $("[data-hp-continue]"), bar = $("[data-hp-bar]");
    if (bar) bar.style.width = pct + "%";
    if (!head || !btn) return;
    var nx = nextLesson(done);
    if (!nx) {
      head.textContent = "You finished all " + COURSE.total + " lessons. Now build the milestone projects.";
      btn.textContent = "Open the roadmap →";
      btn.href = $(".topbar-link").href;
      return;
    }
    var l = nx.lesson;
    btn.href = l.url;
    if (n === 0 && !nx.resume) {
      head.textContent = "Start with lesson 00. It takes about 45 minutes.";
      btn.textContent = "Start lesson 00 →";
    } else {
      head.textContent = (nx.resume ? "Pick up where you left off: " : "Up next: ") + "Lesson " + l.num + ", " + l.title + ".";
      btn.textContent = (nx.resume ? "Continue lesson " : "Start lesson ") + l.num + " →";
    }
  }

  /* ── home page: tag lesson links so they can show a ✓ ── */
  function tagLessonLinks() {
    var map = {};
    COURSE.lessons.forEach(function (l) { map[normPath(l.url)] = l.id; });
    $$(".prose a[href]").forEach(function (a) {
      var id = map[normPath(new URL(a.href, location.href).pathname)];
      if (id) a.dataset.lessonLink = id;
    });
  }

  /* ── lesson page ── */
  var completeBtn, finish, doneBadge;

  function refreshLesson(d) {
    if (!completeBtn) return;
    completeBtn.setAttribute("aria-pressed", String(d));
    completeBtn.textContent = d ? "✓ Completed · undo" : "Mark lesson complete";
    completeBtn.classList.toggle("btn-primary", !d);
    completeBtn.classList.toggle("btn-ghost", d);
    finish.toggleAttribute("data-done", d);
    var l = byId[lessonId];
    $(".finish-title").textContent = d ? "Lesson " + l.num + " complete. Nice work." : "Finished lesson " + l.num + "?";
    if (doneBadge) doneBadge.hidden = !d;
    var next = $(".pager-link.next");
    if (next) {
      next.toggleAttribute("data-up-next", d);
      var dir = $(".pager-dir", next);
      if (!dir.dataset.label) dir.dataset.label = dir.textContent;
      dir.textContent = d ? "Up next →" : dir.dataset.label;
    }
    refreshChecks();
  }

  function refreshChecks() {
    var boxes = $$(".prose li.task-list-item input");
    var sum = $("[data-check-summary]");
    if (!boxes.length || !sum) return;
    var k = boxes.filter(function (b) { return b.checked; }).length;
    sum.hidden = false;
    var all = k === boxes.length;
    sum.textContent = all ? "Checklist: all " + k + " ticked ✓" : "Checklist: " + k + " of " + boxes.length + " ticked";
    if (completeBtn) completeBtn.toggleAttribute("data-nudge", all && !isDone(lessonId));
  }

  function setupIntro(prose) {
    var h1 = $("h1", prose);
    if (!h1) return;
    var ps = [], el = h1.nextElementSibling;
    while (el && el.tagName === "P" && /^(Goal|Time)\b/.test((el.textContent || "").trim())) { ps.push(el); el = el.nextElementSibling; }
    if (!ps.length) return;
    var card = document.createElement("div");
    card.className = "intro-card";
    h1.after(card);
    ps.forEach(function (p) { card.appendChild(p); });
    var ex = $$("h2", prose).filter(function (h) { return /Exercises/.test(h.textContent); })[0];
    var cl = $$("h2", prose).filter(function (h) { return /Checklist/.test(h.textContent); })[0];
    var acts = document.createElement("p");
    acts.className = "intro-actions";
    if (ex && ex.id) acts.insertAdjacentHTML("beforeend", '<a href="#' + ex.id + '">Jump to the exercises ↓</a>');
    if (cl && cl.id) acts.insertAdjacentHTML("beforeend", '<a href="#' + cl.id + '">See the checklist ↓</a>');
    if (acts.children.length) card.appendChild(acts);
  }

  function setupChecklist(prose) {
    var items = $$("li.task-list-item", prose);
    if (!items.length) return;
    var saved = checks()[lessonId] || [];
    items.forEach(function (li, i) {
      var input = $("input", li);
      if (!input) return;
      input.disabled = false;
      input.removeAttribute("disabled");
      input.checked = !!saved[i];
      input.id = "check-" + i;
      var label = document.createElement("label");
      label.htmlFor = input.id;
      var span = document.createElement("span");
      while (input.nextSibling) span.appendChild(input.nextSibling);
      label.appendChild(input);
      label.appendChild(span);
      li.appendChild(label);
      input.addEventListener("change", function () {
        var all = checks();
        all[lessonId] = items.map(function (x) { var b = $("input", x); return !!(b && b.checked); });
        write(CHECKS, all);
        refreshChecks();
      });
    });
  }

  function setupToc(prose) {
    var list = $(".toc-list");
    if (!list) return;
    var heads = $$("h2[id]", prose);
    if (heads.length < 3) { $(".toc").hidden = true; return; }
    heads.forEach(function (h) {
      var li = document.createElement("li");
      var a = document.createElement("a");
      a.href = "#" + h.id;
      a.innerHTML = h.innerHTML.replace(/^\s*\d+\.\s*/, "");
      li.appendChild(a);
      list.appendChild(li);
    });
    var links = $$("a", list);
    if (!("IntersectionObserver" in window)) return;
    var io = new IntersectionObserver(function () {
      var current = null;
      for (var i = 0; i < heads.length; i++) {
        if (heads[i].getBoundingClientRect().top < window.innerHeight * 0.35) current = i;
      }
      links.forEach(function (a, i) { a.setAttribute("aria-current", i === current ? "true" : "false"); });
    }, { rootMargin: "0px 0px -60% 0px" });
    heads.forEach(function (h) { io.observe(h); });
  }

  /* kramdown can leave <summary> wrapped in a <p>; a summary must be the first child to work. */
  function fixDetails(prose) {
    $$("details", prose).forEach(function (d) {
      var s = $("summary", d);
      if (s && s.parentElement !== d) {
        var p = s.parentElement;
        d.insertBefore(s, d.firstChild);
        if (p && !p.textContent.trim() && !p.children.length) p.remove();
      }
    });
  }

  function setupCode(prose) {
    $$("div.highlighter-rouge", prose).forEach(function (box) {
      var pre = $("pre", box);
      if (!pre) return;
      var m = /language-(\w+)/.exec(box.className);
      var lang = m && m[1] !== "plaintext" ? m[1] : "";
      if (lang) {
        var tag = document.createElement("span");
        tag.className = "code-lang";
        tag.textContent = { bash: "shell", sh: "shell", go: "Go", python: "Python" }[lang] || lang;
        box.appendChild(tag);
      }
      var btn = document.createElement("button");
      btn.type = "button";
      btn.className = "copy-btn";
      btn.textContent = "Copy";
      btn.addEventListener("click", function () {
        var text = pre.innerText.replace(/\n$/, "");
        var ok = function () {
          btn.textContent = "Copied"; btn.setAttribute("data-copied", "");
          setTimeout(function () { btn.textContent = "Copy"; btn.removeAttribute("data-copied"); }, 1400);
        };
        if (navigator.clipboard) navigator.clipboard.writeText(text).then(ok, function () {});
      });
      box.appendChild(btn);
    });
  }

  function setupReadProgress() {
    var bar = $(".read-progress > i");
    if (!bar) return;
    var ticking = false;
    function update() {
      var h = document.documentElement.scrollHeight - window.innerHeight;
      bar.style.width = (h > 0 ? Math.min(100, (window.scrollY / h) * 100) : 0) + "%";
      ticking = false;
    }
    window.addEventListener("scroll", function () { if (!ticking) { ticking = true; requestAnimationFrame(update); } }, { passive: true });
    update();
  }

  function setupLesson() {
    var prose = $(".prose");
    write(LAST, lessonId);
    /* The README's own "Next →" line is for GitHub readers; the pager below replaces it here. */
    $$("p", prose).slice(-2).forEach(function (p) {
      if (/^\s*(Next|You've finished)/.test(p.textContent) && $("a", p)) p.hidden = true;
    });
    setupIntro(prose);
    setupChecklist(prose);
    setupToc(prose);
    setupReadProgress();
    completeBtn = $(".complete-btn");
    finish = $(".finish");
    doneBadge = $(".done-badge");
    if (completeBtn) completeBtn.addEventListener("click", function () {
      var on = !isDone(lessonId);
      setDone(lessonId, on);
      refresh();
      if (on) { var nx = $(".pager-link.next"); if (nx) nx.focus({ preventScroll: true }); }
    });
  }

  /* ── chrome: mobile menu, theme, reset ── */
  function setupChrome() {
    var menu = $(".menu-btn"), scrim = $(".scrim");
    function setNav(open) {
      body.toggleAttribute("data-nav-open", open);
      if (menu) menu.setAttribute("aria-expanded", String(open));
      if (scrim) scrim.hidden = !open;
    }
    if (menu) menu.addEventListener("click", function () { setNav(!body.hasAttribute("data-nav-open")); });
    if (scrim) scrim.addEventListener("click", function () { setNav(false); });
    document.addEventListener("keydown", function (e) { if (e.key === "Escape") setNav(false); });

    var cur = $('[aria-current="page"]', $(".sidebar"));
    if (cur) cur.scrollIntoView({ block: "center" });

    var themeBtn = $(".theme-btn");
    if (themeBtn) themeBtn.addEventListener("click", function () {
      var root = document.documentElement;
      var dark = root.dataset.theme ? root.dataset.theme === "dark" : matchMedia("(prefers-color-scheme: dark)").matches;
      root.dataset.theme = dark ? "light" : "dark";
      try { localStorage.setItem(THEME, root.dataset.theme); } catch (e) {}
    });

    var reset = $("[data-reset]");
    if (reset) reset.addEventListener("click", function () {
      if (!reset.dataset.armed) {
        reset.dataset.armed = "1";
        reset.textContent = "Click again to clear all progress";
        setTimeout(function () { if (reset.dataset.armed) { delete reset.dataset.armed; reset.textContent = "Reset my progress"; } }, 3000);
        return;
      }
      var s = roadmap();
      COURSE.lessons.forEach(function (l) { delete s.done[l.id]; });
      write(KEY, s);
      write(CHECKS, {});
      try { localStorage.removeItem(LAST); } catch (e) {}
      delete reset.dataset.armed;
      reset.textContent = "Progress cleared";
      $$(".prose li.task-list-item input").forEach(function (b) { b.checked = false; });
      refresh();
    });
  }

  setupChrome();
  var prose = $(".prose");
  if (prose) { fixDetails(prose); setupCode(prose); }
  if (body.hasAttribute("data-home")) tagLessonLinks();
  if (lessonId) setupLesson();
  refresh();

  /* Another tab (or the roadmap) changed progress: stay in sync. */
  window.addEventListener("storage", function (e) {
    if (e.key === KEY || e.key === CHECKS) {
      if (e.key === CHECKS && lessonId) {
        var saved = checks()[lessonId] || [];
        $$(".prose li.task-list-item input").forEach(function (b, i) { b.checked = !!saved[i]; });
      }
      refresh();
    }
  });
})();
