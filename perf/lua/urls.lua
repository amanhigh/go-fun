require "lua.library"

counter = 1
headers = {}

init = function(args)
    fileName = "./data/" .. args[1]
    if (fileName) then
        print("Loading File: " .. fileName)
        urls = non_empty_lines_from(fileName)
        urls = shuffle(urls)
        print("Testing with " .. #urls .. " unique urls")
    else
        print "No Url File Specified Exiting"
        os.exit()
    end
end

-- Lamda Function Called To Procure next Request, Builds a Batch and Provides Next Request
request = function()
    -- Rollover Counter When Input Finishes
    counter = counter + 1
    if counter >= #urls then
        counter = 1
        print("Url Set Completed Rolling Over")
    end

    --    print(urls[counter])
    return wrk.format("GET", urls[counter], headers)
end

response = function(status, headers, body)
    --     print(status)
end