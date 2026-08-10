# Emona Textile

Website for the Emona Textile company — a React/TypeScript frontend and a Go REST API
backed by MongoDB.

## Stack

| Part     | Technologies                                                            |
| -------- | ----------------------------------------------------------------------- |
| Frontend | React 18, TypeScript, MobX, MUI 5, styled-components, webpack 5          |
| Backend  | Go 1.18+, chi router, MongoDB (official Go driver)                       |
| Storage  | MongoDB — collections `carousel` and `catalog` in the `EmonaDB` database |

## Repository layout

```
backend/          Go REST API
  main.go         router setup and entry point
  handlers.go     HTTP handlers and MongoDB configuration
  app/            server command, storage models, logger, mail, templates
frontend/         React application
  src/pages/      screens (Home, About Us, Catalog, Catalog item)
  src/modules/    MobX stores and helpers
  src/components/ shared UI components
  src/routing/    routes
  webpack/        dev and prod build configs
```

## Requirements

- [Go](https://go.dev/dl/) 1.18 or newer
- [Node.js](https://nodejs.org/) 16 or newer, with `yarn` or `npm`
- A running [MongoDB](https://www.mongodb.com/try/download/community) instance

## Configuration

Both parts read their configuration from environment variables. Copy the example files
and fill them in:

```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

`.env` is git-ignored — keep every real connection string, password, and API key out of
version control.

### Backend

Defaults point at a local MongoDB instance, so no credentials are stored in the
repository.

| Variable    | Default                     | Description            |
| ----------- | --------------------------- | ---------------------- |
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection URI |
| `MONGO_DB`  | `EmonaDB`                   | MongoDB database name  |

### Frontend

| Variable                | Description                                          |
| ----------------------- | ---------------------------------------------------- |
| `GOOGLE_MAPS_API_KEY`   | Maps JavaScript API key for the contacts map          |
| `EMAILJS_SERVICE_ID`    | EmailJS service used by the contact form              |
| `EMAILJS_TEMPLATE_ID`   | EmailJS template used by the contact form             |
| `EMAILJS_PUBLIC_KEY`    | EmailJS public key                                    |
| `GA_TRACKING_ID`        | Google Analytics tracking ID; empty disables analytics |

> **These are not secrets.** webpack inlines them into the browser bundle, where anyone
> can read them — moving them into `.env` only makes them easy to rotate per
> environment. Access must be restricted at the provider: HTTP referrer and API
> restrictions for the Maps key (Google Cloud Console → Credentials), allowed origins
> and captcha for EmailJS. Never put a server-side secret in `frontend/.env`.

Variables already set in the shell or in CI take precedence over the `.env` file. To add
a new one, list its name in `CLIENT_ENV_KEYS` in
[webpack/base.config.js](frontend/webpack/base.config.js) and export it from
[src/config.ts](frontend/src/config.ts).

## Getting started

### 1. Start MongoDB

Locally installed, or via Docker:

```bash
docker run -d --name emona-mongo -p 27017:27017 mongo:6
```

The API expects an `EmonaDB` database with a `carousel` and a `catalog` collection.
Both endpoints return an empty array until the collections are populated.

### 2. Run the backend

```bash
cd backend && go run .
```

The server listens on `:4000`. To point it at a different database, set the variables
before starting it:

```bash
MONGO_URI=mongodb://localhost:27017 MONGO_DB=EmonaDB go run .
```

On Windows PowerShell:

```powershell
$env:MONGO_URI = "mongodb://localhost:27017"; go run .
```

A different port can be passed with the `-addr` flag:

```bash
go run . -addr=:8080
```

### 3. Run the frontend

```bash
cd frontend && yarn install && yarn start
```

The dev server runs on [http://localhost:3000](http://localhost:3000) and talks to the
API at `http://localhost:4000` (see [axiosDefaults.ts](frontend/src/axios/axiosDefaults.ts)).
That origin is the only one allowed by the backend's CORS configuration.

## API

| Method | Path        | Description                                |
| ------ | ----------- | ------------------------------------------ |
| `GET`  | `/`         | Health check, returns a plain-text greeting |
| `GET`  | `/carousel` | Carousel slides for the home screen         |
| `GET`  | `/catalog`  | Catalog categories                          |

Responses are JSON. A storage failure returns `500` and logs the cause server-side.

## Frontend scripts

Run from the `frontend` directory:

| Script            | Description                                       |
| ----------------- | ------------------------------------------------- |
| `yarn start`      | Dev server with hot reload on port 3000           |
| `yarn build`      | Production bundle                                 |
| `yarn lint`       | ESLint with `--fix`                               |
| `yarn typecheck`  | TypeScript check without emitting output          |
| `yarn prettier`   | Format sources                                    |
| `yarn stats`      | Build and open the bundle analyzer                |

## Backend commands

Run from the `backend` directory:

```bash
go build ./...   # compile every package
go vet ./...     # static analysis
go run .         # start the API server
```

## License

BSD 3-Clause — see [LICENSE](LICENSE).
