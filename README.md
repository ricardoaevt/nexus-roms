# Nexus ROMs

![Status](https://img.shields.io/badge/Status-Development-orange)
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![Wails](https://img.shields.io/badge/Wails-v2-red)
![Svelte](https://img.shields.io/badge/Frontend-Svelte-ff3e00?logo=svelte)

Nexus ROMs is a high-performance desktop application built with **Wails** and **Svelte**, designed to organize, identify, and rename your retro gaming collection with precision and ease.

## 🌟 Key Features

- **🚀 Concurrent Orchestration**: Multi-threaded scanning engine (configurable worker counts) for blazing fast identification.
- **🔍 Multi-Provider Scraping**: Integrated support for world-class metadata providers:
  - [ScreenScraper](https://www.screenscraper.fr/)
  - [TheGamesDB](https://thegamesdb.net/)
- **🧩 Smart Templates**: Rename your collection using dynamic tokens like `{Name}`, `{Region}`, `{Platform}`, `{Year}`, and `{Developer}`.
- **💾 Session Recovery**: Robust session management allows you to pause, resume, or restart scanning operations without losing progress.
- **💎 Premium UI**: A reactive, dark-mode dashboard providing real-time logs, live identification tables, and progress statistics.
- **🔒 Security First**: API credentials and personal data are securely stored and handled.
- **📦 Wide Compression Support**: Native handling of `.zip`, `.rar`, and `.7z` archives.

## 🚀 Getting Started

### Prerequisites

To compile Nexus ROMs from source, you will need:
- **Go** (1.24 or higher)
- **Node.js** (v20+) & **NPM**
- **Wails CLI** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Installation & Development

1. Clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/nexus-roms.git
   cd nexus-roms
   ```

2. Run in development mode (Live Reload):
   ```bash
   wails dev
   ```

### 🏗️ Building for Production

Generate a production-ready binary for your current OS:
```bash
wails build
```

The compiled binary will be located in the `build/bin` directory.

## 📦 Automated Releases

This repository includes a **GitHub Actions** workflow that automatically compiles and packages the application for **Windows, Linux, and macOS** whenever a new tag (e.g., `v1.0.0`) is pushed.

Detailed instructions for downloading pre-compiled binaries can be found in the [Releases](https://github.com/YOUR_USERNAME/nexus-roms/releases) section.

## 🛠️ Tech Stack

- **Backend**: Go (Wails Bridged)
- **Frontend**: Svelte / TypeScript
- **Database**: SQLite (CGO-free via `modernc.org/sqlite`)
- **Cryptography**: AES-256-GCM for credential protection

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is personal software. See the [LICENSE](LICENSE) file for details.

---
*Developed with focus on speed, aesthetics, and user experience.*
