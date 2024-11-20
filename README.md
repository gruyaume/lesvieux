# LesVieux

LesVieux is a job platform for seniors.

## Getting Started

Clone the repository:

```shell
git clone git@github.com:gruyaume/lesvieux.git
```

Generate (or copy) a certificate and private key to the following location:

```shell
openssl req -newkey rsa:2048 -nodes -keyout /var/snap/hebdo/common/key.pem -x509 -days 1 -out /var/snap/hebdo/common/cert.pem -subj "/CN=example.com"
```

Build the frontend:

```shell
npm install --prefix ui
npm run build --prefix ui
```

View the frontend:

```shell
go run cmd/lesvieux/main.go -config lesvieux.yaml
```

## Reference

### Configuration

The configuration file is a YAML file with the following fields:

```yaml
db_path: "./lesvieux.db"
port: 8000
tls:
  cert: "cert.pem"
  key: "key.pem"
```

### API

| Endpoint                | HTTP Method | Description             | Parameters         |
| ----------------------- | ----------- | ----------------------- | ------------------ |
| `/api/v1/posts`         | GET         | Get all job posts       |                    |
| `/api/v1/posts`         | POST        | Create a new job post   |                    |
| `/api/v1/accounts`      | GET         | Get all accounts        |                    |
| `/api/v1/accounts`      | POST        | Create a new account    | username, password |
| `/api/v1/accounts/{id}` | GET         | Get an account by id    |                    |
| `/api/v1/accounts/{id}` | DELETE      | Delete an account by id |                    |
| `/api/v1/login`         | POST        | Login                   | username, password |
| `/metrics`              | Get         | Get Prometheus metrics  |                    |
| `/status`               | Get         | Get service status      |                    |

#### Authentication

The API requires authentication. To authenticate, send a POST request to `/api/v1/admin/login` with the username and password in the body. The response will contain a JWT token. Include this token in the `Authorization` header of subsequent requests.

### Metrics

In addition to the Go runtime metrics, the following custom metrics are exposed:
* `http_requests_total`: The total number of HTTP requests.
* `http_request_duration_seconds`: The duration of HTTP requests in seconds.
* `job_posts_total`: The total number of job posts.
