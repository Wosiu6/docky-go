
# docky-go: Modern CLI Docker Container Visualization

[![License][license-shield]][license-url] [![GitHub][github-shield]][github-url]

> A modular, extensible, and beautiful TUI for visualizing Docker containers and their stats, written in Go.

---

## Features

- **Live TUI Dashboard:** See all your Docker containers in a modern, responsive terminal UI.
- **Container Stats:** View CPU, memory, and status for all containers, with special details for supported types.
- **Extensible Architecture:** Add new container types or UI workflows by simply implementing a strategy interface—no core changes needed.
- **Cross-Platform:** Works on Windows (Docker Desktop) and Linux (Docker socket).
- **Fast & Efficient:** Uses Go concurrency for fast stats collection.
- **Beautiful UI:** Built with Bubble Tea and Lip Gloss for a polished look.

---

## Supported Container Types (with custom details)

<details>
  <summary>List of supported containers</summary>
  - PostgreSQL<br />
  - Minecraft<br />
  - Portainer<br />
  - Traefik<br />
  - Immich<br />
  - OwnCloud<br />
  - Nginx<br />
  - Redis<br />
  - MySQL<br />
  - MongoDB<br />
  - Grafana<br />
  - Prometheus<br />
  - Nextcloud<br />
  - Minio<br />
  - MariaDB<br />
  - RabbitMQ<br />
  - Elasticsearch<br />
  - Kibana<br />
  - Jenkins<br />
  - WordPress<br />
  - Vaultwarden<br />
  - Mosquitto<br />
  - Plex<br />
  - Jellyfin<br />
  - Home Assistant<br />
  - Sonarr<br />
  - Radarr
</details>

Want to add your own? Just implement a strategy and a detail renderer—no need to touch the core!

---

## Architecture

- **docker/**: Docker client abstraction (interface + implementation)
- **fetcher/**: Fetches and classifies containers, uses strategy pattern for extensibility
- **fetcher/strategies/**: One file per container type, easy to add more
- **ui/**: Modular Bubble Tea TUI (model, view, renderers, styles, logo, helpers)

---

## Screenshots

<details>
  <summary>Loading</summary>
  <img width="1112" height="629" alt="image" src="https://github.com/user-attachments/assets/36a951c1-964d-44f5-bd38-14e78f3c26b8" />

</details>

<details>
  <summary>3 containers, postgres, nginx, owncloud</summary>
  <img width="1112" height="623" alt="image" src="https://github.com/user-attachments/assets/7c48e732-b332-403a-8b32-9676cde9981a" />

</details>

---

## Roadmap

- [x] Modularize and deduplicate codebase
- [x] Add unit tests for strategies and fetcher
- [ ] Add more container-specific strategies (PRs welcome soon!)
- [ ] Export stats to file or API
- [ ] More advanced filtering and sorting

---

## Getting Started

1. **Install Go** (>=1.25)
2. Clone this repo
3. Run: `go run main.go`

---

## Contribution

Feel free to throw a PR if you're sure you know what you're doing but likely message me first to ensure we are on the same page.

---

[paypal-shield]: https://img.shields.io/static/v1?label=PayPal&message=Donate&style=flat-square&logo=paypal&color=blue
[paypal-url]: https://www.paypal.com/donate/?hosted_button_id=MTY5DP7G8G6T4

[coffee-shield]: https://img.shields.io/static/v1?label=BuyMeCoffee&message=Donate&style=flat-square&logo=buy-me-a-coffee&color=orange
[coffee-url]: https://www.buymeacoffee.com/wosiu6

[license-shield]: https://img.shields.io/badge/license-Apache%20License%202.0-purple
[license-url]: https://opensource.org/license/apache-2-0

[github-shield]: https://img.shields.io/static/v1?label=&message=GitHub&style=flat-square&logo=github&color=grey

[paypal-shield]: https://img.shields.io/static/v1?label=PayPal&message=Donate&style=flat-square&logo=paypal&color=blue
[paypal-url]: https://www.paypal.com/donate/?hosted_button_id=MTY5DP7G8G6T4

[coffee-shield]: https://img.shields.io/static/v1?label=BuyMeCoffee&message=Donate&style=flat-square&logo=buy-me-a-coffee&color=orange
[coffee-url]: https://www.buymeacoffee.com/wosiu6

[license-shield]: https://img.shields.io/badge/license-Apache%20License%202.0-purple
[license-url]: https://opensource.org/license/apache-2-0

[github-shield]: https://img.shields.io/static/v1?label=&message=GitHub&style=flat-square&logo=github&color=grey
[github-url]: https://github.com/Wosiu6/docky-go
