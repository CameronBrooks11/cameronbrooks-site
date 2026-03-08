package content

// Posts is the full list of writing entries, loaded from content/writing/*.md at startup.
// Draft posts (Published: false) are never exposed via any route.
// To add a post: create a new .md file in internal/content/writing/ and restart the server.
var Posts []Post
