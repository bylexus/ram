# RAM

> My brain's external RAM extension

This little project fills an a annoying gap in my daily workflow: When I'm reading news or working on a task,
I occasionally want to note something down to examine later: an URL, a random thought, a tool a command line etc.
But where should I put that small piece of information? In my notes tool? Yeah, not available on my mobile. As Google Keep note?
nah, I never gonna find it... 

So I need a tool to simple and fast take a note and examine it later. This should be this project.


## MVP ideas

* A Note consists of:
  * a texfield
  * an URL
  * some tags
  * a done checkbox (boolean)
* Notes should be kept in a simple SQLite DB
* A small Go Webserver should offer a UI to enter / list notes
* Simple, fast, fast and easy to develop - a first version must not be sophisticated! It should be done soon enough!
* First version: only 1 user, no (or basic) auth

### MVP Developing Path

1. create the db part:
   * find sqlite golang driver
   * create a structure to run db migrations on startup
   * create simple utility functions for the app
2. create the static web server:
   * should serve static files to deliver the basic web page
   * The web page must be minimal - html and some tiny ajax functions
3. create backend endpoints to interact with frontend
4. create authentication

## More ideas

* Authenticated by Passkey (webauth) only!
* Multiple users - each with its own note pool
* Notes export
* Notes API


