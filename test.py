import socket
import json
import time
import random


HOST = "127.0.0.1"
PORT = 9000
SENSORS = [("sensor:IPT", "PSI"), ("sensor:ITC", "C"), ("sensor:MPT", "PSI"), ("sensor:MTC", "C")]

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((HOST, PORT))
    for sensor in SENSORS:
        (name, unit) = sensor
        val = random.randint(0, 1500)
        out = {
        "t": time.time_ns(),
        "l": name,
        "u": unit,
        "v": val
    }   
        
    out = json.dumps(out).encode()
    length = len(out)
    s.sendall(length.to_bytes(4, byteorder="big"))
    s.sendall(out)
    print(name + " | " + str(val))
    
print("Done sending test data!")