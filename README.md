# youtube-dl-multiconfig
youtube-dl has some really great and powerful auto config options to automate downloading YouTube videos.
However, there is not apparent support to configure youtube-dl on a per-playlist basis, so this project aims
to fill that gap.

This project was partially inspired by the nice guide by Erik Ellsinger on [auto-downloading videos](https://erik.ellsinger.me/automatically-download-youtube-videos-to-plex-on-truenas-using-youtube-dl/) on TrueNAS

# config setup
Configuration is done within named "config groups". All parameters are entered as they would be in the command
line, and you can view youtube-dl documentation for info on each of them. The only option straying from this is
`batch-file` which will instead take an array of strings. These will be written to a temp file at runtime before
the file name is passed into `youtube-dl`

View youtube-dl options: [https://github.com/ytdl-org/youtube-dl#description](https://github.com/ytdl-org/youtube-dl#description)

By default the program looks at `~/.config/youtube-dl-multiconfig/config` but this may be changed via the `--config` option.

# example
See `example.json` to view a sample config.

# todo
 - Add functionality to remove existing videos older than N days
 - Potentially allow constants to be shared across config groups
