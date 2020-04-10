#!/bin/python

import os
import xml.etree.ElementTree as ET
FLATPAK_REMOTE_URL="https://flathub.org/repo/flathub.flatpakrepo"
FLATPAK_ARCHS=[ 'x86_64']

FLATPAK_LOC = '/var/lib/appstream-extractor/export-data/appstream-flathub-x86_64-2020-02-26.04:28:37.+0000/'


class Flatpak:
    def refData(self):
        os.system('flatpak remoter-add --user --if-not-exists flathub %s' % FLATPAK_REMOTE_URL)
        for i in FLATPAK_ARCHS:
            os.system('flatpak --user update --appstream --arch=%s' % i)

    def getCache(self):
        self._tree = ET.parse(FLATPAK_LOC + '/appstream.xml')
        self._root = self._tree.getroot()

        self._components = self._root.findall('component')
        self._categories = []
        self._keywords = []
        self._appdata = []

        for c in self._components:
            icons = []
            for i in c.findall('icon'):
                try:
                    icon = {
                        'size': i.attrib['width'],
                        'file': i.text,
                        'type': i.attrib['type']
                    }

                    icons.append(icon)
                except KeyError:
                    pass
            
            categories = ['flatpak']
            try:
                for d in c.findall('categories'):
                    for cd in d.findall('category'):
                        categories.append(cd.text)
                        if cd.text not in self.categories:
                            self.categories.append(cd.text)

            except AttributeError:
                pass

            a = {
                'name': c.find('name').text,
                'icons': icons,
                'category': categories,
                'type': 'flatpak'
            }

            self._appdata.append(a)

        return self._appdata