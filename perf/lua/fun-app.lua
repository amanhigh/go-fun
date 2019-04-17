--
-- Created by IntelliJ IDEA.
-- User: amanpreet.singh
-- Date: 20/07/17
-- Time: 4:32 PM
-- To change this template use File | Settings | File Templates.
--
require "lua.lib.library"
JSON = (loadfile "./lua/lib/JSON.lua")()

local header = { ["Content-Type"] = "application/json" }
local personCount = 1000

request = function()
    local api = "/person"
    local payload = file_get_contents("./request/person.json")

    --local age = math.random(1,personCount)
    --local req = JSON:decode(payload);
    --req.name="Aman-" .. age
    --req.age=age
    --payload=JSON:encode(req)

    --print(api .. payload)
    return wrk.format("POST", api, header, payload)
end