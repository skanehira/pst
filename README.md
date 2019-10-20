# pst
This is TUI process monitor written in Go.

![](https://i.imgur.com/TUcGDfk.gif)

## Installation
```sh
$ git clone https://github.com/skanehira/pst
$ go install
```

## Support OS
- Mac
- Linux

## Usage
```sh
$ pst -h
Usage of pst:
  -log
        enable output log

# run tui
$ pst
```

Default, log file will generate `$HOME/pst.log` if it's not exist.

## Keybindings
### common keybinds
| key      | description |
|----------|-------------|
| Ctrl + c | stop pst    |

### input
| key         | description          |
|-------------|----------------------|
| Enter       | next process         |
| Tab         | focus next panel     |
| Shift + Tab | focus previous panel |

### processes panel
| key         | description          |
|-------------|----------------------|
| j           | next process         |
| k           | previous process     |
| g           | first process        |
| G           | last process         |
| K           | kill select process  |
| Tab         | focus next panel     |
| Shift + Tab | focus previous panel |

### process tree panel
| key         | description          |
|-------------|----------------------|
| j           | next process         |
| k           | previous process     |
| g           | first process        |
| G           | last process         |
| K           | kill select process  |
| Enter       | expand child process |
| Tab         | focus next  panel    |
| Shift + Tab | focus previous panel |
