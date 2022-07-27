% podman-generate-systemd(1)

## NAME
podman\-generate\-systemd - Generate systemd unit file(s) for a container or pod

## SYNOPSIS
**podman generate systemd** [*options*] *container|pod*

## DESCRIPTION
**podman generate systemd** will create a systemd unit file that can be used to control a container or pod.
By default, the command will print the content of the unit files to stdout.

Generating unit files for a pod requires the pod to be created with an infra container (see `--infra=true`).  An infra container runs across the entire lifespan of a pod and is hence required for systemd to manage the life cycle of the pod's main unit.

_Note: If you use this command with the remote client, including Mac and Windows (excluding WSL2) machines, you would still have to place the generated units on the remote system.  Moreover, please make sure that the XDG_RUNTIME_DIR environment variable is set.  If unset, you may set it via `export XDG_RUNTIME_DIR=/run/user/$(id -u)`._

### Kubernetes Integration

A Kubernetes YAML can be executed in systemd via the `podman-kube@.service` systemd template.  The template's argument is the path to the YAML file.  Given a `workload.yaml` file in the home directory, it can be executed as follows:

```
$ escaped=$(systemd-escape ~/sysadmin.yaml)
$ systemctl --user start podman-kube@$escaped.service
$ systemctl --user is-active podman-kube@$escaped.service
active
```

## OPTIONS

#### **--after**=*dependency_name*

Add the systemd unit after (`After=`) option, that ordering dependencies between the list of dependencies and this service. This option may be specified more than once.

User-defined dependencies will be appended to the generated unit file, but any existing options such as needed or defined by default (e.g. `online.target`) will **not** be removed or overridden.

#### **--container-prefix**=*prefix*

Set the systemd unit name prefix for containers. The default is *container*.

#### **--files**, **-f**

Generate files instead of printing to stdout.  The generated files are named {container,pod}-{ID,name}.service and will be placed in the current working directory.

Note: On a system with SELinux enabled, the generated files will inherit contexts from the current working directory. Depending on the SELinux setup, changes to the generated files using `restorecon`, `chcon`, or `semanage` may be required to allow systemd to access these files. Alternatively, use the `-Z` option when running `mv` or `cp`.

#### **--format**=*format*

Print the created units in specified format (json). If `--files` is specified, the paths to the created files will be printed instead of the unit content.

#### **--name**, **-n**

Use the name of the container for the start, stop, and description in the unit file

#### **--new**

Using this flag will yield unit files that do not expect containers and pods to exist.  Instead, new containers and pods are created based on their configuration files.  The unit files are created best effort and may need to be further edited; please review the generated files carefully before using them in production.

Note that `--new` only works on containers and pods created directly via Podman (i.e., `podman [container] {create,run}` or `podman pod create`).  It does not work on containers or pods created via the REST API or via `podman kube play`.

#### **--no-header**

Do not generate the header including meta data such as the Podman version and the timestamp.

#### **--pod-prefix**=*prefix*

Set the systemd unit name prefix for pods. The default is *pod*.

#### **--requires**=*dependency_name*

Set the systemd unit requires (`Requires=`) option. Similar to wants, but declares a stronger requirement dependency.

#### **--restart-policy**=*policy*

Set the systemd restart policy.  The restart-policy must be one of: "no", "on-success", "on-failure", "on-abnormal",
"on-watchdog", "on-abort", or "always".  The default policy is *on-failure*.

#### **--restart-sec**=*time*

Set the systemd service restartsec value. Configures the time to sleep before restarting a service (as configured with restart-policy).
Takes a value in seconds.

#### **--separator**=*separator*

Set the systemd unit name separator between the name/id of a container/pod and the prefix. The default is *-*.

#### **--start-timeout**=*value*

Override the default start timeout for the container with the given value in seconds.

#### **--stop-timeout**=*value*

Override the default stop timeout for the container with the given value in seconds.

#### **--template**

Add template specifiers to run multiple services from the systemd unit file.

Note that if `--new` was not set to true, it is set to true by default. However, if `--new` is set to `false` explicitly, the command will fail.

#### **--wants**=*dependency_name*

Add the systemd unit wants (`Wants=`) option, that this service is (weak) dependent on. This option may be specified more than once. This option does not influence the order in which services are started or stopped.

User-defined dependencies will be appended to the generated unit file, but any existing options such as needed or defined by default (e.g. `online.target`) will **not** be removed or overridden.

## EXAMPLES

### Generate and print a systemd unit file for a container

Generate a systemd unit file for a container running nginx with an *always* restart policy and 1-second stop timeout to stdout. Note that the **RequiresMountsFor** option in the **Unit** section ensures that the container storage for both the GraphRoot and the RunRoot are mounted prior to starting the service. For systems with container storage on disks like iSCSI or other remote block protocols, this ensures that Podman is not executed prior to any necessary storage operations coming online.

```
$ podman create --name nginx nginx:latest
$ podman generate systemd --name --restart-policy=always --stop-timeout 1 nginx
# container-nginx.service
# autogenerated by Podman 4.1.1
# Tue Jul 26 07:38:14 CEST 2022

[Unit]
Description=Podman container-nginx.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/run/user/1000/containers

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=always
TimeoutStopSec=61
ExecStart=/usr/bin/podman start nginx
ExecStop=/usr/bin/podman stop -t 1 nginx
ExecStopPost=/usr/bin/podman stop -t 1 nginx
PIDFile=/run/user/1000/containers/overlay-containers/cbb15974fec717acd3723fc1508e95694ce1653d696f9591580d8a76ce911de3/userdata/conmon.pid
Type=forking

[Install]
WantedBy=default.target
```

### Generate systemd unit file for a container with `--new` flag

The `--new` flag generates systemd unit files that create and remove containers at service start and stop commands (see ExecStartPre and ExecStopPost service actions). Such unit files are not tied to a single machine and can easily be shared and used on other machines.

```
$ podman create --name foobar docker.io/library/alpine sleep inf
$ podman generate systemd --name --new foobar
# container-foobar.service
# autogenerated by Podman 4.1.1
# Tue Jul 26 09:43:24 CEST 2022

[Unit]
Description=Podman container-foobar.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=%t/containers

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/%n.ctr-id
ExecStart=/usr/bin/podman run \
	--cidfile=%t/%n.ctr-id \
	--cgroups=no-conmon \
	--rm \
	--sdnotify=conmon \
	-d \
	--replace \
	--name foobar docker.io/library/alpine sleep inf
ExecStop=/usr/bin/podman stop --ignore --cidfile=%t/%n.ctr-id
ExecStopPost=/usr/bin/podman rm -f --ignore --cidfile=%t/%n.ctr-id
Type=notify
NotifyAccess=all

[Install]
WantedBy=default.target
```

### Generate systemd unit files for a pod with two simple alpine containers

Note `systemctl` should only be used on the pod unit and one should not start or stop containers individually via `systemctl`, as they are managed by the pod service along with the internal infra-container.

You can still use `systemctl status` or `journalctl` to examine container or pod unit files.

```
$ podman pod create --name systemd-pod
$ podman create --pod systemd-pod alpine top
$ podman create --pod systemd-pod alpine top
$ podman generate systemd --files --name systemd-pod
/home/user/pod-systemd-pod.service
/home/user/container-amazing_chandrasekhar.service
/home/user/container-jolly_shtern.service
$ cat pod-systemd-pod.service
# pod-systemd-pod.service
# autogenerated by Podman 4.1.1
# Tue Jul 26 07:44:38 CEST 2022

[Unit]
Description=Podman pod-systemd-pod.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=
Requires=container-amazing_chandrasekhar.service container-jolly_shtern.service
Before=container-amazing_chandrasekhar.service container-jolly_shtern.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStart=/usr/bin/podman start 3ad28bb347bc-infra
ExecStop=/usr/bin/podman stop -t 10 3ad28bb347bc-infra
ExecStopPost=/usr/bin/podman stop -t 10 3ad28bb347bc-infra
PIDFile=/run/user/1092/containers/overlay-containers/2500c510c5e3f56323fcc849ca7a7737f426e22c01308da501494291e73f435e/userdata/conmon.pid
Type=forking

[Install]
WantedBy=default.target
```

### Installation of generated systemd unit files.

Podman-generated unit files include an `[Install]` section, which carries installation information for the unit. It is used by the enable and disable commands of systemctl(1) during installation.

To install a generated unit file as a systemd system service, copy it to the directory ```/etc/systemd/system``` and enable it with
`sudo systemctl enable ...`. The service will then be started at boot.

To install a generated unit file as a systemd user service, copy it to the directory ```$HOME/.config/systemd/user``` and enable it with
`systemctl --user enable ...`.  If lingering has been enabled for the user (`loginctl enable-linger <user-name>`), the service will be
started at boot, otherwise the service will be started on user login but killed after the last session for the user is closed.

```
# Generated systemd files.
$ podman pod create --name systemd-pod
$ podman create --pod systemd-pod alpine top
$ podman generate systemd --files --name systemd-pod

# Copy all the generated files.

$ sudo cp pod-systemd-pod.service container-great_payne.service /etc/systemd/system
$ sudo systemctl enable pod-systemd-pod.service
Created symlink /etc/systemd/system/multi-user.target.wants/pod-systemd-pod.service → /etc/systemd/system/pod-systemd-pod.service.
Created symlink /etc/systemd/system/default.target.wants/pod-systemd-pod.service → /etc/systemd/system/pod-systemd-pod.service.
$ systemctl is-enabled pod-systemd-pod.service
enabled
```

### Use `systemctl` to perform operations on generated installed unit files.



#### Managing a systemd user service for another user

To check the status of the systemd user service _foobar.service_ belonging to the user _testuser_

`sudo systemctl --user -M testuser@ status foobar.service`

To run the podman command `podman container list` as the user _testuser_

```
sudo systemd-run --machine=testuser@ --quiet --user --collect --pipe --wait podman container list
```





Create and enable systemd unit files for a pod using the above examples as reference and use `systemctl` to perform operations.

Since systemctl defaults to using the root user, all the changes using the systemctl can be seen by appending sudo to the podman cli commands. To perform `systemctl` actions as a non-root user, use the `--user` flag when interacting with `systemctl`.

Note: If the previously created containers or pods are using shared resources, such as ports, make sure to remove them before starting the generated systemd units.

```
$ systemctl --user start pod-systemd-pod.service
$ podman pod ps
POD ID         NAME          STATUS    CREATED          # OF CONTAINERS   INFRA ID
0815c7b8e7f5   systemd-pod   Running   29 minutes ago   2                 6c5d116f4bbe
$ sudo podman ps # 0 Number of pods on root.
CONTAINER ID  IMAGE  COMMAND  CREATED  STATUS  PORTS  NAMES
$ systemctl stop pod-systemd-pod.service
$ podman pod ps
POD ID         NAME          STATUS   CREATED          # OF CONTAINERS   INFRA ID
272d2813c798   systemd-pod   Exited   29 minutes ago   2                 6c5d116f4bbe
```

Create a simple alpine container and generate the systemd unit file with `--new` flag.
Enable the service and control operations using the systemctl commands.

Note: When starting the container using `systemctl start` rather than altering the already running container, it spins up a "new" container with similar configuration.

```
# Enable the service.

$ sudo podman ps -a
CONTAINER ID  IMAGE                            COMMAND  CREATED        STATUS     PORTS  NAMES
bb310a0780ae  docker.io/library/alpine:latest  /bin/sh  2 minutes ago  Created           busy_moser
$ sudo systemctl start container-busy_moser.service
$ sudo podman ps -a
CONTAINER ID  IMAGE                            COMMAND  CREATED        STATUS            PORTS      NAMES
772df2f8cf3b  docker.io/library/alpine:latest  /bin/sh  1 second ago   Up 1 second ago              distracted_albattani
bb310a0780ae  docker.io/library/alpine:latest  /bin/sh  3 minutes ago  Created                      busy_moser
```
## SEE ALSO
**[podman(1)](podman.1.md)**, **[podman-container(1)](podman-container.1.md)**, **systemctl(1)**, **systemd.unit(5)**, **systemd.service(5)**, **[conmon(8)](https://github.com/containers/conmon/blob/main/docs/conmon.8.md)**

## HISTORY
April 2020, Updated details and added use case to use generated .service files as root and non-root, by Sujil Shah (sushah at redhat dot com)

August 2019, Updated with pod support by Valentin Rothberg (rothberg at redhat dot com)

April 2019, Originally compiled by Brent Baude (bbaude at redhat dot com)
