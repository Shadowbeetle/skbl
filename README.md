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
$ systemctl --user daemon-reload
$ systemctl --user start skbl.service
```

**If the keyboard backlight is not turned back on after hitting any key, read the Iitial setup section.**

Enable the service to start `skbl` on startup

```sh
$ systemctl --user enable skbl.service
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

TODO


