# Hookbot

Hookbot is a simple Telegram bot that integrates webhook functionality to send PDF files to everyone that is subscribed to the bot functionality. 
The bot was written for a specific use case but I have tried to make it as general purpose as I could without going overboard.
Adapting the bot to send other kind of files should be trivial if you want to.
The bot can handles zipped files because it's been thought to be easily integrated with a service like Zapier/IFTTT. My use case it a Zapier task that receives an email with a PDF attachment and calls the webhook passing the zipped attachments as parameter (better leave the headache of handling emails to someone else).
Hookbot supports both direct conversations and being added to a group. In a group setting anyone can turn on/off the bot functionality which may not be what you want.  

## Configurations
A default configuration file is provided. The bot commands and confirmation messages can be customized from there. 
A configuration file (config.json) should be placed in /etc/hookbot or in the same folder as the executable. 
The bot API token **must** be placed in the configuration file otherwise the application will fail on start.

## Security 
The webhook is protected by secret keys that can be configured in the config file.  The authentication key should be sent in the *authkey* header. Just make sure you're hosting your webhook on HTTPS to avoid leaking the key(s).
I feel like this solution should be provide ”ok“ security without adding a lot of complexity. If you want to prevent brute force attacks on the key configuring fail2ban to prevent this should be trivial.

## Persitance
Subscribers ID are persisted in a SQLite database. The location of this file can be configured in the config file. SQLite has been chosen because nobody really wants to go through the trouble of setting up a DBMs for something so simple. If you already one set up and you want to used that switching from SQLite to that should't be more than a couple lines of code.

## Why Go
It's been almost two years since I written any go code (and it shows), so I wanted to pick it up again. I really liked the idea of having easy multithreading with channels and go routines. Something really similar could probably be written in node using Observables, but, I like how with Go you can easily produce an executable without having to worry about setting up node and installing npm dependencies.

## TODO
- Better modularization
- Add documentation (duh!)
- Extend functionality to handle more files types/text
- Add channel support

## Credits 
Numerous open source projects have made this bot trivial to develop. First and foremost [Telebot](https://github.com/tucnak/telebot) made it as easy as possible to get the bot up and running in no time. [Gin](https://github.com/gin-gonic/gin) really easy to setup the we server and to add a middleware for “authentication”. [Viper](https://github.com/spf13/viper) made it really easy to handle configuration files.