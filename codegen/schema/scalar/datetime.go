package scalar

import (
	"errors"

	"io"
	"time"
)

type Datetime struct {
	Value time.Time
}

// UnmarshalGQL implements the graphql.Marshaler interface
func (d *Datetime) UnmarshalGQL(v interface{}) error {
	ds, err := v.(string)
	if !err {
		return errors.New("fail unmarshal datetime: convert to string")
	}

	var errParse error
	d.Value, errParse = time.Parse(time.RFC3339, ds)
	if errParse != nil {
		return errParse
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (d Datetime) MarshalGQL(w io.Writer) {
	w.Write([]byte(d.Value.Format(time.RFC3339)))
	//res, _ := json.Marshal(d.Value.String())
	//w.Write(res)
}
