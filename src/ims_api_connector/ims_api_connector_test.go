package ims_api_connector

import (
	"testing"
	"time"
	"net/http"
	"reflect"
)

func assertStringsEqual(t *testing.T, was, expected string) {
	if was != expected {
		t.Errorf("'%v' != '%v'", was, expected)
	}
}

func assertDurationsEqual(t *testing.T, was, expected time.Duration) {
	if was != expected {
		t.Errorf("%v != %v", was, expected)
	}
}

func assertErrorIsNil(t *testing.T, was error) {
	if was != nil {
		t.Errorf("%v != nil", was)
	}
}

func assertBoolEqual(t *testing.T, was, expected bool) {
	if was != expected {
		t.Errorf("%v != %v", expected, was)
	}
}

func assertDeepEqual(t *testing.T, was, expected interface{}) {
	if !reflect.DeepEqual(was, expected) {
		t.Errorf("%v != %v", was, expected)
	}
}

func setUp() *IMSAPIConnector {
	return New("some_username", "some_password", "192.168.137.253:8000", 5)
}

func Test_New(t *testing.T) {
	subject := setUp()

	assertStringsEqual(t, subject.Username, "some_username")

	assertStringsEqual(t, subject.Password, "some_password")

	assertStringsEqual(t, subject.BaseURL, "http://192.168.137.253:8000/api/")

	assertDurationsEqual(t, subject.client.Timeout, time.Second*time.Duration(5))
}

func Test_prepareResource(t *testing.T) {
	subject := setUp()
	was := subject.buildResource("auth/login")
	expected := "http://192.168.137.253:8000/api/auth/login/"

	assertStringsEqual(t, was, expected)
}

func Test_buildRequest(t *testing.T) {
	subject := setUp()

	req1, err1 := subject.buildRequest(http.MethodGet, "some url", "some body")
	assertStringsEqual(t, req1.URL.String(), "some%20url")
	assertErrorIsNil(t, err1)
	assertStringsEqual(t, req1.Header.Get("Content-Type"), "application/x-www-form-urlencoded")

	subject.Authenticated = true
	subject.Key = "some_key"

	req2, err2 := subject.buildRequest(http.MethodGet, "some url", "some body")
	assertStringsEqual(t, req2.URL.String(), "some%20url")
	assertErrorIsNil(t, err2)
	assertStringsEqual(t, req2.Header.Get("Content-Type"), "application/json")
	assertStringsEqual(t, req2.Header.Get("Authorization"), "Token some_key")
}

func Test_buildUsernameAndPasswordParams(t *testing.T) {
	subject := setUp()

	body := subject.buildUsernameAndPasswordParams()

	assertStringsEqual(t, string(body), "password=some_password&username=some_username")
}

func Test_handleAuthenticationResponse(t *testing.T) {
	subject := setUp()
	body := []byte("{\"key\": \"1c1a552f5b013bd76b7d6acd731a8e46955f4b13\"}")

	authenticated := subject.handleAuthenticationResponse(body)
	assertBoolEqual(t, authenticated, true)

	assertBoolEqual(t, subject.Authenticated, true)

	assertStringsEqual(t, subject.Key, "1c1a552f5b013bd76b7d6acd731a8e46955f4b13")
}

func Test_handleGetAssetsResponse(t *testing.T) {
	subject := setUp()
	var body_str string
	body_str = "[{\"id\":1,\"name\":\"Asset 1\",\"is_deleted\":false,\"last_updated\":\"1991-02-06T00:00:00.000000+00:00\",\"note\":null,\"json_data\":null,\"type_id\":3,\"primary_ip_device_id\":5,\"site_id\":1,\"tags\":[7]},"
	body_str += "{\"id\":2,\"name\":\"Asset 2\",\"is_deleted\":false,\"last_updated\":\"1991-02-06T00:00:00.000000+00:00\",\"note\":null,\"json_data\":null,\"type_id\":4,\"primary_ip_device_id\":6,\"site_id\":1,\"tags\":[8]}]"
	var body []byte = []byte(body_str)

	was, _ := subject.handleGetAssetsResponse(body)

	expected := make([]Asset, 2)

	lastUpdated := time.Date(1991, 2, 6, 0, 0, 0, 0, time.UTC)

	tags1 := make([]int, 1)
	tags1[0] = 8
	expected[1] = Asset{ID: 2, Name: "Asset 2", IsDeleted: false, LastUpdated: lastUpdated, Note: "", JSONData: nil, TypeID: 4, PrimaryIPDeviceID: 6, SiteID: 1, Tags: tags1}

	tags2 := make([]int, 1)
	tags2[0] = 7
	expected[0] = Asset{ID: 1, Name: "Asset 1", IsDeleted: false, LastUpdated: lastUpdated, Note: "", JSONData: nil, TypeID: 3, PrimaryIPDeviceID: 5, SiteID: 1, Tags: tags2}

	for i := 0; i < 2; i++ {
		if expected[i].ID != was[i].ID {
			t.Errorf("ID not equal")
		}
		if expected[i].Name != was[i].Name {
			t.Errorf("Name not equal")
		}
		if expected[i].IsDeleted != was[i].IsDeleted {
			t.Errorf("IsDeleted not equal")
		}
		if expected[i].LastUpdated.Unix() != was[i].LastUpdated.Unix() {
			t.Log(expected[i].LastUpdated, was[i].LastUpdated)
			t.Errorf("LastUpdated not equal")
		}
		if expected[i].Note != was[i].Note {
			t.Errorf("Note not equal")
		}
		if expected[i].JSONData != was[i].JSONData {
			t.Errorf("JSONData not equal")
		}
		if expected[i].TypeID != was[i].TypeID {
			t.Errorf("TypeID not equal")
		}
		if expected[i].PrimaryIPDeviceID != was[i].PrimaryIPDeviceID {
			t.Errorf("PrimaryIPDeviceID not equal")
		}
		if expected[i].SiteID != was[i].SiteID {
			t.Errorf("SiteID not equal")
		}
		if !reflect.DeepEqual(expected[i].Tags, expected[i].Tags) {
			t.Errorf("Tags not equal")
		}
	}

}
