package proxy

import (
	"github.com/devopsfaith/flatmap/tree"
	"melody/config"
	"strings"
)

const flatmapKey = "flatmap_filter"

type EntityFormatter interface {
	Format(Response) Response
}

type propertyFilter func(*Response)

type EntityFormatterFunc func(Response) Response

func (e EntityFormatterFunc) Format(entity Response) Response { return e(entity) }

type entityFormatter struct {
	Target         string
	Prefix         string
	PropertyFilter propertyFilter
	Mapping        map[string]string
}

type flatmapFormatter struct {
	Target string
	Prefix string
	Ops    []flatmapOp
}

type flatmapOp struct {
	Type string
	Args [][]string
}

func (e entityFormatter) Format(resp Response) Response {
	if e.Target != "" {
		extractTarget(e.Target, &resp)
	}
	if len(resp.Data) > 0 {
		e.PropertyFilter(&resp)
	}
	if len(resp.Data) > 0 {
		for formerKey, newKey := range e.Mapping {
			if v, ok := resp.Data[formerKey]; ok {
				resp.Data[newKey] = v
				delete(resp.Data, formerKey)
			}
		}
	}
	if e.Prefix != "" {
		resp.Data = map[string]interface{}{e.Prefix: resp.Data}
	}
	return resp
}

func NewEntityFormatter(remote *config.Backend) EntityFormatter {
	/**
	1. 返回结果为key : 数组
	{
		"a": []
	}
	 */
	// 创建flatMapFormatter对象
	if entityFormatter := newFlatmapFormatter(remote); entityFormatter != nil {
		return entityFormatter
	}
	/**
	2. 返回结果为key : 对象
	{
		"a" : {...}
	}
	 */
	var propertyFilter propertyFilter
	if len(remote.Whitelist) > 0 {
		propertyFilter = newWhitelistingFilter(remote.Whitelist)
	} else {
		propertyFilter = newBlacklistingFilter(remote.Blacklist)
	}

	mappings := make(map[string]string, len(remote.Mapping))
	for k, v := range remote.Mapping {
		mappings[k] = v
	}

	return entityFormatter{
		Target:         remote.Target,
		Prefix:         remote.Group,
		PropertyFilter: propertyFilter,
		Mapping:        mappings,
	}
}

func newFlatmapFormatter(remote *config.Backend) EntityFormatter {
	if v, ok := remote.ExtraConfig[Namespace]; ok {
		if e, ok := v.(map[string]interface{}); ok {
			if a, ok := e[flatmapKey].([]interface{}); ok {
				if len(a) == 0 {
					return nil
				}
				options := []flatmapOp{}
				for _, o := range a {
					m, ok := o.(map[string]interface{})
					if !ok {
						continue
					}
					op := flatmapOp{}
					if t, ok := m["type"].(string); ok {
						op.Type = t
					} else {
						continue
					}
					if args, ok := m["args"].([]interface{}); ok {
						op.Args = make([][]string, len(args))
						for k, arg := range args {
							if t, ok := arg.(string); ok {
								op.Args[k] = strings.Split(t, ".")
							}
						}
					}
					options = append(options, op)
				}
				if len(options) == 0 {
					return nil
				}
				return &flatmapFormatter{
					Target: remote.Target,
					Prefix: remote.Group,
					Ops:    options,
				}
			}
		}
	}
	return nil
}

func newWhitelistingFilter(whitelist []string) propertyFilter {
	wlDict := make(map[string]interface{})
	for _, k := range whitelist {
		wlFields := strings.Split(k, ".")
		d := buildDictPath(wlDict, wlFields[:len(wlFields)-1])
		d[wlFields[len(wlFields)-1]] = true
	}
	
	return func(response *Response) {
		if whitelistPrune(wlDict, response.Data) {
			for k := range response.Data {
				delete(response.Data, k)
			}
		}
	}
}

func buildDictPath(accumulator map[string]interface{}, fields []string) map[string]interface{} {
	ok := true
	var c map[string]interface{}
	var fIdx int
	fEnd := len(fields)
	p := accumulator
	for fIdx = 0; fIdx < fEnd; fIdx++ {
		if c, ok = p[fields[fIdx]].(map[string]interface{}); !ok {
			break
		}
		p = c
	}
	for ; fIdx < fEnd; fIdx++ {
		c = make(map[string]interface{})
		p[fields[fIdx]] = c
		p = c
	}
	return p
}

func whitelistPrune(wlDict map[string]interface{}, inDict map[string]interface{}) bool {
	canDelete := true
	var deleteSibling bool
	for k, v := range inDict {
		deleteSibling = true
		if subWl, ok := wlDict[k]; ok {
			if subWlDict, okk := subWl.(map[string]interface{}); okk {
				if subInDict, isDict := v.(map[string]interface{}); isDict && !whitelistPrune(subWlDict, subInDict) {
					deleteSibling = false
				}
			} else {
				// whitelist leaf, maintain this branch
				deleteSibling = false
			}
		}
		if deleteSibling {
			delete(inDict, k)
		} else {
			canDelete = false
		}
	}
	return canDelete
}

func newBlacklistingFilter(blacklist []string) propertyFilter {
	bl := make(map[string][]string, len(blacklist))
	for _, key := range blacklist {
		keys := strings.Split(key, ".")
		if len(keys) > 1 {
			if sub, ok := bl[keys[0]]; ok {
				bl[keys[0]] = append(sub, keys[1])
			} else {
				bl[keys[0]] = []string{keys[1]}
			}
		} else {
			bl[keys[0]] = []string{}
		}
	}

	return func(entity *Response) {
		for k, sub := range bl {
			if len(sub) == 0 {
				delete(entity.Data, k)
			} else {
				if tmp := blacklistFilterSub(entity.Data[k], sub); len(tmp) > 0 {
					entity.Data[k] = tmp
				}
			}
		}
	}
}

func blacklistFilterSub(v interface{}, blacklist []string) map[string]interface{} {
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	for _, key := range blacklist {
		delete(tmp, key)
	}
	return tmp
}

func (f flatmapFormatter) Format(response Response) Response {
	if f.Target != "" {
		extractTarget(f.Target, &response)
	}

	f.processOps(&response)

	if f.Prefix != "" {
		response.Data = map[string]interface{}{f.Prefix: response.Data}
	}
	return response
}

func extractTarget(target string, entity *Response) {
	if v, ok := entity.Data[target]; ok {
		if t, ok := v.(map[string]interface{}); ok {
			entity.Data = t
		} else {
			entity.Data = map[string]interface{}{}
		}
	} else {
		entity.Data = map[string]interface{}{}
	}
}

func (f flatmapFormatter) processOps(entity *Response) {
	tree, err := tree.New(entity.Data)
	if err != nil {
		return
	}

	for _, v := range f.Ops {
		switch v.Type {
		case "move":
			tree.Move(v.Args[0], v.Args[1])
		case "del":
			tree.Del(v.Args[0])
		default:
		}
	}

	entity.Data, _ = tree.Get([]string{}).(map[string]interface{})
}
