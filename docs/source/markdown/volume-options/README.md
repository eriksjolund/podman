# Mount options

Proof of concept

Sketchy design idea for generating documentation about mount options
to fix https://github.com/containers/podman/issues/24249

The file [_mount-options.json_](mount-options.json) contains the data source
for creating [_output.md_](output.md)

The generated file [_output.md_](output.md) contains

* a compatibility table showing how mount types and mount options can be combined
* a compatibility table showing how many arguments each option takes
* a list of mount types
* a list of mount options

## Generate Markdown

1. Build container image
   ```
   podman build -t generate .
   ```
2. Build container image
   ```
   podman build -t generate .
   ```
3. Generate _output.md_
   ```
   cat mount-options.json | podman run --rm -i generate > output.md
   ```

