package log

import (
	"time"
	"fmt"
)

type Data struct {
	UserID int
	Object string
	CreatedAt time.Time
	Action string

}

type Log interface {
	Save (d []Data) error
}

type log struct{}

func (r log) Save(d []Data) error {
	for _, elem := range d {
		fmt.Println(elem.UserID, elem.Object, elem.CreatedAt, elem.Action)
	}
	return nil
}



