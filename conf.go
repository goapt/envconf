package envconf

import (
	"strconv"

	"github.com/goapt/dotenv"
	"github.com/pelletier/go-toml"
)

type Conf struct {
	config *toml.Tree
}

func LoadFile(path string) (*Conf, error) {
	config, err := toml.LoadFile(path)

	if err != nil {
		return nil, err
	}

	return &Conf{
		config,
	}, nil
}

// Env replace toml config items
func (v *Conf) Env(path string) error {
	m, err := dotenv.Read(path)

	if err != nil {
		return err
	}

	for k, val := range m {
		if has := v.config.Has(k); has {
			r, err := convert(v.config.Get(k), val)

			if err != nil {
				return err
			}

			v.config.Set(k, r)
		}
	}

	return nil
}

// Unmarshal to struct
func (v *Conf) Unmarshal(i interface{}) error {
	return v.config.Unmarshal(i)
}

func convert(v interface{}, val string) (r interface{}, err error) {
	switch v.(type) {
	case int, int32, int64:
		r, err = strconv.ParseInt(val, 10, 0)
	case bool:
		r, err = strconv.ParseBool(val)
	case float32, float64:
		r, err = strconv.ParseFloat(val, 32)
	default:
		r = val
	}

	return r, err
}
