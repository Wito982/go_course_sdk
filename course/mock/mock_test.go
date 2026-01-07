package mock

import (
	"testing"

	"github.com/Wito982/go_course_sdk/course"
)

func TestMock_Course(t *testing.T) {

	t.Run("should return an error", func(t *testing.T) {
		var _ course.Transport = (*CourseSDKMock)(nil)
	})

}
