--
-- Created by IntelliJ IDEA.
-- User: amanpreet.singh
-- Date: 20/07/17
-- Time: 4:32 PM
-- To change this template use File | Settings | File Templates.
--
require "lua.lib.library"

local header = { ["Content-Type"] = "application/json" }
local personCount = 1000
local api = "/person"

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
    return wrk.format("POST", api, header, payload)
end