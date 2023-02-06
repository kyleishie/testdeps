# testdeps - AKA Testing Dependencies

### Motivation

Mocking databases sucks. testdeps aims to make testing your database code easier by allowing you to spin up a local
container to test against. No more wondering if your tests are passing simply because you suck at writing mocks. 🙃

### TODO:

- [X] Finish pkg/nats
- [ ] Add examples
- [X] Implement PostgreSQL
- [ ] Implement Redis
- [ ] Implement Kafka
- [ ] Implement DynamoDB

### Contribution

Feel free to submit a pull request. Please try to follow the coding style as demonstrated in pkg/mongo. If you add a pkg
it must have the following.

- A "New..." constructor func.
- Container must auto start.
- Container struct must provide access to connection string if applicable.
- pkg must provide convenience functions for using the dep. For example, pkg/mongo has NewClient, NewDatabase.
- pkg must provide "...WithContext" versions of methods and respect context state. ctx must be first parameter.
- pkg must provide "NewTest..." methods that accept a *testing.T and use it to perform post test clean up. See pkg/mongo
  for demonstration.
- tests... If you omit tests or any of the tests are failing your PR will be rejected.
