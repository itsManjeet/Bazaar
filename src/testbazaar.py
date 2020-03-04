from bazaar import Bazaar

b = Bazaar()
b.getApps()

result = b.install('steam')
print(result)