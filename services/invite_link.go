package services

func GenerateInviteLink(groupId string) (string, bool) {
	if groupId == "" {
		return "", false
	}
	inviteLink := "https://explit.com/invite/" + groupId
	return inviteLink, true
}