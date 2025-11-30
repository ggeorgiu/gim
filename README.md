## GIM

Vim-like text editor built in Go. Lightweight, for-fun, personal project.

### Modes

Just like vim, the editor has modes of operating it, enabling/disabling features.

- NORMAL
    - In this mode you can move the cursor throughout the editor. This is the default mode.
    - Pressing `Esc` will enter this mode from whatever mode you're into.
- COMMAND
    - Allows executing commands. See [Commands](#Commands)
    - Press `:` while in `NORMAL` mode to enter it.
- INSERT
    - This is the mode that allows inserting/deleting text.
    - Press `i` while in `NORMAL` mode to enter it.

### Components

#### Status line

- current mode
- cursor position
- show error

#### Line number

- current line number
- line number relative to cursor

## Commands

- `:q` - exit
- `:w` - write contents to file
