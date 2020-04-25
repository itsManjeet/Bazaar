import threading
from time import sleep
from subprocess import Popen, PIPE
import sys
def loop1():
    while True:
        print('loop1 is running')
        sleep(1)

def sbp():
    process = Popen(['sleep', '5'], stdout = PIPE)
    while True:
        print('i am also running')
        out = process.stdout.readline().decode('utf-8')
        if out == '' and process.poll() != None:
            break
        if out != '':
            print(out)


thread1 = threading.Thread(target=loop1)
thread2 = threading.Thread(target=sbp)
thread1.start()
thread2.start()

thread1.join()
thread2.join()