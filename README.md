find iput file:

```sh
$ ls -lah /dev/input/by-path/
```

find the ones that might be keyboards (usually have kbd in them)

test with `cat /dev/input/by-path/<kbd-input>`

if you see gibberish appear in the terminal, you're good

```sh
$ git clone git@github.com:Shadowbeetle/skbl.git
$ make install
$ systemctl --user daemon-reload
$ systemctl --user start skbl.service
```
