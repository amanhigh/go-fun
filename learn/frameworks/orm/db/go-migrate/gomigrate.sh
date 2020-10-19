CMD="/usr/local/bin/migrate"
URL="mysql://aman:aman@tcp(docker:3306)/compute?charset=utf8&parseTime=True&loc=Local"
PATH="./migration"

#Generate Migrations
#migrate create -dir ./migration -ext sql -seq init_schema

# Clear Database
echo -en "\033[1;32m Drop \033[0m \n"
$CMD -path $PATH -database $URL drop

# Apply Migrations
echo -en "\033[1;32m Apply \033[0m \n"
$CMD -path $PATH -database $URL up

# Check Version
echo -en "\033[1;32m Version Check \033[0m \n"
$CMD -path $PATH -database $URL version

# Remove Migration
echo -en "\033[1;32m Down \033[0m \n"
$CMD -path $PATH -database $URL down

#Version Check
echo -en "\033[1;32m Version Check (Post Down) \033[0m \n"
$CMD -path $PATH -database $URL version
