# Bot

The Twitch bot itself

## Running locally

* Make a `.env` file based on `example.env`
* Get an OAuth secret from `https://twitchapps.com/tmi/`
* Add commands
    * To create an invokable command (i.e., activated by typing `!hello` in chat or something similar) create a file
      called `command_name.command.json` based on the example files given
    * To create a message that sends after a certain amount of messages, create a file
      called `command_name.interval.json` based on the example files given. **NB:** These messages aren't guaranteed to
      send. There is a ~30 second limit on each command so if you had an interval message set to send every 10 messages
      and then received 40 messages in 20 seconds, the interval message would not be sent twice. Note also that the
      bot's responses in chat do not count towards the message count
* Run the command `go run .`

## Reserved keywords

You cannot use the following keywords as parameter names in commands

* `username`
    * The username of the user who invoked the command

## Open Source Libraries Used 

### [go-twitch-irc](https://github.com/gempir/go-twitch-irc)

Used for connecting to and interacting with Twitch chat

MIT License

### [godotenv](https://github.com/joho/godotenv)

Used to load environment variables

MIT License
