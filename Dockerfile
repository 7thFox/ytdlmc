FROM golang:1.17-buster
VOLUME [ "/opt/appdata", "/opt/log", "/opt/downloads" ]

# Install packages needed for container
RUN apt-get update && \
    apt-get install -y --no-install-recommends cron ffmpeg python3
RUN git clone https://github.com/7thFox/ytdlmc /opt/src

RUN touch /var/log/cron.log

# Pull yt-dlp binary release
ADD "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp" /bin/yt-dlp
RUN chmod a+rx /bin/yt-dlp

# TODO: Verify GPG Signature pending https://github.com/yt-dlp/yt-dlp/issues/1886

ADD ./dockerbin/startup.sh /bin/startup.sh
ADD ./dockerbin/getvideos.sh /bin/getvideos.sh
RUN chmod a+rx /bin/startup.sh
RUN chmod a+rx /bin/getvideos.sh

# ensure we ALWAYS do anything after this uncached:
ADD "https://www.random.org/cgi-bin/randbyte?nbytes=10&format=h" /opt/skipcache

# Build ytdlmc
WORKDIR /opt/src
RUN git pull --rebase --force
RUN go build -o /opt/src/bin/ytdlmc /opt/src/main.go


ENV CRONTIME="*/30 * * * *"
CMD startup.sh >> /opt/log/ytdlmc.log 2>&1 && tail -f /opt/log/ytdlmc.log
