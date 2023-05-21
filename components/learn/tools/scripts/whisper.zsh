set -e
if [ $# -ne 2 ]; then
    echo "Usage: \$0 whisper_dir journal_dir"
    exit 1
fi

whisper_dir=$1
journal_dir=$2

for file in "$whisper_dir"/*.m4a; do
    filename=$(basename "$file")
    # Strip the .m4a extension
    filename_without_ext="${filename%.*}"

    # Extract date and time
    date=$(echo "$filename_without_ext" | cut -d "_" -f 1-3)
    time=$(echo "$filename_without_ext" | cut -d "_" -f 4-6)
    
    #Journal Detection
    journal_name="$date.md"
    journal_path="$journal_dir/$journal_name";
    whisper_output="$whisper_dir/$filename_without_ext.txt"

    #Convert Note
    echo "\033[1;32m JournalName: $journal_name Time: $time  \033[0m \n";
    whisper "$file" --output_format txt -o $whisper_dir --language en >/dev/null 2>/dev/null;
    
    #Append Note
    echo "\n- VoiceNote: $time" >> $journal_path;
    cat $whisper_output >> $journal_path
done