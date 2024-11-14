import time
import sys
import string
import random
import pyaes


def generate(length):
    letters = string.ascii_lowercase + string.digits
    return ''.join(random.choice(letters) for i in range(length))

def f(n):

    start = round(time.time(),6)

    length_of_message = int(n)
    num_of_iterations = int(n)

    message = generate(length_of_message)

    # 128-bit key (16 bytes)
    KEY = b'\xa1\xf6%\x8c\x87}_\xcd\x89dHE8\xbf\xc9,'

    for loops in range(num_of_iterations):
        aes = pyaes.AESModeOfOperationCTR(KEY)
        ciphertext = aes.encrypt(message)
        print(ciphertext)

        aes = pyaes.AESModeOfOperationCTR(KEY)
        plaintext = aes.decrypt(ciphertext)
        aes = None
    end = round(time.time(),6)
    #print("{} {}, {}, {}, {}, {}".format(id_n, n, pred, start, end, end-start))
    #print(result)
def handler(event, context):
    # n = event.get("num", 0)
    return f(1)