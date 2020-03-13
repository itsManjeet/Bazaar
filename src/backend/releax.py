import os
import sys
import subprocess

class Releax:
    def __init__(self, recipieDir):
        self.recipieDir = recipieDir

    def getApps(self):
        self.categories = []
        self.keywords = []
        self.appdata = []
        for repo in os.listdir(self.recipieDir):
            self.categories.append(repo)
            if not os.path.isdir(os.path.join(self.recipieDir, repo)):
                continue
            for a in os.listdir(os.path.join(self.recipieDir,repo)):
                if not os.path.exists(os.path.join(self.recipieDir, repo, a, 'recipie')):
                    continue
                recipieDir = os.path.join(self.recipieDir, repo, a)
                a = self.__read_recipie(recipieDir)
                self.appdata.append(a)
        
        return self.appdata

    def getApp(self, app):
        for i in self.appdata:
            if i['name'] == app:
                return i
        return None


    def __read_recipie(self, recipieDir):
        #iconFile = os.path.join(recipieDir, 'icon.png')
        recipieFile = os.path.join(recipieDir, 'recipie')
        a = {
            'name': self.__get_var(recipieFile, 'name'),
            'id': self.__get_var(recipieFile, 'name'),
            'summary': self.__get_comment(recipieFile, 'Description'),
            'description': self.__get_comment(recipieFile, 'Description'),
            'license': self.__get_comment(recipieFile, 'License'),
            'developer': ' ',
            'icons': [],
            'repo': recipieDir,
            'categories': [
                'releax',
                recipieDir.split('/')[-2]
            ],
            'keywords': ['releax'],

            'url': [
                {
                    'type': 'home',
                    'address': self.__get_comment(recipieFile, 'URL')
                }
            ],
            'type': 'releax'
            }
        return a

    def __get_comment(self, recipieFile, valof):
        with open(recipieFile, 'r') as openFile:
            for line in openFile.readlines():
                line = line.lstrip()
                try:
                    if line[0] == '#':
                        var = line[1:line.find(':')].lstrip()
                        if var == valof:
                            return line[line.find(':')+1:].lstrip().strip()
                except IndexError:
                    pass
        
        return ''

    def __get_var(self, recipieFile, valof):
        with open(recipieFile, 'r') as openFile:
            for line in openFile.readlines():
                line = line.lstrip()
                var = line[0:line.find('=')].lstrip()
                if var == valof:
                    return line[line.find('=')+1:].lstrip().strip()
        
        return ''

    def getDepends(self, app):
        if type(app) == str:
            app = self.getApp(app)
            if app is None:
                return -1

        result = subprocess.run(['sys-app','dp', app['name']], stdout=subprocess.PIPE)
        deps = result.stdout.decode('utf-8').splitlines()
        a = []
        for i in deps:
            i = i[i.find(' ')+1:]
            a.append(i)
        return a

    def getInstallCMD(self, app):
        if type(app) == str:
            app = self.getApp(app)
            if app is None:
                return -1

        return ['/bin/pkexec','/bin/sys-app','in', app['name'], '--no-ask']

    def getUnInstallCMD(self, app):
        if type(app) == str:
            app = self.getApp(app)
            if app is None:
                return -1

        return ['/bin/pkexec','/bin/sys-app','rm', app['name'], '--no-ask']

    def isInstall(self, app):
        if type(app) == str:
            app = self.getApp(app)
            if app is None:
                return False
        
        result = subprocess.run(['ls','/var/lib/app/'], stdout=subprocess.PIPE)
        result = result.stdout.decode('utf-8')
        for a in result.splitlines():
            if a == app['name']:
                return True
        
        return False

    def install(self, app):
        if type(app) == str:
            app = self.getApp(app)
            if app is None:
                return -1
        
        result = subprocess.run(['/bin/pkexec','/bin/sys-app','in', app['name'], '--no-ask'])
        return result.returncode


    
