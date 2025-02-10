package model

// types defined in this file are used to represent the data structures of the WhatsApp API
// example payload can be found at https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks/payload-examples/

type Profile struct {
	Name string `json:"name"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WAID    string  `json:"wa_id"`
}

type Text struct {
	Body string `json:"body"`
}

type Image struct {
	MimeType string `json:"mime_type"`
	Sha256   string `json:"sha256"`
	ID       string `json:"id"`
}

type Audio struct {
	MimeType string `json:"mime_type"`
	Sha256   string `json:"sha256"`
	ID       string `json:"id"`
	Voice    bool   `json:"voice"`
}

type Message struct {
	From      string `json:"from"` // mapped using alias `from` in Python
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Text      *Text  `json:"text,omitempty"`  // Optional field
	Image     *Image `json:"image,omitempty"` // Optional field
	Audio     *Audio `json:"audio,omitempty"` // Optional field
	Type      string `json:"type"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts,omitempty"` // Optional field
	Messages         []Message `json:"messages,omitempty"` // Optional field
}

type Change struct {
	Value    Value            `json:"value"`
	Field    string           `json:"field"`
	Statuses []map[string]any `json:"statuses,omitempty"` // Using map for dynamic dict-like fields
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Payload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

type UserMessage struct {
	User    User    `json:"user"`
	Message *string `json:"message,omitempty"` // Optional field
	Image   *Image  `json:"image,omitempty"`   // Optional field
	Audio   *Audio  `json:"audio,omitempty"`   // Optional field
}
