### Config ###
export RESTIC_REPOSITORY=./test
export RESTIC_PASSWORD=aman

# Init Repo
# echo "\033[1;33m Init Repo (Once) \033[0m"
# restic init

# Backup Files
echo "\033[1;33m Backup Files \033[0m"
restic --exclude-file ./exclusions --tag jq,origin backup ../jq 
echo "\033[1;34m Add Another Backup \033[0m"
restic --tag plantuml,origin backup ../plantuml

# Check Snapshot
restic snapshots

#Modify Jq Folder add new.txt
echo "\033[1;34m Modifying Source \033[0m"
touch ../jq/new.txt

# Backup Post Modification
restic --tag new.txt backup ../jq

echo "\033[1;34m Snapshots After Modification \033[0m"
restic snapshots

echo "\033[1;34m Undoing Modifications \033[0m"
rm ../jq/new.txt

echo "\033[1;33m Restoring to Last Snapshot \033[0m"
read
restic restore latest --target .. #-t works for target

echo "\033[1;33m Cleaning older Snapshots \033[0m"
restic forget --keep-last 3 --prune


################# Useful Commands ############
# Repo Health: restic check (--read-data Compares Checksum)
# restic find "*.json"
# Dry Run: restic -n backup
# External Backup: restic -r /srv/restic-repo-copy copy --from-repo /srv/restic-repo
# Include List: restic backup /home/user/specific_file /home/user/specific_directory
## Browsing
# restic ls latest - File List
# restic diff 5093dca3 53486dfc (Using Snapshot ids)
# restic cat snapshot bbed3ad3 
## Schedule Backup
# /etc/systemd/system/backup.timer
# systemctl enable restic-backup.timer
## Snapshot Management
# restic forget <snapshotid> (--prune to compress)
## Check Size
# restic stats latest
## Key Management
# restic key list
# restic key passwd (change password)
## Compress
# --compression=max (auto,off)


################# Exclusions ############
# --exclude-file ./exclusions OR --exclude="*.c" --exclude-file=excludes.txt --exclude-larger-than 2048
# Exclude a single file
# echo "/home/user/specific_file" >> ./exclusions
# Exclude a directory
# echo "/home/user/specific_directory" >> ./exclusions
# Exclude files with a specific extension
# echo "*.mp4" >> ./exclusions
# Exclude files that start with a specific string
# echo "temp_*" >> ./exclusions
# Exclude files that match a regular expression
# echo "/home/user/[0-9]*" >> ./exclusions