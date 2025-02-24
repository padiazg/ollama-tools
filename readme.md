# Ollama tools
A few tools to iteract with your local Ollama server.

This project started as a couriosity to know how much memory do I need to run the models I had downloaded as some where almost unusable no matter the specs from my hardware.

I must mention that I'm not an expert in this matter, and most of the logics and formulas were inspired by the code from this repo [ollama-gpu-calculator](https://github.com/aleibovici/ollama-gpu-calculator)

## The theory
Check the theory [here](theory.md)

## Install
Clone the repo
```shell
git clone https://github.com/padiazg/ollama-tools.git
cd ollama-tools
```
Install dependencies and build
```shell
go mod tidy
go build
```

## Usage
**List downloaded models**
```shell
$ ollama-tools list-models --help
List models using the Ollama api

If no model-name is espified all models will be retieved and listed.
You can pass the model-name as an argument or using the --model-name flag

Usage:
  ollama-tools list-models [model-name] [flags]

Flags:
  -h, --help                help for list-models
  -m, --model-name string   Model to list
  -t, --table               Print as table

Global Flags:
      --config string   config file (default is $HOME/.ollama-tools.yaml)
```
Example
```shell
# list all models downloaded by ollama
$ ollama-tools list-models
Available models:
----------------------------------------------------
Model: nomic-embed-text:latest
  Parameters: 136.73M (136727040)
  Quantization: F16
  Context Length: 2048 tokens
  Embedding Length: 768

  Memory Breakdown:
    Model Weights Memory: 0.25 MB
    KV Cache (for context): 0.07 MB
    GPU VRAM: 0.35 MB
    System RAM: 0.71 MB

Model: llama3.1:latest
  Parameters: 8.03B (8030261312)
  Quantization: Q4_K_M
  Context Length: 131072 tokens
  Embedding Length: 4096

  Memory Breakdown:
    Model Weights Memory: 3.74 MB
    KV Cache (for context): 8.93 MB
    GPU VRAM: 13.04 MB
    System RAM: 14.35 MB

Note: This model has a large context length (131072 tokens).
Reducing max_context in your Ollama request can significantly lower memory usage.

Model: deepseek-r1:14b
  Parameters: 14.77B (14770033664)
  ...
```
List an especific model
```shell
# pass the model as a parameter
$ ollama-tools list-models phi4:latest
# or with a flag
$ ollama-tools list-models --model-name phi4:latest
Model: phi4:latest
  Parameters: 14.66B (14659507200)
  Quantization: Q4_K_M
  Context Length: 16384 tokens
  Embedding Length: 5120

  Memory Breakdown:
    Model Weights Memory: 6.83 MB
    KV Cache (for context): 1.51 MB
    GPU VRAM: 9.02 MB
    System RAM: 9.92 MB

Note: This model has a large context length (16384 tokens).
Reducing max_context in your Ollama request can significantly lower memory usage.
```
There's also an option to list the models in a table using the `--table` or `-t` flags
```shell
$ ollama-tools list-models --table
+-------------------------+-----------------------------+-------------------+----------------+------------------+-----------------+------------+--------------+--------------+
| MODEL                   |          PARAMETERS         |    QUANTIZATION   | CONTEXT LENGTH | EMBEDDING LENGTH | BASE MODEL SIZE | KV CACHE   | GPU RAM      | SYSTEM RAM   |
|                         | BILLIONS | UNITS            | LEVEL    | BITS   |                |                  |                 |            |              |              |
+-------------------------+----------+------------------+----------+--------+----------------+------------------+-----------------+------------+--------------+--------------+
| nomic-embed-text:latest |  136.73M |        136727040 | F16      |     16 |           2048 |              768 |         0.25 Gb |    0.07 Gb |      0.35 Gb |      0.71 Gb |
| llama3.1:latest         |    8.03B |       8030261312 | Q4_K_M   |      4 |         131072 |             4096 |         3.74 Gb |    8.93 Gb |     13.04 Gb |     14.35 Gb |
| deepseek-r1:14b         |   14.77B |      14770033664 | Q4_K_M   |      4 |         131072 |             5120 |         6.88 Gb |   12.11 Gb |     19.68 Gb |     21.65 Gb |
| deepseek-r1:latest      |    7.62B |       7615616512 | Q4_K_M   |      4 |         131072 |             3584 |         3.55 Gb |    8.70 Gb |     12.60 Gb |     13.86 Gb |
| phi4:latest             |   14.66B |      14659507200 | Q4_K_M   |      4 |          16384 |             5120 |         6.83 Gb |    1.51 Gb |      9.02 Gb |      9.92 Gb |
+-------------------------+----------+------------------+----------+--------+----------------+------------------+-----------------+------------+--------------+--------------+

$ ollama-tools list-models --model-name phi4:latest --table
Using config file: /Users/pato/.ollama-tools.yaml
+-------------+-----------------------------+-------------------+----------------+------------------+-----------------+------------+--------------+--------------+
| MODEL       |          PARAMETERS         |    QUANTIZATION   | CONTEXT LENGTH | EMBEDDING LENGTH | BASE MODEL SIZE | KV CACHE   | GPU RAM      | SYSTEM RAM   |
|             | BILLIONS | UNITS            | LEVEL    | BITS   |                |                  |                 |            |              |              |
+-------------+----------+------------------+----------+--------+----------------+------------------+-----------------+------------+--------------+--------------+
| phi4:latest |   14.66B |      14659507200 | Q4_K_M   |      4 |          16384 |             5120 |         6.83 Gb |    1.51 Gb |      9.02 Gb |      9.92 Gb |
+-------------+----------+------------------+----------+--------+----------------+------------------+-----------------+------------+--------------+--------------+
```

**Estimate** 
We can estimate the RAM requirement without downloading the model. You must get some values from the model's page and feed it to the app. This comes in handy to download only those models our setup would handle.
Values:
- _parameter count_: in units, not in billions or millions. 
- _context length_: in units, not in kilos.
- _quantization level_: the string as found in the page (Q4_K_M, Q4_K_S, F16, F32, etc). The app will try to trabslate it to bits.
```shell
$ ollama-tools estimate --help
Estimates the RAM rwquirement based on few parameters without the need to download any model

Usage:
  ollama-tools estimate [flags]

Flags:
  -c, --context-length int          Context length
  -h, --help                        help for estimate
  -p, --parameter-count int         Parameters count
  -q, --quantization-level string   Quantization level

Global Flags:
      --config string   config file (default is $HOME/.ollama-tools.yaml)
```
Example
```shell
$ ollama-tools estimate -p 8030261312 -c 131072 -q Q4_K_M
  Memory Breakdown:
    Model Weights Memory: 3.74 MB
    KV Cache (for context): 8.93 MB
    GPU VRAM: 13.04 MB
    System RAM: 14.35 MB
```

## Config
There is no need to config anything as long as your ollama is running in your localhost, but if it's running elsewhere you can set the app to look fo it like this
```shell
OT_OLLAMAURL=http:192.168.1.100:11434 ./ollama-tools list-models
```
Or use a configuration file. Create a file `~/.ollama-tools.yaml` and add this to it
```yaml
ollamaurl: http:192.168.1.100:11434
```
Then use the app as usual

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Buy me a coffee
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://buymeacoffee.com/padiazgy)

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
