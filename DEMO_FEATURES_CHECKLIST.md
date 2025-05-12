# Conference Demo: Features Checklist & Roadmap

**Purpose:**
Every time development resumes, check this file and `IMPLEMENTATION.md` to decide where to continue.

---

## 1. Schema Management & Version Control
- [x] Schema evolution/history: Track schema changes, expose diffs
- [x] Schema validation on write: Enforce/log errors
- [x] Field-level metadata: Comments, tags, types
- [x] Version history & rollback: Commit logs, version compare, rollback (30 days)
- [x] Time travel: CLI/API for version/timestamp query

## 2. Automated Profiling & Data Quality
- [x] Statistical summaries: min, max, mean, median (Go/Python)
- [x] Distribution analysis: HTML (Python), stats (Go), add histogram if possible
- [x] Null %, unique ratio, type inference: Go
- [ ] Consistency checks: Compare schema/types across versions
- [x] Pattern recognition: Regex, email, URL, digit-only (Go)
- [ ] Format validation: Add more (dates, IDs, custom regex)
- [ ] Trend deviation: Compare current/previous runs, alert on changes

## 3. Anomaly Detection
- [x] Outlier detection (z-score): Go profiler
- [ ] IQR-based outliers: Add IQR method
- [ ] Sudden change detection: Compare stats across runs

## 4. Rule Validation
- [x] Null/range checks: Go rules
- [ ] Predefined rule library: Add more rules (length, regex, enum, etc)
- [ ] Custom rule authoring: YAML/SQL/Python-based rules
- [ ] Rule execution history: Store/display validation results over time

## 5. Monitoring & Alerts
- [ ] Real-time dashboard: Web UI or Grafana
- [ ] Metric retention: Store metrics for 30 days
- [ ] Alerting: Email/CLI notifications

## 6. Reports & Visualization
- [x] HTML/JSON reports: Python profiler
- [ ] Trend/failure charts: Add to HTML or Grafana
- [ ] Export options: JSON/CSV export from Go/Python

## 7. Platform, Security, Developer Experience
- [x] Dockerized deployment
- [ ] Basic security: SSL, API keys, user auth
- [ ] CLI tools: One-line validation/profiling
- [ ] Python API: Expose programmatic access
- [ ] Docs: Interactive examples, best practices

---

**How to use:**
- Check off completed features.
- Add notes/links as you go.
- Use this with `IMPLEMENTATION.md` to track progress and next steps.
