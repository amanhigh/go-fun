# Tail a log and send filtered log to another file.
# Can be installed as service. Refer test-script.service under source-deb
tail -f -n0 /var/log/system.log | grep --line-buffered -E 'error' >> /tmp/syslog_error.log

