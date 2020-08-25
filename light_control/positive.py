import time
import sys

import yeelight


def light_control(addr):
    bulb = yeelight.Bulb(addr)

    bulb.turn_on()
    bulb.set_rgb(255, 128, 0)
    print('Start lighting...')
    time.sleep(10)

    for i in range(0, 128, 5):
        bulb.set_rgb(255, 128 + i, 0)
        print('R:255 G:%d B:0' % (128 + i))
        time.sleep(1)

    for i in range(0, 128, 5):
        bulb.set_rgb(255, 255 - i, 0)
        print('R:255 G:%d B:0' % (255 - i))
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
