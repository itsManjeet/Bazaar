#!/bin/sh

DESTDIR=${DESTDIR:-"/"}
PREFIX=${PREFIX:-"/usr"}
BINDIR=${BINDIR:-"$PREFIX/bin"}
DATADIR=${DATADIR:-"$PREFIX/share"}

install -vDm755 data/Bazaar.desktop -t $DESTDIR/$DATADIR/applications/
install -vDm644 data/ui.glade -t $DESTDIR/$DATADIR/bazaar/

install -vDm755 src/bazaar-genrepo.sh $DESTDIR/$BINDIR/bazaar-genrepo
install -vDm755 src/bazaar-ui.py $DESTDIR/$BINDIR/bazaar-ui

install -vDm644 src/backend/*.py -t $DESTDIR/$PREFIX/lib/python3.8/site-packages/backend/
