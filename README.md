Fergulator
==========
[![Build Status](https://travis-ci.org/scottferg/Fergulator.png?branch=master)](https://travis-ci.org/scottferg/Fergulator)

This is an NES emulator, written in Go. It's fairly new and very much a work in progress, so not all games run yet and not all features are implemented. Details are below.

![alt text](http://i.imgur.com/QGwdl.png "Metroid")

## To build on Linux

Requires Go 1.1

From your GOPATH:

        $ sudo apt-get install libsdl1.2-dev libsdl-gfx1.2-dev libsdl-image1.2-dev libglew1.6-dev libxrandr-dev
        $ go get github.com/scottferg/Fergulator

## To build on OSX

You'll need to install [XQuartz](http://xquartz.macosforge.org/landing/) in order
to run on OSX.

Requires Go 1.1

From your GOPATH:

        $ brew install sdl --with-x11-driver
        $ brew install sdl_gfx sdl_image glew
        $ go get github.com/scottferg/Fergulator

## Run the emulator

        $ Fergulator path/to/game.nes

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

        Toggle pause - P
        Frame Advance - \
        Reload Debug JS - D
        Trigger Mode Event - M

## Supported Mappers

* NROM
* UNROM
* CNROM
* MMC1
* MMC2
* MMC3
* MMC5
* ANROM

## Tested games that run well or are playable

[List is in the wiki](https://github.com/scottferg/Fergulator/wiki/Tested-Games)
