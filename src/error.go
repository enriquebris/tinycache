package tinycache

import "fmt"

const (
	errorEmptyKey    = "Empty key"
	errorMissingKey  = "The key '%s' is not present"
	errorExceedLimit = "The value for '%s' exceeds the limits"
)

type TinyCacheError struct {
	msg         string
	missingkey  bool
	emptykey    bool
	exceedlimit bool
}

func (tce *TinyCacheError) Error() string {
	return tce.msg
}

func missingKeyError(key string) string {
	return fmt.Sprintf(errorMissingKey, key)
}
