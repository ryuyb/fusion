# Task Completion Checklist

## Before Committing/Pushing Code

### 1. Code Quality Checks
- [ ] Run linter: `make lint`
- [ ] Fix all linter errors and warnings
- [ ] Ensure code is formatted (gofmt)
- [ ] Check for any golangci-lint issues

### 2. Testing
- [ ] Run all tests: `make test`
- [ ] Ensure all unit tests pass: `make test-unit`
- [ ] Run integration tests: `make test-integration`
- [ ] Run E2E tests: `make test-e2e`
- [ ] Check test coverage: `make test-coverage`
- [ ] Add tests for new functionality

### 3. Build Verification
- [ ] Clean build: `go clean -testcache && go build -o bin/fusion cmd/server/main.go`
- [ ] Verify binary works: `./bin/fusion --help`
- [ ] Test in Docker: `make docker-build`

### 4. Database Changes
- [ ] If schema changed: `make generate-ent`
- [ ] Run migrations: `go run cmd/server/main.go migrate up`
- [ ] Verify schema in database

### 5. Documentation
- [ ] Update Swagger docs: `make generate-swagger`
- [ ] Add/update docstrings for new public functions
- [ ] Update API documentation if needed
- [ ] Add comments for complex logic

### 6. Configuration
- [ ] Update config.yaml if new config fields added
- [ ] Update docker-compose.yml if needed
- [ ] Check environment variable requirements

### 7. Dependencies
- [ ] Check for unused dependencies: `go mod tidy`
- [ ] Update go.mod if new dependencies added
- [ ] Run `go mod vendor` if vendoring is used

### 8. Security
- [ ] Review for security vulnerabilities
- [ ] Check authentication/authorization is properly implemented
- [ ] Verify input validation is in place
- [ ] Ensure sensitive data is not logged
- [ ] Check for SQL injection vulnerabilities
- [ ] Review JWT secret configuration

### 9. Performance
- [ ] Check for N+1 queries in database code
- [ ] Verify proper database connection pool settings
- [ ] Review for memory leaks
- [ ] Check concurrent access patterns

### 10. Code Review
- [ ] Self-review the changes
- [ ] Check for proper error handling
- [ ] Ensure Clean Architecture boundaries are respected
- [ ] Verify dependency injection is properly configured
- [ ] Check naming conventions are followed
- [ ] Review for code duplication

## Post-Commit Actions

### 11. CI/CD
- [ ] Check CI pipeline passes
- [ ] Verify tests run in CI environment
- [ ] Check build succeeds in CI
- [ ] Review test coverage in CI

### 12. Documentation Updates
- [ ] Update CHANGELOG if applicable
- [ ] Update README if needed
- [ ] Notify team of breaking changes
- [ ] Update API documentation

## Docker and Deployment

### 13. Docker Changes
- [ ] Rebuild Docker image: `make docker-build`
- [ ] Test with Docker Compose: `make docker-up`
- [ ] Check logs: `make docker-logs`
- [ ] Test healthchecks work
- [ ] Multi-architecture build if needed: `make docker-build-multiarch`

### 14. Database Migration Checklist
- [ ] Create migration for schema changes
- [ ] Test migration in staging environment
- [ ] Create rollback script if needed
- [ ] Verify migration doesn't cause downtime
- [ ] Update seed data if needed

## Environment-Specific Checks

### 15. Development
- [ ] Run `air` for hot reload works
- [ ] Test with different config files
- [ ] Verify environment variables work

### 16. Staging/Production
- [ ] Test with production-like data
- [ ] Check performance with production load
- [ ] Verify monitoring and logging
- [ ] Check error rates and alerts
- [ ] Test backup and restore procedures

## Common Issues to Watch For

### Architecture Violations
- [ ] Domain layer has no dependencies on outer layers
- [ ] Repository interfaces are in domain, implementations in infrastructure
- [ ] Application layer only depends on domain
- [ ] No circular dependencies between packages

### EntGO Specific
- [ ] Soft delete mixin is used where appropriate
- [ ] Ent code is regenerated after schema changes
- [ ] No direct Ent client usage outside infrastructure layer
- [ ] Use repository pattern for data access

### Error Handling
- [ ] No panic() in production code
- [ ] All errors are properly wrapped
- [ ] Error messages don't expose sensitive info
- [ ] Error codes are meaningful

### Concurrency
- [ ] Goroutines are properly managed
- [ ] No race conditions in tests: `go test -race`
- [ ] Channel operations are safe
- [ ] Context is used for cancellation

## Final Verification
- [ ] All TODOs and FIXMEs are addressed or have tracking issues
- [ ] No debug statements left in code
- [ ] No hardcoded credentials or secrets
- [ ] Environment variables are properly documented
- [ ] API versioning is consistent
- [ ] Response formats are consistent
- [ ] HTTP status codes are appropriate
- [ ] Logging is at appropriate levels
