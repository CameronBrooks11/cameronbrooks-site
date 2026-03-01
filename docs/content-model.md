# Content Model

## Principles

- MVP content is code-backed in `internal/content`.
- No database and no request-time file IO.
- Slugs are canonical URL keys.
- `Body` is stored as `string`; trust conversion happens in services.

## Types

### `Project`

Fields:

- `Slug string`
- `Title string`
- `Description string`
- `Body string` (raw HTML)
- `Tags []string`
- `Date time.Time`
- `Links []Link`
- `Featured bool`

### `Post`

Fields:

- `Slug string`
- `Title string`
- `Summary string`
- `Body string` (raw HTML)
- `Tags []string`
- `Date time.Time`
- `Published bool`

## Publication rules

- `Published=false` posts are never listed or routable.
- Unknown slugs return `404`.
- Home page pulls only featured projects and recent published posts.

## Editing workflow

1. Update `internal/content/data.go`.
2. Keep slug uniqueness and URL-safe format.
3. Run:

   ```sh
   go test ./...
   go vet ./...
   go build ./...
   ```

4. Smoke locally (`make dev` + `make smoke`).
5. Deploy.

## Future migration path

If content volume grows:

- Load markdown files via embed at startup
- Render once to HTML
- Populate same `Project`/`Post` view source

This keeps handler/template contracts unchanged.
