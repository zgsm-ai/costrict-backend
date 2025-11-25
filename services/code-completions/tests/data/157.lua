local function isPalindrome(str)
    local reversed = str:reverse()
    <｜fim▁hole｜>return str == reversed
end

local testString = "level"
if isPalindrome(testString) then
    print(testString .. " is a palindrome")
else
    print(testString .. " is not a palindrome")
end