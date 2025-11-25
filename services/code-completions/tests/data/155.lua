local function bubbleSort(arr)
    local n = #arr
    for i = 1, n do
        for j = 1, n - i do
            if arr[j] > arr[j + 1] then
                <｜fim▁hole｜>arr[j], arr[j + 1] = arr[j + 1], arr[j]
            end
        end
    end
    return arr
end

local numbers = {5, 2, 9, 1, 7}
local sorted = bubbleSort(numbers)

for i, v in ipairs(sorted) do
    print(i, v)
end