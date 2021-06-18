package envconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
)

type database struct {
	Host     string
	Username string
	Password string
	Maxidle  int
}

type logconf struct {
	LogModel    string `toml:"log_mode" yml:"log_mode"`
	LogLevel    string `toml:"log_level" yml:"log_level"`
	LogMaxFiles int    `toml:"log_max_files" yml:"log_max_files"`
}

type app struct {
	Env       string
	Debug     string
	FooEnv    string  `toml:"foo_env" yml:"foo_env"`
	TimeAfter float32 `toml:"time_after" yml:"time_after"`
	Databases map[string]database
	Log       logconf
}

func TestConf_Env(t *testing.T) {

	files := []string{"./testdata/app.toml", "./testdata/app.yml"}

	_ = os.Setenv("FOO_ENV", "abc")

	for _, file := range files {
		conf, err := New(file)
		if err != nil {
			t.Fatal(err)
		}

		err = conf.Env("./testdata/.env.local")

		if err != nil {
			t.Fatal(err)
		}

		app := &app{}

		err = conf.Unmarshal(app)
		if err != nil {
			t.Fatal(err)
		}

		if app.Databases["alpha"].Host != "127.0.0.1" {
			t.Fatal(errors.New(fmt.Sprintf("env parse error, must get %s but get %s", "127.0.0.1", app.Databases["alpha"].Host)))
		}

		if app.Databases["beta"].Password != "123456" {
			t.Fatal(errors.New(fmt.Sprintf("env parse error, must get %s but get %s", "123456", app.Databases["beta"].Password)))
		}

		if app.Databases["beta"].Maxidle != 100 {
			t.Fatal(errors.New(fmt.Sprintf("env parse error, must get %d but get %d", 100, app.Databases["beta"].Maxidle)))
		}

		if app.TimeAfter != 0.125 {
			t.Fatal(errors.New(fmt.Sprintf("env parse error, must get %s but get %f", "0.125", app.TimeAfter)))
		}

		if app.FooEnv != "abc" {
			t.Fatal(errors.New(fmt.Sprintf("env parse error, must get %s but get %s", "abc", app.FooEnv)))
		}

		b, _ := json.Marshal(app)

		t.Log(string(b))
	}
}

func TestConf_EnvWithReader(t *testing.T) {
	files := []string{"./testdata/app.toml", "./testdata/app.yml"}

	for _, file := range files {
		conf, err := New(file)
		if err != nil {
			t.Fatal(err)
		}
		file, err := os.Open("./testdata/.env.local")
		if err != nil {
			return
		}
		defer file.Close()

		err = conf.EnvWithReader(file)

		if err != nil {
			t.Fatal(err)
		}

		app := &app{}

		err = conf.Unmarshal(app)
		if err != nil {
			t.Fatal(err)
		}

		if app.Databases["alpha"].Host != "127.0.0.1" {
			t.Fatal(errors.New(fmt.Sprintf("env parse error, must get %s but get %s", "127.0.0.1", app.Databases["alpha"].Host)))
		}

		if app.Databases["beta"].Password != "123456" {
			t.Fatal(errors.New(fmt.Sprintf("env parse error, must get %s but get %s", "123456", app.Databases["beta"].Password)))
		}

		if app.Databases["beta"].Maxidle != 100 {
			t.Fatal(errors.New(fmt.Sprintf("env parse error, must get %d but get %d", 100, app.Databases["beta"].Maxidle)))
		}

		if app.TimeAfter != 0.125 {
			t.Fatal(errors.New(fmt.Sprintf("env parse error, must get %s but get %f", "0.125", app.TimeAfter)))
		}

		b, _ := json.Marshal(app)

		t.Log(string(b))
	}
}

func TestViper(t *testing.T) {
	// viper.MergeConfig()
}
