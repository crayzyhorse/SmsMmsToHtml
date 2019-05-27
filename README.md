Command Line Tool.

Convert xml to html with SMS and MMS media embedded.

The problem is that many backup apps do not save MMS media along with 
the SMS messages. The app (Google Play Store - SMS Backup & Restore) 
does this if you enable it under the app preferences. However it builds
an xml file that does not easily display well in a browser. This tool
means to remidy that by converting the xml file into an html file.

This version will only convert single contact xml files.  Attempting
to convert multiple contact backups will most likely make a mess of 
the html if it even parses it at all.

The CSS in template.html file can be edited to adjust the color, layout,
font, etc. of the html.

Under app preferences be sure to:
-Enable Add Readable Date
-Enable Add Contact Name
-Enable Include MMS Messages
-(Optionable) Include Emoji/Special Characters
Tested with xml file from version 9.74.1.


A new folder "html" is created to hold your new .html file.
Optional -n flag for setting the name or number you want to associate
with outgoing messages.

Example: SmsMmsToHtml.exe 2019-06-30.xml

Example: SmsMmsToHtml.exe -n=Jennifer 2019-06-30.xml

Example: SmsMmsToHtml.exe -n=432-555-8765 2019-06-30.xml

Concerning emojis and special characters.
The backup app saves these as html encoded entities which do not
always display well in a browser.  SmsMmsToHtml.exe does not attmpt
to decode these at this time (Perhaps in a future release) but does
put them in the html.  Often they show up as black diamond with a 
white question mark.
