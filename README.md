# [joeburgess.dev](https://joeburgess.dev)

My personal website built with Go, serving as a basic snapshot of what I'm up to. Its pretty lightweight and just shows some GitHub activity. I didn't use a frontend framework, just Go's `html/template` package.

<p align="center">
    <img src="https://github.com/user-attachments/assets/49e4a867-a6c3-437a-ab7b-fb5156265f92" width="600" />
</p>

## Features

- 🌓 Dark/light theme (rose pine)
- 📊 GitHub integration showing recent repos/activity
- 🌤️ Weather served by my own [breeze api](https://github.com/josephburgess/breeze)
- 🔄 Data on periodic refresh every 15mins (rather than on refreshing page - keep that traffic down!)
- 📱 Responsive

## Technical Details

- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, JavaScript
- **Hosting**: Deployed w/ Docker on a very minimal [DigitalOcean](https://www.digitalocean.com/) VPS (details in the Terraform setup [here](https://github.com/josephburgess/backstage))

## Development

If for some reason you wanted to run it locally:

```bash
git clone https://github.com/josephburgess/joeburgess.dev.git
cd joeburgess.dev
go run main.go
```

## Weather Widget

The site features a weather widget, mainly because I wanted to integrate it with [breeze](https://github.com/josephburgess/breeze), a lightweight API service I've set up for [gust](http://github.com/josephburgess/gust), another small project I'm working on. At the moment its showing current conditions in Buenos Aires, Argentina (where I'm currently based). 

For the widget to work running locally you'll need a breeze API key, for which you'll need to sign up through gust at the moment. I'm working on making this more accessible :)

## Future Plans

- I'm quite interested to see if I can find a way (that doesn't suck) to update the weather widget dynamically depending where I am in the world!

## Contact

Feel free to reach out to me at joe@joeburgess.dev or connect via the social links on the site.
