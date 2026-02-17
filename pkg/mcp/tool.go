package mcp

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	miniflux "miniflux.app/v2/client"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

type handler struct {
	minifluxUrl string
}

func RegisterTools(server *mcp.Server, minifluxUrl string) {
	h := handler{minifluxUrl}

	mcp.AddTool(server, &mcp.Tool{
		Name: "get_feeds", Description: "Get list of subcribed feeds on Miniflux instance",
		InputSchema: getFeedSchema,
	}, h.getFeeds)

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "get_entry",
			Description: "Get an entry of subcribed feeds on Miniflux instance by its id",
			InputSchema: getEntrySchema,
		},
		h.getEntry)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_entries",
		Description: "Get entries of subcribed feeds on Miniflux instance",
		InputSchema: getEntriesSchema,
	}, h.getEntries)
}

func (h *handler) getFeeds(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, *getFeedsResult, error) {
	minifluxApiKey := req.Extra.Header.Get("X-Api-Key")
	minifluxCli := miniflux.NewClient(h.minifluxUrl, minifluxApiKey)

	feeds, err := minifluxCli.FeedsContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	output := make([]feed, len(feeds))
	for i, f := range feeds {
		output[i] = feed{
			ID:    f.ID,
			Title: f.Title,
			URL:   f.SiteURL,
		}
	}

	textContent, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(textContent),
				},
			},
		}, &getFeedsResult{
			Feeds: output,
		}, nil
}

func (h *handler) getEntry(ctx context.Context, req *mcp.CallToolRequest, args getEntryParams) (*mcp.CallToolResult, *getEntryResult, error) {
	minifluxApiKey := req.Extra.Header.Get("X-Api-Key")
	minifluxCli := miniflux.NewClient(h.minifluxUrl, minifluxApiKey)

	e, err := minifluxCli.EntryContext(ctx, args.EntryID)
	if err != nil {
		return nil, nil, err
	}

	markdown, err := htmltomarkdown.ConvertString(e.Content)
	if err != nil {
		log.Fatal(err)
	}

	output := entry{
		ID:        e.ID,
		Title:     e.Title,
		URL:       e.URL,
		Content:   markdown,
		CreatedAt: e.CreatedAt,
	}

	textContent, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(textContent),
				},
			},
		}, &getEntryResult{
			Entry: output,
		}, nil
}

func (h *handler) getEntries(ctx context.Context, req *mcp.CallToolRequest, args getEntriesParams) (*mcp.CallToolResult, *getEntriesResult, error) {
	minifluxApiKey := req.Extra.Header.Get("X-Api-Key")
	minifluxCli := miniflux.NewClient(h.minifluxUrl, minifluxApiKey)

	var categoryID int64
	if args.Category != "" {
		categories, err := minifluxCli.CategoriesContext(ctx)
		if err != nil {
			return nil, nil, err
		}

		for _, category := range categories {
			if strings.EqualFold(category.Title, args.Category) {
				categoryID = category.ID
				break
			}
		}
	}

	var feedID int64
	if args.Feed != "" {
		feeds, err := minifluxCli.FeedsContext(ctx)
		if err != nil {
			return nil, nil, err
		}
		for _, feed := range feeds {
			if strings.EqualFold(feed.Title, args.Feed) {
				feedID = feed.ID
			}
		}
	}

	var publishedBefore int64
	if !args.PublishedBefore.IsZero() {
		publishedBefore = args.PublishedBefore.Unix()
	}

	var publishedAfter int64
	if !args.PublishedAfter.IsZero() {
		publishedAfter = args.PublishedAfter.Unix()
	}

	filter := &miniflux.Filter{
		Status:          args.Status,
		Order:           "published_at",
		Direction:       "asc",
		Search:          args.Search,
		Limit:           args.Limit,
		CategoryID:      categoryID,
		FeedID:          feedID,
		PublishedBefore: publishedBefore,
		PublishedAfter:  publishedAfter,
	}

	entries, err := minifluxCli.EntriesContext(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	output := make([]entry, len(entries.Entries))

	for i, e := range entries.Entries {
		markdown, err := htmltomarkdown.ConvertString(e.Content)
		if err != nil {
			log.Fatal(err)
		}

		output[i] = entry{
			ID:        e.ID,
			Title:     e.Title,
			URL:       e.URL,
			Content:   markdown,
			CreatedAt: e.CreatedAt,
		}
	}

	textContent, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(textContent),
				},
			},
		}, &getEntriesResult{
			Entries: output,
		}, nil
}
