package mock

import (
	"testing"

	"github.com/Wito982/go_course_sdk/user"
)

func TestMock_User(t *testing.T) {

	t.Run("should return an error", func(t *testing.T) {
		var _ user.Transport = (*UserSDKMock)(nil)
	})

}
