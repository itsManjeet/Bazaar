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
ERR_TEXT = '<span foreground="red"><b>😱 Error</b></span>'
OK_TEXT  = '<span foreground="green"><b>✔️ ok</b></span>'

class errorDialog(Gtk.Dialog):

    def __init__(self, parent, msg_text):
        Gtk.Dialog.__init__(self, 'Error', parent, 0,
            (Gtk.STOCK_OK, Gtk.ResponseType.OK))
        
        self.set_default_size(150, 100)

        label = Gtk.Label.new(msg_text)
        box = self.get_content_area()
        box.add(label)
        self.show_all()


class Bazaar:
    def __init__(self, ui_file):
        self.ui_file = ui_file
        self.icon_list_store = Gtk.ListStore(GdkPixbuf.Pixbuf, str)
        self.sysapp = sysapp.SysApp('/etc/sysconfig/app.conf')

    def __get_basic_widgets__(self):
        self.icon_view = self.get_widget('iconView')
        self.stack_page = self.get_widget('stack_page')
        self.install_button = self.get_widget('install_button')

    def load(self):
        self.builder = Gtk.Builder()
        self.builder.add_from_file(self.ui_file)
        self.builder.connect_signals(self)
        self.window = self.get_widget('mainWindow')

        self.screen = Gdk.Screen.get_default()
        self.provider = Gtk.CssProvider()
        self.__get_basic_widgets__()

        css = b"""
        #ScreenShotImage {
            box-shadow: 2px 2px 2px grey, 0 0 2px black
        }
        """

        self.provider.load_from_data(css)

        Gtk.StyleContext.add_provider_for_screen(
            Gdk.Screen.get_default(),
            self.provider,
            Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION
        )

        self.icon_view.set_model(self.icon_list_store)
        self.icon_view.set_pixbuf_column(0)
        self.icon_view.set_text_column(1)


        threading.Thread(target=self.loadApps, args=[self.sysapp.get_cache(),]).start()
        for c in categories:
            label = Gtk.Label.new(c)
            label.set_margin_top(16)
            label.set_margin_bottom(16)

            label.set_margin_start(45)
            label.set_margin_end(45)

            self.get_widget('categoryListBox').add(label)

        self.terminal = Vte.Terminal()
        self.get_widget('terminalBox').add(self.terminal)
        self.terminal.set_scrollback_lines(1000)
        self.vte_pty = self.terminal.pty_new_sync(Vte.PtyFlags.DEFAULT)
        self.terminal.set_colors(
            Gdk.RGBA(0.0, 0.0, 0.0, 1.0),
            Gdk.RGBA(1.0, 1.0, 1.0, 1.0),
            None
        )
        self.terminal.set_cursor_blink_mode(Vte.CursorBlinkMode.OFF)
        self.terminal.set_color_cursor(
            Gdk.RGBA(1.0, 1.0, 1.0, 1.0)
        )
        self.terminal.set_pty(self.vte_pty)
        self.terminal.set_audible_bell(False)
        self.terminal.connect('child-exited', self.__terminal_finished__)
        self.terminal.show()

    def __terminal_finished__(self, *args):
       self.__post_exec_func__()
       self.vte_pty = self.terminal.pty_new_sync(Vte.PtyFlags.DEFAULT)
       self.terminal.set_pty(self.vte_pty)

    def loadApps(self, app_data):     
        self.icon_list_store.clear()

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

            
            self.icon_list_store.append([pixbuf, a['name']])

        
    def run(self):
        self.window.show_all()
        Gtk.main()

    def _errorDialog(self, msg_text, exit_func = None):
        dialog = errorDialog(self.window, msg_text)
        dialog.run()
        if exit_func is not None:
            exit_func()

        dialog.destroy()
        Gtk.main_quit()

    def get_widget(self, widget_name):
        widget = self.builder.get_object(widget_name)
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
            self.loadApps(self.sysapp.app_data)
        else:
            for i in self.sysapp.app_data:
                if category in i['category']:
                    app_data.append(i)

            self.loadApps(app_data)
    
    def onAppSelected(self, *args):
        item = args[0].get_selected_items()[0][0]
        app_name = self.icon_list_store[item][1]
        self.loadAppPage(app_name)
    

    def __set_widget_content__(self, widget_name, text, MARKUP = False):
        if MARKUP:
            self.get_widget(widget_name).set_markup(text)
        else:
            self.get_widget(widget_name).set_text(text)


    def __set_page__(self, page_name):
        self.stack_page.set_visible_child_name(page_name)

    def __clear_signal__(self, widget, signals):
        if type(widget) == str:
            widget = self.get_widget(widget)

        for i in signals:
            try:
                widget.disconnect_by_func(i)
            except TypeError:
                continue


    def __set_text_buffer__(self, widget, text):
        widget = self.get_widget(widget)
        widget.delete(
            widget.get_start_iter(),
            widget.get_end_iter()
        )
        widget.set_text(text, len(text))

    def load_url(self, widget, extra , url):
        os.system('exo-open --launch WebBrowser %s' % url)

    def loadAppPage(self, app_name):
        #print("getting app %s" % app_name)
        app_data = self.sysapp.get_app(app_name)

        update_button = self.get_widget('update_button')
        self.__clear_signal__(update_button, [self.on_update_click])
        new_version = self.sysapp.latest_version(app_data)
        if new_version != app_data['version']:
            self.get_widget('update_button').set_visible(True)
            update_button.connect('clicked', self.on_update_click, app_data)
            self.__set_widget_content__('app_name_label', '😋 %s [%s -> %s]' % (app_data['name'], app_data['version'], new_version))
        else:
            self.get_widget('update_button').set_visible(False)
            self.__set_widget_content__('app_name_label', '%s %s' % (app_data['name'], app_data['version']))

        
        self.__set_widget_content__('app_desc_label', app_data['detail'])

        self.__set_page__('app_info_page')

        self.__set_widget_content__('lic_label', app_data['license'])
        
        self.__clear_signal__('url_label', [self.load_url])

        self.get_widget('url_label').set_markup('<a href="">%s</a>' % app_data['url'])        	
        self.get_widget('url_label').connect('activate-link', self.load_url, app_data['url'])

        is_visible = False

        if app_data['status'] == 'installed':
            is_visible  = True
            if app_data['dependencies'] == 'satisfied':
                self.__set_widget_content__('dep_label', OK_TEXT, True)
            else:
                self.__set_widget_content__('dep_label', ERR_TEXT, True)
            if app_data['integrity'] == 'ok':
                self.__set_widget_content__('integrity_label', OK_TEXT, True)
            else:
                self.__set_widget_content__('integrity_label', ERR_TEXT, True)

            
            self.__set_widget_content__('install_label', app_data['installed on'])

        else:
            is_visible  = False

        for i in ['dep_label', 'install_label', 'integrity_label']:
            self.get_widget(i).set_visible(is_visible)


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
            
        self.get_widget('app_image').set_from_pixbuf(pixbuf)
        
        
        self.__clear_signal__(self.install_button, [self.on_uninstall_click,self.on_install_click])

        if app_data['status'] == 'installed':
            self.install_button.set_label('uninstall')
            self.install_button.connect('clicked', self.on_uninstall_click, app_data)

        else:
            self.install_button.set_label('install')
            self.install_button.connect('clicked', self.on_install_click, app_data)

        self.__set_text_buffer__('files_list_buffer', ''.join(self.sysapp.get_files(app_data)))

        with open('/%s/%s/%s/recipie' % (self.sysapp.repodir, app_data['repo'], app_name)) as fptr:
            data = fptr.read()
            self.__set_text_buffer__('recipie_buffer', data)
    

        scfile = '/%s/%s/%s/screenshot' % (self.sysapp.repodir, app_data['repo'], app_name)
        if os.path.exists(scfile):
            #print('adding image from',scfile)
            pixbuf = GdkPixbuf.Pixbuf.new_from_file_at_scale(
                filename = scfile,
                width = 500,
                height = 500,
                preserve_aspect_ratio = True
            )
            self.get_widget('screenshot_image').set_from_pixbuf(pixbuf)    


    def on_install_click(self, btn, appdata):
        btn.set_sensitive(False)
        btn.set_label('installing....')
        self.get_widget('back_button').set_sensitive(False)
        self.__cur_app = appdata['name']
        self.__exec_func__(self.sysapp.get_install_cmd(appdata))


    def __exec_func__(self, args):
        self.child_pid, _, _, _ = GLib.spawn_async(
            working_directory = os.getcwd(),
            argv = ['/bin/pkexec'] + args,
            envp = [
                'no_ask=1'
            ] ,
            flags = (GLib.SpawnFlags.DO_NOT_REAP_CHILD),
            child_setup = self._child_setup,
            user_data = self.vte_pty
        )

        self.terminal.watch_child(self.child_pid)

    def _child_setup(self, vte_pty):
        vte_pty.child_setup()
    
    def __post_exec_func__(self):
        self.get_widget('back_button').set_sensitive(True)
        self.__clear_signal__(self.install_button,
            [self.on_install_click, self.on_uninstall_click]
        )
        self.loadAppPage(self.__cur_app)
        self.install_button.set_sensitive(True)

    def on_uninstall_click(self, btn, appdata):
        btn.set_sensitive(False)
        btn.set_label('uninstalling....')
        self.get_widget('back_button').set_sensitive(False)
        self.__cur_app = appdata['name']
        self.__exec_func__(self.sysapp.get_uninstall_cmd(appdata))

    def on_update_click(self, btn, appdata):
        btn.set_sensitive(False)
        self.install_button.set_sensitive(False)
        btn.set_label('updating....')
        self.get_widget('back_button').set_sensitive(False)
        self.__cur_app = appdata['name']
        self.__exec_func__(self.sysapp.get_update_cmd(appdata))
        

    def onBackBtnClick(self, *args):
        self.__set_page__('market_page')

    def onSearchChanged(self, searchbox):
        search_text = searchbox.get_text()
        app_data = []
        for i in self.sysapp.app_data:
            if search_text in i['name']:
                app_data.append(i)

        self.loadApps(app_data)

    def _exec_process_refresh(self):
        self.get_widget('refreshBtn').set_label('refreshing....')
        subprocess = Gio.Subprocess.new(['pkexec', 'sys-app', 'refresh'],0)
        subprocess.wait_check_async(None, self._postRefreshProcess)

    def _postRefreshProcess(self, subprocess, result):
        subprocess.wait_check_finish(result)
        self.get_widget('refreshBtn').set_sensitive(True)
        self.__clear_signal__(self.install_button, [self.on_install_click, self.on_uninstall_click])
        self.get_widget('refreshBtn').set_label('refresh')
        self.sysapp.getCache()

        self.loadApps(self.sysapp.app_data)

    def onRefreshClick(self, *args):
        self.get_widget('refreshBtn').set_sensitive(False)
        self._exec_process_refresh()



if __name__ == '__main__':
    bazaar = Bazaar('data/ui.glade')
    bazaar.load()
    bazaar.run()
