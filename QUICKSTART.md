# Nessi CLI Quickstart Guide

Welcome to Nessi! This guide will help you get started with Nessi for Delta Lake data quality management, Python extensions, telemetry, and more.

---

## Installation

1. **Clone the repo:**
   ```sh
   git clone https://github.com/nessi-dev/nessi.git
   cd nessi
   ```
2. **Build the CLI:**
   ```sh
   go build -o nessi ./cmd/nessi
   ```
3. **(Optional) Install Python 3.8+** for advanced features:
   ```sh
   python3 --version
   # If needed: brew install python@3.8
   ```

---

## Enabling Python Extensions

Enable advanced analytics and ML features:

```sh
./nessi python-extensions enable
```

Check status:
```sh
./nessi python-extensions status
```

---

## Validating a Delta Lake Table

```sh
./nessi validate --table /path/to/table
```

## Profiling a Table

```sh
./nessi profile --table /path/to/table
```

## Monitoring Table Metrics

```sh
./nessi monitor --table /path/to/table
```

## Generating a Quality Report

```sh
./nessi report --table /path/to/table
```

---

## Telemetry

- By default, telemetry is **disabled**.
- To enable:
  ```sh
  ./nessi telemetry enable
  ```
- To view what is collected:
  ```sh
  ./nessi telemetry info
  ```
- To disable:
  ```sh
  ./nessi telemetry disable
  ```

---

## Configuration Validation

Check your config for errors:
```sh
./nessi config validate
```

---

## Need More Features?

- See `nessi info` and `nessi team-features` for enterprise options.
- Explore custom extensions: see `docs/python_extensions.md` and `docs/plugin_system.md`

---

## Troubleshooting

- Ensure Python is installed and in your PATH for Python extensions
- For help on any command:
  ```sh
  ./nessi <command> --help
  ```
- For more, see the full user guide: `docs/USER_GUIDE.md`
