### Network performance for rootless Podman

* Using the network driver _bypass4netns_. It has almost the same performance characteristics as the normal network on the host.
  Compared to _socket activation_ this alternative also supports outgoing connections without slirp4netns.
  Note the security concern mentioned on the [bypass4netns](https://github.com/rootless-containers/bypass4netns) home page that
  _it's probably possible to connect to host loopback IPs by exploiting [TOCTOU](https://elixir.bootlin.com/linux/v5.9/source/include/uapi/linux/seccomp.h#L81) of `struct sockaddr *` pointers_.
  Other [concerns](https://github.com/rootless-containers/bypass4netns/issues/1#issuecomment-1027948113)
  have been raised, but there is an [implementation idea](https://github.com/rootless-containers/bypass4netns/issues/21)
  of how to solve the security problem.

### Mapping UIDs and GIDs instead of using `:U` volume flag

When `podman run` is passed a __--volume__ (__-v__) option with the `:U` volume flag, Podman
will change the ownership of the files and directories within that directory so that it matches
the UID/GID of the running container user. The volume flag `:U` is especially useful when used together with `--userns=auto`
because a file in the volume might have an ownership UID/GID that is not mapped into the container.
Without any `chown()`, such a file ownership would be represented as _nobody:nobody_ (and as _nobody:nogroup_ in Debian and Ubuntu) inside the container.

When `--auto` is used and the container expects a file to be owned by the running container user, using the `:U` volume flag is required.

When `--auto` is _not_ used and the file ownership is mapped to a container user, but not to the expected running container user, using the `:U` volume flag is not needed. In this case the UID/GID mapping can be adjusted by [__--uidmap__](https://docs.podman.io/en/latest/markdown/podman-run.1.html#uidmap-container-uid-from-uid-amount) and [__--gidmap__](https://docs.podman.io/en/latest/markdown/podman-run.1.html#gidmap-container-gid-host-gid-amount).

See the troubleshooting tips of how to remap UIDs and GIDs:

* [Container creates a file that is not owned by the user's regular UID](https://github.com/containers/podman/blob/main/troubleshooting.md#33-container-creates-a-file-that-is-not-owned-by-the-users-regular-uid)

* [Passed-in devices or files can't be accessed in rootless container (UID/GID mapping problem)](https://github.com/containers/podman/blob/main/troubleshooting.md#34-passed-in-devices-or-files-cant-be-accessed-in-rootless-container-uidgid-mapping-problem)

Recursively changing ownership in a directory could take time if the directory contains
many files and subdirectories. Avoiding the `:U` volume flag could thus save some time.

### Using smaller container images

Container image sizes can vary a lot. Podman needs to perform less work when starting smaller container images.
To build a smaller container image, clean up the package cache after installing a package.
When installing a package, consider also if it can be installed without weaker dependencies, install recommendations and documentation.

packaging tool | command to install a package more space efficiently
-----------    | -----------------
dnf            | `RUN dnf update && dnf install -y --nodocs --setopt install_weak_deps=False packageName && dnf clean`
apt-get        | `RUN apt-get update && apt-get install -y --no-install-recommends packageName && rm -rf /var/lib/apt/lists/*`
apk            | `RUN apk --no-cache add packageName`

In case your Containerfile makes use of build tools to build the application, consider writing a
multi-stage Containerfile so that build tools can be avoided in the resulting container image.

### Whether to use --new

Using --new option is useful to be sure that the container is always started in the same well-defined prestine state.
Creating a new container takes some time.

### Option --pull

The `podman run` option [__--pull__](https://docs.podman.io/en/latest/markdown/podman-run.1.html#pull-policy) specifies when to pull the container image.

Using `--pull=always` and `--pull=newer` will connect to the container registry.
This is useful if you want to run the newest version of the container image (e.g. for staying up to date to the latest security updates).
In case it's not necessary consider using `--pull=never` to skip the pulling of the container image.
