package packer

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"

	log "github.com/cihub/seelog"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template"
)

func Build(t io.ReadCloser) {
	config := &config{
		PluginMinPort: 10000,
		PluginMaxPort: 25000,
	}

	if err := config.Discover(); err != nil {
		log.Errorf("%v", err)
	}

	defer t.Close()

	tplBuf := bytes.NewBuffer([]byte{})
	tplBuf.ReadFrom(t)
	tpl, err := template.Parse(tplBuf)
	if err != nil {
		log.Errorf("%v", err)
	}

	log.Debugf("Template: %v", string(tpl.RawContents))

	coreConfig := packer.CoreConfig{
		Components: packer.ComponentFinder{
			Builder:       config.LoadBuilder,
			Hook:          config.LoadHook,
			PostProcessor: config.LoadPostProcessor,
			Provisioner:   config.LoadProvisioner,
		},
		Template: tpl,
		Variables: map[string]string{
			"AWS_ACCESS_KEY_ID": "AKI123",
			"AWS_SECRET_KEY":    "123",
		},
	}

	core, err := packer.NewCore(&coreConfig)
	if err != nil {
		log.Errorf("Unable to create core: %v", err)
	}

	var builds []packer.Build
	for _, name := range core.BuildNames() {
		log.Debugf("Found build: %v", name)

		b, err := core.Build(name)
		if err != nil {
			log.Errorf("Unable to create build: %v", err)
		}

		builds = append(builds, b)
	}

	for _, b := range builds {
		log.Infof("Preparing build: %s", b.Name())
		// b.SetDebug(cfgDebug)
		// b.SetForce(cfgForce)

		warnings, err := b.Prepare()
		if err != nil {
			log.Errorf("%v", err)
		}
		if len(warnings) > 0 {
			for _, warning := range warnings {
				log.Infof("Warning: %s", warning)
			}
		}
	}

	cacheDir := os.Getenv("PACKER_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = "packer_cache"
	}

	cacheDir, err = filepath.Abs(cacheDir)
	if err != nil {
		log.Errorf("Error preparing cache directory: %s", err)
	}

	log.Infof("Setting cache directory: %s", cacheDir)
	cache := &packer.FileCache{CacheDir: cacheDir}

	artifacts := map[string][]packer.Artifact{}
	errors := make(map[string]error)

	var wg sync.WaitGroup
	for _, b := range builds {
		log.Infof("Starting to process: %s", b.Name())
		wg.Add(1)

		go func(b packer.Build) {
			defer wg.Done()

			name := b.Name()
			log.Infof("Starting build of: %s", name)

			runArtifacts, err := b.Run(&Ui{}, cache)
			if err != nil {
				log.Errorf("Build '%s' errored: %s", name, err)
				errors[name] = err
			} else {
				log.Infof("Build '%s' finished.", name)
				artifacts[name] = runArtifacts
			}
		}(b)
	}

	log.Infof("Waiting on builds to complete...")
	wg.Wait()

	if len(errors) > 0 {
		log.Error("There were some problems building")
		for n, e := range errors {
			log.Infof("%s: %v", n, e)
		}
	}

	if len(artifacts) > 0 {
		log.Infof("Builds finished, there are some artifacts")

		for _, ba := range artifacts {
			for _, a := range ba {
				log.Infof("%#v", a)
			}
		}
	}
}
