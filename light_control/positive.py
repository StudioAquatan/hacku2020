import time
import sys
import random

import yeelight


def light_control(addr):
    bulb = yeelight.Bulb(addr)

    bulb.turn_on()
    bulb.set_rgb(255, 128, 0)
    print('Start lighting...')

    for i in range(40):
        r = random.randint(0, 256)
        g = random.randint(0, 256)
        b = random.randint(0, 256)
        bulb.set_rgb(r, g, b)
        print('R:%d G:%d B:%d' % (r, g, b))
        time.sleep(0.5)

    bulb.turn_off()
    print('Stop lighting')
    return ()


if __name__ == '__main__':
    args = sys.argv
    if 2 <= len(args):
        light_control(args[1])
    else:
        print('Arguments are too short')
