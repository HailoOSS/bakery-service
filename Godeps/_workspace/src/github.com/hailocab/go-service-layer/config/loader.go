package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"launchpad.net/tomb"

	nsqlib "github.com/hailocab/go-nsq"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	// Worst-case config loader interval -- hard attempt to reload config
	configPollInterval = 300 * time.Second
	// Sleep between load failure
	configRetryDelay = 1 * time.Second
	// For file loader, how often to inspect mtime
	filePollInterval = 10 * time.Second
)

var defaultLoader *Loader

type reader func() (io.ReadCloser, error)

// Loader represents a loader of configuration. It automatically reloads when it receives notification that it should
// do so on its changes channel, and also every configPollInterval
type Loader struct {
	tomb.Tomb
	c          *Config
	changes    <-chan bool
	r          reader
	reloadLock sync.Mutex
}

// Load will go and grab the config via the reader and then load it into the config
func (ldr *Loader) Load() error {
	r, err := ldr.r()
	if err != nil {
		return err
	}
	// make sure we close this read-closer that we've just got, once we're done with it
	defer r.Close()
	err = ldr.c.Load(r)
	if err != nil {
		return err
	}

	return nil
}

func (ldr *Loader) reload() {
	ldr.reloadLock.Lock()
	defer ldr.reloadLock.Unlock()

	for {
		if err := ldr.Load(); err != nil {
			log.Warnf("[Config] Failed to reload config: %v", err)
			time.Sleep(configRetryDelay)
			continue
		}
		break
	}
}

func NewLoader(c *Config, changes chan bool, r reader) *Loader {
	ldr := &Loader{
		c:       c,
		changes: changes,
		r:       r,
	}

	go func() {
		defer ldr.Done()

		tick := time.NewTicker(configPollInterval)
		defer tick.Stop()

		for {
			select {
			case <-changes:
				ldr.reload()
			case <-tick.C:
				ldr.reload()
			case <-ldr.Dying():
				log.Tracef("[Config] Loader dying: %s", ldr.Err().Error())
				return
			}
		}
	}()

	// Spit a change down the pipe to load now
	changes <- true

	return ldr
}

// NewFileLoader returns a loader that reads config from file fn
func NewFileLoader(c *Config, fn string) (*Loader, error) {
	log.Infof("[Config] Initialising config loader to load from file '%s'", fn)

	rdr := func() (io.ReadCloser, error) {
		file, err := os.Open(fn)
		if err != nil {
			return nil, fmt.Errorf("Error opening file %v: %v", fn, err)
		}
		return file, nil
	}

	changesChan := make(chan bool)
	l := NewLoader(c, changesChan, rdr)
	go func() {
		tick := time.NewTicker(filePollInterval)
		defer tick.Stop()

		var lastMod time.Time
		for {
			select {
			case <-tick.C:
				if fi, err := os.Stat(fn); err == nil {
					if fi.ModTime().After(lastMod) {
						lastMod = fi.ModTime()
						changesChan <- true
					}
				}
			case <-l.Dying():
				// When the loader dies, we should too
				return
			}
		}
	}()

	return l, nil
}

// NewServiceLoader returns a loader that reads config from config service
func NewServiceLoader(c *Config, addr, service, region, env string) (*Loader, error) {
	// define our hierarchy:
	// H2:BASE
	// H2:BASE:<service-name>
	// H2:REGION:<aws region>
	// H2:REGION:<aws region>:<service-name>
	// H2:ENV:<env>
	// H2:ENV:<env>:<service-name>

	hierarchy := []string{
		"H2:BASE",
		fmt.Sprintf("H2:BASE:%s", service),
		fmt.Sprintf("H2:REGION:%s", region),
		fmt.Sprintf("H2:REGION:%s:%s", region, service),
		fmt.Sprintf("H2:ENV:%s", env),
		fmt.Sprintf("H2:ENV:%s:%s", env, service),
	}

	// construct URL
	if !strings.Contains(addr, "://") {
		addr = "https://" + addr
	}
	addr = strings.TrimRight(addr, "/") + "/compile"
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config service address: %v", err)
	}
	q := u.Query()
	q.Set("ids", strings.Join(hierarchy, ","))
	u.RawQuery = q.Encode()

	configUrl := u.String()

	log.Infof("[Config] Initialising service loader for service '%s' in region '%s' in '%s' environment via URL %s", service, region, env, configUrl)

	rdr := func() (io.ReadCloser, error) {
		rsp, err := http.Get(configUrl)
		if err != nil {
			log.Errorf("[Config] Failed to load config via %s: %v", configUrl, err)
			return nil, fmt.Errorf("Failed to load config via %s: %v", configUrl, err)
		}
		defer rsp.Body.Close()
		if rsp.StatusCode != 200 {
			log.Errorf("[Config] Failed to load config via %s - status code %v", configUrl, rsp.StatusCode)
			return nil, fmt.Errorf("Failed to load config via %s - status code %v", configUrl, rsp.StatusCode)
		}
		b, _ := ioutil.ReadAll(rsp.Body)

		loaded := make(map[string]interface{})
		err = json.Unmarshal(b, &loaded)
		if err != nil {
			log.Errorf("[Config] Unable to unmarshal loaded config: %v", err)
			return nil, fmt.Errorf("Unable to unmarshal loaded config: %v", err)
		}

		b, err = json.Marshal(loaded["config"])
		if err != nil {
			log.Errorf("[Config] Unable to unmarshal loaded config: %v", err)
			return nil, fmt.Errorf("Unable to unmarshal loaded config: %v", err)
		}
		rdr := ioutil.NopCloser(bytes.NewReader(b))
		return rdr, nil
	}

	changesChan := make(chan bool)
	l := NewLoader(c, changesChan, rdr)

	go func() {
		// wait until loaded
		l.reload()

		// look out for config changes PUBbed via NSQ -- subscribe via a random ephemeral channel
		channel := fmt.Sprintf("g%v#ephemeral", rand.Uint32())
		consumer, err := nsqlib.NewConsumer("config.reload", channel, nsqlib.NewConfig())
		if err != nil {
			log.Warnf("[Config] Failed to create NSQ reader to pickup config changes (fast reload disabled): ch=%v %v", channel, err)
			return
		}
		consumer.AddHandler(nsqlib.HandlerFunc(func(m *nsqlib.Message) error {
			changesChan <- true
			return nil
		}))

		// now configure it -- NOT via lookupd! There is a bug we think!
		subHosts := AtPath("hailo", "service", "nsq", "subHosts").AsHostnameArray(4150)
		if len(subHosts) == 0 {
			log.Warnf("[Config] No subHosts defined for config.reload topic (fast reload disabled)")
			return
		}

		log.Infof("[Config Load] Subscribing to config.reload (for fast config reloads) via NSQ hosts: %v", subHosts)
		if err := consumer.ConnectToNSQDs(subHosts); err != nil {
			log.Warnf("[Config Load] Failed to connect to NSQ for config changes (fast reload disabled): %v", err)
			return
		}

		// Wait for the Loader to be killed, and then stop the NSQ reader
		l.Wait()
		consumer.Stop()
	}()

	return l, nil
}

// LoadFromFile will load config from a flat text file containing JSON into the default instance
func LoadFromFile(fn string) (err error) {
	if defaultLoader != nil {
		defaultLoader.Killf("Replaced by LoadFromFile")
	}

	defaultLoader, err = NewFileLoader(DefaultInstance, fn)
	return err
}

// LoadFromService will load config from the config service into the default instance
func LoadFromService(service string) (err error) {
	addr := os.Getenv("H2_CONFIG_SERVICE_ADDR")
	region := os.Getenv("EC2_REGION")
	env := os.Getenv("H2O_ENVIRONMENT_NAME")

	if len(addr) == 0 {
		log.Critical("[Config] Config service address not set")
		log.Flush()
		os.Exit(1)
	}

	if defaultLoader != nil {
		defaultLoader.Killf("Replaced by LoadFromService")
	}

	defaultLoader, err = NewServiceLoader(DefaultInstance, addr, service, region, env)
	return err
}
