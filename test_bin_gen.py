import struct
import sys

"""
Generates a binary file with increasing 4-byte integers.
It adds 4 0xE bytes after each integer.
The script argument is the resulting file size in KB.

Args:
    filename: The name of the output binary file.
    num_integers: The number of integers to write.

`hexdump -C test_bin_gen.bin | less`
"""
def generate_increasing_binary(filename, num_integers):
    try:
        with open(filename, "wb") as f:
            for i in range(num_integers):
                f.write(struct.pack(">I", i))  # Write i as a 4-byte unsigned int (big-endian)
                f.write(struct.pack("BBBB", 0xEE, 0xEE, 0xEE, 0xEE))  # Scrive 4 byte con valore 0xE

        print(f"Binary file '{filename}' generated with {num_integers} integers.")

    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    generate_increasing_binary("test_bin_gen.bin", int(1024*int(sys.argv[1])/4))  # 4 bytes per integer