import xml.etree.ElementTree as ET

class FlatPak:
    def __init__(self, data_file):
        self.data_file = data_file

    def GenAppsData(self):
        self.tree = ET.parse(self.data_file)
        self.root = self.tree.getroot()
        self.components = self.root.findall('component')
        self.categories = []
        self.keywords = []
        self.appdata = []

        for component in self.components:
            a = {
                'name': component.find('name').text,
                'id': component.find('id').text,
                'summary': component.find('summary').text,
                'description': self.__get_desc(component),
                #'license': component.find('project_license').text,
                'developer': self.__get_developer_name(component),
                'icons': self.__get_icons(component),
                'categories': self.__get_categories(component),
                'keywords': self.__get_keywords(component),
                'url': self.__get_urls(component),
                'depends': self.__get_depends(component),
                'type': 'flatpak'
            }

            self.appdata.append(a)

        return self.appdata

    def __get_developer_name(self, component):
        try:
            devel_name = component.find('developer_name').text
            return devel_name
        except AttributeError:
            return None

    def __get_desc(self, component):
        try:
            desc = ""
            for d in component.findall('description'):
                for p in d.findall('p'):
                    desc += p.text

            return desc
        except AttributeError:
            return 'None'

    def __get_icons(self, component):
        icons = []
        for i in component.findall('icon'):
            try:
                icon = {
                    'size': i.attrib['width'],
                    'file': i.text,
                    'type': i.attrib['type']
                }
                icons.append(icon)
            except KeyError:
                pass

        
        return icons

    def __get_categories(self, component):
        categories = []
        try:
            for c in component.findall('categories'):
                for cd in c.findall('category'):
                    categories.append(cd.text)
                    if cd.text not in self.categories:
                        self.categories.append(cd.text)

        except AttributeError:
            return None

    def __get_screenshots(self, component):
        screenshoots = []
        try:
            for s in component.findall('screenshots'):
                for sc in s.findall('screenshot'):
                    screenshoot = {
                        'type': sc.attrib['type'],
                        'images': []
                    }

                    for img in sc.findall('image'):
                        screenshoot['images'].append({
                            'height': img.attrib['height'],
                            'width': img.attrib['width'],
                            'source': img.text
                        })
                    screenshoots.append(screenshoot)

            return screenshoots
        except AttributeError:
            return None

    def __get_keywords(self, component):
        keywords = []
        try:
            for k in component.findall('keywords'):
                for kw in k.findall('keyword'):
                    keywords.append(kw.text)
                    if kw.text not in self.keywords:
                        self.keywords.append(kw.text)

        except AttributeError:
            return None

    def __get_urls(self, component):
        urls = []
        for u in component.findall('url'):
            url = {
                'type': u.attrib['type'],
                'address': u.text
            }

            urls.append(url)

        return urls

    def __get_depends(self, component):
        depends = []
        try:
            depends.append(component.attrib['runtime'])
        except KeyError:
            pass

        try:
            depends.append(component.attrib['sdk'])
        except KeyError:
            pass
    
        return depends

    def GetApp(self, app):
        for i in self.appdata:
            if i['name'] == app:
                return i
        return None

    def Install(self, app):
        try:
            print(self.GetApp(app))
            appID = self.GenApp(app)['id']
            return ['flatpak','install',appID, '-y']
        
        except AttributeError:
            return None

        except KeyError:
            return None
