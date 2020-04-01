#!/bin/sh

FLATPAK_REMOTE_URL="https://flathub.org/repo/flathub.flatpakrepo"
FLATPAK_ARCHS=(x86_64 i386)

function Main {
    flatpak remote-add --user --if-not-exists flathub $FLATPAK_REMOTE_URL
    for i in $FLATPAK_ARCHS ; do
        flatpak --user update --appstream --arch=$i
    done

    sys-app rf
}

Main