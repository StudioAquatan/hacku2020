import time
import sys

import yeelight


def light_control(addr):
    bulb = yeelight.Bulb(addr)

    bulb.turn_on()
    bulb.set_rgb(0, 128, 255)
    print('Start lighting...')
    time.sleep(10)

    for i in range(0, 256, 15):
        bulb.set_rgb(i, 128 - (int)(i / 5), 255 - i)
        print('R:%d G:%d B:%d' % (i, 128-i/5, 255 - i))
        time.sleep(1)

    time.sleep(10)
    bulb.turn_off()
    print('Stop lighting')
    return ()


if __name__ == '__main__':
    args = sys.argv
    if 2 <= len(args):
        light_control(args[1])
    else:
        print('Arguments are too short')
