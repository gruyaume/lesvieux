# LesVieux

LesVieux is a job platform for seniors.

## Getting Started

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

| Endpoint                       | HTTP Method | Description                     | Parameters                                  |
| ------------------------------ | ----------- | ------------------------------- | ------------------------------------------- |
| `/api/v1/published_posts`      | GET         | Get all published posts         |                                             |
| `/api/v1/published_posts/{id}` | GET         | Get a published blog post by id |                                             |
| `/api/v1/posts`                | GET         | Get all blog posts              |                                             |
| `/api/v1/posts`                | POST        | Create a new blog post          |                                             |
| `/api/v1/me/posts/{id}`        | PUT         | Update a blog post by id        | title, content, status ("draft, published") |
| `/api/v1/me/posts/{id}`        | GET         | Get a blog post by id           |                                             |
| `/api/v1/me/posts/{id}`        | DELETE      | Delete a blog post by id        |                                             |
| `/api/v1/accounts`             | GET         | Get all accounts                |                                             |
| `/api/v1/accounts`             | POST        | Create a new account            | username, password                          |
| `/api/v1/accounts/{id}`        | GET         | Get an account by id            |                                             |
| `/api/v1/accounts/{id}`        | DELETE      | Delete an account by id         |                                             |
| `/api/v1/login`                | POST        | Login                           | username, password                          |
| `/metrics`                     | Get         | Get Prometheus metrics          |                                             |
| `/status`                      | Get         | Get service status              |                                             |

#### Authentication

The API requires authentication. To authenticate, send a POST request to `/api/v1/admin/login` with the username and password in the body. The response will contain a JWT token. Include this token in the `Authorization` header of subsequent requests.

### Metrics

In addition to the Go runtime metrics, the following custom metrics are exposed:
* `http_requests_total`: The total number of HTTP requests.
* `http_request_duration_seconds`: The duration of HTTP requests in seconds.
* `blog_posts_total`: The total number of blog posts.
