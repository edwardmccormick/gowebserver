# Repository Guidelines

## Project Structure & Module Organization
This repository contains three active app areas plus shared project notes:

- `webserverpoc/`: Go backend using Gin, GORM, MySQL, MongoDB, WebSockets, and AI integration. Main entry point is `main.go`; request handlers and middleware live alongside it.
- `react-vite-poc/`: React 19 + Vite web client. UI code is under `src/components/`, shared helpers under `src/utils/`, and static assets in `public/`.
- `react-native-expo/urmid/`: Expo-based mobile prototype with `App.tsx` as the entry point and images in `assets/`.
- `docs/` and top-level `*.md`, `*.json`, and `*.ps1` files: notes, configs, and one-off scripts.

## Build, Test, and Development Commands
- `cd webserverpoc && go run .`: run the Go API locally.
- `cd webserverpoc && go build .`: compile the backend binary.
- `docker compose up --build`: start the Go service plus MySQL and MongoDB.
- `cd react-vite-poc && npm install && npm run dev`: start the Vite frontend.
- `cd react-vite-poc && npm run build`: produce a production web build.
- `cd react-vite-poc && npm run lint`: run ESLint on JS/JSX files.
- `cd react-native-expo/urmid && npm install && npm start`: launch the Expo app.

## Coding Style & Naming Conventions
Use idiomatic formatting for each stack:

- Go: format with `gofmt`, keep package-level files small, and prefer CamelCase for exported symbols.
- React/Vite: 2-space indentation, semicolons, functional components, and lowercase component filenames such as `login.jsx` and `matchlist.jsx`.
- Expo/TypeScript: follow the existing component-first structure in `App.tsx` and keep assets named descriptively.

## Testing Guidelines
There are currently no committed automated tests in the repository. When adding features:

- add Go tests as `*_test.go` files in `webserverpoc/` and run `go test ./...`;
- add frontend tests only if you also add the supporting test tooling;
- include manual verification steps in the PR when automation is missing.

## Commit & Pull Request Guidelines
Recent commit messages are short and informal (`Styling?`, `small but important`). For new work, prefer clear imperative subjects such as `Add profile photo upload validation`.

PRs should include:

- a concise summary of user-visible or API-visible changes;
- linked issue or task reference when available;
- screenshots for UI work in `react-vite-poc/` or `react-native-expo/`;
- notes about config, schema, or Docker changes.

## Configuration & Security Tips
Do not commit secrets. Backend config is split between `webserverpoc/configlocal.json` for local runs and container-mounted `config.json` for Docker. Review CORS origins and database credentials before deploying beyond local development.
