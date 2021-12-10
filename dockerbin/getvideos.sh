#!/bin/sh

echo "[$(date --rfc-3339=seconds)] Checking for new videos..." >> /opt/log/ytdlmc.log
flock --nonblock --verbose /opt/appdata/ytdlmc.lock \
    /opt/src/bin/ytdlmc \
        -downloader=yt-dlp \
        --config /opt/appdata/config.json >> /opt/log/ytdlmc.log 2>&1
