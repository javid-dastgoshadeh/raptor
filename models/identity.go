package models

type PublicMetadata struct {
	Status string `json:"status,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

type AdminMetadata struct {
	HashedPassword     interface{} `json:"hashed_password,omitempty"`
	CreatedAtOldSystem interface{} `json:"created_at_old_system"`
	CloseBackup        interface{} `json:"close_backup"`
	Trash              interface{} `json:"trash"`
}

type Credential struct {
	Oidc     interface{} `json:"oidc"`
	Password interface{} `json:"password"`
}

type Traits struct {
	Email               interface{} `json:"email,omitempty"`
	PhoneNumber         interface{} `json:"phone_number,omitempty"`
	Username            interface{} `json:"username,omitempty"`
	EmailVerified       interface{} `json:"email_verified,omitempty"`
	PhoneNumberVerified interface{} `json:"phone_number_verified,omitempty"`
	Name                *Name       `json:"name,omitempty"`
	DisplayName         interface{} `json:"display_name,omitempty"`
	Nickname            interface{} `json:"nickname,omitempty"`
	Avatar              interface{} `json:"avatar,omitempty"`
	Socials             *Social     `json:"socials,omitempty"`
}

type Name struct {
	First interface{} `json:"first"`
	Last  interface{} `json:"last"`
}

type Social struct {
	Instagram interface{} `json:"instagram"`
	Twitter   interface{} `json:"twitter"`
	Telegram  interface{} `json:"telegram"`
}

type Identity struct {
	Credential     *Credential     `json:"credential,omitempty"`
	Traits         *Traits         `json:"traits,omitempty"`
	State          interface{}     `json:"state,omitempty"`
	SchemaID       interface{}     `json:"schema_id"`
	MetadataAdmin  *AdminMetadata  `json:"metadata_admin"`
	MetadataPublic *PublicMetadata `json:"metadata_public"`
}

type State string

const (
	Active   State = "active"
	Inactive State = "inactive"
)

type MessageSender string

const (
	Sms      MessageSender = "sms"
	Email    MessageSender = "email"
	Username MessageSender = "username"
)

type Device string

const (
	App Device = "app"
	Web Device = "web"
)
