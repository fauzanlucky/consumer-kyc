# Forky API Testing Framework

## Mocha

> Mocha is a feature-rich JavaScript test framework running on Node.js and in the browser, making asynchronous testing simple and fun. Mocha tests run serially, allowing for flexible and accurate reporting, while mapping uncaught exceptions to the correct test cases

## Chai

> Chai is a assertion library for Node and the browser and can be used with any testing framework (like Mocha)

## Prerequisites

- Install Node.js and npm, we will use npm to install Mocha and Chai
- POSTMAN for making fast HTTP requests to the API

## Getting Started

```sh
$ cd qa
$ npm install
$ cp env.example .env
```

### This is directory structure within `qa`

```bash
.
├── api
│   ├── health
│   └── v1
│       └── ...
├── data
│   └── v1
│       └── ...
├── reports
│   ├── assets
│   ├── mochawesome.html
│   └── mochawesome.json
├── seeds
├── .env.example
├── .eslintignore
├── .eslintrc.js
├── package.json
└── README.md
```

| Directory | Description                                                                         |
| --------- | ----------------------------------------------------------------------------------- |
| data      | For test data, file extension should be .json                                       |
| helpers   | Common code, for general needed such as assertion, http-status-codes, token, & etc. |

### Foldering and Naming Convention

1. Filename using `snake_case`
2. Variable name using `camelCase`
3. Foldername using `kebab-case`
4. The file (`.js` and `.json`) should be located in a path that matches the endpoint path in a folder named 'data'
5. The filename of (`schema.json`) using suffix `.schema.json`

## Run the test

You can specify the command that you want to run from `package.json` file.

```sh
$ yarn test:api | to run all tests
$ yarn test:api -- --grep @tag | to run test with specific tag
$ yarn test:api -- --grep @skip --invert | to exclude @skip tag

$ yarn test:reports | to generate mochawesome report
```

## "No amount of testing can prove a software right, a single test can prove a software wrong." - Amir Ghahrai
