from releax.config import Config
import os
class SysApp:
    def __init__(self,config_file):
        self.config = Config(config_file)
        self.config.parse()

        self.repodir = self.config.get('repo_dir')
        self.repos = self.config.get('repo')

    def get_cache(self):
        self.app_data = []
        for repo in self.repos.split():
            for app in os.listdir(os.path.join(self.repodir, repo)):
                appdir = os.path.join(self.repodir, repo, app)
                
                if not os.path.exists(appdir + '/recipie'):
                    continue

                app_conf = {
                    'name': app,
                    'category': self.__get_category__(appdir),
                    'icon': self.__get_icon__(appdir),
                    'repo': repo
                }

                self.app_data.append(app_conf)

        return self.app_data

    def get_app(self, app_name):
        for app in self.app_data:
            if app['name'] == app_name:
                for i in os.popen("sys-app info %s" % app_name).readlines():
                    for j in ['\x1b[1;35m', '\x1b[1;37m','\x1b[0m','\x1b[1;32m','\x1b[1;31m']:
                        i = i.replace(j,'')
                    var = i[:i.find(':')].strip().lower()
                    val = i[i.find(':') + 1:].strip().lower()
                    app[var] = val
                return app
        return None


    def __get_category__(self, appdir):
        with open(appdir + '/recipie') as fptr:
            for l in fptr.readlines():
                if l.startswith('#') and 'Category' in l:
                    category = l.split(':')[1].strip().split()
                    return category
        
        return []

    def __get_icon__(self, appdir):
        if os.path.exists(appdir + '/icon'):
            return appdir + '/icon'
        return None

    def latest_version(self, app):
        recipie_path = os.path.join(self.repodir, app['repo'], app['name'], 'recipie')
        with open(recipie_path) as fptr:
            for i in fptr.readlines():
                if i.startswith('version='):
                    return i.split('=')[1].strip()

        return app['version']


    def get_files(self, app):
        try:
            files_path = os.path.join(self.config.get('data_dir'), app['name'], 'files')
            with open(files_path) as fptr:
                files_list = fptr.readlines()
                return files_list
        except FileNotFoundError:
            return ''


    def get_install_cmd(self, app):
        return ['sys-app','install', app['name'], '--no-ask']
    
    def get_uninstall_cmd(self, app):
        return ['sys-app','remove', app['name'], '--no-ask']

    def get_update_cmd(self, app):
        return ['sys-appupdate', app['name']]
