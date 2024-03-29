#!/usr/bin/python3

from pyroute2 import IPRoute

ip = IPRoute()
octavia_interface = ip.link_lookup(ifname='octavia')

if len(octavia_interface):
    current_addresses = ip.get_addr(label='octavia')
    if '172.99.0.32' not in current_addresses:
        ip.addr('add', index = octavia_interface[0], address='172.99.0.32',
                mask=32)
ip.close()
