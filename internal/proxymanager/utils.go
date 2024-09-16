package proxymanager

import (
	"fmt"
	"math/rand"

	"github.com/fsnotify/fsnotify"
	"github.com/kitabisa/mubeng/common/errors"
	"github.com/kitabisa/mubeng/pkg/helper"
)

// NextProxy will navigate the next proxy to use
func (p *ProxyManager) NextProxy() (string, error) {
	var proxy string

	if p.Length <= 0 {
		return proxy, errors.ErrNoProxyLeft
	}

	p.CurrentIndex++
	if p.CurrentIndex > p.Length-1 {
		p.CurrentIndex = 0
	}

	return p.Proxies[p.CurrentIndex], nil
}

// RandomProxy will choose a proxy randomly from the list
func (p *ProxyManager) RandomProxy() (string, error) {
	var proxy string

	if p.Length <= 0 {
		return proxy, errors.ErrNoProxyLeft
	}

	return p.Proxies[rand.Intn(p.Length)], nil
}

// RemoveProxy removes target proxy from proxy pool
func (p *ProxyManager) RemoveProxy(target string) error {
	for i, v := range p.Proxies {
		if v == target {
			p.Proxies = append(p.Proxies[:i], p.Proxies[i+1:]...)
			p.Length -= 1

			return nil
		}
	}

	return fmt.Errorf("unable to find %q in the proxy pool", target)
}

// Rotate proxy based on method
//
// Valid methods are "sequent" and "random", default return empty string.
func (p *ProxyManager) Rotate(method string) (string, error) {
	var proxy string
	var err error

	switch method {
	case "sequent":
		proxy, err = p.NextProxy()
	case "random":
		proxy, err = p.RandomProxy()
	}

	if proxy != "" {
		proxy = helper.EvalFunc(proxy)
	}

	return proxy, err
}

// Watch proxy file from events
func (p *ProxyManager) Watch() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return watcher, err
	}

	if err := watcher.Add(p.filepath); err != nil {
		return watcher, err
	}

	return watcher, nil
}

// Reload proxy pool
func (p *ProxyManager) Reload() error {
	i := p.CurrentIndex

	p, err := New(p.filepath)
	if err != nil {
		return err
	}
	p.CurrentIndex = i

	return nil
}
