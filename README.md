# Best Match

A service to match customers with partners, by the best match.

A `match` is calculated as:

- A partner able to work will all the materials specified by the customer
- A customer within the operating radius of the partner

Results are ordered by best partner rating, and then by distance.

## Operations

Get info about a specific partner

`http://localhost:8080/partner/b276cb54ac524f8cadb1afce5ced67c4`