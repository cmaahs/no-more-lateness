# no-more-lateness

A simple utility to go through your Google Calendar and find the next upcoming meetings, if it is within 5 minutes of start time, launch the online meeting tool

Clearly not super simple, as it turns out, and one of those things that has just made it to "It works for me" status.  Clearly a great deal of clean up needs to happen.  So enter the "we will see how it works" phase and hopefully I'll get time to work on it a bit more.

## TODO

I mistakenly placed the code to extract details from calendar items in the calendar provider.  Perhaps
a better way will be to pass a common format to a unit that is simply for scanning the text and extracting
the information.
