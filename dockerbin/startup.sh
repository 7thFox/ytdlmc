#!/bin/sh

echo "[$(date --rfc-3339=seconds)] Container Starting..." >> /opt/log/ytdlmc.log

echo "[$(date --rfc-3339=seconds)] Writing crontab..." >> /opt/log/ytdlmc.log
echo "\
SHELL=/bin/sh\n\
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin\n\
$CRONTIME root getvideos.sh" > /etc/crontab

echo "[$(date --rfc-3339=seconds)] Starting Cron..." >> /opt/log/ytdlmc.log
cron

echo "[$(date --rfc-3339=seconds)] Container Started" >> /opt/log/ytdlmc.log