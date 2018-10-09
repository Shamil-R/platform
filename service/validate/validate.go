package validate

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Data struct {
	Rule      string
	RuleValue string
	Value     string
	Type      string
}

type Validate interface {
	Validate(d []Data) error
}

type validate struct{}

func (v validate) Validate(d []Data) error {
	for _, elem := range d {
		switch elem.Rule {
		case "min":

			if elem.Type == "String" {
				temp, err := strconv.Atoi(elem.RuleValue)
				if err != nil {
					return err
				}

				err = validation.Validate(elem.Value,
					validation.Length(temp, math.MaxInt64),
				)
				if err != nil {
					return err
				}

			} else if elem.Type == "Int" {
				err := validation.Validate(elem.Value,
					validation.Min(elem.RuleValue),
				)
				if err != nil {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("Валидация директивы %s для типа %s не реализована в исходном коде", elem.Rule, elem.Type))
			}

		case "max":
			if elem.Type == "String" {
				temp, err := strconv.Atoi(elem.RuleValue)
				if err != nil {
					return err
				}

				err = validation.Validate(elem.Value,
					validation.Length(0, temp), // not empty
				)
				if err != nil {
					return err
				}

			} else if elem.Type == "Int" {
				err := validation.Validate(elem.Value,
					validation.Max(elem.RuleValue), // not empty
				)
				if err != nil {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("Валидация директивы %s для типа %s не реализована в исходном коде", elem.Rule, elem.Type))
			}

		default:
			return errors.New(fmt.Sprintf("Директива %s не реализована в исходном коде", elem.Rule))
		}
	}
	return nil
}
