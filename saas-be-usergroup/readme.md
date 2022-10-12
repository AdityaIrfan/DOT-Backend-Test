## How To Migrate Databae
1. Create New Migration  
`sql-migrate new -config=migration.yaml -env="local" engagement`
2. Migrate database  
`sql-migrate up -config=migration.yaml -env="local"`

## models entities  generation
`https://github.com/volatiletech/sqlboiler#pro-tips`

## mock ports 
`https://github.com/vektra/mockery`
go dir ports and try mocks

## test using
`https://github.com/stretchr/testify`