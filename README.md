# quay-yala

**quay-yala** is a simple Go-based CLI tool designed to analyze Red Hat Quay Debug logs. It automates the tedious process of digging through massive log files by extracting critical error tracebacks and access log patterns.

## 🚀 Features

* **Log Parsing:** Automatically detects and extracts error tracebacks and structured access logs.
* **Batch Processing:** Analyze multiple log files in a single command.
* **Automatic Reporting:** Generates individual report files for every input log analyzed.
* **Cross-Platform:** Available for Linux, macOS, and Windows (Intel & ARM).

---

## 📦 Installation

You can download the latest binaries from the [Releases](https://github.com/nasonawa/quay-yala/releases) page.

### From Source

If you have Go installed:

```bash
go install github.com/nasonawa/quay-yala@latest

```

---

## 🛠 Usage

To analyze one or more Quay debug logs, use the `-i` (input) flag:

```bash
quay-yala -i quay-log1.log quay-log2.log

```

### Output

For every input file provided, the tool generates a corresponding report in the current directory:

* `quay-log1.log` ➡️ `report-quay-log1.log`
* `quay-log2.log` ➡️ `report-quay-log2.log`

---

## 📊 Report Details

The generated reports provide a structured view of:

1. **Error Tracebacks:** Grouped exceptions found within the debug stream.
2. **Access Logs:** Parsed request metadata (Method, Path, Status Codes) for quick auditing.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
