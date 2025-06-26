package seeder

import userEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"

func mapEmails(emailSeeds []EmailSeed) []userEntity.Email {
	var emails []userEntity.Email
	for _, e := range emailSeeds {
		emails = append(emails, userEntity.Email{
			Email:      e.Email,
			IsVerified: e.IsVerified,
		})
	}
	return emails
}

func mapPhones(phoneSeeds []PhoneSeed) []userEntity.Phone {
	var phones []userEntity.Phone
	for _, p := range phoneSeeds {
		phones = append(phones, userEntity.Phone{
			Number:     p.Number,
			IsVerified: p.IsVerified,
		})
	}
	return phones
}

// ADDED: Helper function to convert status
func mapUserStatus(status string) userEntity.UserStatus {
	switch status {
	case "active":
		return userEntity.UserStatusActive
	case "invited":
		return userEntity.UserStatusInvited
	case "suspended":
		return userEntity.UserStatusSuspended
	case "removed":
		return userEntity.UserStatusRemoved
	default:
		return userEntity.UserStatusInvited // Default fallback
	}
}
