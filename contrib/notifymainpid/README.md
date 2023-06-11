
### systemd service with User= (a proof of concept)

Run a socket-activated systemd system service with the systemd
directive `User=`. 

Relevant feature request:

* https://github.com/containers/podman/issues/12778

Demo:



On a Fedora 38 system run these commands
```
dest=/var/tmp/notify-mainpid
podman build -t notify .
podman create --name tmpctr localhost/notify
podman cp tmpctr:/notify-mainpid $dest
chmod 755 $dest
sudo chown root:root $dest
sudo chcon --reference /usr/bin/cp $dest
```


Disable SELinux
```
sudo setenforce 0
```
(Maybe disabling SELinux is not needed ?)
In any case, remember to enable SELinux after trying out the demo.


Set up the service
```
sudo useradd test
sudo cp echo.service /etc/systemd/system
sudo cp echo.socket /etc/systemd/system
sudo systemctl daemon-reload
```


Open a shell to create the podman user namespace
```
sudo machinectl shell test@
podman unshare /bin/true
```
(let this shell keep on running)

Rationale: __notify-mainpid__ uses the environment variable `SYSTEMD_EXEC_PID`. If the podman user namespace
is not already available `SYSTEMD_EXEC_PID` will not be equal to the normal Podman PID (see https://github.com/containers/podman/discussions/18842). That is why we would like to have the podman user namespace already created before starting the echo service.


open another shell
```
sudo systemctl start echo.socket
sudo rm /var/tmp/conmon-pidfile
```


```
$ echo hello | socat  -t 60 - tcp4:127.0.0.1:908
echo
```

If the 60 seconds is not enough, increase the value. The first time the service is started, the container image
might be pulled which could take some time.


Instead of connecting to the socket, another way to start the echo service is to run

```
sudo systemctl start echo.service
```

## possible future improvements

Just some ideas.

* rewrite notify-mainpid.cpp . Probably C is more suitable than C++.  notify-mainpid should check the ownership of the PID and ownership of conmon-pidfile matches the user.
* remove the racy solution with a sleep in container_internal.go. Introduce some other synchronization. (Maybe using `OpenFile=` of a root-owned file that podman can detect that the file has been edited by notify-mainpid)
* maybe make use of OpenFile= that podman could write the conmon-pid to.
* add another example that handles `--sdnotify=container` (for example with `/socket-activate-echo --sdnotify`)

Solving  `--sdnotify=container` will be more complicated.
