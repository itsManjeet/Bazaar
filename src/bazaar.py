#!/bin/python
import gi
gi.require_version('Gtk','3.0')
from gi.repository import Gtk, GLib, Gdk, GdkPixbuf, Gio
from time import sleep
from backend import sysapp
import os
import threading
import subprocess


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
        self._sysapp = sysapp.SysApp('/etc/sysconfig/app.conf')

    def load(self):
        self._builder = Gtk.Builder()
        self._builder.add_from_file(self._ui_file)
        self._builder.connect_signals(self)
        self._window = self._getWidget('mainWindow')

        self._screen = Gdk.Screen.get_default()
        self._provider = Gtk.CssProvider()

        css = b"""
        """

        self._provider.load_from_data(css)

        Gtk.StyleContext.add_provider_for_screen(
            Gdk.Screen.get_default(),
            self._provider,
            Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION
        )
        self._getWidget('iconView').set_model(self._icon_list_store)
        self._getWidget('iconView').set_pixbuf_column(0)
        self._getWidget('iconView').set_text_column(1)



        threading.Thread(target=self._load_apps, args=[self._sysapp.getCache(),]).start()
        for c in os.listdir(self._sysapp._repo_dir):
            label = Gtk.Label(c)
            label.set_padding(30, 15);
            self._getWidget('categoryListBox').add(label)


    def _load_apps(self, app_data):     
        self._icon_list_store.clear()

        for a in app_data:
            pixbuf = Gtk.IconTheme.get_default().load_icon('application-x-pak', 64, 0)
            if a['icon'] is not None:
                pixbuf = GdkPixbuf.Pixbuf.new_from_file(a['icon'])

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

    def _getWidget(self, widget_name):
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
        if category == 'all':
            self._load_apps(self._sysapp._appdata)
        else:
            for i in self._sysapp._appdata:
                if category == i['category']:
                    app_data.append(i)

            self._load_apps(app_data)
    
    def onAppSelected(self, *args):
        item = args[0].get_selected_items()[0][0]
        app_name = self._icon_list_store[item][1]
        self._loadAppPage(app_name)
    
    def _loadAppPage(self, app_name):
        print("getting app %s" % app_name)
        app_data = self._sysapp.getApp(app_name)

        self._getWidget('appNameLbl').set_text('%s %s' % (app_data['name'], app_data['version']))
        self._getWidget('appDescLbl').set_text(app_data['detail'])
        self._getWidget('stackPage').set_visible_child_name('appInfo')
        try:
            pixbuf = Gtk.IconTheme.get_default().load_icon(app_name, 64, 0)
        except:
            pixbuf = Gtk.IconTheme.get_default().load_icon('application-x-pak', 64, 0)
            
        self._getWidget('appImage').set_from_pixbuf(pixbuf)
        self._clearClickButton()
        clickbtn = self._getWidget('clickButton')
        if app_data['status'] == 'installed':
            clickbtn.set_label('uninstall')
            clickbtn.connect('clicked',self.onUninstallBtnClick, app_name)
        else:
            clickbtn.set_label('install')
            clickbtn.connect('clicked',self.onInstallBtnClick, app_name)

        
        try:
            with open('/var/lib/app/%s/files' % app_name) as fptr:
                data = fptr.read()
                self._getWidget('textbuffer1').set_text(data, len(data))

        except FileNotFoundError:
            start = self._getWidget('textbuffer1').get_start_iter()
            end = self._getWidget('textbuffer1').get_end_iter()
            self._getWidget('textbuffer1').delete(start, end)

        
        recipieBuffer = self._getWidget('recipieBuffer')
        with open('/%s/%s/%s/recipie' % (self._sysapp._repo_dir, app_data['repo'], app_name)) as fptr:
            start = recipieBuffer.get_start_iter()
            end = recipieBuffer.get_end_iter()
            recipieBuffer.delete(start, end)
            data = fptr.read()
            recipieBuffer.set_text(data, len(data))

        self._getWidget('licenseLbl').set_text(app_data['license'])
        self._getWidget('urlLbl').set_markup('<a href="">%s</a>' % app_data['home url'])
        self._getWidget('urlLbl').connect('activate-link', self.onLoadUrl, app_data['home url'])
        is_visible = False
        if app_data['status'] == 'installed':
            is_visible  = True
            self._getWidget('depLbl').set_text(app_data['dependencies'])
            self._getWidget('installOnLbl').set_text('installed on')
            self._getWidget('integLbl').set_text('integrity')
        else:
            is_visible  = False

        self._getWidget('depLbl').set_visible(is_visible)
        self._getWidget('installOnLbl').set_visible(is_visible)
        self._getWidget('integLbl').set_visible(is_visible)

    def _clearClickButton(self):
        btn = self._getWidget('clickButton')
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
        self._getWidget('backBtn').set_sensitive(False)
        self._sysExecFunc("install",appname)

    def _sysExecFunc(self, method , appname):
        self._getWidget('clickButton').set_label('%sing....' % method)
        subprocess = Gio.Subprocess.new(['pkexec', 'sys-app', method, appname],0)
        subprocess.wait_check_async(None, self._postExecFunc)

    def _postExecFunc(self, subprocess, result):
        subprocess.wait_check_finish(result)
        self._getWidget('backBtn').set_sensitive(True)
        self._clearClickButton()
        self._loadAppPage(self.__cur_app)
        self._getWidget('clickButton').set_sensitive(True)

    def onUninstallBtnClick(self, btn, appname):
        btn.set_sensitive(False)
        self.__cur_app = appname
        self._getWidget('backBtn').set_sensitive(False)
        self._sysExecFunc("remove",appname)
        

    def onBackBtnClick(self, *args):
        self._getWidget('stackPage').set_visible_child_name('marketPage')

    def onSearchChanged(self, searchbox):
        search_text = searchbox.get_text()
        app_data = []
        for i in self._sysapp._appdata:
            if search_text in i['name']:
                app_data.append(i)

        self._load_apps(app_data)

    def _exec_process_refresh(self):
        self._getWidget('refreshBtn').set_label('refreshing....')
        subprocess = Gio.Subprocess.new(['pkexec', 'sys-app', 'refresh'],0)
        subprocess.wait_check_async(None, self._postRefreshProcess)

    def _postRefreshProcess(self, subprocess, result):
        subprocess.wait_check_finish(result)
        self._getWidget('refreshBtn').set_sensitive(True)
        self._clearClickButton()
        self._getWidget('refreshBtn').set_label('refresh')
        self._sysapp.getCache()

        self._load_apps(self._sysapp._appdata)

    def onRefreshClick(self, *args):
        self._getWidget('refreshBtn').set_sensitive(False)
        self._exec_process_refresh()



if __name__ == '__main__':
    bazaar = Bazaar('data/ui.glade')
    bazaar.load()
    bazaar.run()
