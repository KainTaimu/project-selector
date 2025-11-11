### Quick directory jumper
Allows quick jumping to a directory

### Usage
```bash
# Create entries. Opens editor specified by $EDITOR or $VISUAL, whichever is first
$ project-selector -e

$ project-selector
(1) ~/Projects/my_project
(2) ~/Projects/bla
> 1

$ pwd
~/Projects/my_project
```

### Install
```bash
$ git clone https://github.com/KainTaimu/project-selector.git
...
$ cd project-selector
/project-selector$ go install .
```
