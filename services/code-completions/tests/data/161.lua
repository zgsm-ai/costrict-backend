local function fibonacci(n)
    if n <= 1 then
        return n
    else
        return fibonacci(n - 1) + fibonacci(n - 2)
    end
end

local function main()
    local n = 10
    local result = fibonacci(n)
    print("Fibonacci(" .. n .. ") = " .. result)<｜fim▁hole｜>
end

main()