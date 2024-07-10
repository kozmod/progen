package core

import (
	"testing"
)

func Test(t *testing.T) {
	e := Engin{}
	e.AddFilesActions(
		WithPriority(
			NewAction("files",
				File{
					Path: "/file_1",
					Data: []byte("xxxxxxxxxx"),
				},
			),
			1,
		),
	).
		AddCreateDirsActions()

	_ = e.Run()
}
