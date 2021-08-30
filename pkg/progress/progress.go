package progress

import "github.com/schollz/progressbar/v3"

func Bar(max int, callback func(bar *progressbar.ProgressBar) error) error {
	bar := progressbar.NewOptions(max)

	if err := callback(bar); err != nil {
		return err
	}

	if err := bar.Finish(); err != nil {
		return err
	}

	return nil

}
