# Generate GPG Key: gpg --gen-key
# List GPG Keys: gpg --list-key
###### Export: 
#PublicKey: gpg --output public.pgp --armor --export <email>
#PrivateKey: gpg --output private.pgp --armor --export-secret-key <email>
###### Delete: 
#PrivateKey: gpg --delete-secret-key <email> (Need to be done First)
#PublicKey: gpg --delete-key <email>
###### Importing: 
# gpg --import <key_file>
# Trust Key (Enter 5 for ultimiate) , q to exit: : gpg --edit-key <key_id> trust
####### Backup Management:
# Trigger FullBackup if older than 5 Months: duplicity --full-if-older-than 5M /source/path file:///destination/path
# Force Full: duplicity full  /source/path file:///destination/path
# Force Incremental: duplicity incremental  /source/path file:///destination/path
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

echo "\033[1;32m Import Keys \033[0m \n";
gpg --import public.pgp
gpg --import private.pgp

echo "\033[1;32m Full Backup \033[0m \n";
duplicity --full-if-older-than 10s ../jq $BACKUP_FILE --encrypt-key=$KEY_ID

echo "\033[1;34m Modifying Source \033[0m \n";
echo "Test" > ../jq/test.txt


echo "\033[1;32m Incremental Backup \033[0m \n";
duplicity ../jq $BACKUP_FILE --encrypt-key=$KEY_ID

echo "\033[1;32m Incremental Backup (No Change) \033[0m \n";
duplicity ../jq $BACKUP_FILE --encrypt-key=$KEY_ID

echo "\033[1;34m Removing Modifications \033[0m \n";
rm ../jq/test.txt

echo "\033[1;32m List Files \033[0m \n";
duplicity list-current-files --encrypt-key=$KEY_ID $BACKUP_FILE

echo "\033[1;32m Remove Older Backups \033[0m \n";
duplicity remove-older-than 10s $BACKUP_FILE --force
 
echo "\033[1;33m Restore Backup \033[0m \n";
rm -rf $RESTORE_PATH
duplicity restore --encrypt-key=$KEY_ID $BACKUP_FILE $RESTORE_PATH

echo "\033[1;33m Verifying Backup \033[0m \n";
duplicity verify $BACKUP_FILE $RESTORE_PATH --compare-data --encrypt-key=$KEY_ID
