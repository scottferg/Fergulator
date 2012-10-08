Fergulator
==========
![alt text](https://secure.travis-ci.org/scottferg/Fergulator.png "Travis build status")

This is an NES emulator, written in Go. It's fairly new and very much a work in progress, so not all games run yet and not all features are implemented. Details are below.

![alt text](http://i.imgur.com/QGwdl.png "Metroid")

## To build on Linux

        $ sudo apt-get install libsdl1.2-dev libsdl-gfx1.2-dev libsdl-image1.2-dev
        $ go get -u github.com/0xe2-0x9a-0x9b/Go-SDL/sdl
        $ go test
        $ go build

## To build on OSX

        $ brew install sdl sdl_image sdl_gfx
        $ PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go get -u github.com/0xe2-0x9a-0x9b/Go-SDL/sdl
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

## Supported Mappers

* NROM
* UNROM
* CNROM
* MMC1
* MMC3

## Tested games that run well or are playable

[List is in the wiki](https://github.com/scottferg/Fergulator/wiki/Tested-Games)

## What isn't working

* Sound
* Second controller
* Scrolling and palettes on a number of MMC1 games
* Save states for some MMC1 games
* Some minor graphical glitches on screen boundary

## Next planned mappers

* MMC5
* MMC2
