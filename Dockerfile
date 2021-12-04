FROM golang:1.17-buster
VOLUME [ "/opt/appdata", "/opt/log", "/opt/downloads" ]

# Install packages needed for container
RUN apt-get update && \
    apt-get install -y --no-install-recommends cron ffmpeg python3
RUN git clone https://github.com/7thFox/ytdlmc /opt/src

# Setup Cron
RUN touch /var/log/cron.log
ARG CRONTIME="*/15 * * * *"
ARG FLOCKPATH=/opt/appdata/ytdlmc.lock
RUN echo "${CRONTIME} root      echo \"[\$(date --rfc-3339=seconds)] Checking for new videos\" >> /opt/log/ytdlmc.log && flock --nonblock --verbose $FLOCKPATH /opt/src/bin/ytdlmc -downloader=yt-dlp --config /opt/appdata/config.json >> /opt/log/ytdlmc.log 2>&1" >> /etc/crontab

# we name it youtube-dl so it aliases
ADD "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp" /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp

# ensure we ALWAYS do anything after this uncached:
ADD "https://www.random.org/cgi-bin/randbyte?nbytes=10&format=h" /opt/skipcache

WORKDIR /opt/src
RUN git pull --rebase --force
RUN go build -o /opt/src/bin/ytdlmc /opt/src/main.go
# Run
CMD echo "[$(date --rfc-3339=seconds)] Container Started" >> /opt/log/ytdlmc.log && cron && tail -f /var/log/cron.log
