# TODOs

## Lists

- [x] make the endpoints.go file for the catalog endpoints and use them in test
- [x] add in the catalog part the middleware for authentication and modify the test to have that configured as well.
- [x] In the DeleteUserByAdmin handler (and maybe others), the logger is missing. Check if all handlers have logging setup
- [x] make better helper functions for the unit and integration tests of catalog service
- [x] make sure that the catalog is called by authuser/partner only using imports and not the old client configuration
- [ ] rearrange the test helper and split them per service so that it is easier to organise files and find what you want.
Some files need to renamed and split.
- [ ] make sure that the settings is called only using imports and not the old client configuration
- [ ] Tests to fix:
    - catalog:
        - set the auth middleware for all tests where this is necessary
        - TestGetPrice
        - TestGetPricesByProductID
        - TestUpdatePrice

- [ ] Tests to complete:
    - product:
        - TestGetAdminAllProducts

- [ ] check for the code in app and make it work. You need to find the right names for the services use the right dependency.
use the right environment (common/envmode is for that purpose)




