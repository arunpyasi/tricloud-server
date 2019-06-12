package noti

var disableNoti = false

func notificationProviders(cs CredentialStore) []Provider {

	var ps []Provider
	if disableNoti {
		return ps
	}

	if cs.GetAPIFile("firebase") != "" {
		fb := NewFirebase(cs.GetAPIFile("firebase"))
		ps = append(ps, fb)
	}

	ps = append(ps, &WebHooked{})

	return ps
}
