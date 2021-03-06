# Javascript

Dependencies:

* `nodejs` (v18.4.0)
* `npm` (8.12.1)

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

## Regenerate package-lock.json

To regenerate `package-lock.json` with docker, first remove it:

```
$ rm javascript/package-lock.json
```

Now, regenerate it using docker:

```
$ docker run -it --rm --name node-package-lock -v "$PWD":/usr/src/app -w /usr/src/app/javascript node:18-alpine npm install
```
