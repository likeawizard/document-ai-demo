package store_test

import (
	"testing"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/store"
	"github.com/stretchr/testify/assert"
)

// These test will most likely fail on non Linux/Unix systems. The app itself might run on Windows but it has been not tested there.
func TestSystemStorePath(t *testing.T) {
	type testCase struct {
		name       string
		location   string
		filename   string
		want       string
		shouldFail bool
	}

	tcs := []testCase{
		{
			name:       "Missing filename",
			location:   "",
			filename:   "",
			want:       "",
			shouldFail: true,
		},
		{
			name:       "Filename no location",
			location:   "",
			filename:   "somefile.png",
			want:       "somefile.png",
			shouldFail: false,
		},
		{
			name:       "Filename and location",
			location:   "subdirectory",
			filename:   "somefile.png",
			want:       "subdirectory/somefile.png",
			shouldFail: false,
		},
		{
			name:       "Filename and location with trailing /",
			location:   "subdirectory/",
			filename:   "somefile.png",
			want:       "subdirectory/somefile.png",
			shouldFail: false,
		},
		{
			name:       "Filename and location at root",
			location:   "/",
			filename:   "somefile.png",
			want:       "/somefile.png",
			shouldFail: false,
		},
		{
			name:       "Filename and location with depth",
			location:   "dir/childDir",
			filename:   "somefile.png",
			want:       "dir/childDir/somefile.png",
			shouldFail: false,
		},
		{
			name:       "Filename with traversal up",
			location:   "..",
			filename:   "somefile.png",
			want:       "../somefile.png",
			shouldFail: false,
		},
		{
			name:       "Filename with redundant traversal removed",
			location:   "whateverImGoingBackAnyway/../dir",
			filename:   "somefile.png",
			want:       "dir/somefile.png",
			shouldFail: false,
		},
	}

	for _, tc := range tcs {
		ss := store.NewSystemStore(config.StorageCfg{Location: tc.location})
		got, err := ss.GetPath(tc.filename)
		if tc.shouldFail {
			assert.Error(t, err, tc.name)
		}
		assert.Equal(t, tc.want, got, tc.name)
	}

}
