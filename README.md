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

- Run the container-shifter binary to pull down the desired images with `./container-shifter pull --config-file /path/to/config.yml`. Pulls will happen in parallel.

- Push images to the private registries with `./container-shifter push --config-file /path/to/config.yml --docker-config /path/to/docker/config.json`. Pushes will happen in parallel for each registry, i.e. all images are pushed at once to the first registry, then the second, and so on.

Both of the config flags are optional. Docker one defaults to `$HOME/.docker/config.json`. The push/pull steps will eventually be combined into an all-in-one command.

### Dealing With Airgapped Environments

One of the reasons I wrote this tool is to try and help ease the burden of dealing with an airgapped environment. A system without access to the public internet needs to have a fully private registry. Here's how you might use the built-in commands to accomplish this

From a system with internet access:

- Create your config.yml to specify the public images.  Also include your internal registry, we'll use it in the next step.
- Pull the docker images you specify in your config file with `./container-shifter pull --config-file /path/to/config.yml`.
- Save the images to a local path with `./container-shifter save --config-file /path/to/config.yml --image-directory /desired/path/for/image/tars/` (ensure the trailing slash in img directory for now)
- Throw the tars directory and container-shifter on a USB

From an airgapped system with access to the internal registry:

- Mount up the USB
- Import the docker images with `./container-shifter load --config-file /path/to/config.yml --image-directory /desired/path/for/image/tars/` (ensure the trailing slash in img directory for now)
- Docker login to the internal registry if needed.
- Push to the internal registry with `./container-shifter push --config-file /path/to/config.yml --docker-config /path/to/docker/config.json`

## Building

This tool can be easily built using the provided Makefile. Simply issue one of the following:
- `make darwin` - Builds 64-bit Mac client in the bin/ directory
- `make linux` - Builds 64-bit Linux client in the bin/ directory
- `make all` - Builds both of the above in the bin/ directory

## Next Steps

- Combination command for download/push
- Improve interaction with registries, search for tags and that kind of thing
- Provide a "daemon mode" that monitors the desired containers and if new versions pop up, pull/push
- Support <,>,= notations for container version tags?