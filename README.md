To get the bot up and running:

1. Make sure you have set up a new application as a bot properly on https://discordapp.com/developers/applications/ and grab the token from the Bot tab.

1. Add the bot to your server using the OAuth tab.

1. Copy and customize the following files:
- token.example.txt -> token.txt

1. Run the bot to generate the default config. Configs are set up for VSCode Run, or you can manually:
```
go get ./...
go run main.go
```

To update packages and dependencies:
```
go get -u ./...
```
