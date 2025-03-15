## Next steps

- [ ] Better validation for current routes
  - [Go Validator](https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme) seems to be a good choice
- [ ] Better error responses
- [ ] Add handlers/repositories for other entities
  - [ ] expense
  - [ ] category
- [ ] Find a way to generate Swagger docs from the code
- [ ] Deploy
  - Google Cloud Run might be a good choice
  - See [discussion](https://www.reddit.com/r/golang/comments/15zgudv/where_would_you_host_a_go_app/)
  - [ ] Enforce HTTPS
  - [ ] Add rate limiting
- [ ] Add Stack authentication
- [ ] Set up CORS
- [ ] Create CI/CD pipeline

## Future improvements

- [ ] Write tests
- [ ] Cache API responses when appropriate
- [ ] Error tracking (Sentry)
- [ ] Add logging
- [ ] Add metrics
