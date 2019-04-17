--
-- Created by IntelliJ IDEA.
-- User: amanpreet.singh
-- Date: 27/11/14
-- Time: 5:17 PM
-- To change this template use File | Settings | File Templates.
--
JSON = (loadfile "./lua/lib/JSON.lua")()

function shuffle(array)
    local n = #array
    for i = 1, n do
        j, k = math.random(n), math.random(n)
        array[j], array[k] = array[k], array[j]
    end
    return array
end

function file_exists(file)
    local f = io.open(file, "rb")
    if f then
        f:close()
    end
    return f ~= nil
end

function non_empty_lines_from(file)
    lines = {}

    if not file_exists(file) then
        print("Specified file " .. file .. " not found")
        os.exit()
        return lines
    end

    -- Read All Content
    content = file_get_contents(file)
    -- Break Into Lines
    lines = mysplit(content, "\n")
    -- Filter Empty Lines
    non_empty_lines = filter_non_empty(lines)
    -- print(table.concat(non_empty_lines, " "))
    return non_empty_lines
end

-- String Operations
function filter_non_empty(lines)
    non_empty_lines = {}
    for i, s in pairs(lines) do
        if (not is_empty(s)) then
            non_empty_lines[#non_empty_lines + 1] = s
        end
    end
    return non_empty_lines
end

function is_empty(s)
    return s == nil or s == '-'
end

function mysplit(inputstr, sep)
    if sep == nil then
        sep = "%s"
    end
    local t = {};
    i = 1
    for str in string.gmatch(inputstr, "([^" .. sep .. "]+)") do
        t[i] = str
        i = i + 1
    end
    return t
end

position = 0

function batch(array, size)
    if position > #array then
        position = 1
        print("Rolling Over Batch Set")
    end
    --print('Array Length: ' .. #array .. ' Range: ' .. position .. ' to ' .. position+size)

    local result = {}
    for i = position, position + size do
        result[#result + 1] = array[i]
    end
    position = position + size

    --print(JSON:encode(result))
    return result
end

function batchString(array, position, size)
    --print('Array Length: ' .. #array .. ' Range: ' .. position .. ' to ' .. position + size)
    --print('Values: ' .. array[position] .. '-' .. array[position + size])
    local result = "#"
    for i = position, position + size do
        local p = array[i]
        --print('Index: ' .. i .. ' Result: ' .. result .. ' Next Val: ' .. p)
        result = result .. ",\"" .. p .. "\"";
        --print("Concat Result: " .. result)
    end
    return string.gsub(result, "#,", "")
end

function filterByPattern(tbl, pattern)
    local newtbl = {}
    j = 0
    for i, v in pairs(tbl) do
        for n, p in pairs(pattern) do
            if string.match(v, p) then
                newtbl[j] = v
                j = j + 1
                break
            end
        end
    end
    return newtbl
end

-- Helper function which reads the contents of a file(This function is from the helloworld.lua example above)
function file_get_contents(filename)
    local file = io.open(filename, "r")
    if not file then
        return nil
    end

    local contents = file:read("*all") -- See Lua manual for more information
    file:close() -- GC takes care of this if you would've forgotten it

    return contents
end

-- Array related functions
function addToSet(set, key)
    set[key] = true
end

function removeFromSet(set, key)
    set[key] = nil
end

function setContains(set, key)
    return set[key] ~= nil
end

function arrayToSet(arr)
    local set = {}
    for k, v in pairs(arr) do
        addToSet(set, v)
    end
    return set
end
