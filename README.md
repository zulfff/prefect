# <p align="center">Prefect</p>

<p align="center">
    <img src="https://raw.githubusercontent.com/daniswastaken/prefect/main/ui/media/banner.png" alt="Prefect Banner" width="100%">
    <br/>
    <i>"It's perfect but it's also not perfect at the same time, so it's Prefect."</i>
</p>

---

**Prefect** is a modern, lightweight, and visually stunning dashboard and file management system designed for homelabs. Built with a high-performance Go backend and a sleek, glassmorphic Web UI, it provides real-time system insights and seamless file operations in a unified experience.

## ✨ Features

- 🖥️ **Live System Stats**: Real-time monitoring of CPU usage (with temperature and power tracking), RAM availability, and Disk health via high-speed WebSockets.
- 📂 **Integrated File Explorer**: A fully-featured web-based explorer with support for copy, cut, paste, rename, and delete operations.
- 🦄 **Glassmorphic UI**: A premium, state-of-the-art design inspired by modern operating systems, featuring vibrant gradients, blur effects, and smooth animations.
- ☁️ **SUDO-Aware Home Detection**: Intelligently detects the actual user's home directory even when running with root privileges, ensuring a personalized experience.
- 🚀 **SPA Architecture**: Navigates between system insights and applications like the File Explorer without page reloads, using a sophisticated iframe overlay system.

## 🛠️ Technology Stack

- **Backend**: [Go](https://go.dev/) (Golang) — Leveraging its concurrency and system-level performance.
- **Frontend**: 
  - Pure **HTML5** & **Vanilla JavaScript** for high performance and low overhead.
  - **Vanilla CSS** for the exquisite glassmorphic design system.
- **Communication**: **WebSockets** for real-time, bidirectional data streaming.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher installed on your system.

### Running Prefect

1. **Clone the repository**:
   ```bash
   git clone https://github.com/daniswastaken/prefect.git
   cd prefect
   ```

2. **Run the application**:
   ```bash
   go run main.go
   # Or with sudo if system-level access is required
   sudo go run main.go
   ```

3. **Open the Dashboard**:
   Navigate to `http://localhost:8080` in your web browser.

## 📜 License

This project is licensed under the [MIT License](LICENSE).

## 🤝 Credits

Created with ❤️ for the homelab community. Reference [CREDITS.md](CREDITS.md) for a full list of contributors and inspirations.
