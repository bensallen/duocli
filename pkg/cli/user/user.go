package user

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bensallen/duocli/pkg/duocli"
	"github.com/duosecurity/duo_api_golang/admin"
	"github.com/urfave/cli/v2"
)

func Create(c *cli.Context) error {
	username := c.String("username")
	if username == "" {
		return fmt.Errorf("username argument required")
	}

	groups := c.StringSlice("group")
	email := c.String("email")
	firstName := c.String("firstName")
	lastName := c.String("lastName")
	status := c.String("status")

	if status != "" {
		switch status {
		case "active":
		case "bypass":
		case "disabled":
		default:
			return fmt.Errorf("status not set to active, bypass or disabled: %s", status)

		}
	}

	adm, err := duocli.LoadAdminConfig(c.String("config"))
	if err != nil {
		return err
	}

	log.Printf("adding user %s", username)

	user := admin.User{
		Username:  username,
		Email:     email,
		FirstName: &firstName,
		LastName:  &lastName,
		Status:    status,
	}

	result, err := adm.CreateUser(user.URLValues())
	if err != nil {
		return err
	}

	if result.Stat != "OK" {
		return fmt.Errorf("Duo API returned non-ok status response when creating user %s, message: %s", username, *result.Message)
	}

	return associateGroupsWithUser(result.Response, groups, adm)
}

func Modify(c *cli.Context) error {
	username := c.String("username")
	addgroups := c.StringSlice("addgroup")
	delgroups := c.StringSlice("delgroup")
	email := c.String("email")
	firstName := c.String("firstName")
	lastName := c.String("lastName")
	status := c.String("status")
	create := c.Bool("create")

	if status != "" {
		switch status {
		case "active":
		case "bypass":
		case "disabled":
		default:
			return fmt.Errorf("status not set to active, bypass or disabled: %s", status)

		}
	}

	adm, err := duocli.LoadAdminConfig(c.String("config"))
	if err != nil {
		return err
	}

	getUser, err := adm.GetUsers(admin.GetUsersUsername(username))
	if err != nil {
		return err
	}
	if len(getUser.Response) == 0 && !create {
		return fmt.Errorf("user not found %s", username)
	}
	if len(getUser.Response) > 1 {
		return fmt.Errorf("more than one user found with this username or alias")
	}

	user := admin.User{
		Username:  username,
		Email:     email,
		FirstName: &firstName,
		LastName:  &lastName,
		Status:    status,
	}

	var result *admin.GetUserResult
	if create && len(getUser.Response) == 0 {
		log.Printf("adding user %s", username)
		result, err = adm.CreateUser(user.URLValues())
	} else {
		log.Printf("updating user %s", username)
		result, err = adm.ModifyUser(getUser.Response[0].UserID, user.URLValues())
	}
	if err != nil {
		return err
	}

	if result.Stat != "OK" {
		return fmt.Errorf("Duo API returned non-ok status response when modifying user %s, message: %s", username, *result.Message)
	}

	if err := associateGroupsWithUser(result.Response, addgroups, adm); err != nil {
		return err
	}

	// Don't try to delete groups if the user was just created
	if create && len(getUser.Response) == 0 {
		return nil
	}

	return disassociateGroupsWithUser(result.Response, delgroups, adm)
}

func associateGroupsWithUser(user admin.User, groups []string, adm *admin.Client) error {
	if len(groups) > 0 {
		duoGroups, err := adm.GetGroups()
		if err != nil {
			return err
		}

	GROUP:
		for _, group := range groups {
			for _, userDuoGroup := range user.Groups {
				if group == userDuoGroup.Name {
					log.Printf("group %s is already associated with user %s, skipping", group, user.Username)
					continue GROUP
				}
			}
			var grpFound string
			for _, duoGroup := range duoGroups.Response {
				if group == duoGroup.Name {
					grpFound = duoGroup.GroupID
				}
			}
			if grpFound != "" {
				log.Printf("associating group %s with user %s", group, user.Username)
				result, err := adm.AssociateGroupWithUser(user.UserID, grpFound)
				if err != nil {
					return err
				}
				if result.Stat != "OK" {
					return fmt.Errorf("Duo API returned non-ok status response when associating user with group %s, message: %s", group, *result.Message)
				}
			} else {
				log.Printf("warning, group %s not found in Duo, skipping", group)
				continue
			}
		}
	}
	return nil
}

func disassociateGroupsWithUser(user admin.User, groups []string, adm *admin.Client) error {
	if len(groups) > 0 {
		duoGroups, err := adm.GetGroups()
		if err != nil {
			return err
		}

		for _, group := range groups {
			var existGroupFound bool
			for _, userDuoGroup := range user.Groups {
				if group == userDuoGroup.Name {
					existGroupFound = true
				}
			}

			if !existGroupFound {
				log.Printf("group %s is not associated with user %s, skipping", group, user.Username)
				continue
			}

			var grpFound string
			for _, duoGroup := range duoGroups.Response {
				if group == duoGroup.Name {
					grpFound = duoGroup.GroupID
				}
			}
			if grpFound != "" {
				log.Printf("disassociating group %s from user %s", group, user.Username)
				result, err := adm.DisassociateGroupFromUser(user.UserID, grpFound)
				if err != nil {
					return err
				}
				if result.Stat != "OK" {
					return fmt.Errorf("Duo API returned non-ok status response when disassociating user with group %s, message: %s", group, *result.Message)
				}
			} else {
				log.Printf("warning, group %s not found in Duo, skipping", group)
				continue
			}
		}
	}
	return nil
}

func Delete(c *cli.Context) error {
	usernames := c.StringSlice("username")
	devices := c.Bool("devices")

	if len(usernames) == 0 {
		return fmt.Errorf("username must be specified")
	}

	adm, err := duocli.LoadAdminConfig(c.String("config"))
	if err != nil {
		return err
	}

	for _, user := range usernames {
		result, err := adm.GetUsers(admin.GetUsersUsername(user))
		if err != nil {
			return err
		}
		if result.Stat != "OK" {
			return fmt.Errorf("Duo API returned non-ok status response on searching for user: %s, message: %s", user, *result.Message)
		}
		if len(result.Response) == 0 {
			log.Printf("warning, user %s not found, skipping", user)
			continue
		}
		if len(result.Response) > 1 {
			log.Printf("warning, searching for user %s returned more than one result, skipping", user)
			continue
		}

		if devices {
			for _, phone := range result.Response[0].Phones {
				log.Printf("deleting phone %s from user %s", phone.Name, user)
				phoneResult, err := adm.DeletePhone(phone.PhoneID)
				if err != nil {
					return err
				}
				if phoneResult.Stat != "OK" {
					log.Printf("warning, Duo API returned non-ok status response when deleting phone: %s for user: %s, message: %s", phone.Name, user, *result.Message)
				}
			}
		}

		log.Printf("deleting user %s", user)
		deleteResult, err := adm.DeleteUser(result.Response[0].UserID)
		if err != nil {
			return err
		}
		if deleteResult.Stat != "OK" {
			log.Printf("warning, Duo API returned non-ok status response when deleting user: %s, message: %s", user, *result.Message)
		}
	}
	return nil
}

func Get(c *cli.Context) error {
	usernames := c.StringSlice("username")

	if len(usernames) == 0 {
		return fmt.Errorf("username must be specified")
	}

	adm, err := duocli.LoadAdminConfig(c.String("config"))
	if err != nil {
		return err
	}

	userResponses := []admin.User{}

	for _, user := range usernames {
		result, err := adm.GetUsers(admin.GetUsersUsername(user))
		if err != nil {
			return err
		}
		if result.Stat != "OK" {
			return fmt.Errorf("Duo API returned non-ok status response on searching for user: %s, with message: %s", user, *result.Message)
		}
		if len(result.Response) == 0 {
			log.Printf("warning, user %s not found, skipping", user)
			continue
		}
		if len(result.Response) > 1 {
			log.Printf("warning, searching for user %s returned more than one result, skipping", user)
			continue
		}

		userResponses = append(userResponses, result.Response[0])
	}

	if len(userResponses) > 0 {
		jsonOut, err := json.Marshal(&userResponses)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", jsonOut)
	} else {
		log.Printf("No results")
	}

	return nil
}
