# mtbcal

[![Build Status](https://travis-ci.org/slashk/mtbcal.svg?branch=master)](https://travis-ci.org/slashk/mtbcal)
[![Code Climate](https://codeclimate.com/github/slashk/mtbcal/badges/gpa.svg)](https://codeclimate.com/github/slashk/mtbcal)

## Documentation

To view generated docs for mtbcal, run the below command and point your browser to http://127.0.0.1:6060/pkg/

    godoc -http=:6060 2>/dev/null &

## Install

```bash
$ buffalo db migrate up
$ buffalo build -z
$ ./bin/mtbcal
```

## Run Tests

   buffalo test

## Run in dev

    buffalo dev

[Powered by Buffalo](http://gobuffalo.io)

