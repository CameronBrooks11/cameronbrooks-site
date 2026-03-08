# Content Model

## Principles

- Posts are loaded from Markdown files in `internal/content/writing/` at startup.
- No database and no request-time file IO; the embed is compiled into the binary.
- Slugs are derived from filenames and are canonical URL keys.
- `Body` is rendered from Markdown to HTML at load time; trust conversion happens in services.

## Adding a new post

1. Create a new file: `internal/content/writing/my-post-slug.md`
2. Write the frontmatter, a `---` separator, then the Markdown body:

   ```markdown
   title: My Post Title
   date: 2026-03-07
   summary: One sentence description shown in lists.
   tags: go, web
   published: true

   ---

   Your Markdown content here. Full CommonMark syntax supported.
   ```

3. Restart the server (`make dev` locally, `make deploy` for production).

The slug is always the filename without the `.md` extension.
Set `published: false` to keep a draft unroutable and out of all lists.

## Frontmatter fields

| Field       | Required | Format            | Notes                           |
| ----------- | -------- | ----------------- | ------------------------------- |
| `title`     | yes      | free text         |                                 |
| `date`      | yes      | `YYYY-MM-DD`      | Used for sorting and display    |
| `summary`   | yes      | one sentence      | Used in list view and meta tag  |
| `tags`      | no       | comma-separated   | e.g. `go, web, tools`           |
| `published` | yes      | `true` or `false` | `false` = draft; never routable |

## Post type

Fields available in templates via `PostView`:

- `Slug string`
- `Title string`
- `Summary string`
- `Body template.HTML` (rendered Markdown, trust-converted in services)
- `Tags []string`
- `Date string` (formatted: `January 2006`)
- `Published bool`

## Publication rules

- `published: false` posts are never listed or routable.
- Unknown slugs return `404`.
- Home page pulls the 5 most recent published posts.
- Posts are sorted newest-first everywhere.
