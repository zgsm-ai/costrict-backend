local function processTable(t)
    local result = {}
    for key, value in pairs(t) do
        <｜fim▁hole｜>if type(value) == "number" then
            result[key] = value * 2
        else
            result[key] = value
        end
    end
    return result
end

local data = {
    name = "Lua",
    version = 5.4,
    type = "scripting"
}

local processed = processTable(data)
for key, value in pairs(processed) do
    print(key, value)
end