# Javascript

Dependencies on Ubuntu 20.04:

* `nodejs` (v10.19.0)
* `npm` (6.14.4)

Install node dependencies:

```
$ npm ci
```

Running tests:

```
$ npm run test
```

Uplifting minor dependencies:

```
$ npm update
```

## Docker

To build with docker:

```
$ docker build -t jox-algorithm-javascript -f Dockerfile.javascript .
```

Run tests:

```
$ docker run jox-algorithm-javascript
```
