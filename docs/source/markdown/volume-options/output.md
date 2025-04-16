| |artifact | bind | devpts | glob | image | ramfs | tmpfs | volume |
| --- |--- | --- | --- | --- | --- | --- | --- | --- |
| bind-nonrecursive |   | x |   | x |   |   |   |   |
| bind-propagation |   | x |   | x |   |   |   |   |
| chown(U) |   | x |   | x |   |   |   |   |
| destination(dst,target) | x | x | x | x | x | x | x | x |
| digest | x |   |   |   |   |   |   |   |
| gid |   |   | x |   |   |   |   |   |
| idmap |   | x |   | x |   |   |   | x |
| max |   |   | x |   |   |   |   |   |
| mode |   |   | x |   |   |   |   |   |
| no-dereference |   | x |   | x |   |   |   |   |
| notmpcopyup |   |   |   |   |   | x | x |   |
| ptmxmode |   |   | x |   |   |   |   |   |
| readonly(ro) |   | x |   | x |   |   |   |   |
| readwrite(rw) |   | x |   | x |   |   |   |   |
| relabel |   | x |   | x |   |   |   |   |
| rprivate |   | x |   | x |   |   |   |   |
| rshared |   | x |   | x |   |   |   |   |
| rslave |   | x |   | x |   |   |   |   |
| shared |   | x |   | x |   |   |   |   |
| slave |   | x |   | x |   |   |   |   |
| source(src) | x | x |   | x | x |   |   | x |
| subpath |   |   |   |   | x |   |   | x |
| title | x |   |   |   |   |   |   |   |
| tmpcopyup |   |   |   |   |   | x | x |   |
| tmpfs-mode |   |   |   |   |   | x | x |   |
| tmpfs-size |   |   |   |   |   | x | x |   |


| | 0 arg allowed | 1 arg allowed |
| - | - | - |
| bind-nonrecursive | x | x | 
| bind-propagation |  | x | 
| chown | x | x | 
| destination |  | x | 
| digest |  | x | 
| gid |  | x | 
| idmap | x | x | 
| max |  | x | 
| mode |  | x | 
| no-dereference | x | x | 
| notmpcopyup | x | x | 
| ptmxmode |  | x | 
| readonly | x | x | 
| readwrite | x | x | 
| relabel | x | x | 
| rprivate | x |  | 
| rshared | x |  | 
| rslave | x |  | 
| shared | x |  | 
| slave | x |  | 
| source |  | x | 
| subpath |  | x | 
| title |  | x | 
| tmpcopyup | x | x | 
| tmpfs-mode |  | x | 
| tmpfs-size |  | x | 


### Types for `--mount`:

#### artifact

artifact documentation.
Available options: _destination_, _digest_, _source_, _title_

#### bind

bind documentation.
Available options: _bind-nonrecursive_, _bind-propagation_, _chown_, _destination_, _idmap_, _no-dereference_, _readonly_, _readwrite_, _relabel_, _rprivate_, _rshared_, _rslave_, _shared_, _slave_, _source_

#### devpts

devpts documentation.
Available options: _destination_, _gid_, _max_, _mode_, _ptmxmode_

#### glob

globs documentation.
Available options: _bind-nonrecursive_, _bind-propagation_, _chown_, _destination_, _idmap_, _no-dereference_, _readonly_, _readwrite_, _relabel_, _rprivate_, _rshared_, _rslave_, _shared_, _slave_, _source_

#### image

image documentation.
Available options: _destination_, _source_, _subpath_

#### ramfs

ramfs documentation.
Available options: _destination_, _notmpcopyup_, _tmpcopyup_, _tmpfs-mode_, _tmpfs-size_

#### tmpfs

tmpfs documentation.
Available options: _destination_, _notmpcopyup_, _tmpcopyup_, _tmpfs-mode_, _tmpfs-size_

#### volume

volume documentation.
Available options: _destination_, _idmap_, _source_, _subpath_


### Options for `--mount`:

#### bind-nonrecursive

Do not set up a recursive bind mount. By default it is recursive.
Available for types: _bind_, _glob_

#### bind-propagation

See also mount(2).
Available for types: _bind_, _glob_

#### chown(U)

Recursively change the owner and group of the source volume based on the UID and GID of the container.
Available for types: _bind_, _glob_

#### destination(dst,target)

Mount destination spec. When the destination is specified, the files and directories matching the glob on the base file name on the destination directory are mounted. The option `type=glob,src=/foo*,destination=/tmp/bar` tells container engines to mount host files matching /foo* to the /tmp/bar/ directory in the container.
Available for types: _artifact_, _bind_, _devpts_, _glob_, _image_, _ramfs_, _tmpfs_, _volume_

#### digest

If the artifact source contains multiple blobs a digest can be specified to only mount the one specific blob with the digest.
Available for types: _artifact_

#### gid

numeric GID of the file owner
Available for types: _devpts_

#### idmap

*true* or *false* (default if unspecified: *false*).  If true, create an idmapped mount to the target user namespace in the container. The idmap option is only supported by Podman in rootful mode.
Available for types: _bind_, _glob_, _volume_

#### max

maximum number of PTYs
Available for types: _devpts_

#### mode

octal permission mask for the file (default: 600).
Available for types: _devpts_

#### no-dereference

do not dereference symlinks but copy the link source into the mount destination.
Available for types: _bind_, _glob_

#### notmpcopyup

Disable copying files from the image to the tmpfs/ramfs.
Available for types: _ramfs_, _tmpfs_

#### ptmxmode


Available for types: _devpts_

#### readonly(ro)


Available for types: _bind_, _glob_

#### readwrite(rw)


Available for types: _bind_, _glob_

#### relabel


Available for types: _bind_, _glob_

#### rprivate

The same as `bind-propagation=rprivate`.
Available for types: _bind_, _glob_

#### rshared

The same as `bind-propagation=rshared`.
Available for types: _bind_, _glob_

#### rslave

The same as `bind-propagation=rslave`.
Available for types: _bind_, _glob_

#### shared

The same as `bind-propagation=shared`.
Available for types: _bind_, _glob_

#### slave

The same as `bind-propagation=slave`.
Available for types: _bind_, _glob_

#### source(src)

mount source spec
Available for types: _artifact_, _bind_, _glob_, _image_, _volume_

#### subpath


Available for types: _image_, _volume_

#### title

If the artifact source contains multiple blobs a title can be set which is compared against `org.opencontainers.image.title` annotation.
Available for types: _artifact_

#### tmpcopyup

Enable copyup from the image directory at the same location to the tmpfs/ramfs. Used by default.
Available for types: _ramfs_, _tmpfs_

#### tmpfs-mode

Octal file mode of the tmpfs/ramfs (e.g. 700 or 0700.).
Available for types: _ramfs_, _tmpfs_

#### tmpfs-size

Size of the tmpfs/ramfs mount, in bytes. Unlimited by default in Linux.
Available for types: _ramfs_, _tmpfs_
