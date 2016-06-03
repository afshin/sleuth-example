# `sleuth-example`

This is an example of two web services: an `article-service` and a `comment-service`. The article service is dependent on the comment service.

The `http` branch solves this dependency by making HTTP requests, assuming that the comment service is on the same machine and its port never changes.

The `master` branch solves this dependency by creating a [`sleuth`](https://github.com/ursiform/sleuth) network between the two services.

This codebase is the basis for a tutorial on how to use `sleuth`: "[Service autodiscovery in Go with sleuth](http://darian.af/post/master-less-peer-to-peer-micro-service-autodiscovery-in-golang-with-sleuth/#naive-implementation-using-http-requests:7997f7408f245e3f1e7de9f602b588e5)".
