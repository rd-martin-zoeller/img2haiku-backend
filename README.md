# img2haiku-backend

A simple Google Cloud Function taking an image as an input and creating a Japanese haiku out of it.

## Why did you build this?

At university back in 2015, it took us weeks to build an algorithm that could detect nuts in a simple image. And it sucked. If you changed the image, the algorithm would fail. 10 years later, we can give an AI an image and tell it to write a poem about it. And it just works. That workflow would not have existed ten years ago. Nobody _needs_ it, but it's fun to see how far we've come.

TL;DR: I've built it because it's possible.

## How does it work?

The function needs a language string (like "English", "German", etc.) and a base64 JPEG image. This input is then sent to OpenAI's ChatGPT 4o along with a prompt instructing the AI to respond in a specific JSON format. ChatGPT's response is then interpreted as such JSON, sanitized, and returned to the caller.

This Google Cloud Function implementation is intended to be used with an iOS client from which people can upload their images. In a real-world scenario, the JWT used to authenticate against this API may be provided by a separate, small auth server that only issues tokens to legitimate clients. Such a validation may be based on Device Check or similar mechanisms.

## How to use the demo

1. Create an OpenAI API key
2. Clone the repository
3. In the directory of the repository, run `go mod download`
4. Optional: Copy the path to a photo or an image you want to create a haiku from
5. Edit the `server.sh` script and add your OpenAI API key
6. Run `./server.sh`
7. Copy the JWT it puts out
8. Edit the `client.sh` script 
    - Add your JWT
    - Optional: Add your image path
    - Optional: Update the language you want your haiku to be in
    - Optional: Add some tags if you want to get a haiku in a specific mood
9. Run `./client.sh`
