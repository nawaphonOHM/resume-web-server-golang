# resume-web-server-golang

High-performance static single page application (SPA) web server in Go with direct MongoDB metadata resolution and cloud object storage streaming using the `github.com/nawaphonOHM/whatever/v3` framework.

---

## HTTP Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/` | Resolves and streams the entrypoint `index.html` asset directly from cloud storage. |
| `GET` | `/{fileName}` | Resolves `{fileName}` from MongoDB; if found, streams object from storage; if not found, transparently falls back to `index.html` (SPA routing). |
| `GET` | `/health` | Framework readiness / liveness probe returning HTTP 200. |
| `GET` | `/ready` | Framework readiness probe returning HTTP 200. |

> **Routing limitation:** Only single-segment request paths are routed. Multi-segment HTTP paths such as `/assets/app.js` or `/projects/42` return HTTP 404 and do **not** fall back to `index.html`. This applies to HTTP routing only; the MongoDB `location` metadata may still reference nested object keys in cloud storage.

---

## MongoDB Collection Schema

The application queries the `resume-single-page-application` MongoDB collection to look up static file locations.

Documents must follow this schema:

```json
{
  "_id": { "$oid": "..." },
  "fileName": "index.html",
  "location": "gs://bucket-name/index.html",
  "hashValueFile": "abc123hash"
}
```

### Schema Rules & Normalization

- **`fileName`** (string, required): Lookup key for static assets. Paths are stored without leading slashes (e.g., `index.html`, `app.js`). HTTP requests with leading slashes (such as `/index.html` or `///app.js`) are stripped and normalized prior to querying MongoDB.
- **`location`** (string, required): Cloud object storage URI or object key indicating the location of the static file. Supported formats:
  - **Cloud storage URIs (`gs://` or `s3://`)**:
    - `gs://bucket-name/path/to/file.ext` -> target bucket `bucket-name`, object key `path/to/file.ext`.
    - `s3://bucket-name/path/to/file.ext` -> target bucket `bucket-name`, object key `path/to/file.ext`.
    - `gs:///path/to/file.ext` -> target bucket configured default bucket (`OHM9996_STORAGE_BUCKET`), object key `path/to/file.ext`.
    - Note: URIs without an object key (e.g., `gs://bucket-name` or `gs://bucket-name/`) are invalid. Non-cloud schemes (e.g., `http://`, `https://`) are not supported.
  - **Non-URI relative or absolute object keys**:
    - `index.html` or `/index.html` -> target bucket configured default bucket (`OHM9996_STORAGE_BUCKET`), object key `index.html`.
    - `assets/main.js` or `/assets/main.js` -> target bucket configured default bucket (`OHM9996_STORAGE_BUCKET`), object key `assets/main.js`.
    - Non-URI paths always resolve strictly as object keys in the default bucket (leading slashes are stripped) and are never split into bucket and key segments. Note that nested object keys in cloud storage are supported via the metadata `location`, while the HTTP route itself serves single-depth requests (`GET /{fileName}`).
- **`hashValueFile`** (string, required): Content hash value for validation and cache identification.

---

## Environment Configuration

Configuration values are loaded through standard `OHM9996_` environment variables or an optional `.env` file.

### Application Settings
| Environment Variable | Default | Description |
|---|---|---|
| `OHM9996_INDEX_FILE` | `index.html` | SPA entrypoint file name. |
| `OHM9996_MONGODB_COLLECTION` | `resume-single-page-application` | MongoDB collection storing static file metadata. |

### Server Settings
| Environment Variable | Default | Description |
|---|---|---|
| `OHM9996_SERVER_HOST` | *(empty)* | Host address interface to bind HTTP listener (empty binds to all interfaces). |
| `OHM9996_SERVER_PORT` | `8080` | Port for HTTP listener. |
| `OHM9996_SERVER_API_GROUP` | `raw` | API route grouping mode (`standard`, `raw`). Framework default is `standard`; application defaults to `raw` when unset for root SPA routing. |
| `OHM9996_GIN_MODE` | `release` | Server execution mode (`release`, `debug`, `test`). |
| `OHM9996_SERVER_DISPLAY_NAME` | `application` | Service display name reported in server logs. |
| `OHM9996_SERVER_READ_TIMEOUT` | `10s` | HTTP server read timeout duration. |
| `OHM9996_SERVER_WRITE_TIMEOUT` | `10s` | HTTP server write timeout duration. |
| `OHM9996_SERVER_IDLE_TIMEOUT` | `60s` | HTTP server idle timeout duration. |
| `OHM9996_SERVER_SHUTDOWN_TIMEOUT` | `10s` | Graceful shutdown timeout duration upon receiving SIGINT/SIGTERM. |

### MongoDB Settings
| Environment Variable | Default | Description |
|---|---|---|
| `OHM9996_MONGODB_HOST` | *(required)* | MongoDB host address or Google Cloud Firestore endpoint (e.g. `*.firestore.goog`). |
| `OHM9996_MONGODB_PORT` | `27017` | Port for MongoDB connection (use `443` for Google Cloud Firestore). |
| `OHM9996_MONGODB_DATABASE` | *(empty)* | Target database name. |
| `OHM9996_MONGODB_USERNAME` | *(empty)* | Authentication username. |
| `OHM9996_MONGODB_PASSWORD` | *(empty)* | Authentication password. |
| `OHM9996_MONGODB_AUTH_SOURCE` | *(empty)* | Authentication database source. |
| `OHM9996_MONGODB_PROTOCOL` | `mongodb` | Protocol scheme (`mongodb` or `mongodb+srv`). |
| `OHM9996_MONGODB_CONNECT_TIMEOUT` | `10s` | MongoDB connection timeout duration. |
| `OHM9996_MONGODB_SERVER_SELECTION_TIMEOUT` | `5s` | Server selection timeout duration. |
| `OHM9996_MONGODB_SOCKET_TIMEOUT` | `10s` | Socket read/write timeout duration. |
| `OHM9996_MONGODB_MAX_CONN_IDLE_TIME` | `10m` | Maximum connection idle duration. |
| `OHM9996_MONGODB_MAX_POOL_SIZE` | `100` | Maximum connection pool size. |
| `OHM9996_MONGODB_MIN_POOL_SIZE` | `5` | Minimum connection pool size. |

*Note: Google Cloud Firestore MongoDB endpoints (`firestore.goog`) are automatically detected and configured with TLS, SCRAM-SHA-256 authentication, `loadBalanced=true`, `retryWrites=false`, and Firestore-safe collection probes.*

### Object Storage Settings
| Environment Variable | Default | Description |
|---|---|---|
| `OHM9996_STORAGE_BUCKET` | *(required)* | Target cloud storage bucket name (must be 3–222 characters). Also used as the default bucket when metadata `location` is a relative object key. |
| `OHM9996_STORAGE_REGION` | `us-east-1` | Storage bucket region (e.g. `us-east-1`, `asia-southeast1`). |
| `OHM9996_STORAGE_ENDPOINT` | *(empty)* | Custom endpoint URL for S3/GCS compatibility (e.g. `https://storage.googleapis.com` for GCS XML API or `http://localhost:9000` for MinIO). |
| `OHM9996_STORAGE_ACCESS_KEY_ID` | *(empty)* | HMAC access key ID for S3/GCS authentication. |
| `OHM9996_STORAGE_SECRET_ACCESS_KEY` | *(empty)* | HMAC secret access key for S3/GCS authentication. |
| `OHM9996_STORAGE_USE_PATH_STYLE` | `false` | Enable path-style URL addressing (`http://endpoint/bucket/key`). |
| `OHM9996_STORAGE_CONNECT_TIMEOUT` | `10s` | Storage client connection timeout duration. |
| `OHM9996_STORAGE_REQUEST_TIMEOUT` | `30s` | Storage get object request timeout duration. |
| `OHM9996_STORAGE_PING_TIMEOUT` | `5s` | Storage bucket verification timeout duration. |

### Logging & OpenTelemetry Settings
| Environment Variable | Default | Description |
|---|---|---|
| `OHM9996_LOGGING_LEVEL` | `info` | Structured log level (`debug`, `info`, `warn`, `error`). |
| `OHM9996_LOGGING_FORMAT` | `text` | Structured log output format (`text`, `json`). |
| `OHM9996_OTEL_ENABLED` | `true` | Enable OpenTelemetry tracing exporter. Set to `false` in local development if no collector is running. |
| `OHM9996_OTEL_EXPORTER_OTLP_ENDPOINT` | `localhost:4317` | OpenTelemetry OTLP gRPC collector endpoint. |

---

## Deployment

The application is containerized with Cloud Native Buildpacks and deployed to Google Cloud Run via Google Cloud Build.

### 1. Framework Dependency
The application directly depends on upstream `github.com/nawaphonOHM/whatever/v3` (pinned at `v3.7.0`). The framework provides native support for raw unversioned routing (`rest.APIGroupRaw`), request context access, streaming responses, and managed MongoDB/Storage integrations.

### 2. Google Secret Manager & IAM Setup
In Cloud Run deployments, secrets are injected as environment variables via Secret Manager. Ensure the following IAM roles and secrets are configured:

#### IAM Permissions
- **Cloud Run Runtime Service Account (`--service-account=$_SERVICE_ACCOUNT`)**: Must be granted the `roles/secretmanager.secretAccessor` role on each configured secret so Cloud Run can access and inject them at revision startup.
- **Object Storage HMAC Principal**: The service account or user associated with the HMAC credentials (`OHM9996_STORAGE_ACCESS_KEY_ID` / `OHM9996_STORAGE_SECRET_ACCESS_KEY`) must have read permissions on the target bucket (e.g. `roles/storage.objectViewer` or `storage.objects.get` on Google Cloud Storage).

#### Secret Manager Mappings
Ensure the following secrets are created and versioned in Google Secret Manager:

- **MongoDB / Firestore Secrets**:
  - `OHM9996_MONGODB_HOST`: MongoDB host or Firestore endpoint (`firestore.goog`).
  - `OHM9996_MONGODB_PORT`: MongoDB port (`27017` or `443` for Firestore).
  - `OHM9996_MONGODB_DATABASE`: Target database name.
  - `OHM9996_MONGODB_USERNAME`: Database username.
  - `OHM9996_MONGODB_PASSWORD`: Database password.
  - `OHM9996_MONGODB_AUTH_SOURCE`: Auth database (e.g. `admin` or project ID for Firestore).
  - `OHM9996_MONGODB_MAX_CONN_IDLE_TIME`: e.g. `10m`.
  - `OHM9996_MONGODB_SOCKET_TIMEOUT`: e.g. `10s`.
- **Cloud Storage Secrets (GCS / S3 Interoperability)**:
  - `OHM9996_STORAGE_BUCKET`: Target GCS or S3 bucket name.
  - `OHM9996_STORAGE_ENDPOINT`: `https://storage.googleapis.com` (for Google Cloud Storage) or AWS S3 endpoint.
  - `OHM9996_STORAGE_REGION`: Storage region (e.g. `asia-southeast1` or `us-east-1`).
  - `OHM9996_STORAGE_ACCESS_KEY_ID`: Google Cloud Storage HMAC Access Key ID.
  - `OHM9996_STORAGE_SECRET_ACCESS_KEY`: Google Cloud Storage HMAC Secret Access Key.
  - `OHM9996_STORAGE_USE_PATH_STYLE`: `false` (or `true` for local MinIO).
- **Server & Observability Secrets**:
  - `OHM9996_SERVER_HOST`: `""` or `0.0.0.0`.
  - `OHM9996_SERVER_PORT`: `8080`.
  - `OHM9996_SERVER_API_GROUP`: `raw`.
  - `OHM9996_OTEL_ENABLED`: `false` (or `true` if OTLP collector is configured).
  - `OHM9996_LOGGING_FORMAT`: `json`.

### 3. Health & Startup Probes
Cloud Run is configured in `cloudbuild.yaml` with a startup probe targeting the built-in HTTP health check:
```yaml
--startup-probe=httpGet.port=8080,httpGet.path=/health,timeoutSeconds=1,periodSeconds=10,failureThreshold=3
```
The `/health` endpoint responds with HTTP 200 immediately once the REST listener and connection pools are initialized.

---

## Framework Routing & API Architecture

This repository uses upstream `github.com/nawaphonOHM/whatever/v3` (pinned at `v3.7.0`).

### Root SPA Routing (`rest.APIGroupRaw`)
The static SPA routes (`/`, `/:fileName`) are registered with `Group: rest.APIGroupRaw`, allowing endpoints to mount directly at the root URL path without `/api` or version prefixes.

### Dependency Boundary
To maintain a strict dependency boundary, `github.com/gin-gonic/gin` and `go.mongodb.org/mongo-driver/v2` are permitted only as transitive dependencies via `github.com/nawaphonOHM/whatever/v3`. They are strictly banned from direct imports across the repository by `depguard` lint rules, and are excluded from direct `require` entries in the root `go.mod`.

### Routing Constraints & Shadowing

- **Gin Radix Tree Wildcards**: Gin's tree-based router disallows a single root catch-all wildcard (`/*fileName`) because it conflicts with framework static probe routes (`/health` and `/ready`). The handler registers `/` and `/:fileName` to serve single-level asset paths without wildcard collisions.
- **Probe Route Shadowing**: Static assets named literally `health` or `ready` are shadowed by framework health check endpoints.

---

## Build & Test

### Run Unit Tests
```bash
go test -race -v ./...
```

### Run Linter
```bash
golangci-lint run ./...
```

### Build Binary
```bash
go build -o bin/server .
```
