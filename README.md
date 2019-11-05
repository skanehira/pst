# pst
This is TUI process monitor written in Go.

![](https://i.imgur.com/TsrokJ7.gif)

## Features
- Monitor process's list, info, tree, open files,
- Kill process

## Support OS
- Mac
- Linux

## Requirements
- ps
- lsof

## Installation
```sh
$ git clone https://github.com/skanehira/pst
$ cd pst
$ go install
```

## Options
You can change the process info to be displayed with environment `PS_ARGS`.

Default `PS_ARGS` value is `pid,ppid,%cpu,%mem,lstart,user,command`.

e.g make alias and use it.

```sh
alias pst="env PS_ARGS=%cpu,%mem,lstart pst"
```

## Usage
```sh
$ pst -h
Usage of pst:
  -log
        enable output log
  -proc string
        use word to filtering process name when starting

# run tui
$ pst
```

Default, log file will generate `$HOME/pst.log` if it's not exist.

## Keybindings
### common keybindings
| key         | description          |
|-------------|----------------------|
| Ctrl + c    | stop pst             |
| j           | move down            |
| k           | move up              |
| h           | move left            |
| l           | move right           |
| g           | move to top          |
| G           | move to bottom       |
| Ctrl + f    | next page            |
| Ctrl + b    | previous page        |
| Tab         | focus next panel     |
| Shift + Tab | focus previous panel |

### input
| key         | description          |
|-------------|----------------------|
| Enter       | next process         |

### processes panel
| key         | description          |
|-------------|----------------------|
| K           | kill select process  |

### process tree panel
| key         | description          |
|-------------|----------------------|
| K           | kill select process  |
| Enter       | expand child process |
