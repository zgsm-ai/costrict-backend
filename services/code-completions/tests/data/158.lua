<｜fim▁hole｜>local function calculateArea(shape, ...)
    local args = {...}
    
    if shape == "rectangle" then
        local width, height = args[1], args[2]
        return width * height
    elseif shape == "circle" then
        local radius = args[1]
        return math.pi * radius * radius
    else
        return 0
    end
end

local rectArea = calculateArea("rectangle", 5, 3)
local circleArea = calculateArea("circle", 2)

print("Rectangle area:", rectArea)
print("Circle area:", circleArea)