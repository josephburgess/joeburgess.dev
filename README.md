# joeburgess.dev

My personal website built with Go, serving as a lightweight snapshot of what I'm up to.

<p align="center">
    <img src="https://github.com/user-attachments/assets/49e4a867-a6c3-437a-ab7b-fb5156265f92" width="600" />
</p>

## Overview

This is quite a simple site showing recent Github activity. Built using Go with no frontend framework, just Go's templating ability.

## Features

- 🌓 Dark/light theme toggle with persistent preferences
- 📊 GitHub integration showing my recent repositories and activity
- 🌤️ Weather widget served by my own [Breeze API](https://github.com/josephburgess/breeze)
- 🔄 Data only refreshes every 15 mins to keep traffic down
- 📱 Responsive design

## Technical Details

- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, JavaScript
- **Hosting**: Deployed w/ Docker on a very minimal [DigitalOcean](https://www.digitalocean.com/) VPS (details in the Terraform setup [here](https://github.com/josephburgess/backstage))
- **Theme**: Rose Pine - same as my editor :)

## Development

If for some reason you wanted to run it locally:

```bash
git clone https://github.com/josephburgess/joeburgess.dev.git
cd joeburgess.dev
go run main.go
```

For the weather to work you'll need an API key for [breeze](https://github.com/josephburgess/breeze)... But to do that right now you'll need to use [gust](http://github.com/josephburgess/gust).

## Weather Widget

The site features a weather widget, mainly because I wanted to integrate it with [breeze](https://github.com/josephburgess/breeze), a lightweight API service I've set up for [gust](http://github.com/josephburgess/gust), another small project I'm working on. At the moment its showing current conditions in Buenos Aires, Argentina (where I'm currently based). For the widget to work running locally you'll need a breeze API key, for which you'll need to sign up through gust at the moment. I'm working on making this more accessible :)

## Future Plans

- I'm quite interested to see if I can find a way that doesn't suck to update the weather widget dynamically depending where I am in the world!

## Contact

Feel free to reach out to me at joe@joeburgess.dev or connect via the social links on the site.
