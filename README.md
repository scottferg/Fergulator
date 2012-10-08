Fergulator
==========
![alt text](https://secure.travis-ci.org/scottferg/Fergulator.png "Travis build status")

This is an NES emulator, written in Go. It's fairly new and very much a work in progress, so not all games run yet and not all features are implemented. Details are below.

![alt text](http://i.imgur.com/QGwdl.png "Metroid")

## To build on Linux

        $ sudo apt-get install libsdl1.2-dev libsdl-gfx1.2-dev libglfw-dev libglew1.6-dev libxrandr-dev
        $ go get -u github.com/0xe2-0x9a-0x9b/Go-SDL/sdl
        $ go test
        $ go build

## To build on OSX

        $ brew install sdl sdl_gfx glfw glew
        $ PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go get -u github.com/0xe2-0x9a-0x9b/Go-SDL/gfx
        $ PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go get -u github.com/banthar/gl
        $ PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go get -u github.com/jteeuwen/glfw
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

* Super Mario Bros
* Super Mario Bros 2
* Contra
* Bionic Commando
* Blaster Master
* Double Dragon
* Final Fantasy
* Tecmo Bowl
* Ninja Gaiden
* Legend of Zelda
* Metroid
* Tetris
* Balloon Fight
* Adventure Island
* Castlevania
* Castlevania 2
* Dig Dug
* Donkey Kong
* Donkey Kong Jr.
* Duck Tales
* Duck Tales 2
* Excitebike
* Galaga
* Gun Smoke
* Hydlide
* Ice Climber
* Ice Hockey
* Kung Fu
* Lode Runner
* Mega Man
* Mega Man 2
* Mega Man 3
* Mega Man 4
* Mega Man 5
* Mega Man 6
* Metal Gear
* Prince of Persia
* A Boy and His Blob
* Snake Rattle 'n' Roll
* Bart vs. The Space Mutants
* Kid Icarus
* Mighty Bomb Jack
* Bubble Bobble
* Adventures of Link

## What isn't working

* Sound
* Second controller
* Scrolling and palettes on a number of MMC1 games
* Save states for some MMC1 games
* Some minor graphical glitches on screen boundary

## Next planned mappers

* MMC3
* MMC5
* MMC2
