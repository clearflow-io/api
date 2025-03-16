## Next steps

- [x] Better validation for current routes
  - [Go Validator](https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme) seems to be a good choice
- [x] Better error responses
- [x] Add handlers/repositories for other entities
  - [x] expense
  - [x] category
- [ ] Add Stack authentication
- [ ] Set up CORS
- [ ] Deploy
  - Google Cloud Run might be a good choice
  - See [discussion](https://www.reddit.com/r/golang/comments/15zgudv/where_would_you_host_a_go_app/)
  - [ ] Enforce HTTPS
  - [ ] Add rate limiting
- [ ] Create CI/CD pipeline

## Future improvements

- [ ] Write tests
- [ ] Cache API responses when appropriate
- [ ] Error tracking (Sentry)
- [ ] Add logging
- [ ] Add metrics
