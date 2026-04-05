# [joeburgess.dev](https://joeburgess.dev)

My personal website built with Go, serving as a basic snapshot of what I'm up to. Its pretty lightweight and just shows some GitHub activity. I didn't use a frontend framework, just Go's `html/template` package.

<p align="center">
    <img src="https://github.com/user-attachments/assets/2442bb25-0fc1-4049-a1f0-42a17be08ee4" width="600" />
</p>

## Features

- Dark/light theme (rose pine)
- GitHub integration — recent repos and activity
- Weather via my own [breeze API](https://github.com/josephburgess/breeze)
- Background data refresh (hourly, not on page load)
- Blog powered by [glogger](https://github.com/josephburgess/glogger)
- Responsive

## Stack

- **Backend**: Go with stdlib `net/http`
- **Templating**: Go `html/template`
- **Blog**: [glogger](https://github.com/josephburgess/glogger) — my own lightweight markdown blog engine
- **Hosting**: Docker on a DigitalOcean VPS (infra via [backstage](https://github.com/josephburgess/backstage))

## Running locally

If for some reason you wanted to run it locally:

```bash
git clone https://github.com/josephburgess/joeburgess.dev.git
cd joeburgess.dev
cp .env.example .env  # add your BREEZE_API_KEY if you want the weather widget
go run main.go
```

## Weather Widget

I added a widget mainly because I wanted to integrate it with [breeze](https://github.com/josephburgess/breeze), a lightweight API service I've set up for [gust](http://github.com/josephburgess/gust), another small project I'm working on. I am now based back home in London, so that's where it shows the weather for.

For the widget to work running locally you'll need an API key. You can either install gust and use the one-click signup to get a key, or use your own OpenWeatherMap 3.0 key.

## Blog

The site now includes a blog powered by [glogger](https://github.com/josephburgess/glogger), a lightweight blog engine package I built in go. It supports simple markdown content (no database), multiple themes, and simple integration with existing go sites.

## Future Plans

I'm quite interested to see if I can find a way (that doesn't suck) to update the weather widget dynamically depending where I am in the world! I also have a few things I want to add to [glogger](https://github.com/josephburgess/glogger) too - you can see the vague roadmap in the project's README.

## Contact

Feel free to reach out to me at hi@joeburgess.dev or connect via the social links on the site.
