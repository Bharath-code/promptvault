package mcp

import (
	"context"
	"fmt"
	"html"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/model"
	"golang.org/x/time/rate"
)

// sanitizeOutput escapes HTML special characters to prevent injection
func sanitizeOutput(s string) string {
	return html.EscapeString(s)
}

// rateLimiter manages rate limiting for MCP tools
var (
	rateLimiter   = rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests per second
	rateLimiterMu sync.Mutex
)

// checkRateLimit returns an error if rate limit is exceeded
func checkRateLimit() error {
	rateLimiterMu.Lock()
	defer rateLimiterMu.Unlock()
	
	if !rateLimiter.Allow() {
		return fmt.Errorf("rate limit exceeded. Please wait before making another request")
	}
	return nil
}

func Serve(database *db.DB) error {
	s := server.NewMCPServer("PromptVault", "0.1.0")

	// 1. Tool: search_prompts
	searchTool := mcp.NewTool("search_prompts",
		mcp.WithDescription("Full-text search across all prompts in PromptVault"),
		mcp.WithString("query", mcp.Required(), mcp.Description("The search query")),
	)
	s.AddTool(searchTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Check rate limit
		if err := checkRateLimit(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		query, err := req.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError("Argument query is missing or invalid"), nil
		}
		prompts, err := database.Search(ctx, query)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var b strings.Builder
		for _, p := range prompts {
			b.WriteString(fmt.Sprintf("---\nID: %s\nTitle: %s\nStack: %s\n\n%s\n\n", 
				sanitizeOutput(p.ID), 
				sanitizeOutput(p.Title), 
				sanitizeOutput(p.Stack), 
				sanitizeOutput(p.Content)))
		}
		res := b.String()
		if res == "" {
			res = "No prompts found for query: " + sanitizeOutput(query)
		}
		return mcp.NewToolResultText(res), nil
	})

	// 2. Tool: list_stacks
	listStacksTool := mcp.NewTool("list_stacks",
		mcp.WithDescription("List all available tech stack categories"),
	)
	s.AddTool(listStacksTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Check rate limit
		if err := checkRateLimit(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		_, stacks, _ := database.Stats(ctx)
		res := fmt.Sprintf("There are %d active stacks. Default taxonomy includes:\n%s", stacks, strings.Join(model.DefaultStacks, "\n"))
		return mcp.NewToolResultText(res), nil
	})

	// 3. Tool: get_stack_prompts
	getStackTool := mcp.NewTool("get_stack_prompts",
		mcp.WithDescription("Get all prompts for a specific tech stack"),
		mcp.WithString("stack", mcp.Required(), mcp.Description("The tech stack (e.g. frontend/react/hooks)")),
	)
	s.AddTool(getStackTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Check rate limit
		if err := checkRateLimit(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		stack, err := req.RequireString("stack")
		if err != nil {
			return mcp.NewToolResultError("Argument stack is missing or invalid"), nil
		}
		prompts, err := database.List(ctx, stack)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var b strings.Builder
		for _, p := range prompts {
			b.WriteString(fmt.Sprintf("---\nID: %s\nTitle: %s\n\n%s\n\n", 
				sanitizeOutput(p.ID), 
				sanitizeOutput(p.Title), 
				sanitizeOutput(p.Content)))
		}
		res := b.String()
		if res == "" {
			res = "No prompts found in stack: " + sanitizeOutput(stack)
		}
		return mcp.NewToolResultText(res), nil
	})

	// Expose standard resource templates for cursor integrations
	s.AddResourceTemplate(mcp.NewResourceTemplate(
		"promptvault://stack/{category}",
		"Tech Stack Prompts",
		mcp.WithTemplateDescription("A curated set of system prompts for {category}"),
	), func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Check rate limit
		if err := checkRateLimit(); err != nil {
			return nil, err
		}
		
		stack := strings.TrimPrefix(req.Params.URI, "promptvault://stack/")
		prompts, err := database.List(ctx, stack)
		if err != nil {
			return nil, err
		}
		var b strings.Builder
		for _, p := range prompts {
			b.WriteString(fmt.Sprintf("---\n%s\n---\n%s\n\n", 
				sanitizeOutput(p.Title), 
				sanitizeOutput(p.Content)))
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "text/markdown",
				Text:     b.String(),
			},
		}, nil
	})

	// Start standard I/O server
	return server.ServeStdio(s)
}
