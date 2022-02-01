# gmtree

Generate treeview for command "`go mod graph`" .

## Requirements

* Go: 1.16+
* Go Mods: [go.mod](go.mod)

## Install

```bash
$ go get -u https://github.com/rujax/gmtree
```

## Argument

| Name | Alias | Type | Default Value |
| --- | --- | --- | --- |
| --indent | -i | int | 2 |

## Usage

Print treeview on Stdout

```bash
$ go mod graph | gmtree # Indent: 2
$ go mod graph | gmtree -i n # Indent: n
```

Save treeview to file

```bash
$ go mod graph | gmtree > treeview_file_path
```

## Example

Mac

![example_mac.png](example_mac.png)
