[![Imgur](https://i.imgur.com/z7wBEwj.png)](https://i.imgur.com/z7wBEwj.png)
:trollface:

## don't let your memes be dreams
![erni](erni.jpg)

## how to use
# DON'T

## i insist, how to use?
- either build from the tools and run it on dir you want to change
- put as _ import for side effect, but changes is only reflected after subsequent run

## why no test
[![Imgur](https://i.imgur.com/AlKoUmy.jpg)](https://i.imgur.com/AlKoUmy.jpg)

i have no freaking idea how to test this at the moment, the fact that it works is already a magic to me

## roadmap
- works on init (some sort of JIT change maybe idk i don't understand this thing)
as of now initial `go run`/`go build` does not reflect changes, but subsequent one does. need to look under the hood.
- works on imported packages
as of now `erni` only works on function from the same package
- works based on return types of function
as of now `erni` only fixes
- works on whole repo
as of now `erni` only works on file in dir
- tweakable
with comment or a ala `go generate` maybe

### disclaimer
i have every reason to believe that this tool **will break** your code
