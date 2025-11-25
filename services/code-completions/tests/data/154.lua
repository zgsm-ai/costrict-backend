function fibonacci(n)
    if n <= 1 then
        return n
    else
        <｜fim▁hole｜>return fibonacci(n - 1) + fibonacci(n - 2)
    end
end

for i = 1, 10 do
    print(string.format("F(%d) = %d", i, fibonacci(i)))
end