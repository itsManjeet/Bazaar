import gi
gi.require_version('Gtk','3.0')
gi.require_version('Vte', '2.91')
from gi.repository import Gtk, GdkPixbuf, Gdk, Vte, GLib
from bazaar import Bazaar

resourceFile = 'data/ui.glade'
listStore = Gtk.ListStore(GdkPixbuf.Pixbuf, str)

def UpdateAppPage(data):
    for app in data:
        if len(app['icons']) < 1:
            img = Gtk.IconTheme.get_default().load_icon('gtk-missing-image',64, 0)
        elif app['icons'][0]['type'] == 'cached':
            img = GdkPixbuf.Pixbuf.new_from_file(bazaar.getAppIcon(app))
        
        listStore.append([img,app['name']])


def BuildCategories(container):
    defaultCategories = [ 'Accessories', 'Development', 'Games', 'Graphics', 'Internet', 'Multimedia', 'Office', 'System', 'Settings' ]
    for c in defaultCategories:
        btn = Gtk.Button.new_with_label(c)
        btn.connect('clicked', UpdateCategories, c)
        container.add(btn)

def UpdateCategories(widget, category):
    listStore.clear()
    apps = bazaar.getAppsFromCategory(category)
    print(apps)
    UpdateAppPage(apps)

class Hander:
    def OnMainDestroy(self, *args):
        Gtk.main_quit()

    def SelectionChanged(self, iconView):
        selected = iconView.get_selected_items()
        try:
            appName = listStore[selected[0][0]][1]
            appIcon = listStore[selected[0][0]][0]

        except IndexError:
            return

        selectedApp = bazaar.GetApp(appName)
        self.__updateAppPage(selectedApp, appIcon)
        bazaarPage.set_visible_child_name('AppInfoPage')

    def SearchChanged(self, *args):
        print('text changed')    
    
    def BackButtonClicked(self, *args):
        bazaarPage.set_visible_child_name('MainPage')

    def __updateAppPage(self, app, icon):
        appName_label.set_text(app['name'])
        appSum_label.set_text(app['summary'])
        appDesc_label.set_text(app['description'])
        appLogo_image.set_from_pixbuf(icon)

        appInstall_button.connect('clicked', self.__install_app, app['name'])

    def __install_app(self, widget, app):
        print(bazaar.Install(app))


        
bazaar = Bazaar()
builder = Gtk.Builder()
builder.add_from_file(resourceFile)

window = builder.get_object('mainWindow')
bazaarPage = builder.get_object('bazaarPage')

appList_iconview = builder.get_object('appList_iconview')

appName_label = builder.get_object('appName_label')
appSum_label = builder.get_object('appSum_label')
appDesc_label = builder.get_object('appDesc_label')
appLogo_image = builder.get_object('appLogo_image')
appInstall_button = builder.get_object('appInstall_button')

categoryBox = builder.get_object('CategoryBox')

builder.connect_signals(Hander())

appList_iconview.set_model(listStore)
appList_iconview.set_pixbuf_column(0)
appList_iconview.set_text_column(1)

UpdateAppPage(bazaar.getApps())
BuildCategories(categoryBox)

window.show_all()
Gtk.main()
