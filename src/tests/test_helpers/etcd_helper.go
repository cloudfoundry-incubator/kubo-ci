package test_helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	. "github.com/onsi/gomega"
)

func GetEtcdLeaderID(etcdIP string) string {
	var statsResp struct {
		LeaderInfo struct {
			Leader string `json:"leader"`
		} `json:"leaderInfo"`
	}

	statsEndpoint := fmt.Sprintf("http://%s:4001/v2/stats/self", etcdIP)
	resp, err := http.Get(statsEndpoint)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	json.Unmarshal(body, &statsResp)
	Expect(statsResp.LeaderInfo.Leader).NotTo(BeEmpty())
	return statsResp.LeaderInfo.Leader
}

func GetEtcdLeaderClientURL(etcdIP string) (string, error) {
	var etcdMembersResp struct {
		Members []struct {
			ID         string   `json:"id"`
			ClientURLs []string `json:"clientURLs"`
		} `json:"members"`
	}
	membersEndpoint := fmt.Sprintf("http://%s:4001/v2/members", etcdIP)
	etcdLeaderID := GetEtcdLeaderID(etcdIP)

	resp, err := http.Get(membersEndpoint)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())

	json.Unmarshal(body, &etcdMembersResp)
	for _, member := range etcdMembersResp.Members {
		if member.ID == etcdLeaderID {
			return member.ClientURLs[0], nil
		}
	}
	return "", errors.New("cannot find etcd leader ip")
}

func PutKeyToEtcd(etcdIP string, key string, value string) {
	etcdLeaderClientUrl, err := GetEtcdLeaderClientURL(etcdIP)
	Expect(err).NotTo(HaveOccurred())

	keysEndpoint := fmt.Sprintf("%s/v2/keys/%s", etcdLeaderClientUrl, key)

	httpClient := &http.Client{}
	data := url.Values{}
	data.Set("value", value)
	req, err := http.NewRequest("PUT", keysEndpoint, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	Expect(err).NotTo(HaveOccurred())
	Expect(strconv.Itoa(resp.StatusCode)).Should(MatchRegexp(`20[01]`))
}

func GetKeyFromEtcd(etcdIP string, key string) string {
	var etcdKeysResp struct {
		Node struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"node"`
	}

	etcdLeaderClientUrl, err := GetEtcdLeaderClientURL(etcdIP)
	Expect(err).NotTo(HaveOccurred())

	keysEndpoint := fmt.Sprintf("%s/v2/keys/%s", etcdLeaderClientUrl, key)
	resp, err := http.Get(keysEndpoint)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &etcdKeysResp)
	return etcdKeysResp.Node.Value
}

func DeleteKeyFromEtcd(etcdIP string, key string) {
	etcdLeaderClientUrl, err := GetEtcdLeaderClientURL(etcdIP)
	Expect(err).NotTo(HaveOccurred())

	keysEndpoint := fmt.Sprintf("%s/v2/keys/%s", etcdLeaderClientUrl, key)
	httpClient := &http.Client{}
	req, err := http.NewRequest("DELETE", keysEndpoint, nil)
	resp, err := httpClient.Do(req)
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(200))
}
