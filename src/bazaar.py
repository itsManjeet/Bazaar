#!/bin/python
import gi
gi.require_version('Gtk','3.0')
gi.require_version('Vte','2.91')
from gi.repository import Gtk, GLib, Gdk, GdkPixbuf, Gio, Vte
from time import sleep
from backend import sysapp
import os
import threading
import subprocess
from dbus import SystemBus, Interface


categories = [
    'All', 'Accessories', 'Development', 'Graphics', 'Internet', 'Games', 'Multimedia', 'Office', 'Customizations', 'Plugins', 'System', 'Libraries'
]

class errorDialog(Gtk.Dialog):

    def __init__(self, parent, msg_text):
        Gtk.Dialog.__init__(self, 'Error', parent, 0,
            (Gtk.STOCK_OK, Gtk.ResponseType.OK))
        
        self.set_default_size(150, 100)

        label = Gtk.Label(msg_text)
        box = self.get_content_area()
        box.add(label)
        self.show_all()

def loadWebsite(url):
    os.system('exo-open --launch WebBrowser %s' % url)

class Bazaar:
    def __init__(self, ui_file):
        self._ui_file = ui_file
        self._icon_list_store = Gtk.ListStore(GdkPixbuf.Pixbuf, str)
        self.sysapp = sysapp.SysApp('/etc/sysconfig/app.conf')

    def load(self):
        self._builder = Gtk.Builder()
        self._builder.add_from_file(self._ui_file)
        self._builder.connect_signals(self)
        self._window = self.getWidget('mainWindow')

        self._screen = Gdk.Screen.get_default()
        self._provider = Gtk.CssProvider()

        css = b"""
        #ScreenShotImage {
            box-shadow: 2px 2px 2px grey, 0 0 2px black
        }
        """

        self._provider.load_from_data(css)

        Gtk.StyleContext.add_provider_for_screen(
            Gdk.Screen.get_default(),
            self._provider,
            Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION
        )
        self.getWidget('iconView').set_model(self._icon_list_store)
        self.getWidget('iconView').set_pixbuf_column(0)
        self.getWidget('iconView').set_text_column(1)


        threading.Thread(target=self._load_apps, args=[self.sysapp.getCache(),]).start()
        for c in categories:
            label = Gtk.Label(c)
            label.set_padding(40, 15);
            self.getWidget('categoryListBox').add(label)

        self.terminal = Vte.Terminal()
        self.getWidget('terminalBox').add(self.terminal)
        self.terminal.set_scrollback_lines(1000)
        self.vte_pty = self.terminal.pty_new_sync(Vte.PtyFlags.DEFAULT)
        self.terminal.set_colors(
            Gdk.RGBA(0.0, 0.0, 0.0, 1.0),
            Gdk.RGBA(1.0, 1.0, 1.0, 1.0),
            None
        )
        self.terminal.set_pty(self.vte_pty)
        self.terminal.set_audible_bell(False)
        self.terminal.connect('child-exited', self.__terminal_finished__)
        self.terminal.show()

    def __terminal_finished__(self, *args):
       self.__post_exec_func__()
       self.vte_pty = self.terminal.pty_new_sync(Vte.PtyFlags.DEFAULT)
       self.terminal.set_pty(self.vte_pty)

    def _load_apps(self, app_data):     
        self._icon_list_store.clear()

        for a in app_data:
            pixbuf = Gtk.IconTheme.get_default().load_icon('application-x-pak', 64, 0)
            if a['icon'] is not None:
                pixbuf = GdkPixbuf.Pixbuf.new_from_file_at_scale(
                    filename = a['icon'],
                    height = 64,
                    width = 64,
                    preserve_aspect_ratio = True)

            try:
                pixbuf = Gtk.IconTheme.get_default().load_icon(a['name'], 64, 0)
            except Exception as e:
                pass

            
            self._icon_list_store.append([pixbuf, a['name']])

        
    def run(self):
        self._window.show_all()
        Gtk.main()

    def _errorDialog(self, msg_text, exit_func = None):
        dialog = errorDialog(self._window, msg_text)
        dialog.run()
        if exit_func is not None:
            exit_func()

        dialog.destroy()
        Gtk.main_quit()

    def getWidget(self, widget_name):
        widget = self._builder.get_object(widget_name)
        if widget == None:
            self._errorDialog('Unable to get widget %s' % widget_name)
        else:
            return widget


    def onDestroy(self, *args):
        Gtk.main_quit()

    def onCategoryClick(self, listbox, listboxrow):
        category = listboxrow.get_child().get_text()
        app_data = []
        if category == 'All':
            self._load_apps(self.sysapp._appdata)
        else:
            for i in self.sysapp._appdata:
                if category in i['category']:
                    app_data.append(i)

            self._load_apps(app_data)
    
    def onAppSelected(self, *args):
        item = args[0].get_selected_items()[0][0]
        app_name = self._icon_list_store[item][1]
        self.loadAppPage(app_name)
    
    def loadAppPage(self, app_name):
        #print("getting app %s" % app_name)
        app_data = self.sysapp.getApp(app_name)

        self.getWidget('appNameLbl').set_text('%s %s' % (app_data['name'], app_data['version']))
        self.getWidget('appDescLbl').set_text(app_data['detail'])
        self.getWidget('stackPage').set_visible_child_name('appInfo')
        pixbuf = Gtk.IconTheme.get_default().load_icon('application-x-pak', 128, 0)
        if app_data['icon'] is not None:
            pixbuf = GdkPixbuf.Pixbuf.new_from_file_at_scale(
                filename = app_data['icon'],
                height = 128,
                width = 128,
                preserve_aspect_ratio = True)

        try:
            pixbuf = Gtk.IconTheme.get_default().load_icon(app_data['name'], 128, 0)
        except Exception as e:
            pass
            
        self.getWidget('appImage').set_from_pixbuf(pixbuf)
        self.clearClickButton()
        clickbtn = self.getWidget('clickButton')
        if app_data['status'] == 'installed':
            clickbtn.set_label('uninstall')
            clickbtn.connect('clicked',self.onUninstallBtnClick, app_name)
        else:
            clickbtn.set_label('install')
            clickbtn.connect('clicked',self.onInstallBtnClick, app_name)

        
        try:
            with open('/var/lib/app/%s/files' % app_name) as fptr:
                data = fptr.read()
                self.getWidget('textbuffer1').set_text(data, len(data))

        except FileNotFoundError:
            start = self.getWidget('textbuffer1').get_start_iter()
            end = self.getWidget('textbuffer1').get_end_iter()
            self.getWidget('textbuffer1').delete(start, end)

        
        recipieBuffer = self.getWidget('recipieBuffer')
        with open('/%s/%s/%s/recipie' % (self.sysapp._repo_dir, app_data['repo'], app_name)) as fptr:
            start = recipieBuffer.get_start_iter()
            end = recipieBuffer.get_end_iter()
            recipieBuffer.delete(start, end)
            data = fptr.read()
            recipieBuffer.set_text(data, len(data))

        self.getWidget('licenseLbl').set_text(app_data['license'])
        self.getWidget('urlLbl').set_markup('<a href="">%s</a>' % app_data['url'])
        try:
        	self.getWidget('urlLbl').disconnect_by_func(self.onLoadUrl)
        except:
        	pass
        	
        self.getWidget('urlLbl').connect('activate-link', self.onLoadUrl, app_data['url'])
        is_visible = False

        ErrText = '<span foreground="red"><b>😱 Error</b></span>'
        OkText = '<span foreground="green"><b>✔️ ok</b></span>'

        if app_data['status'] == 'installed':
            is_visible  = True
            if app_data['dependencies'] == 'satisfied':
                self.getWidget('depLbl').set_markup(OkText)
            else:
                self.getWidget('depLbl').set_markup(ErrText)

            if app_data['integrity'] == 'ok':
                self.getWidget('integLbl').set_markup(OkText)
            else:
                self.getWidget('integLbl').set_markup(ErrText)

            self.getWidget('installOnLbl').set_text(app_data['installed on'])
        else:
            is_visible  = False

        self.getWidget('depLbl').set_visible(is_visible)
        self.getWidget('installOnLbl').set_visible(is_visible)
        self.getWidget('integLbl').set_visible(is_visible)

        sc = self.getWidget('screenShotBox')
        scfile = '/%s/%s/%s/screenshot' % (self.sysapp._repo_dir, app_data['repo'], app_name)
        if os.path.exists(scfile):
            #print('adding image from',scfile)
            pixbuf = GdkPixbuf.Pixbuf.new_from_file_at_scale(
                filename = scfile,
                width = 500,
                height = 500,
                preserve_aspect_ratio = True
            )
            img = Gtk.Image.new_from_pixbuf(pixbuf)
            img.set_name('ScreenShotImage')
            img.show()
            sc.add(img)

    def clearClickButton(self):
        btn = self.getWidget('clickButton')
        try:
            btn.disconnect_by_func(self.onInstallBtnClick)
        except TypeError:
            pass
        
        try:
            btn.disconnect_by_func(self.onUninstallBtnClick)
        except TypeError:
            pass
        

    def onLoadUrl(self, *args):
        os.system('exo-open --launch WebBrowser %s' % args[2])

    def onInstallBtnClick(self, btn, appname):
        btn.set_sensitive(False)
        self.__cur_app = appname
        self.getWidget('backBtn').set_sensitive(False)
        self.__exec_func__("install",appname)

    def __exec_func__(self, method , appname):
        self.getWidget('clickButton').set_label('%sing....' % method)
        self.child_pid, _, _, _ = GLib.spawn_async(
            working_directory = os.getcwd(),
            argv = ['/bin/pkexec','sys-app',method, appname, '--no-ask'],
            flags = (GLib.SpawnFlags.DO_NOT_REAP_CHILD),
            child_setup = self._child_setup,
            user_data = self.vte_pty
        )

        self.terminal.watch_child(self.child_pid)
        #GLib.spawn_close_pid(self.child_pid)

    def _child_setup(self, vte_pty):
        vte_pty.child_setup()
    
    def __post_exec_func__(self):
        self.getWidget('backBtn').set_sensitive(True)
        self.clearClickButton()
        self.loadAppPage(self.__cur_app)
        self.getWidget('clickButton').set_sensitive(True)

    def onUninstallBtnClick(self, btn, appname):
        btn.set_sensitive(False)
        self.__cur_app = appname
        self.getWidget('backBtn').set_sensitive(False)
        self.__exec_func__("remove",appname)
        

    def onBackBtnClick(self, *args):
        self.getWidget('stackPage').set_visible_child_name('marketPage')

    def onSearchChanged(self, searchbox):
        search_text = searchbox.get_text()
        app_data = []
        for i in self.sysapp._appdata:
            if search_text in i['name']:
                app_data.append(i)

        self._load_apps(app_data)

    def _exec_process_refresh(self):
        self.getWidget('refreshBtn').set_label('refreshing....')
        subprocess = Gio.Subprocess.new(['pkexec', 'sys-app', 'refresh'],0)
        subprocess.wait_check_async(None, self._postRefreshProcess)

    def _postRefreshProcess(self, subprocess, result):
        subprocess.wait_check_finish(result)
        self.getWidget('refreshBtn').set_sensitive(True)
        self.clearClickButton()
        self.getWidget('refreshBtn').set_label('refresh')
        self.sysapp.getCache()

        self._load_apps(self.sysapp._appdata)

    def onRefreshClick(self, *args):
        self.getWidget('refreshBtn').set_sensitive(False)
        self._exec_process_refresh()



if __name__ == '__main__':
    bazaar = Bazaar('data/ui.glade')
    bazaar.load()
    bazaar.run()
