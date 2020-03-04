from flatpakBazaar import FlatPak
from releaxBazaar import Releax

flatpakFile = '/var/lib/appstream-extractor/export-data/appstream-flathub-x86_64-2020-02-26.04:28:37.+0000/appstream.xml'
iconsData = '/var/lib/appstream-extractor/export-data/appstream-flathub-x86_64-2020-02-26.04:28:37.+0000/icons/'


class Bazaar:
    def __init__(self):
        self.flatpak = FlatPak(flatpakFile)
        self.releax = Releax('/usr/recipies/')
    
    def getApps(self):
        self.appdata = self.flatpak.getApps()
        self.appdata += self.releax.getApps()

        self.categories = self.flatpak.categories

        return self.appdata

    def getApp(self, app):
        for i in self.appdata:
            if i['name'] == app:
                return i
        return None

    def getAppIcon(self, app):
        return iconsData + '/64x64/' + app['icons'][0]['file']

    def getAppsFromCategory(self, categoryArr):
        apps = []
        for a in self.appdata:
            if a['categories'] is None:
                continue
            for c in a['categories']:
                if c in categoryArr:
                    if a not in apps:
                        apps.append(a)

        return apps

    def getDepends(self, app):
        if type(app) == str:
            app = self.getApp(app)
            if app is None:
                return -1

        if app['type'] == 'releax':
            return self.releax.getDepends(app)
        elif app['type'] == 'flatpak':
            return self.flatpak.getDepends(app)
        else:
            return -3

    def getInstallCMD(self, app):
        if type(app) == str:
            app = self.getApp(app)
            if app is None:
                return -1

        if app['type'] == 'releax':
            return self.releax.getInstallCMD(app)
        elif app['type'] == 'flatpak':
            return self.flatpak.getInstallCMD(app)
        else:
            return ['/usr/bin/sh','echo','"dont know how to install %s"' % app['id']]

    def isInstall(self, app):
        if type(app) == str:
            app = self.getApp(app)
            if app is None:
                return -1

        if app['type'] == 'releax':
            return self.releax.isInstall(app)
        elif app['type'] == 'flatpak':
            return self.flatpak.isInstall(app)
        else:
            return False

    def install(self, app):
        if type(app) == str:
            app = self.getApp(app)
            if app is None:
                return -1

        if app['type'] == 'releax':
            return self.releax.install(app)
        else:
            print('not yet implemented')