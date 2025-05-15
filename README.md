# aerospace-marks

**Beta**: I use it daily, but it's a WIP. Please report any issue or ideas in [issues](https://github.com/cristianoliveira/aerospace-marks/issues)

This is a cli for AeroSpace WM to manage marks on windows. 

I's heavily based on sway marks [sway](https://man.archlinux.org/man/sway.5.en) but following the `aerospace` style of commands
```text
mark --add|--replace [--toggle] <identifier>

Marks are arbitrary labels that can be used to identify certain
windows and then jump to them at a later time. Each identifier
can only be set on a single window at a time since they act as
a unique identifier. By default, mark sets identifier as the only
mark on a window. --add will instead add identifier to the list
of current marks for that window. If --toggle is specified mark will
remove identifier if it is already marked.
```
### Why not use named workspaces?

> TL;DR: Dynamic mappings to specific windows, not workspaces.

I wanted something more granular than workspaces, I want to jump to a specific window, given a context.
Sometimes my "video" context is youtube, sometimes it's a video player, sometimes it's a browser. I want to be able to jump to the right window, no matter what workspace it is on.

## Usage

Allow one to mark a window with the given identifier. 
```text
aerospace-marks mark [--add|--replace] [--toggle] <identifier>
```
And then focus to a window with the given mark.
```text
aerospace-marks focus <identifier>
```

### Config Usage

```toml
# ~/.config/aerospace/config.toml
[mode.main.binding]
# Mark
cmd-shift-a = ["exec-and-forget aerospace-marks focus a 2> /tmp/marks.log", "mode main"]
cmd-shift-b = ["exec-and-forget aerospace-marks focus browser 2> /tmp/marks.log", "mode main"]
# Focus
cmd-ctrl-a = ["exec-and-forget aerospace-marks focus a 2> /tmp/marks.log", "mode main"]
cmd-ctrl-b = ["exec-and-forget aerospace-marks focus browser 2> /tmp/marks.log", "mode main"]
```

See more in [documentation](docs/aerospace-marks)

### Packages

 - AeroSpace Socket Client - [aerospacecli](pkgs/aerospacecli)
