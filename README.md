# ytdlmc (youtube-dl-multiconfig)
youtube-dl has some really great and powerful auto config options to automate downloading YouTube videos.
However, there is not apparent support to configure youtube-dl on a per-playlist basis, so this project aims
to fill that gap.

This project was partially inspired by the nice guide by Erik Ellsinger on [auto-downloading videos](https://erik.ellsinger.me/automatically-download-youtube-videos-to-plex-on-truenas-using-youtube-dl/) on TrueNAS

# options
```
  -config string
        The path to your config file (default "~/.config/ytdlmc/config")
  -downloader string
        The downloader to use. Technically only build to work for 'youtube-dl', but forks like 'yt-dlp' work for most things (default "youtube-dl")
  -simulate
        Print command and don't execute
```

# config
NOTE: Go is very pedantic about it's JSON, and will not except trailing commas or comments!

Configuration is done within named "config groups". All parameters are entered as they would be in the command
line, and you can view [youtube-dl documentation](https://github.com/ytdl-org/youtube-dl#options) for info 
on each of them. The only option straying from this is `batch-file`: 
```
-a, --batch-file FILE                File containing URLs to download ('-'
                                     for stdin), one URL per line. Lines
                                     starting with '#', ';' or ']' are
                                     considered as comments and ignored.
```
Which will instead take an array of strings.

## additional json elements
In addition to the youtube-dl options, some extra elements are available for the json config:
 - `disable` - Skips the group, allowing you to effectively comment it out
 - `comment` - Field is ignored, but allows you to add a comment describing the config group

# example
See `example.json` to view a sample config.

# todo
 - Add support for [yt-dlp](https://github.com/yt-dlp/yt-dlp) and change default downloader
 - Allow sub-config groups to inherit settings
 - Add functionality to remove existing videos older than N days
