# mtbcal

[![Build Status](https://travis-ci.org/slashk/mtbcal.svg?branch=master)](https://travis-ci.org/slashk/mtbcal)

## Documentation

To view generated docs for mtbcal, run the below command and point your brower to http://127.0.0.1:6060/pkg/

    godoc -http=:6060 2>/dev/null &

### Buffalo

http://gobuffalo.io/docs/getting-started

### Pop/Soda

http://gobuffalo.io/docs/db

## Database Configuration

 	development:
 		dialect: postgres
 		database: mtbcal_development
 		user: <username>
 		password: <password>
 		host: 127.0.0.1
 		pool: 5

 	test:
 		dialect: postgres
 		database: mtbcal_test
 		user: <username>
 		password: <password>
 		host: 127.0.0.1

 	production:
 		dialect: postgres
 		database: mtbcal_production
 		user: <username>
 		password: <password>
 		host: 127.0.0.1
 		pool: 25

 ### Running Migrations

    buffalo soda migrate

 ## Run Tests

    buffalo test

 ## Run in dev

    buffalo dev

[Powered by Buffalo](http://gobuffalo.io)

