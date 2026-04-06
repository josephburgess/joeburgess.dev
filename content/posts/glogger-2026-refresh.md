---
title: "glogger in 2026: rss, stdlib routing and a clean up"
date: 2026-04-06
description: "decided to use some of my easter hols to dust off my golang"
tags: ["go", "glogger"]
---

It's been a while since I first built glogger and wrote about it here. The package has sat quietly, used exclusively by this site and parsing its extremely low number of posts in the background, but I had some extra time off over Easter so I've gone back and done a bit of a refresh.

## killing gorilla/mux

The original version of glogger used [gorilla/mux](https://github.com/gorilla/mux) for routing since I'd already used it for my personal site and was familiar with it. Coming back though I found it had been archived and is no longer maintained. Thankfully it also seems go 1.22 has added method+path pattern matching directly to the stdlib `ServeMux` so a replacement didn't take long to find.

The migration was pretty clean, patterns were similar (arguably more readable tbh), and one less dependency which is an extra win when trying to keep the package as lightweight as possible!

```go
// before (gorilla/mux)
func (b *Blog) RegisterHandlers(router *mux.Router) {
    blogRouter := router.PathPrefix(b.config.URLPrefix).Subrouter()
    blogRouter.HandleFunc("/", b.handleListPosts).Methods("GET")
    blogRouter.HandleFunc("/{slug}", b.handleSinglePost).Methods("GET")
}

// after (stdlib)
func (b *Blog) Handler() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", b.handleListPosts)
    mux.HandleFunc("GET /{slug}", b.handleSinglePost)
    return mux
}
```

### syntax highlighting

The original post mentioned syntax highlighting as a feature I wanted to add at some point. I initially looked at using [goldmark-highlighting](https://github.com/yuin/goldmark-highlighting) using [chroma](https://github.com/alecthomas/chroma) but I found it was a bit temperamental with golang, not to mention adding another dependency and more boilerplate.

I insteadopted for using [highlight.js](https://highlightjs.org) loaded from a CDN, which means no additional go dependencies.

Each glogger theme has a sensible default syntax theme but you can override it with any highlight.js theme name from their [examples](https://highlightjs.org/examples):

```go
glogger.Config{
    Theme:       glogger.ThemeRosePine,
    SyntaxTheme: "tokyo-night-dark",
}
```

### yaml frontmatter

Posts now require YAML frontmatter. Previously glogger was not doing much about post metadata and I just had it grabbing a title from the first heading. Now posts look like this:

```markdown
---
title: "my post title"
date: 2026-01-15
description: "shown in post list and RSS feed"
tags: [go, blogging]
draft: false
---

content here...
```

The filename (minus `.md`) becomes the URL slug. Setting `draft: true` hides the post and prevents it being served. Tags are filterable and clicking one opens a filtered listview.


### rss 2.0

RSS isnt something I've used extensively but was actually very simple to add. Just needed a new `GET /feed.xml` route, three new config fields, and some `encoding/xml` structs. There's also now an RSS link in the heading.

```go
glogger.Config{
    Title:       "My Blog",
    Description: "writing about things",
    BaseURL:     "https://example.com",
}
```


### mount helper

Previously mounting glogger on your mux looked like:

```go
 blog, err := glogger.New(glogger.Config{
    URLPrefix: "/blog",
}

prefix := blog.URLPrefix()
mux.HandleFunc("GET "+prefix, func(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, prefix+"/", http.StatusMovedPermanently)
})
mux.Handle(prefix+"/", http.StripPrefix(prefix, blog.Handler()))
```

now its:

```go
blog.Mount(mux)
```

`URLPrefix` is set once in `Config` and `Mount` handles everything — including the redirect from `/blog` to `/blog/` which was easy to forget.

## still on the list...

I still want to add pagination, and perhaps make theming more intuitive, but I'm not much of a frontend person. It's been really nice to pick up go agian after a long-ish hiatus. Its a relatively intuitive language but doesn't hold your hand.
