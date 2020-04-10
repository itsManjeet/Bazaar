from releax.config import Config
import os
class SysApp:
    def __init__(self,config_file):
        self._conf = Config(config_file)
        self._conf.parse()

        self._repo_dir = self._conf.get('repo_dir')
        self._repos = self._conf.get('repo')

    def getCache(self):
        self._appdata = []
        for repo in self._repos.split():
            for app in os.listdir(os.path.join(self._repo_dir, repo)):
                try:
                    if not os.path.exists(os.path.join(self._repo_dir,repo,app,'recipie')):
                        continue
                except Exception as e:
                    continue
                icon = None
                icon_path = os.path.join(self._repo_dir,repo,app,'icon')
                if os.path.exists(icon_path):
                    icon = icon_path
                self._appdata.append(
                    {
                        'name': app,
                        'category': repo,
                        'icon': icon,
                        'type': 'releax'
                    }
                )
        return self._appdata

    def getApp(self, app_name):
        app_data = {}
        for i in os.popen("sys-app info %s" % app_name).readlines():
            for j in ['\x1b[1;35m', '\x1b[1;37m','\x1b[0m','\x1b[1;32m']:
                i = i.replace(j,'')
            var = i[:i.find(':')].strip().lower()
            val = i[i.find(':') + 1:].strip().lower()
            app_data[var] = val

        for i in self._appdata:
            if i['name'] == app_name:
                app_data['repo'] = i['category']


        return app_data

