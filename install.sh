#!/bin/sh

DESTDIR=${DESTDIR:-"/"}
PREFIX=${PREFIX:-"/usr"}
BINDIR=${BINDIR:-"$PREFIX/bin"}
DATADIR=${DATADIR:-"$PREFIX/share"}

install -vDm755 data/bazaar.desktop -t $DESTDIR/$DATADIR/applications/
install -vDm644 data/ui.glade -t $DESTDIR/$DATADIR/bazaar/

install -vDm755 src/bazaar.py $DESTDIR/$BINDIR/bazaar

sed -i "s|data/ui.glade|$DATADIR/bazaar/ui.glade|g" $DESTDIR/$BINDIR/bazaar
install -vDm644 src/backend/*.py -t $DESTDIR/$PREFIX/lib/python3.8/site-packages/backend/
