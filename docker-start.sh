#!/bin/bash

/go/bin/buffalo db migrate up
# /go/bin/grift fixtures
# /go/bin/grift counts
/go/src/github.com/slashk/mtbcal/bin/app
