package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// 向etcd中注册的Server信息
type Server struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Weight  int    `json:"weight"`
	Version string `json:"version"`
	Ttl     int64  `json:"ttl"`
}


// 创建注册key
func (s Server) BuildRegisterKey() string {
	// 如果服务没有版本 (Version)，则注册键为 "/服务名称/服务地址"。
	if len(s.Version) == 0 {
		return fmt.Sprintf("/%s/%s", s.Name, s.Addr)
	}
	// 如果服务有版本，则注册键为 "/服务名称/服务版本/服务地址"。
	return fmt.Sprintf("/%s/%s/%s", s.Name, s.Version, s.Addr)
}

func ParseValue(val []byte) (Server, error) {
	server := Server{}
	if err := json.Unmarshal(val, &server); err != nil {
		return server, err
	}
	return server, nil
}

func ParseKey(key string) (Server, error) {
	strs := strings.Split(key, "/")
	if len(strs) == 2 {
		// no version
		return Server{
			Name: strs[0],
			Addr: strs[1],
		}, nil
	}

	if len(strs) == 3 {
		// has version
		return Server{
			Name:    strs[0],
			Addr:    strs[1],
			Version: strs[2],
		}, nil
	}

	return Server{}, errors.New("invalid key")
}
