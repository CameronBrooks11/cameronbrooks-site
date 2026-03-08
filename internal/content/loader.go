package content

import (
	"bufio"
	"bytes"
	"io/fs"
	"log/slog"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

func init() {
	posts, err := loadPosts(writingFS)
	if err != nil {
		slog.Error("failed to load writing posts", "err", err)
		return
	}
	Posts = posts
}

func loadPosts(fsys fs.FS) ([]Post, error) {
	matches, err := fs.Glob(fsys, "writing/*.md")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, path := range matches {
		raw, err := fs.ReadFile(fsys, path)
		if err != nil {
			return nil, err
		}
		post, parseErr := parsePostFile(path, raw)
		if parseErr != nil {
			slog.Warn("skipping post with parse error", "path", path, "err", parseErr)
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// parsePostFile parses a Markdown file with simple key: value frontmatter.
//
// Format:
//
//	title: My Post
//	date: 2026-01-15
//	summary: One sentence.
//	tags: go, web
//	published: true
//	---
//
//	Markdown body here...
//
// The slug is derived from the filename: "writing/my-post.md" → "my-post".
func parsePostFile(path string, raw []byte) (Post, error) {
	slug := strings.TrimSuffix(strings.TrimPrefix(path, "writing/"), ".md")

	meta := map[string]string{}
	var bodyLines []string
	inBody := false

	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		line := scanner.Text()
		if !inBody {
			if line == "---" {
				inBody = true
				continue
			}
			if k, v, ok := strings.Cut(line, ": "); ok {
				meta[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return Post{}, err
	}

	var date time.Time
	if ds := meta["date"]; ds != "" {
		if t, err := time.Parse("2006-01-02", ds); err == nil {
			date = t
		}
	}

	var tags []string
	if ts := meta["tags"]; ts != "" {
		for _, t := range strings.Split(ts, ",") {
			if t = strings.TrimSpace(t); t != "" {
				tags = append(tags, t)
			}
		}
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(strings.Join(bodyLines, "\n")), &buf); err != nil {
		return Post{}, err
	}

	return Post{
		Slug:      slug,
		Title:     meta["title"],
		Summary:   meta["summary"],
		Body:      buf.String(),
		Tags:      tags,
		Date:      date,
		Published: meta["published"] == "true",
	}, nil
}
