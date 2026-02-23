package mcp

import "time"

type getEntryParams struct {
	EntryID int64 `json:"entryId"`
}

type getEntriesParams struct {
	Status          string    `json:"status"`
	Search          string    `json:"search"`
	Limit           int       `json:"limit"`
	Category        string    `json:"category"`
	Feed            string    `json:"feed"`
	PublishedAfter  time.Time `json:"publishedAfter"`
	PublishedBefore time.Time `json:"publishedBefore"`
}

type entry struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type getEntryResult struct {
	Entry entry `json:"entries"`
}

type getEntriesResult struct {
	Entries []entry `json:"entries"`
}

type feed struct {
	ID    int64  `json:"id"`
	Title string `json:"string"`
	URL   string `json:"url"`
}

type getFeedsResult struct {
	Feeds []feed `json:"feeds"`
}
