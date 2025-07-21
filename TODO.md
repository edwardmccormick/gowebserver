Todos:

Front-end:

[ x ] a signup flow that isn't (as) overwhelming
 - I think this is going to end up being, at least at first, a carousel with detail 'buttons' - pick one and go to the next screen, with maybe a back button enabled but no 'forward' button? Or maybe we let you go back and forward, but alert before you submit without every field filled. Working concept
[ x ] a split in the signup flow between a 'basic' profile and a 'complete' or 'advanced' profile, with all the metrics
[ x ] the ability to update your profile
- This can probably be almost identical as the signup flow, except we show you what you currently have
[ x ] upload photograph user flow, if not an actual implementation
- Actual implementation - Presigned URL to upload files, S3 keys get saved, presigned URLs to display the image files, everybody goes home happy.
- Drag and drop functionality, add from your photo gallery on phone (this depends a lot on the implementation for )
[ ] - Crop a photo to get a square 'profile' picture, otherwise it's going to look weird lol
[ ] the ability to sign up from Facebook, or Google, or Apple, or....LinkedIn?!?!?! Github!??!?! lol
- This might be the most easy to implement, although probably some troubling possibilities like, long term.
[ ] - password reset/account recovery flow
[ ] the ability to like photographs, or prompts
[ ] the ability to like prompts
[ ] the ability to send a message with a like, in response to a prompt or photo
[ ] The element and message that the person 'liked' when deciding whether to respond to the person or not
[ ] sane defaults for finding 'matches'
[ ] Some intutive way of displaying the information from 'details' or the 'advanced' profile
[ ] Some intuitive way of displaying 'metrics' about the user - how long, how many matches, how many chats, etc.
[ ] view matches by geographic location on Mapbox (or something similar, although let's be honest, it's going to be Mapbox)
[ ] Long term, a total rework/groom of front-end display elements. Not necessarily logic and functionality, but just how everything *looks*, to incorporate more brand elements. And a snazzier ui, frankly. 😒
[ ] emoji keyboard lol
[ ] front-end Bedrock->Lambda->ALB POC, in lieu of using a Gemini key.

Back-end:

[ ] Data structures to support above, specifically prompts
[ ] Database connection and data persistance
[ ] Database seed/scorch functionality
[ ] Data persistance for chat
[ ] Password reset, including email user(?) - this might actually be easier once we go cloud-native
[ ] Chat monitoring for urls, email@addresses or email (at) addresses, and...addresses?
[ ] First message monitoring. No swear words, no racial slurs, that kind of thing.
[ ] Chat monitoring, with *some* of that functionality.
[ ] geo-encoding through zip codes
[ ] some sort of 'share your location for better matches' kind of functionality
[ ] Gemini Flash 'create a date' functionality. Although actually, I would rather that be a front-end function, but I don't think there's going to be a way to use an API key from the front end without exposing it, and....exposing it would be *not good*.
[ ] Could probably do it from the front end with a lambda and an ALB/lambda address or something. Might be fun to mess with. Probably need POCs all directions there.
[ ] Gemini Flash 'ice breaker' functionality.

Architectural:

[ ] A not-stupid implementation of user stats
[ x ] Containerize everything, and docker compose to spin up the right databases with the right connection strings
[ x ] some idea of how well react will scale to an iOS/Android app, and what those might possibly look like.
[ ] Postgres vs MySQL, although this will ultimately come down (like most of my decisions) to which is cheaper. Picking up pennies in front off the train tracks!
[ ] Alerts, at least something besides server-side events, since those are pretty much a no-go with serverless. Message queue/delivery service?
[ ] Probably some hard decisions about vibe chats, for or against
[ ] Host a bedrock model, or use an API key for AI? And a rough cost estimate for each. Although I'm actually willing to die a little on this hill...I think the AI element makes things 'nicer' in a way that will resonate
[ ] How to stop scammers, spammers, catfish, et. al.

Operational:

[ ] Some idea of what the elements of a trust and safety 'policy' might look like
[ ] Some plan for spreading the word/advertising/word of mouth. Reddit and 'stickers' is not the answer. 'If you build it, they will come' is really not an answer, even if it's the best one I've got at this point.
[ ] Terms of service
[ ] Privacy statement, or whatever other legally required documents we're going to need. These are more 'go live' than 'alpha' pieces, but still....they might be easy wins.
[ ] Would I consider a cofounder? Would I consider an angel? Would I consider some sort of equity arrangement to help us get 'off the ground'? Do I have the stones like those little Asian girls from Coffee Meets Bagel to tell Mark Cuban (rightly!) 

#FFFFFF - White
#000000 - Black
#FF0000 - Red
#0E6BA8 - Bice Blue
#717744 - Reseda Green
#ED8B00