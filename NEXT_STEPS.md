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
  - [ ] Enforce HTTPS
  - [ ] Add rate limiting
- [ ] Create CI/CD pipeline with tests

## Future improvements

- [ ] Write tests
- [ ] Cache API responses when appropriate
- [ ] Error tracking (Sentry)
- [ ] Add logging
- [ ] Add metrics
