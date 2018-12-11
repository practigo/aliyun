package aliyun

import "fmt"

// CanonicalizedError defines the common request response
// if any error occurs.
// It also implements the error interface.
type CanonicalizedError struct {
	Code      string `json:"Code,omitempty"`
	Message   string `json:"Message,omitempty"`
	RequestID string `json:"RequestId,omitempty"`
	HostID    string `json:"HostId,omitempty"`
	Status    int    `json:"status,omitempty"` // from HTTP
}

func (ce *CanonicalizedError) Error() string {
	return fmt.Sprintf("Status: %d - Error: %s (%s); Host: %s, RequestID: %s",
		ce.Status, ce.Code, ce.Message, ce.HostID, ce.RequestID)
}

// some common error codes for all products
// more on https://error-center.aliyun.com/status/product/Public
const (
	ErrCodeForbidden             = "Forbidden"
	ErrCodeInternalError         = "InternalError"
	ErrCodeInvalidParameter      = "InvalidParameter"
	ErrCodeUnknownError          = "UnknownError"
	ErrCodeSignatureNonceUsed    = "SignatureNonceUsed"
	ErrCodeUnsupportedHTTPMethod = "UnsupportedHTTPMethod"
	ErrCodeAPINotFound           = "InvalidApi.NotFound"
	ErrCodeMissingSecurityToken  = "MissingSecurityToken"
	ErrCodeSignatureDoesNotMatch = "SignatureDoesNotMatch"
)
