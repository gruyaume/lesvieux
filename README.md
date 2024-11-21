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

| Endpoint                          | HTTP Method | Description                   | Parameters      |
| --------------------------------- | ----------- | ----------------------------- | --------------- |
| `/api/v1/posts`                   | GET         | List public job posts         |                 |
| `/api/v1/employers`               | GET         | List employers                |                 |
| `/api/v1/employers`               | POST        | Create employer               | email, password |
| `/api/v1/employers/{id}`          | GET         | Get employer by id            |                 |
| `/api/v1/employers/{id}`          | DELETE      | Delete employer by id         |                 |
| `/api/v1/employers/accounts`      | GET         | List employer accounts        |                 |
| `/api/v1/employers/accounts`      | POST        | Create employer account       | employer_id     |
| `/api/v1/employers/accounts/{id}` | GET         | Get employer account by id    |                 |
| `/api/v1/employers/accounts/{id}` | PUT         | Update employer account by id | employer_id     |
| `/api/v1/employers/accounts/{id}` | DELETE      | Delete employer account by id |                 |
| `/api/v1/employers/login`         | POST        | Employer Login                | email, password |
| `/api/v1/admin/login`             | POST        | Admin Login                   | email, password |
| `/api/v1/admin/accounts`          | GET         | List admin accounts           | email, password |
| `/api/v1/admin/accounts`          | POST        | Create admin account          | email, password |
| `/api/v1/admin/accounts/{id}`     | GET         | Get admin account by id       |                 |
| `/api/v1/admin/accounts/{id}`     | PUT         | Update admin account by id    | email, password |
| `/api/v1/admin/accounts/{id}`     | DELETE      | Delete admin account by id    |                 |
| `/metrics`                        | Get         | Get Prometheus metrics        |                 |
| `/status`                         | Get         | Get service status            |                 |

#### Authentication

The API requires authentication. To authenticate, send a POST request to `/api/v1/admin/login` with the email and password in the body. The response will contain a JWT token. Include this token in the `Authorization` header of subsequent requests.

### Metrics

In addition to the Go runtime metrics, the following custom metrics are exposed:
* `http_requests_total`: The total number of HTTP requests.
* `http_request_duration_seconds`: The duration of HTTP requests in seconds.
* `job_posts_total`: The total number of job posts.
