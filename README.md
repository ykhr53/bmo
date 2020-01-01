# BMO
[![Build Status](https://travis-ci.org/ykhr53/bmo.svg?branch=master)](https://travis-ci.org/ykhr53/bmo)
[![codecov](https://codecov.io/gh/ykhr53/bmo/branch/master/graph/badge.svg)](https://codecov.io/gh/ykhr53/bmo)

## What's this
BMO is a slack bot.

## Required package
Thank you nlopes@, for providing this great library.  
https://github.com/nlopes/slack  

and some others... please see go.mod.

## Environment variables
BMO requires three environment variables.

- **SLACKTOKEN**  
you can find it in your app's settings under *Install App* > *Bot User OAuth Access Token*  
It starts with "xoxb-".  
- **VTOKEN**  
It's in *Basic Information* setting, showed as "Verification Token".  
- **BOTUNAME**  
Unique user ID of the bot, starts with "U"  

## Functions
### vote
**Syntax**

Increment *name*'s the number of votes.
```
name++ <discription>
```

Decrement *name*'s the number of votes.
```
name-- <discription>
```

BMO can hook and combine multiple votes.
```
name++ name-- name++ foo++ <discription>
```

## ToDo
- add new functions 🆕
- add ToDo things 🤔

## Done
- enable counting number of votes 🔢
- combine multiple increments for the same name 🤝
- decrement 👎
- be unit testable 📝