#!/bin/sh

DESTDIR=${DESTDIR:-"/"}
PREFIX=${PREFIX:-"usr"}
BINDIR=${BINDIR:-"$PREFIX/bin/"}
LIBEXECDIR=${LIBEXECDIR:-"$PREFIX/libexec"}
DATADIR=${DATADIR:-"$PREFIX/share"}

mkdir -pv build/
cd src

mv const{.go,.old}

echo "
package main
const uiFile = \"/$DATADIR/bazaar/ui.glade\"
" > const.go

go build -o ../build/bazaar
mv const{.old,.go}
cd ..

strip build/bazaar
install -vDm755 build/bazaar $DESTDIR/$BINDIR/bazaar
install -vDm644 data/ui.glade -t $DESTDIR/$DATADIR/bazaar
install -vDm755 data/bazaar.desktop -t $DESTDIR/$DATADIR/applications/
install -vDm755 data/bazaar-pkexec.sh $DESTDIR/$BINDIR/bazaar-pkexec
install -vDm644 data/org.freedesktop.policykit.bazaar.policy -t $DESTDIR/$DATADIR/polkit-1/actions/

