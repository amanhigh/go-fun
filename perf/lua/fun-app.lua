--
-- Created by IntelliJ IDEA.
-- User: amanpreet.singh
-- Date: 20/07/17
-- Time: 4:32 PM
-- To change this template use File | Settings | File Templates.
--
require "lua.lib.library"

local personCount = 1000
local api = "/person"
local logfile = io.open("wrk.log", "w");

init = function()
    payload = file_get_contents("./request/person.json")
    req = JSON:decode(payload);
end

request = function()
    --local age = math.random(1,personCount)
    --req.name="Aman-" .. age
    --req.age=age
    --payload=JSON:encode(req)

    --print(api .. payload)
    return wrk.format("POST", api, std_header, payload)
end

response = function(status, headers, body)
    --print(status)
    --if status ~= 200 then
        resp = JSON:decode(body)
        logfile:write(status .. " " .. resp.Name .. "\n")
    --end
end
