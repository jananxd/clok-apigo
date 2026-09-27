Sept 27, 2026
- [x] Setup a Golang project
- [ ] Create an endpoint that will return a URL, hold the endpoint, create a redis sub, response to the HTTP request once we receive a pub

Specs of the endpoint above:
- the endpoint should hold the request until the timeout or we receive an authentication pub in the redis sub
- it should also create a secret key so that we can match the session from the CLI and the one who authenticates using web.
- check if we really need to use redis pub sub since it might not be needed since we are sure that we will only have a single instance of the web server and it might be too much to do it on redis.

Sept 28, 2026
- [ ] Vibecode a Frontend so that the user can login using a webpage
- [ ] Setup authentication, whether it's on the Frontend and the Frontend will hit an endpoint to verify the token + publish to the long running HTTP request initially.
