from flatpakBazaar import FlatPak


flatpakFile = '/var/lib/appstream-extractor/export-data/appstream-flathub-x86_64-2020-02-26.04:28:37.+0000/appstream.xml'
iconsData = '/var/lib/appstream-extractor/export-data/appstream-flathub-x86_64-2020-02-26.04:28:37.+0000/icons/'


class Bazaar:
    def __init__(self):
        self.flatpak = FlatPak(flatpakFile)
    
    def getApps(self):
        self.appData = self.flatpak.GenAppsData()

        return self.appData

    def GetApp(self, app):
        flatpakApp = self.flatpak.GetApp(app)
        if flatpakApp is None:
            return None
        
        return flatpakApp

    def getAppIcon(self, app):
        return iconsData + '/64x64/' + app['icons'][0]['file']

    def getAppsFromCategory(self, categoryName):
        apps = []
        for a in self.appData:
            if a['categories'] is None:
                continue
            elif categoryName in a['categories']:
                apps.append(a)

        return apps

    def Install(self, app):
        return self.flatpak.Install(app)