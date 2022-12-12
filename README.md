This project aims to demonstrate following basic concepts:
* RESTful API (CRUD with proper HTTP status codes)
* Data mapping to objects
* Validation (limited to request for simplicity)
* Objects with enforced valid state via constructor
* Integration with external API ([FakeStoreAPI](https://fakestoreapi.com))
* Testing
* Usage of [Echo framework](https://echo.labstack.com)

Notes:
* Usually there's middle layer (business/domain, for example ProductService), that acts as a glue between input
  (ie: REST API) and output (ie: external API, database) layers. However, for this demo it wouldn't bring any value,
  so was skipped for simplicity.
* `AddProductCommand` or `UpdateProductCommand` (not to be mistaken with CLI commands) are meant to be immutable objects
  that are always in valid state, thus the extra logic in constructor. "Command" suffix means that it will be used to
  change state of the data, while "Query" would mean reading only. Given that the input and output form of the data is
  basically the same, it feels like overkill in this project - so is there for demo purposes.
* Similar concept to the above is called `ValueObject` for single value, like "Email", "URL" or something with domain
  meaning.
