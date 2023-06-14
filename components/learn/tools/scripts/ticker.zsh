# Check if the first argument is provided
if [ -z $1 ]; then
  echo "Usage: $0 <directory_path>"
  echo "Example: $0 /path/to/your/directory"
  exit 1
fi

# Get the directory path from the first argument
directory_path=$1

# Check if the directory exists
if [ ! -d $directory_path ]; then
  echo "The specified directory does not exist."
  exit 1
fi

# Store clipboard content in a variable
clipboard_content=$(xclip -selection clipboard -o)

#TODO Check Clipboard Regex "POLYCAB.WDH.CTREND"

# Find files created in the last hour in the specified directory with Name SNAG
files=$(find "$directory_path" -maxdepth 1 -type f -cmin -25 -name '*SNAG*')

# Iterate through the files and rename them
for file in $files; do
  # Replace 'SNAG' with clipboard content (Ticker and Timeframe) in the file name
  new_name=$(echo $file | sed "s/SNAG/$clipboard_content/g")
  echo "\033[1;32m Processing: $new_name \033[0m \n";
  mv "$file" "$new_name"
done

# Mark Processed to avoid reuse
echo "Processed: $clipboard_content" | xclip -sel clip