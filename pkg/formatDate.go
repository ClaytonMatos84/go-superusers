package pkg

import "time"

type CustomDate struct {
	time.Time
}

const dateFormat = "2006-01-02"

func (d *CustomDate) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(`"`+dateFormat+`"`, string(b))
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

func (d CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Format(dateFormat) + `"`), nil
}
