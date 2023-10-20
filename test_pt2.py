import socket
import json
import time
import random

HOST = "127.0.0.1"
PORT = 9000
SENSORS = [
    ("sensor:IPT", "PSI"),
    ("sensor:ITC", "C"),
    ("sensor:MPT", "PSI"),
    ("sensor:MTC", "C"),
]

def send_data(conn, data):
    length = len(data)
    conn.sendall(length.to_bytes(4, byteorder="big"))
    conn.sendall(data)

def main():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        # Send a "START" message to initiate the handshake
        start_message = "START".encode()
        send_data(s, start_message)

        # Receive and process the handshake response
        response_length_bytes = s.recv(2)
        response_length = int.from_bytes(response_length_bytes, byteorder="big")
        response_data = s.recv(response_length)
        response = response_data.decode()
        print("Received handshake response:", response)

        if response == "OK":
            while True:
                for sensor in SENSORS:
                    (name, unit) = sensor
                    val = random.randint(200, 1500)
                    data = {
                        "t": int(time.time() * 1e9),
                        "l": name,
                        "u": unit,
                        "v": val,
                }
                data_json = json.dumps(data).encode()
                send_data(s, data_json)
                print(name + " | " + str(val))
            time.sleep(1)  # Add a delay between data transmissions

        print("Done sending test data!")

if __name__ == "__main__":
    main()
