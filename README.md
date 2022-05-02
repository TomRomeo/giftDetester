# Gift Detester
A Discord bot to remove malicious discord gift links

## Invitation link
[Invite the bot to your server](https://discord.com/api/oauth2/authorize?client_id=939443581603696650&permissions=10243&scope=bot%20applications.commands)  

Alternatively, you can also host this bot yourself!

---

Lately, there have been many scams where hackers get access to Discord accounts and use them to send phishing links of supposed Discord Nitro gifts in every server.
Sadly, a lot of people click on those links by which they as well give those hackers access to their account, creating a vicious cycle.
Most of the time, users falling for this scam don't even realize it until they see messages of their account that they did not send.

<img src="https://i.imgur.com/TbRzMMr.png" alt="twitter post from Hylian Wolf" width="400" />

## Gift De(tester) to the rescue!
I built this bot to mitigate these phishing link attacks.
After you add the bot to your server, it will start analysing all links contained inside new messages.
When a message link looks similar to a discord link (for example `dlscord.com`, where there is a lowercase L in place of an i) but isn't **actually** from discord,
the bot will delete the message and take the configured action. You can also configure it to log the action to a logchannel

## So what exactly happens when a member sends a malicious link?
1. The link get's deleted
2. The user receives a dm, informing them that their account has probably been compromised and how to get it back.
3. The user get's either kicked or timeouted from the server

<img src="https://i.imgur.com/BostBwR.png" alt="twitter post from Hylian Wolf" width="400" />

## Configuration
You can configure different actions for what to do with the compromised user
1. Kicking the user, sending them an invite to re-join after they reset their password (default)
2. Timeouting the user for a configured period of time, giving them the possibility to reset their password without kicking them from the server.

## Commands

| Command | Description |
|-|-|
| /phishing action <kick\|ban> | Configure what the bot will do to a compromised user |
| /phishing logs #channel | Set a logchannel |
| /phishing duration | Set the timeout duration if enabled |
