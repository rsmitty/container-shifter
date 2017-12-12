# container-shifter
A binary to grab docker images and import them into private registries

## Usage

This tool is currently a WIP, but it works for a very basic use case.

- Generate a config.yml file. You can put this anywhere you like, but the location lookup defaults to wherever you're running the container-shifter binary from.
```yaml
containers:
 - "docker.io/rsmitty/routereflector"
 - "docker.io/rsmitty/boinc"

registries:
 - "registry.rsmitty.xyz"
```
It's important to note that you'll currently need to specify full paths and versions for the containers you want to pull. This includes quay.io/gcr.io/docker.io, as well as any non-latest tags you want.

- `docker login` to any registries as necessary. You need the auth info added to your config.json for docker.

- Run the container-shifter binary to pull down the desired images with `./container-shifter pull --config-file /path/to/config.yml`

- Push images to the private registries with `./container-shifter pull --config-file /path/to/config.yml --docker-config /path/to/docker/config.json`

Both of the config flags are optional. Docker one defaults to `$HOME/.docker/config.json`. The push/pull steps will eventually be combined into an all-in-one command.

## Building

This tool can be easily built using the provided Makefile. Simply issue one of the following:
- `make darwin` - Builds 64-bit Mac client in the bin/ directory
- `make linux` - Builds 64-bit Linux client in the bin/ directory
- `make all` - Builds both of the above in the bin/ directory
