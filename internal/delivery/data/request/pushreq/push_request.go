package pushreq

type PushSubscriptionKeys struct {
	P256DH string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type UpsertPushSubscriptionRequest struct {
	Endpoint       string               `json:"endpoint"`
	ExpirationTime *int64               `json:"expirationTime"`
	Keys           PushSubscriptionKeys `json:"keys"`
	Timezone       string               `json:"timezone,omitempty"`
}

type DeletePushSubscriptionRequest struct {
	Endpoint string `json:"endpoint"`
}
