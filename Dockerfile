# I'd use -buster, but youtube-dl isn't working in that version
FROM golang:1.17-bullseye
VOLUME [ "/opt/config", "/opt/log", "/opt/downloads" ]

# Install packages needed for container
RUN apt-get update && \
    apt-get install -y --no-install-recommends youtube-dl cron ffmpeg

RUN youtube-dl --version

# Setup Cron
RUN touch /var/log/cron.log
ARG CRONTIME="*/15 * * * *"
RUN echo "${CRONTIME} root      echo \"[\$(date --rfc-3339=seconds)] Checking for new videos\" >> /opt/log/ytdlmc.log && /opt/src/bin/ytdlmc --config /opt/config/config.json >> /opt/log/ytdlmc.log 2>&1" >> /etc/crontab
# RUN echo "${CRONTIME} root      echo \"[\$(date --rfc-3339=seconds)] Checking for new videos\" >> /opt/log/ytdlmc.log && /opt/src/bin/ytdlmc --simulate --config /opt/config/config.json >> /opt/log/ytdlmc.log 2>&1" >> /etc/crontab
# RUN cat /etc/crontab

# Build youtube-dl-multiconfig
RUN git clone https://github.com/7thFox/youtube-dl-multiconfig /opt/src
# Calls for a random number to break the cahing of the git pull
# (https://stackoverflow.com/questions/35134713/disable-cache-for-specific-run-commands/58801213#58801213)
ADD "https://www.random.org/cgi-bin/randbyte?nbytes=10&format=h" /opt/skipcache
WORKDIR /opt/src
RUN git pull --rebase --force
RUN go build -o /opt/src/bin/ytdlmc /opt/src/main.go
# Run
CMD echo "[$(date --rfc-3339=seconds)] Container Started" >> /opt/log/ytdlmc.log && cron && tail -f /var/log/cron.log
