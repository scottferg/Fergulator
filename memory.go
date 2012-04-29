package main

type Word uint8
type BigWord uint16

type Memory [0xffff]Word

var (
    memory Memory
)
