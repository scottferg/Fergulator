Fergulator
==========
![alt text](https://secure.travis-ci.org/scottferg/Fergulator.png "Travis build status")

This is an NES emulator, written in Go. It's fairly new and very much a work in progress, so not all games run yet and not all features are implemented. Details are below.

![alt text](http://i.imgur.com/QGwdl.png "Metroid")

## To build on Linux

From your GOPATH:

        $ sudo apt-get install libsdl1.2-dev libsdl-gfx1.2-dev libsdl-image1.2-dev libglew1.6-dev libxrandr-dev
        $ go get
        $ go test
        $ go build

## To build on OSX

From your GOPATH:

        $ brew install sdl sdl_gfx sdl_image glew
        $ PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go get
        $ go test
        $ go build

## Run the emulator

        $ ./Fergulator path/to/game.nes

## Controls

        A - Z
        B - X
        Start - Enter
        Select - Right Shift
        Up/Down/Left/Right - Arrows

        Save State - S
        Load State - L

        Reset - R

        1:1 aspect ratio - 1
        2:1 aspect ratio - 2
        3:1 aspect ratio - 3
        4:1 aspect ratio - 4

        Emulate overscan - O
        Toggle audio - I

## Supported Mappers

* NROM
* UNROM
* CNROM
* MMC1
* MMC2
* MMC3
* ANROM

## Tested games that run well or are playable

[List is in the wiki](https://github.com/scottferg/Fergulator/wiki/Tested-Games)

## Next planned mappers

* MMC5
