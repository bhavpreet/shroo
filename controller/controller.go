package controller

import (
	"os"
	"path"
	"plugin"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/bhavpreet/config"
	"github.com/go-bongo/go-dotaccess"
	"github.com/sirupsen/logrus"
)

type peripheral struct {
	Plug *plugin.Plugin
	Defs map[string]interface{}                   `yaml:"defs"`
	Func map[string]func(interface{}) interface{} `yaml:"func,omitempty"`
}

type exec struct {
	Func *string     `yaml:"func,omitempty"`
	Args interface{} `ymal:"args,omitempty"`
	Ret  *string     `ymal:"args,omitempty"`
}

type execs []exec

type actions []map[string]*action
type action struct {
	last  time.Time
	Every *Duration `yaml:"every,omitempty"`
	At    *Time     `yaml:"at,omitempty"`
	For   *Duration `yaml:"duration,omitempty"`
	If    *string   `yaml:"if,omitempty"`
	Exec  *execs    `yaml:"exec,omitempty"`
	Else  *execs    `yaml:"else,omitempty"`
	A     *actions  `yaml:"actions,omitempty"`
}

type conf struct {
	Tick        *Duration              `yaml:"tick"`
	Vars        map[string]interface{} `yaml:"vars,omitempty"` // global vars
	Peripherals map[string]peripheral  `yaml:"peripherals,flow,omitempty"`

	Actions actions `yaml:"actions,flow,omitempty"`
}

var cfg *conf

type Duration struct {
	time.Duration
}

type Time struct {
	time.Time
}

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	var err error
	if err = unmarshal(&s); err != nil {
		return err
	}

	d.Duration, err = time.ParseDuration(s)
	return err
}

func (t *Time) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	var err error
	if err = unmarshal(&s); err != nil {
		return err
	}

	t.Time, err = time.Parse(time.Kitchen, s)
	return err
}

func Init() {
	cfg = new(conf)
	config.Bind(cfg)
	for k := range cfg.Peripherals {
		logrus.Debug(k)
	}
	// config.Init()

	var err error
	var dir string
	dir = os.Getenv("LD_LIBRARY_PATH")
	if dir == "" {
		dir = path.Dir(os.Args[0])
	}

	// for all peripherals, call Init
	for p := range cfg.Peripherals {
		x := cfg.Peripherals[p]
		if x.Plug, err = plugin.Open(dir + "/" + p + "/" + p + ".so"); err != nil {
			logrus.WithFields(logrus.Fields{"err": err}).Fatalf("Unable to open plugin: %s", p)
		}

		// Init
		init, err := x.Plug.Lookup("Init")
		if err != nil {
			logrus.WithFields(logrus.Fields{"err": err}).Fatalf("Function init not found for Plugin %s", p)
		}

		if f, ok := init.(func(interface{}) interface{}); ok {
			logrus.Debugf("%s.Init...", p)
			f(cfg.Peripherals[p].Defs)
		} else {
			logrus.Debugf("Unknown Init for %s", p)
		}

		// Map functions
		for f := range x.Func {
			if s, err := x.Plug.Lookup(f); err != nil {
				logrus.WithFields(logrus.Fields{"err": err}).Fatalf("Function %s not found for Plugin %s", f, p)
			} else {
				x.Func[f] = s.(func(interface{}) interface{})
				logrus.Debugf("Found %s for %s", f, p)
			}
		}

	}
}

func (exe *exec) Run(act *action) {
	fn, err := dotaccess.Get(cfg, *exe.Func)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":    err,
			"action": act,
		}).Fatal("Bad Exec")
	}

	var args interface{}
	if exe.Args != nil {
		if _, ok := exe.Args.([]interface{}); ok {
			ea := exe.Args.([]interface{})
			_args := make([]interface{}, len(ea))
			for idx, a := range ea {
				_args[idx], _ = dotaccess.Get(cfg.Vars, a.(string))
			}
			args = _args
		} else {
			args, _ = dotaccess.Get(cfg.Vars, exe.Args.(string))
		}
	}

	logrus.Debug(*exe.Func)
	ret := fn.(func(interface{}) interface{})(args)

	if exe.Ret != nil {
		dotaccess.Set(cfg.Vars, *exe.Ret, ret)
		logrus.WithFields(logrus.Fields{
			"ret":           ret,
			"conf.vars.bme": cfg.Vars["bme"],
		}).Debug()
	}
}

func (_execs *execs) Run(act *action) {
	if _execs == nil {
		return
	}

	for _, exec := range *_execs {
		exec.Run(act)
	}
}

func (a *action) Run(now time.Time) {
	if a.Every != nil {
		if now.Sub(a.last) >= a.Every.Duration {
			a.last = now
			a.Exec.Run(a)
			// run sub actions
			a.A.Run(now)
		} else if a.For != nil && now.Sub(a.last) < a.For.Duration {
			a.Exec.Run(a)
			// run sub actions
			a.A.Run(now)
		} else {
			a.Else.Run(a)
		}
	}

	if a.At != nil {
		nowK, _ := time.Parse(time.Kitchen, now.Format(time.Kitchen))
		if nowK.Before(a.At.Time) {
			nowK = nowK.Add(24 * time.Hour)
		}

		if nowK.Sub(a.At.Time) < a.For.Duration {
			a.last = a.At.Time
			a.Exec.Run(a)
			// run sub actions
			a.A.Run(now)
		} else {
			a.Else.Run(a)
		}
	}

	if a.If != nil {
		value, err := gval.Evaluate(*a.If, cfg.Vars)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":           err,
				"if expression": *a.If,
			}).Fatalf("Unable to evaluate expression")
		}

		if value == true {
			a.last = now
			a.Exec.Run(a)
			// run sub actions
			a.A.Run(now)
		} else if a.For != nil && now.Sub(a.last) < a.For.Duration {
			// we dont update last here since duration has
			// to end come to and end
			a.Exec.Run(a)
			// run sub actions
			a.A.Run(now)
		} else {
			a.Else.Run(a)
		}
	}

}

func (acts *actions) Run(now time.Time) {
	if acts == nil {
		return
	}
	for _, act := range *acts {
		for a, exe := range act {
			logrus.Debugf("Running: %s", a)
			exe.Run(now)
		}
	}
}

func Run() {
	ticker := time.NewTicker(cfg.Tick.Duration)
	defer ticker.Stop()
	for now := range ticker.C {
		logrus.WithFields(logrus.Fields{
			"now": now,
		}).Debug()
		cfg.Actions.Run(now)
	}
}
