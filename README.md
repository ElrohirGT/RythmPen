# RythmPen

## How to play?

You can technically play it on [itchio](https://elrohirgt.itch.io/rythmpen). But
I don't recommend it. For some reason the music doesn't work great there...

```bash
# The command is: go run ./cmd/player/ -audio <path to audio> -map <path to map>
# Easy difficulty
go run ./cmd/player/ -audio "./songs/easy_Joyful, Фрози, Zachz Winner - Boogie [NCS Release].mp3" -map "./songs/easy_Joyful, Фрози, Zachz Winner - Boogie [NCS Release].map" -duration 30s
```

```bash
# The command is: go run ./cmd/player/ -audio <path to audio> -map <path to map>
# Medium difficulty
go run ./cmd/player/ -audio "./songs/medium_Willie Colon - Oh Qué Será_ (Letra Oficial).mp3" -map "./songs/medium_Willie Colon - Oh Qué Será_ (Letra Oficial).map" -duration 1m14s
```

## How to generate a new map for a song?

```bash
# The command is:
# The generated map will be saved on -dst path
go run ./cmd/recorder/ -src <path to audio> -dst <path to map>
```

## How to build for the web?

```bash
GOOS=js GOARCH=wasm go build -o ./cmd/itchio/build/game.wasm ./cmd/itchio && rm ./cmd/itchio/game.zip; zip -r ./cmd/itchio/game.zip ./cmd/itchio/build/
```

## Songs copyright

```
Song: Zachz, Фрози, Joyful - Boogie
Music provided by NoCopyrightSounds
Free Download/Stream: http://ncs.io/Boogie
Watch: http://ncs.lnk.to/BoogieAT/youtube
```
