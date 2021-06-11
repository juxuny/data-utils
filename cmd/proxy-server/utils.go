package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
)

func checkTcp(remote string) error {
	conn, err := net.Dial("tcp", remote)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	return nil
}

func checkHttp(u string) error {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("wrong http status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	_, _ = ioutil.ReadAll(resp.Body)
	return nil
}
