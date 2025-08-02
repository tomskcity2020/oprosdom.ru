package shared_models

type MsgCode struct {
    Urgent      bool   `json:"urgent"`
    Type        string `json:"type"`        // "sms" или "call"
    PhoneNumber string `json:"phone_number"` // формат e.194
    Message     string `json:"message"`
}