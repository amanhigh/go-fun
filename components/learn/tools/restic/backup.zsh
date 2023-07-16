### One Time Setup for Restic ###
REPO=./test
PASS_FILE=./pass.txt
# Init Repo
# echo "\033[1;33m Init Repo (Once) \033[0m"
# restic --repo $REPO -p $PASS_FILE init

# Backup Files
echo "\033[1;34m Backup Files \033[0m"
restic --repo $REPO -p $PASS_FILE --tag jq,origin backup ../jq 

# Check Snapshot
restic --repo $REPO -p $PASS_FILE snapshots

#Modify Jq Folder add new.txt
echo "\033[1;34m Modifying Source \033[0m"
touch ../jq/new.txt

# Backup Post Modification
restic --repo $REPO -p $PASS_FILE --tag new.txt backup ../jq

echo "\033[1;33m Add Another Backup \033[0m"
restic --repo $REPO -p $PASS_FILE --tag plantuml,origin backup ../plantuml

echo "\033[1;34m Snapshots After Modification \033[0m"
restic --repo $REPO -p $PASS_FILE snapshots

echo "\033[1;34m Undoing Modifications \033[0m"
rm ../jq/new.txt

echo "\033[1;33m Restoring to Last Snapshot \033[0m"
read
restic --repo $REPO -p $PASS_FILE restore latest --target ..

echo "\033[1;33m Cleaning older Snapshots \033[0m"
restic --repo $REPO -p $PASS_FILE forget --keep-last 3 --prune