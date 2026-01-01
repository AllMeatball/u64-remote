# U64 Remote
Remotely send and run a .PRG file on a Commodore 64 Ultimate, (or any 1541u)


## Usage and Setup
1. Setup the your device with the `creds.json` file in the config path (see the Config Path section)
2. Write this into of `creds.json` in this format (replacing the values with your own):
```json
{
    "Address": "http://192.x.x.x",
    "Password": "Password goes here if set on system"
}
```
3. Run the command with your prg file like the following:
```
u64-ultimate [prg file]
```

### Config Path
* Windows: `C:\Users\username\AppData\Roaming\u64-remote`
* Linux/Unix: `~/.config/u64-remote`
* macOS: `~/Library/Application Support/u64-remote` 

## TODO
* Figure out how to use data streams (e.g. for debug and video)

