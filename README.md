# lassie-go
Lassie provides a client for the [REST API](https://api.lora.telenor.io) for
[Telenor LoRa](https://lora.engineering).

## Configuration

The configuration file is located at `${HOME}/.lassie`. The file is a simple
list of key/value pairs. Additional values are ignored. Comments must start
with a `#`:

    #
    # This is the URL of the Congress REST API. The default value is
    # https://api.lora.telenor.io and can usually be omitted.
    address=https://api.lora.telenor.io

    #
    # This is the API token. Create new token by logging in to the Congress
    # front-end at https://lora.engineering and create a new token there.
    token=<your api token goes here>


The configuration file settings can be overridden by setting the environment
variables `LASSIE_ADDRESS` and `LASSIE_TOKEN`. If you only use environment variables
the configuration file can be ignored.

Use the `NewWithAddr` function to bypass the default configuration file and
environment variables when you want to create the client programmatically.
