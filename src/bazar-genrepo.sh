#!/bin/sh

FLATPAK_REMOTE_URL="https://flathub.org/repo/flathub.flatpakrepo"
FLATPAK_ARCHS=(x86_64 i386)


function 

function Main {
    flatpak remote-add --user --if-not-exists flathub https://flathub.org 
}

