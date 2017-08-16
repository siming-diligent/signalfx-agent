// Observers are what observe the environment to discover running services. An
// Observer only has to implement one method, Configure, which receives the
// configuration for that observer.  That Configure method might be called
// multiple times, so an observer should be prepared for its config to change
// and take appropriate steps.
//
// The ultimate output of an observer are ServiceInstances, which represent the
// discovered services.
package observers

import (
	"github.com/signalfx/neo-agent/core/config"
	"github.com/signalfx/neo-agent/utils"
	log "github.com/sirupsen/logrus"
)

type Shutdownable interface {
	Shutdown()
}

// ObserverFactory creates an unconfigured instance of an observer
type ObserverFactory func(*ServiceCallbacks) interface{}

var observerFactories = map[string]ObserverFactory{}

// These are blank (zero-value) instances of the configuration struct for a
// particular observer type.
var configTemplates = map[string]interface{}{}

// Register an observer of _type with the agent.  configTemplate should be a
// zero-valued struct that is of the same type that is accepted by the
// Configure method of the observer.
func Register(_type string, factory ObserverFactory, configTemplate interface{}) {
	if _, ok := observerFactories[_type]; ok {
		log.WithFields(log.Fields{
			"observerType": _type,
		}).Error("Observer type already registered")
		return
	}
	observerFactories[_type] = factory
	configTemplates[_type] = configTemplate
}

// ServiceCallbacks are a set of functions that the observers call upon
// detecting new services and disappearing services.
type ServiceCallbacks struct {
	Added   func(*ServiceInstance)
	Removed func(*ServiceInstance)
}

func configureObserver(observer interface{}, conf *config.ObserverConfig) bool {
	log.WithFields(log.Fields{
		"config": *conf,
	}).Debug("Configuring observer")

	finalConfig := utils.CloneInterface(configTemplates[conf.Type])

	config.FillInConfigTemplate("ObserverConfig", finalConfig, conf)
	if finalConfig == nil {
		return false
	}

	return config.CallConfigure(observer, finalConfig)
}
