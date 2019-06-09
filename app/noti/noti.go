package noti

func notificationProviders(cs CredentialStore) []Provider {

	var ps []Provider
	if cs.GetAPIFile("firebase") != "" {
		fb := NewFirebase(cs.GetAPIFile("firebase"))
		ps = append(ps, fb)
	}

	ps = append(ps, &WebHooked{})

	return ps
}
