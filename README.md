# wicked-wit
a web-based cards game


## misc
### migrations
this project uses goose

```
goose -dir ./migrations up 
```

#### seeding for debug purposes
```
goose -dir ./migrations/seed -no-versioning up
```