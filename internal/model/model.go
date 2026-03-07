package model

import "time"

// Prompt represents a stored AI prompt
type Prompt struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Tags       []string   `json:"tags"`
	Stack      string     `json:"stack"`
	Models     []string   `json:"models"`
	Verified   bool       `json:"verified"`
	UsageCount int        `json:"usage_count"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

// TestResult represents the result of testing a prompt
type TestResult struct {
	ID              string    `json:"id"`
	PromptID        string    `json:"prompt_id"`
	Model           string    `json:"model"`
	Input           string    `json:"input"`
	ExpectedOutput  string    `json:"expected_output"`
	ActualOutput    string    `json:"actual_output"`
	Passed          bool      `json:"passed"`
	Score           float64   `json:"score"` // 0-100 similarity score
	LatencyMs       int       `json:"latency_ms"`
	TokenUsage      int       `json:"token_usage"`
	ErrorMessage    string    `json:"error_message,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// TestSuite represents a collection of tests for a prompt
type TestSuite struct {
	ID        string       `json:"id"`
	PromptID  string       `json:"prompt_id"`
	Name      string       `json:"name"`
	Tests     []TestResult `json:"tests"`
	PassRate  float64      `json:"pass_rate"` // Percentage of tests passed
	AvgScore  float64      `json:"avg_score"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// PromptVersion represents a versioned snapshot of a prompt
type PromptVersion struct {
	ID          string    `json:"id"`
	PromptID    string    `json:"prompt_id"`
	Version     int       `json:"version"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tags"`
	Stack       string    `json:"stack"`
	Models      []string  `json:"models"`
	Verified    bool      `json:"verified"`
	CommitMsg   string    `json:"commit_msg,omitempty"`
	Author      string    `json:"author,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// DefaultStacks is the hierarchical tech stack taxonomy
var DefaultStacks = []string{
	// Frontend
	"frontend/react/hooks",
	"frontend/react/performance",
	"frontend/react/testing",
	"frontend/react/nextjs",
	"frontend/vue/composition",
	"frontend/vue/nuxt",
	"frontend/svelte/sveltekit",
	"frontend/angular",
	"frontend/css/tailwind",
	"frontend/css/animations",
	"frontend/typescript",
	"frontend/accessibility",

	// Backend
	"backend/node/express",
	"backend/node/nestjs",
	"backend/node/fastify",
	"backend/python/django",
	"backend/python/fastapi",
	"backend/python/flask",
	"backend/go/gin",
	"backend/go/grpc",
	"backend/go/stdlib",
	"backend/rust/axum",
	"backend/rust/actix",
	"backend/java/spring",
	"backend/ruby/rails",

	// Database
	"database/postgresql",
	"database/mysql",
	"database/mongodb",
	"database/redis",
	"database/sqlite",
	"database/prisma",
	"database/drizzle",

	// DevOps
	"devops/docker",
	"devops/kubernetes",
	"devops/terraform",
	"devops/aws",
	"devops/gcp",
	"devops/azure",
	"devops/github-actions",
	"devops/gitlab-ci",

	// AI/ML
	"ai/prompting",
	"ai/rag",
	"ai/agents",
	"ai/fine-tuning",
	"ai/evaluation",
	"ai/langchain",
	"ai/llm-ops",

	// Mobile
	"mobile/react-native",
	"mobile/flutter",
	"mobile/swift",
	"mobile/kotlin",

	// General
	"general/code-review",
	"general/debugging",
	"general/testing",
	"general/documentation",
	"general/security",
	"general/performance",
	"general/architecture",
	"general/refactoring",
}

// SeedPrompts returns curated starter prompts for `promptvault init`
func SeedPrompts() []*Prompt {
	return []*Prompt{
		{
			Title:   "Fix React useEffect dependencies",
			Content: "Analyze this React component and fix all useEffect dependency array issues.\n\nFor each fix:\n1. Explain WHY the dependency is needed (what value could change)\n2. Show the corrected useEffect with the full dependency array\n3. Flag any potential infinite loop risks and explain how to resolve them\n4. Suggest useCallback/useMemo wrapping if a function or object dependency is causing re-renders\n5. If an effect should only run once, explain the trade-off of using an empty dependency array\n\nDo NOT just suppress the ESLint warning — fix the root cause.",
			Stack:   "frontend/react/hooks",
			Tags:    []string{"debugging", "hooks", "async"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
		{
			Title:   "TypeScript strict type narrowing",
			Content: "Review this TypeScript code and apply strict type narrowing:\n\n1. Replace ALL `any` types with specific, narrow types\n2. Add discriminated union types where appropriate (use a `type` or `kind` field)\n3. Implement type guard functions (isFoo(x): x is Foo) for runtime checks\n4. Ensure exhaustive switch/case checks on union types using `never`\n5. Convert loose interfaces to strict types where mutation isn't needed\n6. Add JSDoc comments for any complex generic constraints\n7. Use `satisfies` operator instead of type assertions where possible\n\nShow before/after for each change.",
			Stack:   "frontend/typescript",
			Tags:    []string{"typescript", "types", "narrowing"},
			Models:  []string{"claude-sonnet"},
		},
		{
			Title:   "FastAPI endpoint with full validation",
			Content: "Create a FastAPI endpoint following production standards:\n\n1. Use Pydantic v2 models for request AND response schemas\n2. Return proper HTTP status codes (201 for creation, 404 for not found, 422 for validation)\n3. Handle errors with HTTPException and include error detail schemas\n4. Add OpenAPI documentation: summary, description, response_model, tags\n5. Use Depends() for database session injection\n6. Add rate limiting with slowapi\n7. Include proper CORS headers if needed\n8. Add request ID tracking via middleware\n9. Write the corresponding test using httpx.AsyncClient\n\nReturn type-safe responses, never raw dicts.",
			Stack:   "backend/python/fastapi",
			Tags:    []string{"api", "validation", "pydantic"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
		{
			Title:   "Terraform AWS module best practices",
			Content: "Create a Terraform module for this AWS resource following:\n\n1. Use variables.tf with proper descriptions, types, and validation blocks\n2. Add outputs.tf for all created resource ARNs, IDs, and endpoints\n3. Tag ALL resources: environment, project, managed-by=terraform, created-by\n4. Use data sources instead of hardcoding AMIs, VPC IDs, or account IDs\n5. Use for_each (not count) for creating multiple similar resources\n6. Follow naming: {project}-{environment}-{resource}\n7. Add lifecycle { prevent_destroy = true } for stateful resources\n8. Include a README.md with usage example\n9. Pin provider versions in versions.tf\n10. Use locals for computed values to keep resource blocks clean",
			Stack:   "devops/terraform",
			Tags:    []string{"terraform", "aws", "iac"},
			Models:  []string{"claude-sonnet"},
		},
		{
			Title:   "Go idiomatic error handling",
			Content: "Refactor this Go code to use idiomatic error handling:\n\n1. Wrap errors with fmt.Errorf(\"doing X: %w\", err) for context chains\n2. Add context.Context as first parameter to all functions that do I/O\n3. Define sentinel errors (var ErrNotFound = errors.New(...)) for known failure modes\n4. Use structured logging with slog (not log.Printf)\n5. Return early on errors — no nested if-err blocks deeper than 1 level\n6. Add defer for cleanup: file.Close(), rows.Close(), tx.Rollback()\n7. Use errors.Is() and errors.As() for error checking, never string comparison\n8. For HTTP handlers, use a central error handler middleware\n\nShow the full refactored code, not just snippets.",
			Stack:   "backend/go/stdlib",
			Tags:    []string{"golang", "errors", "context"},
			Models:  []string{"claude-sonnet"},
		},
		{
			Title:   "Production Dockerfile multi-stage build",
			Content: "Create a production-ready Dockerfile using multi-stage builds:\n\n1. Use the smallest possible base image (alpine or distroless)\n2. Stage 1: build stage with all dev dependencies\n3. Stage 2: runtime stage with only the compiled binary/artifacts\n4. Set a non-root USER for security\n5. Add HEALTHCHECK instruction\n6. Use .dockerignore to exclude node_modules, .git, tests\n7. Pin ALL base image versions with SHA256 digest\n8. Order layers from least-changed to most-changed for cache efficiency\n9. Use COPY --from=build to cherry-pick artifacts\n10. Add LABEL with version, maintainer, and description\n11. Set proper EXPOSE and ENTRYPOINT (not CMD for services)",
			Stack:   "devops/docker",
			Tags:    []string{"docker", "deployment", "security"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
		{
			Title:   "SQL query performance optimization",
			Content: "Analyze and optimize this SQL query for performance:\n\n1. Run EXPLAIN ANALYZE and identify the slowest operations\n2. Add or modify indexes to eliminate sequential scans\n3. Rewrite subqueries as JOINs where it improves the query plan\n4. Use CTEs only when they improve readability (they're optimization fences in some DBs)\n5. Add appropriate WHERE clause index hints\n6. Check for N+1 query patterns and suggest batch alternatives\n7. Estimate row counts at each step and identify cardinality issues\n8. Suggest partitioning if the table exceeds 10M rows\n9. Check for implicit type casts that prevent index usage\n10. Provide the optimized query with comments explaining each change",
			Stack:   "database/postgresql",
			Tags:    []string{"sql", "performance", "indexing"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
		{
			Title:   "Comprehensive code review prompt",
			Content: "Review this code with the rigor of a senior engineer. Check:\n\n**Correctness:**\n- Edge cases: null, empty, boundary values, concurrent access\n- Error handling: are all error paths covered? Any swallowed errors?\n- Off-by-one errors in loops and slices\n\n**Security:**\n- Input validation and sanitization\n- SQL injection, XSS, CSRF vulnerabilities\n- Secrets or credentials in code\n- Proper authentication/authorization checks\n\n**Performance:**\n- O(n²) or worse algorithms that could be O(n)\n- Unnecessary allocations in hot paths\n- Missing pagination for list endpoints\n- N+1 query patterns\n\n**Maintainability:**\n- Functions exceeding 50 lines (suggest extraction)\n- God objects or god functions\n- Missing or misleading naming\n- Dead code or commented-out code\n\nFor each issue: state severity (critical/major/minor/nit), location, and a concrete fix.",
			Stack:   "general/code-review",
			Tags:    []string{"review", "quality", "standards"},
			Models:  []string{"claude-sonnet", "gpt-4o", "gemini-pro"},
		},
		{
			Title:   "Debug systematic root cause analysis",
			Content: "I'm debugging an issue. Help me find the root cause systematically:\n\n1. **Reproduce**: Define the exact steps to reproduce. What's the expected vs actual behavior?\n2. **Isolate**: What changed recently? Check git log for the last 5 relevant commits.\n3. **Hypothesize**: List the top 3 most likely causes, ranked by probability.\n4. **Test**: For each hypothesis, what's the quickest way to confirm or eliminate it?\n5. **Bisect**: If it's a regression, use git bisect to find the exact commit.\n6. **Fix**: Once root cause is found, explain WHY it broke, not just what to change.\n7. **Prevent**: Suggest a test that would catch this regression in the future.\n\nDon't guess randomly — work through this checklist in order.",
			Stack:   "general/debugging",
			Tags:    []string{"debugging", "troubleshooting", "methodology"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
		{
			Title:   "Write comprehensive unit tests",
			Content: "Write unit tests for this code following best practices:\n\n1. Use the AAA pattern: Arrange, Act, Assert\n2. Test the happy path FIRST, then edge cases\n3. Cover these edge cases:\n   - Empty/null/undefined inputs\n   - Boundary values (0, -1, MAX_INT)\n   - Invalid types if applicable\n   - Concurrent access if applicable\n4. Each test should have a descriptive name: `test_[action]_[condition]_[expected]`\n5. Use proper mocking — don't call real APIs/databases\n6. Assert specific values, not just \"truthy\"\n7. Add a test for error cases — verify the error message/type\n8. Keep each test independent — no shared mutable state\n9. Aim for >90% line coverage on the target function\n10. Group related tests with describe/context blocks",
			Stack:   "general/testing",
			Tags:    []string{"testing", "unit-tests", "tdd"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
		{
			Title:   "API documentation generator",
			Content: "Generate comprehensive API documentation for this endpoint:\n\n1. **Endpoint**: Method, path, and description\n2. **Authentication**: Required auth method and scopes\n3. **Request**:\n   - Headers (Content-Type, Authorization, custom headers)\n   - Path parameters with types and constraints\n   - Query parameters with defaults and valid ranges\n   - Request body schema with all fields, types, and validation rules\n4. **Response**:\n   - Success response (200/201) with full schema and example\n   - Error responses (400, 401, 403, 404, 422, 500) with error schema\n5. **Examples**: At least 2 curl examples (success + error case)\n6. **Rate Limits**: If applicable\n7. **Changelog**: Note if this is new or changed recently\n\nOutput in OpenAPI 3.0 YAML format.",
			Stack:   "general/documentation",
			Tags:    []string{"api", "docs", "openapi"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
		{
			Title:   "React component with accessibility",
			Content: "Build this React component with full accessibility:\n\n1. Use semantic HTML elements (nav, main, section, article, aside) — not divs for everything\n2. Add proper ARIA attributes: aria-label, aria-describedby, role\n3. Ensure keyboard navigation: Tab, Enter, Escape, Arrow keys work correctly\n4. Add focus management: auto-focus on modal open, return focus on close\n5. Use proper heading hierarchy (h1 → h2 → h3, no skipping)\n6. Ensure color contrast ratio ≥ 4.5:1 (AA) or 7:1 (AAA)\n7. Add alt text for all images, aria-hidden for decorative ones\n8. Handle screen reader announcements for dynamic content (aria-live)\n9. Test with: keyboard-only navigation, VoiceOver/NVDA simulation\n10. Use prefers-reduced-motion for animations\n\nInclude the component AND a note on which WCAG 2.1 criteria each decision addresses.",
			Stack:   "frontend/react/hooks",
			Tags:    []string{"react", "accessibility", "a11y", "wcag"},
			Models:  []string{"claude-sonnet"},
		},
		{
			Title:   "Kubernetes deployment manifest",
			Content: "Create a production Kubernetes deployment for this service:\n\n1. Deployment with rolling update strategy (maxSurge: 1, maxUnavailable: 0)\n2. Resource requests AND limits (CPU + memory) — use realistic values\n3. Liveness probe (HTTP GET /healthz, initialDelay: 10s)\n4. Readiness probe (HTTP GET /readyz, initialDelay: 5s)\n5. Pod disruption budget (minAvailable: 1)\n6. Horizontal pod autoscaler (CPU target: 70%, min: 2, max: 10)\n7. Service with proper selector labels\n8. ConfigMap for non-secret config\n9. Secret references for credentials (never inline)\n10. Anti-affinity rules to spread pods across nodes\n11. ServiceAccount with minimal RBAC\n12. Network policy to restrict ingress/egress\n\nUse proper labels: app, version, component, managed-by.",
			Stack:   "devops/kubernetes",
			Tags:    []string{"k8s", "deployment", "production"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
		{
			Title:   "RAG pipeline implementation",
			Content: "Implement a Retrieval-Augmented Generation pipeline:\n\n1. **Document Processing**:\n   - Chunking strategy: 512 tokens with 50-token overlap\n   - Clean and normalize text (remove headers/footers/page numbers)\n   - Preserve document structure (headings, lists)\n\n2. **Embedding**:\n   - Use text-embedding-3-small (or alternative)\n   - Batch embed for efficiency\n   - Store in vector DB (Pinecone/Weaviate/pgvector)\n\n3. **Retrieval**:\n   - Hybrid search: vector similarity + BM25 keyword\n   - Re-rank top-k results with cross-encoder\n   - Return top 5 chunks with relevance scores\n\n4. **Generation**:\n   - System prompt with retrieved context\n   - Citation: each claim must reference source chunk\n   - Hallucination guard: if context doesn't contain the answer, say so\n\n5. **Evaluation**:\n   - Faithfulness: is the answer supported by retrieved context?\n   - Relevance: are the retrieved chunks actually related to the query?\n   - Track latency per stage (retrieval, generation)",
			Stack:   "ai/rag",
			Tags:    []string{"rag", "embeddings", "vector-search"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
		{
			Title:   "Security audit checklist",
			Content: "Perform a security audit on this codebase:\n\n**Authentication & Authorization:**\n- [ ] Passwords hashed with bcrypt/argon2 (NOT MD5/SHA256)\n- [ ] JWT tokens have expiry and proper signing\n- [ ] API keys stored in environment variables, not code\n- [ ] Role-based access control on all endpoints\n- [ ] Session fixation protection\n\n**Input Handling:**\n- [ ] All user input sanitized before database queries\n- [ ] Parameterized queries (no string concatenation for SQL)\n- [ ] File upload type/size validation\n- [ ] Rate limiting on authentication endpoints\n- [ ] CORS configured to specific domains (not *)\n\n**Data Protection:**\n- [ ] HTTPS enforced everywhere\n- [ ] Sensitive data encrypted at rest\n- [ ] PII logged? Remove or mask it\n- [ ] Database backups encrypted\n- [ ] Secrets rotated periodically\n\n**Infrastructure:**\n- [ ] Dependencies scanned for known CVEs\n- [ ] Docker containers running as non-root\n- [ ] No debug/verbose mode in production\n- [ ] Error messages don't leak internal details\n\nFor each finding: severity (Critical/High/Medium/Low), impact, and fix.",
			Stack:   "general/security",
			Tags:    []string{"security", "audit", "owasp"},
			Models:  []string{"claude-sonnet", "gpt-4o"},
		},
	}
}
