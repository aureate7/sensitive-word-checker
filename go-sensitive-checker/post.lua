wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
local base = "神农中医药信息平台提供智能识别与知识整合功能。"
local repeat_times = math.floor(2000 / #base)
local text = ""

for i = 1, repeat_times do
  text = text .. base .. tostring(i)
end

wrk.body = string.format('{"text":"%s"}', text)