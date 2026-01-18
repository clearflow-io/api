## Next steps

- [x] Better validation for current routes
  - [Go Validator](https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme) seems to be a good choice
- [x] Better error responses
- [x] Add handlers/repositories for other entities
  - [x] expense
  - [x] category
- [x] Set up CORS
- [x] Add Clerk authentication
- [x] Change error responses to use `json:"error"` instead of `json:"errors"`
- [x] Deploy
- [x] Enforce HTTPS
- [x] Add rate limiting (using `httprate`)
- [x] Create CI/CD pipeline with tests (Workflow added, branch protection pending)
- [x] Enhance logging (capture status codes and response size)
- [ ] Add security headers (HSTS, CSP, etc.)

## Future improvements

- [ ] Write more tests
- [ ] Cache API responses when appropriate
- [ ] Error tracking (Sentry)
- [ ] Centralized logging (Elasticsearch/Logstash/Kibana or similar)
- [ ] Audit logging for sensitive operations
- [ ] Log rotation and retention policy
- [ ] Performance monitoring (OpenTelemetry)
