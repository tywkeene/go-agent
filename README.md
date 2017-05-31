# go-agent
 
## What is it?
Go-agent is an automation engine, the agent-server sits on the cloud (AWS) and agent-clients log in and ping the server.
What's so special about that? Nothing really. The cool part is hooks.

## Okay. what are hooks?
Think about any modern software development workflow. You push to github, your tests run on travis, it emails some people,
and does anything else that needs to be done, automatically.

Hooks in Go-Agent allow you to automate tasks on machines registered with the agent-server.

## Why would I need this?
Say you have a laptop, and a desktop. You have your gpg keys, your ssh keys, and any other files you wouldnt want lost
or stolen, you don't necessarily want them on your laptop, since it could be lost or stolen, but you need them to do any work. 
Go-Agent would allow you to automatically download these files over https from a secure and encrypted server on login, 
fire up sshd on your desktop, and you're good to go.

After you're done you don't even need to do anything, just logout as you normally would, and the agent-client would securily wipe
your files without you even needing to think about it.

## How is it done?
All of this would be done via hook scripts run by go-agent in response to device events. Super simple, super easy.

This is just one of the many possible use cases of Go-Agent, and it's really up to you, and since you have access to your shell
or scripting language of choice.

## So how do I get and run it?
More on that coming soon, since Go-Agent is still under heavy development

Feel free to contact me with any questions, comments or suggestions or if you would like to contribute to the project.

Cheers.
