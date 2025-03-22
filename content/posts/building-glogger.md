# first post! glogger: a simple go blog engine

2025-03-22

Well, welcome to this blog I guess! I have lately built a few projects - my personal site included - in go. I've found it super fun to code in, it works out of the box and feels like a happy middle ground between very high level languages and lower level systems languages. I have now become obsessed with learning rust but I also wanted to add a blog to my site. I really wanted it to be as barebones as possible - there is some great software out there but I found it all a bit overengineered, so I decided I'd do one last go side project before getting really deep into finishing the rust book and trying to build something in rust!

## MVP

I wanted something super lightweight that would be easy enough to implement without relying on a full CMS or complex external tools. My initial mvp outlined four things I wanted to achieve:

1. Parse simple markdown posts - no CMS just files stored in the directory
2. Support a few themes
3. Simple integration - just plugged into the core router
4. Lightweight/fast

## Key Features

### Simple Integration

Since this was originally just going to be built INTO my personal site, making it into a package instead I really needed it to integrate simply/with a few lines of code. Here's all it takes to add blogging to your Go site:

```go
// Create a new blog with default settings
blog, err := glogger.New(glogger.Config{
  ContentDir: "content/posts",  // Where your markdown files are stored
  URLPrefix: "/blog",           // URL prefix for the blog routes
  Theme: glogger.ThemeRosePine, // Optional theme selection
})

// Register with your router
blog.RegisterHandlers(router)
```

With this minimal configuration, glogger sets up these routes:

- `/blog` - List of all blog posts
- `/blog/{slug}` - Individual post pages
- `/blog/_themes/{theme}.css` - Serves the theme `.css` files

### Go features that I found useful

#### 1. Embedding

I was quite pleased with how simply the embed package works to stitch files directly into the binary:

```go
//go:embed assets/templates/*.html
var templatesFS embed.FS

//go:embed assets/themes/*.css
var themeFS embed.FS
```

This ensures that the templates/css etc are always available at runtime. Before discovering this I was having real trouble allowing users to define their theme and the blog kept loading without the styling.

#### 2. File walking

Finding and parsing all blog posts is handled through the [filepath.Walk](https://pkg.go.dev/path/filepath#Walk) function:

```go
err := filepath.Walk(b.config.ContentDir, func(path string, info fs.FileInfo, err error) error {
    if err != nil {
        return err
    }

    if info.IsDir() || !strings.HasSuffix(path, ".md") {
        return nil
    }

    post, err := parsePost(path)
    if err != nil {
        return err
    }

    // Process post...
    return nil
})
```

I found this to be a really nice way to iterate over all the files and apply the parser to them without lots of boilerplate.

#### 3. Markdown parsing

Sadly this is not a zero dependency package! I used [goldmark](https://github.com/yuin/goldmark). It felt like too much of an undertaking to write a markdown parser from scratch! Especially when goldmark is well maintained, well loved, lightweight, FOSS software.

```go
md := goldmark.New(
    goldmark.WithParserOptions(
        parser.WithAutoHeadingID(),
    ),
)
var buf bytes.Buffer
if err := md.Convert([]byte(filteredContent), &buf); err != nil {
    return Post{}, err
}
```

#### 4. Routing

Another dependency :( ... for now I've relied on [gorilla/mux](https://github.com/gorilla/mux) for the routing. This was simply because its what I already had used for my personal site.

```go
func (b *Blog) RegisterHandlers(router *mux.Router) {
    blogRouter := router.PathPrefix(b.config.URLPrefix).Subrouter()

    blogRouter.HandleFunc("/", b.handleListPosts).Methods("GET")
    blogRouter.HandleFunc("/{slug}", b.handleSinglePost).Methods("GET")
    blogRouter.HandleFunc("/_themes/{theme}.css", b.handleThemeCSS).Methods("GET")
}
```

This is a bit of a limitation though, since users of glogger must be using mux as their router. I'd like to make the package router agnostic by:

1. Creating a router interface to allow for different router implementations
2. Use this to provide adapters for other popular routing packages, priority being the standard library!

## Overall

Quite happy with how this all turned out! Its a super simple/lightweight package that works well for me (at least!). I'm not much of a front-end developer, and am not very design oriented, but I've found go's templating to be quite intuitive and provided all I needed for this simple/fast implementation.

## Other plans

I'm probably going to come back to add the odd feature as time goes on and I find things that don't work quite how I wanted. A few things on the list:

- Syntax highlighting - not looked into this yet but I don't like my code snippets being all one colour!
- Pagination (will probably come once I write enough blog posts to need it!)
- As mentioned above, better router support
- Metadata support - bit more of a stretch goal especially with no database
