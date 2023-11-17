# PiedPiper

PiedPiper, the name is totally original and has nothing to do with any TV shows, is a compression tool using [Huffman Coding Trees](https://en.wikipedia.org/wiki/Huffman_coding).\
I only built it for learning purposes, so use it at your own risk :D

## Usage

- Download the binary for your platform:
  - `https://github.com/MarioBassem/pied_piper/releases`
- or, clone and build the application:
  
    ```bash
    git clone https://github.com/MarioBassem/pied_piper.git
    make build
    ```
  
- To encode a file:

    ```bash
    piedpiper path/to/decompressed_file path/to/generated/compressed_file
    ```

- To decode a file:

    ```bash
    piedpiper -decode path/to/compressed_file path/to/generated/decompressed_file
    ```

## Testing

To run tests:

```bash
make test
```
