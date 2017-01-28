// Package mimic contains a mockup tool with middleware.
//
// The MIT License (MIT)
//
// Copyright (c) 2016 Angel Del Castillo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package mimic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

// Load loads yml file with parameters for mocking.
func Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Load : close config file : err [%s]", err)
		}
	}()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	cache = &Cache{
		List: make(map[string]*Endpoint),
	}

	// parse file

	err = yaml.Unmarshal(b, &cache.Config)
	if err != nil {
		return err
	}

	// set cache list from config file

	err = cache.Prepare()
	if err != nil {
		return err
	}

	for k := range cache.List {
		log.Printf("Load : list k [%s] [%#v]", k, cache.List[k])
	}

	return nil
}

// Config yaml.
type Config struct {
	Endpoints []*Endpoint `yaml:"endpoints"`
}

// Endpoint config.
type Endpoint struct {
	Method   string            `yaml:"method"`
	Headers  map[string]string `yaml:"headers"`
	URI      string            `yaml:"uri"`
	Params   string            `yaml:"params"`
	Response *Response         `yaml:"response"`
}

// Response mock struct.
type Response struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

var (
	catcher Catcher = &JSONCatcher{}
)

// SetCatcher sets the catcher implementation for middleware interference.
// You can implement your custom Catcher, by default JSONCatcher is set.
// In order to set your own catcher this function must be called at init time.
func SetCatcher(c Catcher) error {
	// set default Catcher interface.
	catcher = c

	// load keys into cache list
	return cache.Prepare()
}

// Handler middleware.
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// try catch a mockup query.

		v, ok := catcher.Catch(r)
		if ok {
			// write response

			if err := catcher.Render(w, v); err != nil {
				log.Printf("Mimic : render : err [%s]", err)
			}
			return
		}

		h.ServeHTTP(w, r)
	})
}

var (
	cache *Cache
)

// Cache contains mock data.
type Cache struct {
	Config *Config
	List   map[string]*Endpoint
	sync.RWMutex
}

// Prepare will load the keys that default catcher returns from config file and store as cache map.
func (c *Cache) Prepare() error {
	for i := range cache.Config.Endpoints {
		x := cache.Config.Endpoints[i]

		s, err := catcher.Key(x)
		if err != nil {
			return err
		}

		cache.List[s] = x
	}
	return nil
}

// Get returns cached mockup data.
func (c *Cache) Get(s string) (*Endpoint, error) {
	c.RLock()
	defer c.RUnlock()

	e, ok := c.List[s]
	if ok {
		return e, nil
	}

	return nil, errors.New("mock not found")
}

// Get return Endpoint data from cache if exists.
func Get(s string) (*Endpoint, error) {
	return cache.Get(s)
}

// Catcher interface allows custom request handling. View JSONCatcher and
// XMLCatcher for a sample implementation.
type Catcher interface {
	// Key method returns the key for the cache list.
	Key(interface{}) (string, error)

	// Catch method call Key and if exists returns Endpoint data.
	Catch(*http.Request) (interface{}, bool)

	// Render method allow custom http response when catch handles something.
	Render(http.ResponseWriter, interface{}) error
}

// JSONCatcher implements Catcher.
type JSONCatcher struct{}

// Key implements Catcher.
func (j *JSONCatcher) Key(v interface{}) (string, error) {
	switch e := v.(type) {
	case *Endpoint:
		// means cache is on prepare state and needs the slug from config file.
		s := fmt.Sprintf("%v %v %v", e.Method, e.URI, e.Params)
		log.Printf("Key : endpoint set [%s]", s)
		return s, nil
	case *http.Request:
		// means catcher is finding this key in cache map.
		s := fmt.Sprintf("%v %v %v", e.Method, e.RequestURI, e.Form.Encode())
		log.Printf("Key : http request get [%s]", s)
		_, err := Get(s)
		if err == nil {
			return s, nil
		}

		return "", errors.New("not found")
	default:
		return "", fmt.Errorf("type not found: %v", e)
	}
}

// Catch implements Catcher.
func (j *JSONCatcher) Catch(r *http.Request) (interface{}, bool) {
	if err := r.ParseForm(); err != nil {
		return nil, false
	}

	key, err := j.Key(r)
	if err != nil {
		return nil, false
	}

	s, err := Get(key)
	if err != nil {
		return nil, false
	}
	return s, true
}

// Render implements Catcher.
func (j *JSONCatcher) Render(w http.ResponseWriter, v interface{}) error {
	e, ok := v.(*Endpoint)
	if !ok {
		return errors.New("type not supported")
	}

	res := e.Response
	if res == nil {
		return errors.New("cache response body nil")
	}

	// set status.
	w.WriteHeader(res.Status)

	// set custom headers.

	for k, val := range res.Headers {
		log.Printf("Render : set header k v [%v][%v]", k, val)
		w.Header().Set(k, val)
	}
	w.Header().Set("AnotherHeader", "empty")

	if _, err := fmt.Fprint(w, res.Body); err != nil {
		return err
	}

	return nil
}

// XMLCatcher implements Catcher.
// TODO;
type XMLCatcher struct{}

// Catch implements Catcher.
func (j *XMLCatcher) Catch(r *http.Request) (interface{}, bool) {
	return nil, false
}

// Render implements Catcher.
func (j *XMLCatcher) Render(w http.ResponseWriter, v interface{}) error {
	return nil
}
