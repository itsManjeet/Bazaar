import gi
gi.require_version('Gtk','3.0')
gi.require_version('Vte', '2.91')
from gi.repository import Gtk, GdkPixbuf, Gdk, Vte, GLib
from backend.bazaar import Bazaar
import threading
import os

resourceFile = 'data/ui.glade'
listStore = Gtk.ListStore(GdkPixbuf.Pixbuf, str)

def ErrorDialog(errMessage, toexit = False):
    dialog = Gtk.MessageDialog(
        None, 0, Gtk.MessageType.ERROR,
        Gtk.ButtonsType.CANCEL,
        errMessage,
    )

    dialog.run()
    dialog.destroy()
    if toexit:
        Gtk.main_quit()

def ResponseDialog(message):
    dialog = Gtk.MessageDialog(
        None, 0, Gtk.MessageType.QUESTION,
        Gtk.ButtonsType.YES_NO,
        message,
    )

    response = dialog.run()
    dialog.destroy()

    return response


def UpdateAppPage(data):
    for app in data:
        try:
            if len(app['icons']) < 1:
                try:
                    img = Gtk.IconTheme.get_default().load_icon(app['name'],64, 0)
                except gi.repository.GLib.Error:
                    img = Gtk.IconTheme.get_default().load_icon('package-x-generic',64, 0)
            elif app['icons'][0]['type'] == 'cached':
                img = GdkPixbuf.Pixbuf.new_from_file(bazaar.getAppIcon(app))
        except TypeError:
            continue
        listStore.append([img,app['name']])


def BuildCategories(container):
    defaultCategories = {
        'All': ['All'],
        'Accessories': ['Calendar', 'Viewer', 'Emulator', 'Utility', 'VideoConference', 'FileTransfer', 'FileTools', 'Archiving', 'TextEditor', 'Publishing', 'Presentation', 'Sequencer', 'RemoteAccess', 'OCR', 'Scanning', 'PackageManager', 'Application', 'Calculator', 'Dictionary', 'Maps', 'Clock', 'Chart', 'FileManager', 'Productivity', 'Spreadsheet'],
        'Education': ['Literature', 'LearnToCode', 'Documentation', 'Astronomy','Math', 'Robotics', 'Physics', 'Matrix', 'GUIDesigner', 'NumericalAnalysis', 'Chemistry', 'Geography', 'Geoscience'],
        'Development': ['Development','DataVisualization', 'IDE', 'WebDevelopment', 'RevisionControl', 'ProjectManagement', 'Publishing', 'Java', 'ContactManagment', 'Debugger', 'Profiling', 'GUIDesigner', 'Productivity'],
        'Games': ['Emulator', 'Games', 'ArcadeGame', 'RolePlaying', 'Shooter','Simulation', 'StrategyGame', 'ActionGame', 'AdventureGame', 'CardGame', 'LogicGame', 'BlocksGame', 'BoardGame', 'SportsGame', 'KidsGame', 'Adventure', 'Games', 'Role Playing'],
        'Internet': ['Chat','Feed','InstantMessaging', 'Network', 'P2P', 'Email', 'VideoConference', 'Telephony', 'News', 'WebBrowser', 'IRCClient', 'HamRadio', 'Internet', 'Communication'],
        'Multimedia': ['Audio', 'AudioVideo', 'Music', 'Chat', 'Email', 'Player', 'Video', 'VideoConference', 'AudioVideoEditing', '2DGraphics', 'RasterGraphics', 'TV', 'Graphics', 'GTK', 'Photography', 'Qt', 'Recorder', '3DGraphics', 'Tuner', 'HamRadio', 'ImageProcessing', 'VectorGraphics', 'GUIDesigner', 'Art', 'Mixer'],
        'Office': ['Office', 'DataVisualization', 'Education', 'MedicalSoftware','Science', 'Electronics', 'Engineering', 'WordProcessor', 'VideoConference', 'Literature', 'ProjectManagement', 'Economy', 'Finance', 'Documentation', 'Presentation', 'Productivity', 'Spreadsheet'],
        'System': ['System', 'Security', 'ConsoleOnly', 'Monitor', 'FileSystem', 'Languages', 'Translation', 'core', 'extra'],
        'Settings': ['HardwareSettings', 'DesktopSettings', ],
        'Xfce': ['xfce']
    }
    for c in defaultCategories:
        btn = Gtk.Button.new_with_label(c)
        context = btn.get_style_context()
        context.add_class('category')
        btn.connect('pressed', UpdateCategories, defaultCategories[c])
        container.add(btn)

def UpdateCategories(widget, category):
    listStore.clear()
    if 'All' in category:
        UpdateAppPage(bazaar.appdata)
    else:
        apps = bazaar.getAppsFromCategory(category)
        UpdateAppPage(apps)

class Hander:
    def __init__(self):
        self.isSearching = False

    def OnMainDestroy(self, *args):
        Gtk.main_quit()

    def SelectionChanged(self, iconView):
        selected = iconView.get_selected_items()
        try:
            appName = listStore[selected[0][0]][1]
            appIcon = listStore[selected[0][0]][0]

        except IndexError:
            return

        selectedApp = bazaar.getApp(appName)
        self.__updateAppPage(selectedApp, appIcon)
        bazaarPage.set_visible_child_name('AppInfoPage')

    def SearchChanged(self, searchBox):
        listStore.clear()
        searchBoxText = searchBox.get_text().lower()
        data = []
        for x in bazaar.appdata:
            if searchBoxText in x['name'].lower() or searchBoxText in x['summary'].lower() or searchBoxText in x['description'].lower():
               data.append(x)

        UpdateAppPage(data)

    
    def BackButtonClicked(self, *args):
        bazaarPage.set_visible_child_name('MainPage')

    def __updateAppPage(self, app, icon):
        terminalBox.hide()
        appName_label.set_text(app['name'])
        appSum_label.set_text(app['summary'])
        buffer = appDesc_label.get_buffer()
        DescText = '%s\n\nLicense:%s\nDeveloper:%s\n' % (app['description'], app['license'], app['developer'])
        for u in app['url']:
            DescText += '\n%s:%s' % (u['type'], u['address'])

        buffer.set_text(DescText)
        appLogo_image.set_from_pixbuf(icon)
        if bazaar.isInstall(app):
            appInstall_button.set_label('uninstall')
            appInstall_button.connect('clicked', self.__uninstall_app, app)
        else:
            appInstall_button.set_label('install')
            appInstall_button.connect('clicked', self.__install_app, app)

    def __execute_cmd(self, cmd):
        terminalBox.set_no_show_all(False)
        terminalBox.show_all()
        Terminal.spawn_sync(
                Vte.PtyFlags.DEFAULT,
                os.environ['HOME'],
                cmd,
                [],
                GLib.SpawnFlags.DO_NOT_REAP_CHILD,
                None,
                None
            )


    def __uninstall_app(self, widget, app):
        print('uninstalling %s' % app['name'])

    def __install_app(self, widget, app):
        terminalBox.set_no_show_all(False)
        terminalBox.show_all()
        result = bazaar.getDepends(app)
        if result == -1:
            ErrorDialog('unable to find app "%s"' % app['name'])
        elif result == -2:
            ErrorDialog('unable to process app of unknown type')
        else:
            if len(result) > 1:
                resp = ResponseDialog('%s are going to install' % (','.join(result)))
                if resp == Gtk.ResponseType.NO:
                    ErrorDialog('Ok Progress Exit')
                    return
            self.__execute_cmd(bazaar.getInstallCMD(app['name']))
            



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
appInfo_box = builder.get_object('AppInfoBox')
terminalBox = builder.get_object('terminalBox')
Terminal = Vte.Terminal()
terminalBox.add(Terminal)
terminalBox.set_no_show_all(True)

categoryBox = builder.get_object('CategoryBox')
categoryBox.show()

builder.connect_signals(Hander())

appList_iconview.set_model(listStore)
appList_iconview.set_pixbuf_column(0)
appList_iconview.set_text_column(1)

UpdateAppPage(bazaar.getApps())
BuildCategories(categoryBox)
window.show_all()
Gtk.main()
