package sms

type Sender interface {
	Send(*Config) error
}

func NewConfig() *Config {
	return &Config{
		ExtraData: make(map[string]interface{}),
	}
}

type Config struct {
	Mobile      string
	Content     string
	Template    string
	CallbackURL string
	SignName    string
	ExtraData   map[string]interface{}
}

func (c *Config) Extra(key string) interface{} {
	if v, y := c.ExtraData[key]; y {
		return v
	}
	return nil
}

func (c *Config) Get(key string) interface{} {
	return c.Extra(key)
}

func (c *Config) Add(key string, val interface{}) {
	c.ExtraData[key] = val
}

func (c *Config) Del(key string) {
	if _, y := c.ExtraData[key]; y {
		delete(c.ExtraData, key)
	}
}

func (c *Config) Clear() {
	c.ExtraData = make(map[string]interface{})
}
