# EnvConf 基于viper + env的golang配置库
<a href="https://github.com/goapt/envconf/actions"><img src="https://github.com/goapt/envconf/workflows/build/badge.svg" alt="Build Status"></a>

为了解决配置集中化方案，基于viper格式和env配置自动替换的解决方案

## 理念
在应用程序配置文件中，我们只需要使用toml或者yml,ini等viper支持的配置文件写出一个配置文件骨架，里面不包含任何的敏感数据（如用户名和密码），然后在.env.local文件中写入敏感信息，程序启动的时候会自动将数据覆盖替换，如下

1、创建 app.toml 骨架配置文件
```toml
env = "local"
# debug  off close, on open
debug = "on"

time_after = 0.5

[databases]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [databases.alpha]
  host = ""
  username = ""
  password = ""
  maxidle = 10

  [databases.beta]
  host = ""
  username = ""
  password = ""
  maxidle = 10

[log]
    log_mode = "std"
    log_level = "debug"
    log_max_files = 15

```
2、创建敏感信息替换文件 `.env.local`
```
debug=off
time_after=0.125
# server config
databases.alpha.host=127.0.0.1
databases.alpha.username=root
databases.alpha.password=123456
databases.alpha.maxidle=100
databases.beta.host=127.0.0.2
databases.beta.username=root
databases.beta.password=123456
# use variables
databases.beta.maxidle=${databases.alpha.maxidle}

# log config
log.log_max_files=10
```

env文件遵循 [JSON-Path](https://github.com/pelletier/go-toml/tree/master/query) 规范

此处我们注意到两点

1、env文件支持注释，使用 `#` 开头的则为注释行

2、支持变量应用，比如 `databases.beta.maxidle` 就是对 `databases.alpha.maxidle` 值的引用，变量应用使用 `${KEY}` 格式，当ENV变量不在当前配置项中，则会从系统环境变量中获取

> 因此我们可以仅在本地创建.env.local，而在系统发布的时候由发布系统将ACM配置中心的env线上配置数据写入到项目目录的.env.prod文件，从而实现线上的安全隔离

3、载入配置文件，以及覆盖数据
```go
import "github.com/goapt/envconf"

type Database struct {
	Host     string
	Username string
	Password string
	Maxidle  int
}

type LogConf struct {
	LogModel    string `toml:"log_mode"`
	LogLevel    string `toml:"log_level"`
	LogMaxFiles int    `toml:"log_max_files"`
}

type App struct {
	Env       string
	Debug     string
	TimeAfter float32 `toml:"time_after"`
	Databases map[string]Database
	Log       LogConf
}

conf, err := envconf.New("./testdata/app.toml")

if err != nil {
    return err
}

err = envconf.Env("./testdata/.env.local")

if err != nil {
    return err
}

app := &App{}

err = envconf.Unmarshal(app)
if err != nil {
    return err
}
```

4、最后得到数据如下
```
env = "local"
# debug  off close, on open
debug = "off"

time_after = 0.125

[databases]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [databases.alpha]
  host = "127.0.0.1"
  username = "root"
  password = "123456"
  maxidle = 100

  [databases.beta]
  host = "127.0.0.2"
  username = "root"
  password = "123456"
  maxidle = 100

[log]
    log_mode = "std"
    log_level = "debug"
    log_max_files = 10
```

> Enjoy

## Thanks

https://github.com/joho/godotenv

https://github.com/spf13/viper