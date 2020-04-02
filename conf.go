package envconf

import (
	"path/filepath"
	"strconv"

	"github.com/goapt/dotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Conf struct {
	*viper.Viper
	TagName string
}

type ViperOption func(v *viper.Viper)

func New(cf string, options ...ViperOption) (*Conf, error) {
	conf := viper.New()
	conf.SetConfigFile(cf)
	tagName := "mapstructure"
	ext := filepath.Ext(cf)
	if len(ext) > 1 {
		tagName = ext[1:]
	}

	for _, option := range options {
		option(conf)
	}

	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Conf{
		Viper:   conf,
		TagName: tagName,
	}, nil
}

// Env replace toml config items
func (v *Conf) Env(path string) error {
	m, err := dotenv.Read(path)
	if err != nil {
		return err
	}

	for k, val := range m {
		if has := v.Viper.IsSet(k); has {
			r, err := convert(v.Viper.Get(k), val)

			if err != nil {
				return err
			}

			v.Viper.Set(k, r)
		}
	}

	return nil
}

// Unmarshal to struct
func (v *Conf) Unmarshal(i interface{}) error {
	return v.Viper.Unmarshal(i, func(config *mapstructure.DecoderConfig) {
		config.TagName = v.TagName
	})
}

func convert(v interface{}, val string) (r interface{}, err error) {
	switch v.(type) {
	case int:
		r, err = strconv.ParseInt(val, 10, 0)
	case int8:
		r, err = strconv.ParseInt(val, 10, 8)
	case int16:
		r, err = strconv.ParseInt(val, 10, 16)
	case int32:
		r, err = strconv.ParseInt(val, 10, 32)
	case int64:
		r, err = strconv.ParseInt(val, 10, 64)
	case bool:
		r, err = strconv.ParseBool(val)
	case float32:
		r, err = strconv.ParseFloat(val, 32)
	case float64:
		r, err = strconv.ParseFloat(val, 64)
	default:
		r = val
	}

	return r, err
}
