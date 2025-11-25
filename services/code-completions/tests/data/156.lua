local function reverseString(str)
    local reversed = ""
    for i = #str, 1, -1 do
        <｜fim▁hole｜>reversed = reversed .. string.sub(str, i, i)
    end
    return reversed
end

local original = "Hello World"
local reversed = reverseString(original)

print("Original:", original)
print("Reversed:", reversed)