# star-wars-api
A simple CRUD api for Star Wars planets. See [the docs](https://gugabfigueiredo.github.io/star-wars-api/)

---

### Local

To run the API locally:

```bash
$ make compose-up
$ make run
```

Test it with
```bash
$ curl localhost:8080/sw-api/health?user=jedimaster
```

Fill the database with planets from [swapi](https://swapi.dev/)
```bash
$ curl localhost:8080/sw-api/planets/update-movies
```

Clean everything when you are done
```bash
$ make compose-down
```