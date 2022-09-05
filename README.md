# Best Match

A service to match customers with partners, by the best match.

A `match` is calculated as:

- A partner able to work will all the materials specified by the customer
- A customer within the operating radius of the partner

Results are ordered by best partner rating, and then by distance.

## Operations

### Get the best matches

The best-matching provider depends on the request body. To do such a call, use the endpoint `http://localhost:8080/partners` with a request body like so:

```json
{
  "materials": [
    "07cab731-d981-4915-9444-cc997eec351f",
    "1606f175-3502-4028-9501-6b591c00f1f3",
    "ac47d822-ffc9-48b7-8492-4d49e921d4df"
  ],
  "area": 100.5,
  "phone": "0123456789",
  "address": {
    "lat": "52.3599795",
    "long": "4.8851198"
  }
}
```

The response is a list of partners, order by best match (rating and then distance to the customer):

```json
200 OK

[
  {
    "ranking": 1,
    "distance_km": 1.5015628,
    "Id": "b276cb54-ac52-4f8c-adb1-afce5ced67c4",
    "Name": "Acme Inc.",
    "Address": "(4.8986299,52.3706706)",
    "Radius": 10,
    "Rating": 5
  },
  {
    "ranking": 2,
    "distance_km": 1.4156744,
    "Id": "6360d1e7-ccd0-43d0-8bf5-d7bc807213d3",
    "Name": "De Twee Broers",
    "Address": "(4.8964412,52.3706706)",
    "Radius": 10,
    "Rating": 4.3
  }
]
```


### Get info about a specific partner

To get info about a specific partner, call `http://localhost:8080/partner/` passing a partner ID in uuid form (in any valid form):  

`http://localhost:8080/partner/b276cb54ac524f8cadb1afce5ced67c4`

```json
200 OK

{
  "Id": "b276cb54-ac52-4f8c-adb1-afce5ced67c4",
  "Name": "Acme Inc.",
  "Address": "(4.8986299,52.3706706)",
  "Radius": 10,
  "Rating": 5
}
```

## Design decisions

If this was a production system, some things would look different:

- Instead of seeding the postgres database via a `initdb` file, sql migrations would be used;
- For integration testing, a separate compose file would be used to setup a database in place only for the tests;
- Graceful shutdown of the HTTP service would be in place, so lingering requests would have time to finish before the service shut down.