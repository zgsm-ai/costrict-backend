local Person = {}
Person.__index = Person

function Person.new(name, age)
    local self = setmetatable({}, Person)
    self.name = name
    self.age = age
    return self
end

function Person:introduce()
    <｜fim▁hole｜>return "Hi, I'm " .. self.name .. " and I'm " .. self.age .. " years old."
end

local function main()
    local person = Person.new("Alice", 30)
    print(person:introduce())
end

main()