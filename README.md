# BMO
## What's this
BMO is a slack bot.

## Required package
Thank you nlopes@, for providing this great library.  
https://github.com/nlopes/slack  

## Env
BMO requires three environment variables.

- SLACKTOKEN:  
you can find it in your app's settings under *Install App* > *Bot User OAuth Access Token*  
It starts with "xoxb-".  
- VTOKEN:  
It's in *Basic Information* setting, showed as "Verification Token".  
- BOTUNAME:  
User name of the bot, starts with "U-"  

## Functions
### vote
**Syntax**
```
name++ <discription>
```

Increment *name*'s the number of votes.  

## ToDo
- combine multiple increments for the same name 🤝
- decrement 👎
- be unit testable 📝
- add new functions 🆕

## Done
- enable counting number of votes 🔢