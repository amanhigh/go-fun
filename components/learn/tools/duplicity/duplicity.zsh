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



# Duplicity Tutorial: https://www.youtube.com/watch?v=G8M3GnAkufw
#

KEY_ID=DAEBA62B1F563C062EBB33BA464D9D0E1EEBE051 #Get From List Key

echo "\033[1;32m Import Keys \033[0m \n";
gpg --import public.pgp

echo "\033[1;32m Full Backup \033[0m \n";
duplicity full ../jq file://Backup --encrypt-key=$KEY_ID

echo "\033[1;34m Modifying Source \033[0m \n";
echo "Test" > ../jq/test.txt

echo "\033[1;32m Incremental Backup \033[0m \n";
duplicity incremental ../jq file://Backup --encrypt-key=$KEY_ID

echo "\033[1;34m Removing Modifications \033[0m \n";
rm ../jq/test.txt