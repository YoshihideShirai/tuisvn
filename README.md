# Tuisvn: Terminal user interface for Subversion

## What is tuisvn

Tuisvn is subversion(svn) client.  
It functions mainly as a svn repository browser like tig(https://github.com/jonas/tig).

![](./docs/screenshot.png)

## System requirement

- Linux
- Installed svn command.

## Installation

~~~shell
go install github.com/YoshihideShirai/tuisvn@latest
~~~

## Key bindings

### tree screen

- k : Move-up
- j : Move-down
- Enter : change directory
- l : Enter log screen on selected file and directory.
- q : Back previous screen.

### log screen

- k : Move-up
- j : Move-down
- Enter : Enter rev screen
- q : Back previous screen.

### rev screen

- k : Move-up
- j : Move-down
- Enter : Enter diff screen on change directory.
- q : Back previous screen.

### diff screen

- k : Move-up
- j : Move-down
- q : Back previous screen.

## Usage

~~~shell
tuisvn [repository path or url]
~~~
If argv is omitted, repository path is set the current directory.

## Development status

Tuisvn is under development.  
Implemented status is following...
- [x] file and directory tree. (svn ls)
- [x] revision history. (svn log)
- [x] revision diff
- [ ] revision cat

## License

Copyright (C) 2022 [Yoshihide Shirai](mailto:yoshihide.shirai@gmail.com).

Licensed under the [MIT License](LICENSE).
