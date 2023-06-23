# API Consumption Problems 101 - a hands-on guide

API providers have a LOT to worry about. Between security, performance,
availability, and cost-efficiency, it's no wonder that the UX for API consumers
isn't always the best. What are issues that API consumers face? And how should
you deal with them?

In this article, we'll go over some of the most common issues that API consumers
experience using a hands-on approach - we developed an API provider that will
let you experiment with these issues first-hand.

* [� TL;DR](#-tldr)
* [🌒 An API provider that causes issues](#-an-api-provider-that-causes-issues)
  * [Try it out online](#try-it-out-online)
  * [Try it out locally](#try-it-out-locally)
* [🌓 API consumption issues](#-api-consumption-issues)
  * [1. Rate limiting](#1-rate-limiting)
    * [What is rate limiting?](#what-is-rate-limiting)
    * [Rate limiting by IP address](#rate-limiting-by-ip-address)
    * [Rate limiting by API Key](#rate-limiting-by-api-key)
    * [Exponential backoff](#exponential-backoff)
  * [2. Server-side errors](#2-server-side-errors)
    * [What are server-side errors?](#what-are-server-side-errors)
    * [Dealing with 5xx errors using retries](#dealing-with-5xx-errors-using-retries)
  * [3. Caching](#3-caching)
    * [Why do we need caching?](#why-do-we-need-caching)
    * [Caching data on the client-side](#caching-data-on-the-client-side)
* [🌔 So what's next?](#-so-whats-next)
* [🌕 How can I help?](#-how-can-i-help)

## 🌑 TL;DR

* API providers have a lot to worry about, and sometimes the UX for API consumers
  isn't the best.
* There are many patterns that API consumers must implement to deal with 3rd
  party APIs: rate limiting, server errors, and caching.
* [Lunar.dev](https://lunar.dev) is a platform that helps you manage your 3rd
  party APIs better and deals with many of the issues we discussed in this
  article for you.

## 🌒 An API provider that causes issues

![Server Cover art](/assets/cover.jpg "Server Cover art")

To help you experiment with API consumption issues, [we developed
a server](https://problematic-api-server-thecoreman.koyeb.app/swagger/) with
an API endpoint for each problematic API pattern we want to demonstrate in this
article.

### Try it out online

You can try out the server in the browser by going to the API documentation.

Open [the API
documentation](https://problematic-api-server-thecoreman.koyeb.app/swagger/).
You should see something like this:

![Swagger API - 1](/assets/api-1.png "Swagger API screenshot 1")

Click on the `Try it out` button:

![api-2-try-it-out](/assets/api-2-try-it-out.png "Swagger API screenshot 2")

And then on `Execute`:

![api-3-execute](/assets/api-3-execute.png "Swagger API screenshot 3")

You should see a `200 OK` response with text from a random book. Like this:

![api-4-response](/assets/api-4-response.png "Swagger API screenshot 4")

> This time, we got a quote from "The Iliad" by Homer:
>
> ```"Or lay the spoils of conquest at her feet,"```
>
> telling about the death of Atrides of the Greeks by the hand of Hector of the
> Trojans.

You can also try it in the terminal by running the following command:

```sh
curl -X 'GET' \
  'https://problematic-api-server-thecoreman.koyeb.app/rate-limit/by-ip' \
  -H 'accept: application/json'
```

### Try it out locally

If you prefer to run the server locally, [you can find the code
here](https://github.com/TheCoreMan/problematic-api-server). To run it, clone
the repository and run:

```sh
go run cmd/server.go
```

Then go to [http://localhost:4578/swagger/](http://localhost:4578/swagger/) to
try it out. This will allow you to check out the code to see how a provider
might implement the different patterns we'll discuss in this article.

## 🌓 API consumption issues

### 1. Rate limiting

![429 meme](/assets/429-meme.jpg "429 meme")

#### What is rate limiting?

Rate limiting is a mechanism that limits the number of requests a client can
make to an API. It's used to prevent abuse and to ensure that the API is
available to all consumers. You'll usually encounter rate limiting just when
you don't want to - as your consumption of the API increases since you enjoy
what it has to offer.

Rate limiting (usually) has two parameters:

* 📈 **Limit count**: The number of requests a client can make in a given time
  period.
* 🪟 **Time Window**: The time period in which the rate limit applies. This can
  be a sliding window (e.g. the last 5 minutes) or a fixed window (e.g. the
  current hour).

The question is "what counts as a client". Every provider can implement this
differently. Some providers will count a client as a single IP address, while
others will look at the authorization data and count an "API Key" as a client.
Let's look at both of these cases.

#### Rate limiting by IP address

This is a very common and very basic rate-limiting option for providers. To try
this out, let's call an endpoint! Run this command a few times, _fast_:

```sh
curl -X 'GET' \
  'https://problematic-api-server-thecoreman.koyeb.app/rate-limit/by-ip' \
  -H 'accept: application/json'
```

At first, the API should respond with `200 OK` and a quote from a random book.
But if you send the requests fast enough, you'll started getting rate-limited,
and the server will respond with `429 Too Many Requests`:

```sh
❯ curl -X 'GET' ...
{"book-name":"Random","text":"  And one would pillage, one would burn the place."}
❯ curl -X 'GET' ...
Rate limit exceeded. reason: Rate limit exceeded for IP 10.11.0.16 by 113.777755ms%
```

How to deal with this? There are two options: either you can wait for the rate
limit to expire (wait for the window to expire using the data in the response),
or you can try to get multiple external IP addresses and use them in a round
robin fashion to avoid hitting the rate limit.

#### Rate limiting by API Key

This is a more advanced rate-limiting option. It allows the provider to make
sure that each client is rate-limited separately. This is useful if the provider
wants to offer different rate limits to different clients. For example, offer a
free tier with a low rate limit and a paid tier with a higher rate limit.

To try this out, let's call an endpoint! Our server gets the API key from the
`X-Account-Id` HTTP header. So let's call the endpoint. The server expects an
email address, and will error if you don't pass an email address:

![API screenshot 5 - 500 error](/assets/api-5-internal-server-error.png "Swagger API screenshot 5 - 500 error")

Let's try to call the endpoint with a valid email address. Any email address
will do. Let's use `astronaut@lunar.dev`. Run this command a few times, _fast_:

```sh
curl -X 'GET' \
  'https://problematic-api-server-thecoreman.koyeb.app/rate-limit/by-account' \
  -H 'accept: application/json' \
  -H 'X-Account-Id: astronaut@lunar.dev'
```

Again, the API will respond with `200 OK` at first:

![API screenshot 6 - 200 OK](/assets/api-6-200-ok.png "Swagger API screenshot 6 - 200 OK")

But soon enough it will start rate-limiting the specific account.

![API screenshot 7 - 429 Too Many Requests](/assets/api-7-account-rate-limit.png "Swagger API screenshot 7 - 429 Too Many Requests")

How to deal with this? There are a few options when account-based rate limiting
is in play. First, you can wait for the rate limit to expire, similarly to the
IP-based rate limiting. Second, you can try to get multiple API keys and use
them in a round-robin fashion to avoid hitting the rate limit. Finally, if this
is a paid service, you can reach out to the provider and ask for a higher rate
limit (perhaps at the cost of a higher monthly fee 🤑).

#### Exponential backoff

Rate limiting by fixed windows is a way to prevent abuse, but for high-volume
services, it can be too simple an approach. A good reference to read more about
this from the perspective of a provider is [Exponential Backoff And Jitter](
https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/) from
AWS's architecture blog.

![Exponential Backoff](https://d2908q01vomqb2.cloudfront.net/fc074d501302eb2b93e2554793fcaf50b3bf7291/2017/10/03/exponential-backoff-and-jitter-blog-figure-7.png "Exponential backoff")

As a consumer, when your provider uses exponential backoff, you'll get a
`429 Too Many Requests` response, and if you don't respect the provider's
recommendations, you'll get rate-limited again and again - but every time with
a longer wait time. This can get you locked out of the API for a long time, so
make sure to respect the window which is usually passed in the response. This
can be in a "Retry-After" header, or in the response body.

Let's try it out! Run this command a few times, _fast_ (but not too much!):

```sh
❯ curl -X 'GET' \
  'https://problematic-api-server-thecoreman.koyeb.app/rate-limit/exponential-backoff' \
  -H 'accept: application/json'
{"book-name":"Random","text":"said Dron."}
❯ curl -X 'GET' ...
Rate limit exceeded. reason:Rate limit exceeded for IP 10.11.0.16 by 30s.
Violated 1 time
❯ curl -X 'GET' ...
Rate limit exceeded. reason:Rate limit exceeded for IP 10.11.0.16 by 1m.
Violated 2 times
❯ curl -X 'GET' ...
Rate limit exceeded. reason:Rate limit exceeded for IP 10.11.0.16 by 2m.
Violated 3 times
```

As you can see, the server started with a 30 second wait time, and when we
violated it more and more times, the server increased it by a factor of 2 every
time. This is a very common approach to rate limiting.

How to deal with this? The best way to deal with exponential backoff is to NEVER
hit the rate limit. If you develop your retry code intelligently, you'll be able
to work with the provider in an optimal fashion and never get locked out.

### 2. Server-side errors

![500 meme](/assets/500-meme.jpeg "500 meme")

#### What are server-side errors?

Another issue you may face when consuming 3rd party APIs is server-side errors.
Server side errors in HTTP are represented by 5xx status codes. You can check
out an exhustive list of all 5xx status codes
[here](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#server_error_responses),
but the most common one is `500 Internal Server Error`. This is a generic error
that the server returns when it encounters an error that it doesn't know how to
handle. This can be caused by a bug in the server's code, or by a bug in the
server's infrastructure (e.g. a database that is down).

To emulate this, we've added an endpoint to our server that returns a 500 error
randomly. If we pass the "error_percentage" query parameter, we can control the
probability of the server returning a 500 error. If we pass 50%, the server will
return a 500 error 50% of the time:

![api-8-error-1](/assets/api-8-error-1.png "api-8-error-1")

#### Dealing with 5xx errors using retries

The thing with 5xx errors is that they are not your fault. The best way to deal
with them is to retry the request a bunch of times and hope that the server's
issues are temporary and will be resolved.

```sh
❯ curl --silent --include -X 'GET' \
  'https://problematic-api-server-thecoreman.koyeb.app/errors/percent?error_percent=50' \
  -H 'accept: */*'
HTTP/2 500

"Rolled a 0.48 which is less than 0.50"
❯ curl --silent --include -X 'GET' \
  'https://problematic-api-server-thecoreman.koyeb.app/errors/percent?error_percent=50' \
  -H 'accept: */*'
HTTP/2 200

"Rolled a 0.90 which is greater than or equal to 0.50"
```

As you can see, the server returned a `500` error the first time, and a `200 OK`
the second time. The retry logic worked!

### 3. Caching

![Caching meme](/assets/caching-meme.jpeg "Caching meme")

#### Why do we need caching?

Caching is a mechanism that allows us to store data in a fast storage layer.
Consuming an API takes time and resources. For data that's relatively static,
it's a waste of resources to consume the API every time we need the data. This
is where caching comes in. We can store the data in a local storage layer, and
consume it from there instead of consuming the API every time.

Cacheable data usually has a TTL (Time To Live) - a time period after which the
data is considered stale and should be refreshed. This is usually passed using
the `Cache-Control` HTTP header.

#### Caching data on the client-side

Let's try it out! Our server has an endpoint that returns a the quotes based
on parameters. Since the books don't change, we can cache the quotes on the
client-side. Let's try it out!

![API 9 - with control](/assets/api-9-with-control.png "API 9 - with control")

```sh
curl -X 'GET' \
  'https://problematic-api-server-thecoreman.koyeb.app/cacheable?book-title=iliad.txt&line-number=1234&with-control=true' \
  -H 'accept: application/json'
```

The response will include the data and a `Cache-Control` header with a TTL:

![API 10 - with control response](/assets/api-10-with-control-response.png "API 10 - with control response")

Response body:

```json
{
  "book-name": "iliad.txt",
  "line-number": 1234,
  "text": "  Spleen to mankind his envious heart possess'd,"
}
```

Response headers:

```text
cache-control: max-age=3600
```

How to deal with this? This means that we can cache the data in a local database
such as Redis or in memory for up to 3600 seconds (1 hour). If we need the data,
we can get it from the cache instead of consuming the API. If we need to refresh
the data, we can consume the API again and rehydrate the cache.

## 🌔 So what's next?

As you can see, there are many issues that API consumers can face. If you want
to manage your 3rd party API consumption with no sweat, consider giving
[Lunar.dev](https://lunar.dev) a try! It's a platform that helps you manage your
3rd party APIs better and deals with many of the issues we discussed in this
article for you.

## 🌕 How can I help?

There are many other issues that API consumers face, such as dealing with
Pagination, security considerations, and more. If you want to contribute to this
educational project, check out the code on [GitHub](
https://github.com/TheCoreMan/problematic-api-server) and submit a Pull Request.
