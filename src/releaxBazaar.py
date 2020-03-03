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

            'categories': [
                'releax'
            ],
            'keywords': ['releax'],

            'url': self.__get_comment(recipieFile, 'URL'),
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

    
