package elastic

import (
	"encoding/json"
	"fmt"

	"github.com/hailocab/service-layer/config"

	log "github.com/cihub/seelog"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/olivere/elastic.v2"
)

var (
	// Config holds info about elastic search
	Config *esConfig
)

// Elastic data
type Elastic struct {
	Host  string
	Index string

	client *elastic.Client
}

// NewWithDefaults Create new elastic with default options
func NewWithDefaults() (*Elastic, error) {
	if Config == nil {
		return nil, fmt.Errorf("No defaults to use")
	}

	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("https://%s", Config.Host)))

	e := &Elastic{
		Host:  Config.Host,
		Index: Config.Index,

		client: client,
	}

	_, err = e.CreateIndex(Config.Index)
	if err != nil {
		return nil, fmt.Errorf("Problem with index: %v", err)
	}

	return e, nil
}

// CreateIndex will create an index
func (e *Elastic) CreateIndex(name string) (bool, error) {
	ok, err := e.IndexExists(name)
	if err != nil {
		return false, fmt.Errorf("Unable to create index: %v", err)
	}

	// index doesn't exist, return false
	if !ok {
		return false, nil
	}

	createIndex, err := e.client.CreateIndex(name).Do()
	if err != nil {
		return false, fmt.Errorf("Unable to create index: %v", err)
	}

	return createIndex.Acknowledged, nil
}

// IndexExists checks to see if an index already exists
func (e *Elastic) IndexExists(name string) (bool, error) {
	ok, err := e.client.IndexExists(name).Do()
	if err != nil {
		return false, fmt.Errorf("Couldn't check if index exists: %v", err)
	}

	return ok, nil
}

func (e *Elastic) Write(body string) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	if _, err := e.client.Index().Index(e.Index).Type("message").Id(id.String()).BodyString(body).Do(); err != nil {
		return err
	}

	return nil
}

// Init loads config and sets up elastic search
func Init() {
	conf, err := loadConfig()
	if err != nil {
		panic(err)
	}

	Config = conf
}

func loadConfig() (*esConfig, error) {
	configJSON := config.AtPath(
		"hailo", "service", "bakery", "elastic",
	).AsJson()

	log.Debugf("Elastic Config: %v", string(configJSON))

	var conf esConfig
	if err := json.Unmarshal(configJSON, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

type esConfig struct {
	Index string `json:"index"`
	Host  string `json:"host"`
}
