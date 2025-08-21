# wicked-wit
a web-based cards game


## misc
### migrations
this project uses goose for db migrations

```
goose -dir ./server/migrations up 
```

#### seeding for debug purposes
```
goose -dir ./server/migrations/seed -no-versioning up
```