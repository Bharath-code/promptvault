package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/model"
)

func Serve(database *db.DB) error {
	s := server.NewMCPServer("PromptVault", "0.1.0")

	// 1. Tool: search_prompts
	searchTool := mcp.NewTool("search_prompts",
		mcp.WithDescription("Full-text search across all prompts in PromptVault"),
		mcp.WithString("query", mcp.Required(), mcp.Description("The search query")),
	)
	s.AddTool(searchTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := req.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError("Argument query is missing or invalid"), nil
		}
		prompts, err := database.Search(query)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var b strings.Builder
		for _, p := range prompts {
			b.WriteString(fmt.Sprintf("---\nID: %s\nTitle: %s\nStack: %s\n\n%s\n\n", p.ID, p.Title, p.Stack, p.Content))
		}
		res := b.String()
		if res == "" {
			res = "No prompts found for query: " + query
		}
		return mcp.NewToolResultText(res), nil
	})

	// 2. Tool: list_stacks
	listStacksTool := mcp.NewTool("list_stacks",
		mcp.WithDescription("List all available tech stack categories"),
	)
	s.AddTool(listStacksTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		_, stacks, _ := database.Stats()
		res := fmt.Sprintf("There are %d active stacks. Default taxonomy includes:\n%s", stacks, strings.Join(model.DefaultStacks, "\n"))
		return mcp.NewToolResultText(res), nil
	})

	// 3. Tool: get_stack_prompts
	getStackTool := mcp.NewTool("get_stack_prompts",
		mcp.WithDescription("Get all prompts for a specific tech stack"),
		mcp.WithString("stack", mcp.Required(), mcp.Description("The tech stack (e.g. frontend/react/hooks)")),
	)
	s.AddTool(getStackTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stack, err := req.RequireString("stack")
		if err != nil {
			return mcp.NewToolResultError("Argument stack is missing or invalid"), nil
		}
		prompts, err := database.List(stack)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var b strings.Builder
		for _, p := range prompts {
			b.WriteString(fmt.Sprintf("---\nID: %s\nTitle: %s\n\n%s\n\n", p.ID, p.Title, p.Content))
		}
		res := b.String()
		if res == "" {
			res = "No prompts found in stack: " + stack
		}
		return mcp.NewToolResultText(res), nil
	})

	// Expose standard resource templates for cursor integrations
	s.AddResourceTemplate(mcp.NewResourceTemplate(
		"promptvault://stack/{category}",
		"Tech Stack Prompts",
		mcp.WithTemplateDescription("A curated set of system prompts for {category}"),
	), func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		stack := strings.TrimPrefix(req.Params.URI, "promptvault://stack/")
		prompts, err := database.List(stack)
		if err != nil {
			return nil, err
		}
		var b strings.Builder
		for _, p := range prompts {
			b.WriteString(fmt.Sprintf("---\n%s\n---\n%s\n\n", p.Title, p.Content))
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
