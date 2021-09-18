package gelatin

import "log"

type GelatinClient struct {
	// Export data from this service
	from GelatinService

	// Import data into this service
	into GelatinService
}

func NewGelatinClient(from GelatinService, into GelatinService) *GelatinClient {
	return &GelatinClient{
		from,
		into,
	}
}

func (c *GelatinClient) MigrateUsers(passwords map[string]string) error {
	fromUsers, err := c.from.User().GetUsers(false)
	if err != nil {
		return err
	}

	intoUsers, err := c.into.User().GetUsers(false)
	if err != nil {
		return err
	}

	var createUsers []*GelatinUser
	var deleteUsers []*GelatinUser

	// Find users we need to create
	for i, fromUser := range fromUsers {
		exists := false
		for _, intoUser := range intoUsers {
			if fromUser.Name == intoUser.Name {
				exists = true
				break
			}
		}

		if !exists {
			createUsers = append(createUsers, &fromUsers[i])
		}
	}

	// Find users we need to delete
	for i, intoUser := range intoUsers {
		exists := false
		for _, fromUser := range fromUsers {
			if intoUser.Name == fromUser.Name {
				exists = true
				break
			}
		}

		if !exists {
			deleteUsers = append(deleteUsers, &intoUsers[i])
		}
	}

	for i, user := range createUsers {
		newUser, err := c.into.User().CreateUser(user.Name)
		if err != nil {
			return err
		}

		log.Printf("created %s: %s", newUser.Name, newUser.Id)

		createUsers[i] = newUser
	}

	// TODO: Migrate this to iterate over deleteUsers
	for i, user := range createUsers {
		err := c.into.User().DeleteUser(user.Id)
		if err != nil {
			return err
		}

		log.Printf("Deleted: %s", user.Name)

		createUsers[i] = nil
	}

	return nil
}
