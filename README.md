Turn off keyboard backlight after a certain period of idleness

# Requirements

- linux
- systemd
- go 1.12 or above

# Installation

In order to get a full setup clone the repo than then run make install.

```sh
$ git clone git@github.com:Shadowbeetle/skbl.git
$ make install
```

This installs the binary, creates a user systemd service, adds your user to the group `input` and creates the necessary config files.

To start using `skbl` with the default setup simply run:

```sh
$ systemctl daemon-reload
$ systemctl ser start skbl@$USER.service
```

**If the keyboard backlight is not turned back on after hitting any key, read the [Initial setup](#initial-setup) section.**

Enable the service to start `skbl` on startup

```sh
$ systemctl enable skbl@$USER.service
```

You can of course use `skbl` by simply running go get and set everythin up yourself by hand as well.

```sh
$ go get github.com/Shadowbeetle/skbl
```

# Initial setup

`skbl` listens to events from `/dev/input/mice` and `/dev/input/event4`. Your keyboard might be a mapped to a different event, so you probably need to change that.

1. you'll need to find the approprate iput file:

```sh
$ ls -lah /dev/input/by-path/
```

2. find the ones that might be keyboards (usually have kbd in them)
3. test with `cat /dev/input/by-path/<kbd-input>`
4. if you see gibberish appear in the terminal, you're good.

# Configuration

`skbl` can be configured by the following ways (consecutive modes override the previous ones):

1. `/etc/skbl/config.toml` is the system wide default config
2. `$HOME/.skbl/config.toml` is the config of the given user 
3. using flags

## Conifg file

Installing `skbl` with `make install` creates the default config file in `/etc/skbl`, and the user config in `$HOME/.skbl` for the user running `make`. Additional users need to create their own config files.

```sh
$ mkdir $HOME/.skbl
$ cp /etc/skbl/config.toml $HOME/.skbl
```

The config file currently has two fields: 

```toml
wait-seconds = "10s" # idle time after which backlight should be turned off
inputs = ["/dev/input/mice", "/dev/input/event4"] # input files to listen to
```

for more information on inputs see [Initial setup](#initial-setup).

## Flags

Individual sessions can be configured using flags as well eg.

```sh
$ skbl --wait 1s --input /dev/input/event1 --input /dev/input/mice
# or
$ skbl -w 1s -i /dev/input/event1 -i /dev/input/mice
```

