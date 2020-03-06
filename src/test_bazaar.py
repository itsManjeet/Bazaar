from backend.bazaar import Bazaar

if __name__ == '__main__':
    b = Bazaar()
    b.getApps()
    d = b.getInstallCMD('acl')
    print(d)
