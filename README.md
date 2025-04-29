# GoExchange
A fun implementation of an asset exchange using Go. I am mostly doing this to learn more about different facets of Go. Most of this code is not AI-generated for obvious reasons.

## To-do
- [x] Basic order queuing and matching logic (Max-min heap).
- [x] Handle limit orders
- [x] Handle partial fills.
- [x] Handle market orders.
- [x] Handle multiple assets.
- [x] Create an API layer to handle orders from customers (net/http).
- [ ] Simulate large order flow and experiment with concurrency (Channels + Goroutines).
- [ ] Expose orderbook using websockets.
- [ ] Simple UI to showcase project.
- [ ] Maintain 80%+ code coverage.
