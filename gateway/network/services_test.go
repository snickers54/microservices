package network

import (
	"strconv"
	"testing"
)

func TestServiceString(t *testing.T) {
	s := Service{
		Name: "test",
		IP:   "10.0.0.1",
		Port: "8000",
		Version: Version{
			Name:  "test",
			Major: 0,
			Minor: 0,
			Patch: 1,
		},
		Status: STATUS_ACTIVE,
	}
	expected := "test - test (0.0.1) | 10.0.0.1:8000"
	if s.String() != expected {
		t.Errorf(errorMsg, expected, s.String())
	}
}

func createService(name, ip string) *Service {
	return &Service{
		Name: name,
		IP:   ip,
		Port: "8000",
		Version: Version{
			Name:  "test",
			Major: 0,
			Minor: 0,
			Patch: 1,
		},
		Status: STATUS_ACTIVE,
	}
}

func TestServicesLen(t *testing.T) {
	ss := append(Services{}, createService("test1", "10.0.0.1"), createService("test2", "10.0.0.2"), createService("test3", "10.0.0.3"))
	if ss.Len() != 3 {
		t.Errorf(errorMsg, "3", "3")
	}
}

func TestServicesSwap(t *testing.T) {
	ss := append(Services{}, createService("test1", "10.0.0.1"), createService("test2", "10.0.0.2"), createService("test3", "10.0.0.3"))
	ss.Swap(0, 1)
	if ss[0].Name == "test1" || ss[1].Name == "test2" || ss[0].Name == ss[1].Name {
		t.Errorf(errorMsg, "test2|test1", ss[0].Name+"|"+ss[1].Name)
	}
}

func TestServicesLess(t *testing.T) {
	ss := append(Services{}, createService("test1", "10.0.0.1"), createService("test2", "10.0.0.2"), createService("test3", "10.0.0.3"))
	if ss.Less(0, 1) == false {
		t.Errorf(errorMsg, "true", "false")
	}
}

func TestServicesAdd(t *testing.T) {
	ss := Services{}
	for i := 1; i <= 10; i++ {
		ss.Add(*createService("test"+strconv.Itoa(i), "10.0.0."+strconv.Itoa(i)))
		if len(ss) != i {
			t.Errorf(errorMsg, strconv.Itoa(i), strconv.Itoa(len(ss)))
		}
	}
}

func TestServicesAddNoDuplicates(t *testing.T) {
	ss := Services{}
	for i := 1; i <= 10; i++ {
		ss.Add(*createService("test", "10.0.0.1"))
		if len(ss) != 1 {
			t.Errorf(errorMsg, "1", strconv.Itoa(len(ss)))
		}
	}
}
