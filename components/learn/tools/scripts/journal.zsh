set -e
if [ -z $1 ]; then
  echo "Usage: $0 <capture_path> <trade_asset_path>"
  exit 1
fi

#Source and Target Directories
captures=$1
brain=$2
assets=$brain/assets/trading
journal=$brain/journals
yesterday=$journal/`date -d yesterday +%Y_%m_%d.md`

# Check if the directory exists
if [ ! -d $captures ]; then
  echo "The specified captures directory does not exist."
  exit 1
fi

# Find files with Name trend
files=$(find "$captures" -maxdepth 1 -type f -name '*trend*')

# Check the count of files
count=$(echo "$files" | wc -l)

if [ $count -gt 2 ]; then
    echo "\033[1;32m Processing SNF Journal \033[0m \n";
    echo "\n- SNF Journal #trading-tome" >> $yesterday

    TICKER=""
    PREVIOUS_TICKER=""

    # Iterate through the files and rename them
    for file in $files; do
        # Extract filename from the filepath
        filename=$(basename "$file")

        # Extract the parts from the filename using awk
        TAGS=$(echo "$filename" | awk -F '--' '{print $1}')
        CHRONO=$(echo "$filename" | awk -F '--' '{print $2}')

        TICKER=$(echo "$TAGS" | awk -F '[.]' '{print $1}')
        TIMEFRAME=$(echo "$TAGS" | awk -F '[.]' '{print $2}')
        TREND=$(echo "$TAGS" | awk -F '[.]' '{print $3}')
        TYPE=$(echo "$TAGS" | awk -F '[.]' '{print $4}')
        REASON=$(echo "$TAGS" | awk -F '[.]' '{print $5}')

        YEAR=$(echo "$CHRONO" | awk -F '[-]' '{print $1}')
        MONTH=$(echo "$CHRONO" | awk -F '[-]' '{print $2}')
        DAY=$(echo "$CHRONO" | awk -F '[-]' '{print $3}')

        #Organize Year and Month Wise
        asset_path=$assets/$YEAR/$MONTH
        mkdir -p $asset_path
        asset=$asset_path/$filename
        
        # Check if TICKER has changed from the previous iteration
        if [ "$TICKER" != "$PREVIOUS_TICKER" ]; then
            echo "\t- | \`$TICKER\` | #t.$TIMEFRAME | #t.$TREND | #t.$TYPE | #t.fail |$( [ -n "$REASON" ] && echo " #r.$REASON" || echo "" )|" >> $yesterday
            PREVIOUS_TICKER="$TICKER"
        fi
        
        echo "\033[1;33m Processing: $asset \033[0m \n";

        #Make Journal Entry and Move File
        echo "\t\t- ![$filename](../assets/trading/$YEAR/$MONTH/$filename)" >> $yesterday
        mv "$file" "$asset"
    done
fi