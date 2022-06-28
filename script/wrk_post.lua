-- example HTTP POST script which demonstrates setting the
-- HTTP method, body, and adding a header
wrk.method = "POST"
wrk.body   = '{"id":"0003","code":3,"name":"dubbogo","age":99}'
wrk.headers["Content-Type"] = "application/json"