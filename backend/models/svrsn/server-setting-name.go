package svrsn

// SvrSettingName Server Setting Name
type SvrSettingName string

const (
	SessSecretKey        SvrSettingName = "SessionSecretKey"
	OwnerRegistered      SvrSettingName = "OwnerRegistered"
	HashIDSalt           SvrSettingName = "HashIDSalt"
	Hostname             SvrSettingName = "Hostname"
	RequireEmailActivate SvrSettingName = "RequireEmailActivate"
	SenderEmailServer    SvrSettingName = "SenderEmailServer"
	SenderEmailAddress   SvrSettingName = "SenderEmailAddress"
	SenderEmailUser      SvrSettingName = "SenderEmailUser"
	SenderEmailPassword  SvrSettingName = "SenderEmailPassword"
)
