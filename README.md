Fergulator
==========

This is an NES emulator, written in Go. It's fairly new, so not all games run yet and not all features are implemented. Details are below.

## To build on Linux

        $ sudo apt-get install libsdl1.2-dev libsdl-gfx1.2-dev libsdl-image1.2-dev libsdl-mixer1.2-dev libsdl-sound1.2-dev libsdl-ttf2.0-dev
        $ go get -u github.com/0xe2-0x9a-0x9b/Go-SDL/...
        $ go test
        $ go build

## To build on OSX

        $ brew install sdl sdl_image sdl_sound sdl_gfx sdl_mixer sdl_ttf
        $ export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go get -u github.com/0xe2-0x9a-0x9b/Go-SDL/...
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

## Tested games that run flawlessly or near flawlessly

* Super Mario Bros
* Contra
* Bionic Commando
* Blaster Master
* Double Dragon
* Final Fantasy
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
* Metal Gear
* Prince of Persia

## What isn't working

* Sound
* Second controller

## Next planned mappers

* MMC3
* MMC5
* MMC2
