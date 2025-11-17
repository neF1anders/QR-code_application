# Application QR-generator
The application became my mini-project for consolidation of skills in programming on Golang.

## Content
- [Idea](#Idea)
- [Technologies](#Technologies)
- [Usage](#Usage)
- [Requirements](#Requirements)
- [Approach](#Approach)
- [Issues](#Issues)

## Idea
The main idea is to make a generator of QR-code for fast transfer of text information of text to the smartphone. 
It was a problem for me, when i need to continue watching smth or reading an article on the way from my PC and there is a question "how to transfer the link?". I wanted to implement a transfer without disturbing effect of social networks. Just to copy the link. Only text.
That`s why my idea is to copy the link, tap the button, scan the generated QR-code. Simple problem - simple solution.

## Technologies
- [Golang](https://go.dev/)
- [clipboard](github.com/atotto/clipboard)
- [qrcode lib](github.com/skip2/go-qrcode)
- [hotkey` logging lib](golang.design/x/hotkey)

## Usage
The usage of this application is almost as complicated as its implementation due to some issues, that must be resolved, but they are not yet. 
- Download the code in go file
- Compile it with the libs, downloaded
- Use .exe to generate the code once

### Requirements
See [Technologies](#Technologies)

## Approach
My implementation consists of such steps as:
- Logging the hotkey command (Ctrl + F1 in my code)
- Opening the localhost for listen-serve
- Generating an image of QR-code with lib and serving to localhost
- Closing the host after 30 seconds

## Issues
My application is not ideal at all. I need to tackle these problems:
- The application works once, generates the code and don`t read the hotkeys until the application is reloaded
- The application never stops consuming the RAM after finishing the generation. The architecture must be fixed
- The application has no settings for setting the hotkeys, the 'background' mode and so on
- No autoload with system

There are also some ideas to decorate the application:
- Autoload with system as a auto-enabled function
- API or the extension mode for Chrome (i don`t now is it possible)
- The simple interface to interact with app not by writing code
