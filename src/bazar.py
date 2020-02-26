#!/bin/python

import gi
gi.require_version('Gtk','3.0')
from gi.repository import Gtk, GdkPixbuf, Gdk
import xml.etree.ElementTree as ET

IconsData = '/var/lib/appstream-extractor/export-data/appstream-flathub-x86_64-2020-02-26.04:28:37.+0000/icons/'
resource_file = 'data/ui.glade'

listStore = Gtk.ListStore(GdkPixbuf.Pixbuf, str)

class BazarBack:

    def __init__(self,dataFile):
        self.tree = ET.parse(dataFile)
        self.root = self.tree.getroot()
        self.Components = self.root.findall('component')

    def GetAppData(self):
        self.AppData = []
        for cmpnt in self.Components:
            appData = {
                'name': cmpnt.find('name').text,
                'id': cmpnt.find('id').text,
                'summary': cmpnt.find('summary').text,
                'description': self.__getDesc(cmpnt),
                'icons': self.__getIcons(cmpnt),
                'categories': self.__getcategories(cmpnt),
                'url': self.__getURL(cmpnt)
            }

            self.AppData.append(appData)

        print('Total Apps',len(self.AppData))
        return self.AppData
    def __getDesc(self,cmpnt):
        try:
            allp = cmpnt.find('description').findall('p')
            d = ""
            for p in allp:
                d += p.text

            return d
        except AttributeError:
            return "None"
        

    def __getIcons(self,cmpnt):
        allIcons = cmpnt.findall('icon')
        icons = []
        for a in allIcons:
            if a.attrib['type'] != 'cached':
                continue
            icon = {
                'size': a.attrib['width'],
                'file': a.text
            }

            icons.append(icon)
        return icons
    
    def __getcategories(self,cmpnt):
        try:
            allCategories = cmpnt.find('categories').findall('category')
            Categories = []
            for c in allCategories:
                Categories.append(c.text)
            return Categories 
        except AttributeError:
            return None
          
    def __getURL(self, cmpnt):
        allUrls = cmpnt.findall('url')
        urls = []
        for u in allUrls:
            url = {
                'type': u.attrib['type'],
                'address': u.text
            }

            urls.append(url)

        return urls

class Bazar:
    def __init__(self):
        self.bg = BazarBack('/var/lib/appstream-extractor/export-data/appstream-flathub-x86_64-2020-02-26.04:28:37.+0000/appstream.xml')

    def UpdateAppsList(self,data):
        for app in data:
            try:
                img = GdkPixbuf.Pixbuf.new_from_file(IconsData + '/64x64/' + app['icons'][0]['file'])
            except:
                img = Gtk.IconTheme.get_default().load_icon('gtk-missing-image', 64, 0)
            listStore.append([img,app['name']])


class Handler:
    def OnMainDestroy(self, *args):
        Gtk.main_quit()

    def SelectionChanged(self, IconView):
        selected = IconView.get_selected_items()
        SelectedApp = listStore[selected[0][0]][1]
        SelectedAppIcon = listStore[selected[0][0]][0]
        for i in b.bg.AppData:
            if i['name'] == SelectedApp:
                Stack.set_visible_child_name('AppInfoPage')
                AppNameLabel.set_text(i['name'])
                AppSummaryLabel.set_text(i['summary'])
                DescriptionLabel.set_text(i['description'])
                AppImage.set_from_pixbuf(SelectedAppIcon)
                InstallerButton.connect('clicked', self.InstallApp, i)

    def SearchChanged(self, searchWidget):
        listStore.clear()
        searchText = searchWidget.get_text()
        searchText = searchText.lower()
        finalData = b.bg.AppData
        if len(searchText) != 0:
            d = []
            for x in finalData:
                if searchText in x['name'].lower() or searchText in x['summary'].lower() or searchText in x['description'].lower():
                    d.append(x)
            finalData = d

        b.UpdateAppsList(finalData)



    def InstallApp(self, widget, Data):
        print(widget, Data)

        

    def BackButtonClicked(self, *args):
        Stack.set_visible_child_name('MainPage')

Builder = Gtk.Builder()
Builder.add_from_file(resource_file)
Window = Builder.get_object('mainWindow')

IconListView = Builder.get_object('IconListView')
Stack = Builder.get_object('Stack')
AppNameLabel = Builder.get_object('AppNameLabel')
AppSummaryLabel = Builder.get_object('AppSummaryLabel')
AppImage = Builder.get_object('AppImage')
DescriptionLabel = Builder.get_object('DescriptionLabel')
InstallerButton = Builder.get_object('InstallerButton')

IconListView.set_model(listStore)
IconListView.set_pixbuf_column(0)
IconListView.set_text_column(1)

Builder.connect_signals(Handler())

b = Bazar()
b.bg.GetAppData()
b.UpdateAppsList(b.bg.AppData)

Window.show_all()
Gtk.main()
