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
- **Wails CLI**: Install via `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **C Compiler**: Required for SQLite (e.g., GCC on Linux, MingW on Windows).

### 🔑 API Configuration

To use the identification features, you will need accounts from these providers:

1. **ScreenScraper**: Register at [screenscraper.fr](https://www.screenscraper.fr/).
2. **TheGamesDB**: Get an API key at [thegamesdb.net](https://thegamesdb.net/).

Once you have your credentials, enter them in the **Settings** tab within the application.

### 🏗️ Installation & Development

1. Clone the repository:

   ```bash
   git clone https://github.com/ricardoaevt/nexus-roms.git
   cd nexus-roms
   ```

2. Install dependencies:

   ```bash
   wails doctor
   ```

3. Run in development mode (Live Reload):

   ```bash
   wails dev
   ```

### 🔨 Building for Production

Generate a production-ready binary for your current OS:

```bash
wails build
```

The compiled binary will be located in the `build/bin` directory.

## 📦 Automated Releases

This repository is configured with **GitHub Actions** to automatically build and package the application for **Windows, Linux, and macOS**. 

**To trigger a new release:**

1. Update the version in `wails.json`.
2. Create and push a tag:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

The binaries will appear in the [Releases](https://github.com/ricardoaevt/nexus-roms/releases) section once the workflow finishes.

## 🛠️ Tech Stack

- **Backend**: Go 1.24+
- **Frontend**: Svelte
- **Database**: SQLite (Encrypted)
- **Framework**: [Wails v2](https://wails.io/)

## 🤝 Contributing

1. Fork the Project.
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`).
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the Branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.

---
*Developed with focus on speed, aesthetics, and user experience for the Retro Gaming Community.*
