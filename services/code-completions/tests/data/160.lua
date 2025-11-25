local function readFile(filename)
    local file = io.open(filename, "r")
    if not file then
        return nil, "Could not open file"
    end
    
    local content = file:read("*all")
    file:close()
    return content
end

local function main()
    local filename = "example.txt"
    <｜fim▁hole｜>local content, err = readFile(filename)
    if content then
        print("File content:")
        print(content)
    else
        print("Error:", err)
    end
end

main()