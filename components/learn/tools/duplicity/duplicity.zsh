# Generate GPG Key: gpg --gen-key
# List GPG Keys: gpg --list-key
###### Export: 
#PublicKey: gpg --output public.pgp --armor --export <email/id>
#PrivateKey: gpg --output private.pgp --armor --export-secret-key <email/id>
###### Delete: 
#PrivateKey: gpg --delete-secret-key <email/id> (Need to be done First)
#PublicKey: gpg --delete-key <email/id>
###### Importing: 
# gpg --import <key_file>
# Trust Key (Enter 5 for ultimiate) , q to exit: : gpg --edit-key <key_id> trust
####### Backup Management:
# Trigger FullBackup if older than 5 Months: duplicity --full-if-older-than 5M /source/path file:///destination/path
# Force Full: duplicity full  /source/path file:///destination/path
# Force Incremental: duplicity incremental  /source/path file:///destination/path
# TimeFormate: s, m, h, D, W, M, or Y (indicating seconds, minutes, hours, days, weeks, months, or years respectively)
####### Exclusions:
# Dry Run Testing with Verbosity 7 (max 9): --dry-run -v7
# Exclude Directory (** Represents Base): --exclude=**/Code
# Exclude Cache Only Under Code: **/Code/*Cache*
# Exclude Cache Recursively Under Code: **/Code/**Cache**
# Exclude workspaceStorage Recursively anywhere under Base: **workspaceStorage**
# Exclude Locks ignore case: ignorecase:**/**lock**
####### Restoration:
# Environment Variable PASSPHRASE can be Set to avoid any PROMPTS on DECRYPTION or Verification.
# Set Environment Variable without Logging it to History. 
# read -rs PASSPHRASE; export PASSPHRASE;

# Duplicity Tutorial: https://www.youtube.com/watch?v=G8M3GnAkufw
# Help Page: https://manpages.ubuntu.com/manpages/bionic/man1/duplicity.1.html

KEY_ID=DAEBA62B1F563C062EBB33BA464D9D0E1EEBE051 #Get From List Key (Passphrase: amanps)
BACKUP_FILE=file://Backup
RESTORE_PATH=Restore
export PASSPHRASE=amanps

echo -e "\033[1;32m Import Keys \033[0m";
gpg --import public.pgp
# gpg --import private.pgp

echo -e "\033[1;32m Full Backup \033[0m";
duplicity --full-if-older-than 10s ../jq $BACKUP_FILE --encrypt-key=$KEY_ID

echo -e "\033[1;34m Modifying Source \033[0m";
echo "Test" > ../jq/test.txt

echo -e "\033[1;32m Incremental Backup \033[0m";
duplicity ../jq $BACKUP_FILE --encrypt-key=$KEY_ID

echo -e "\033[1;32m Incremental Backup (No Change) \033[0m";
duplicity ../jq $BACKUP_FILE --encrypt-key=$KEY_ID

echo -e "\033[1;34m Removing Modifications \033[0m";
rm ../jq/test.txt

echo -e "\033[1;32m List Files \033[0m";
duplicity list-current-files --encrypt-key=$KEY_ID $BACKUP_FILE

echo -e "\033[1;32m Remove Older Backups \033[0m";
duplicity remove-older-than 10s $BACKUP_FILE --force
 
echo -e "\033[1;33m Restore Backup \033[0m";
rm -rf $RESTORE_PATH
duplicity restore --encrypt-key=$KEY_ID $BACKUP_FILE $RESTORE_PATH

echo -e "\033[1;33m Verifying Backup \033[0m";
duplicity verify $BACKUP_FILE $RESTORE_PATH --compare-data --encrypt-key=$KEY_ID
