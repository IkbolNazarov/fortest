package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type StructValidator interface {
	ValidateStruct(interface{}) error
	Engine() interface{}
}

type Context struct {
	Writer    http.ResponseWriter
	Request   *http.Request
	Validator StructValidator
}
type HandleFunc func(*Context)

var decoder = schema.NewDecoder()

func (f HandleFunc) ServeHTTP(c *Context) {
	f(c)
}

func (c *Context) Param(key string) string {
	return mux.Vars(c.Request)[key]
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) QuerySet(key string, value string) {
	c.Request.URL.Query().Set(key, value)
}

func (c *Context) ShouldBindParams(i interface{}) error {
	if err := c.ParseRequestParams(&i, c.Request); err != nil {
		return err
	}

	if c.Validator != nil {
		if err := c.Validator.ValidateStruct(i); err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) ShouldBind(i interface{}) error {
	if err := c.ParseRequest(&i, c.Request); err != nil {
		return err
	}

	if c.Validator != nil {
		if err := c.Validator.ValidateStruct(i); err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) ShouldBindQuery(i interface{}) error {
	if err := c.ParseRequestQuery(&i, c.Request); err != nil {
		return err
	}

	if c.Validator != nil {
		if err := c.Validator.ValidateStruct(i); err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) ShouldBindForm(i interface{}) error {
	if err := c.ParseRequestForm(&i, c.Request); err != nil {
		return err
	}

	if c.Validator != nil {
		if err := c.Validator.ValidateStruct(i); err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) ParseRequest(f interface{}, r *http.Request) error {

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(f); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (c *Context) ParseRequestParams(f interface{}, r *http.Request) error {
	params := mux.Vars(r)

	jsonStr, err := json.Marshal(params)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if err = json.Unmarshal(jsonStr, &f); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (c *Context) ParseRequestQuery(f interface{}, r *http.Request) error {
	params := make(map[string]interface{})
	queryParams := r.URL.Query()
	for i, v := range queryParams {
		params[i] = strings.Join(v, "")
	}
	jsonStr, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
		return err
	}
	// Convert json string to struct
	if err := json.Unmarshal(jsonStr, f); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Context) ParseRequestForm(f interface{}, r *http.Request) error {
	params := make(map[string]string)
	queryParams := r.Form
	for i, v := range queryParams {
		params[i] = strings.Join(v, "")
	}
	jsonStr, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := json.Unmarshal(jsonStr, f); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Context) Redirect(code int, location string) {
	if (code < http.StatusMultipleChoices || code > http.StatusPermanentRedirect) && code != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", code))
	}
	http.Redirect(c.Writer, c.Request, location, code)
}

func (c *Context) JSON(data interface{}) error {
	c.Writer.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(c.Writer).Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) String(data interface{}) error {
	c.Writer.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprintln(c.Writer, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) Respond(data interface{}, contentType string, writer func(data interface{})) error {
	c.Writer.Header().Set("Content-Type", contentType)
	writer(data)

	return nil
}

func (c *Context) PaginationMeta(total int64, offset int, limit int) map[string]interface{} {
	return map[string]interface{}{
		"total":  total,
		"offset": offset,
		"limit":  limit,
	}
}

func (c *Context) ErrorsMeta(errors interface{}) map[string]interface{} {
	return map[string]interface{}{
		"errors": errors,
	}
}

func (c *Context) OK(data interface{}) {
	OK(c.Writer, data)
}

func (c *Context) OKMeta(meta map[string]interface{}, data interface{}) {
	OKMeta(c.Writer, meta, data)
}

func (c *Context) BadRequest(data interface{}) {
	BadRequest(c.Writer, data)
}

func (c *Context) BadRequestMeta(meta map[string]interface{}, data interface{}) {
	BadRequestMeta(c.Writer, meta, data)
}

func (c *Context) Unauthorized(data interface{}) {
	Unauthorized(c.Writer, data)
}

func (c *Context) Forbidden(data interface{}) {
	Forbidden(c.Writer, data)
}

func (c *Context) NotFound(data interface{}) {
	NotFound(c.Writer, data)
}

func (c *Context) Internal(data interface{}) {
	Internal(c.Writer, data)
}
