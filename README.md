# CLI bookmark tool

Allows for bookmarking directories on the terminal

## Usage

```bash
# Create entries. Opens editor specified by $EDITOR or $VISUAL, with priority given to $EDITOR
$ bookmark -e

$ bookmark
(1) ~/Projects/my_project
(2) ~/Projects/bla
> 1

$ pwd
~/Projects/bookmark
```

## Install

```bash
$ git clone https://github.com/KainTaimu/bookmark.git
...
$ cd bookmark
/bookmark$ go install .
```
