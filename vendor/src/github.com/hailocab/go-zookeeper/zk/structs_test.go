package zk

import (
	"reflect"
	"testing"
)

func TestServerList(t *testing.T) {
	testAddrs := [][]string{
		[]string{"1.2.3.4"},
		[]string{"1.2.3.4", "4.3.2.1"},
	}

	for _, addrs := range testAddrs {
		sl := &serverList{addrs: addrs}

		if !sl.contains(addrs[0]) {
			t.Errorf("Expected serverList to contain %s", addrs[0])
		}
		if v := sl.changed([]string{"9.9.9.9"}); !v {
			t.Errorf("Expected serverList changed true, got: %t", v)
		}
		if v := sl.changed(addrs); v {
			t.Errorf("Expected serverList changed false, got: %t", v)
		}
		// dupe test
		if v := sl.changed([]string{"1.2.3.4", "1.2.3.4"}); !v {
			t.Errorf("Expected serverList changed true, got: %t", v)
		}

		var next int
		if len(addrs) > 1 {
			next = 1
		}

		if addr := sl.next(); addr != addrs[next] {
			t.Errorf("Expected next addr %s, got %s", addrs[next], addr)
		}

		if v := sl.hasNext(); v {
			t.Errorf("Expected hasNext false, got %t", v)
		}
	}
}

func TestEncodeDecodePacket(t *testing.T) {
	encodeDecodeTest(t, &requestHeader{-2, 5})
	encodeDecodeTest(t, &connectResponse{1, 2, 3, nil})
	encodeDecodeTest(t, &connectResponse{1, 2, 3, []byte{4, 5, 6}})
	encodeDecodeTest(t, &getAclResponse{[]ACL{{12, "s", "anyone"}}, Stat{}})
	encodeDecodeTest(t, &getChildrenResponse{[]string{"foo", "bar"}})
	encodeDecodeTest(t, &pathWatchRequest{"path", true})
	encodeDecodeTest(t, &pathWatchRequest{"path", false})
	encodeDecodeTest(t, &CheckVersionRequest{"/", -1})
	encodeDecodeTest(t, &multiRequest{Ops: []multiRequestOp{{multiHeader{opCheck, false, -1}, &CheckVersionRequest{"/", -1}}}})
}

func encodeDecodeTest(t *testing.T, r interface{}) {
	buf := make([]byte, 1024)
	n, err := encodePacket(buf, r)
	if err != nil {
		t.Errorf("encodePacket returned non-nil error %+v\n", err)
		return
	}
	t.Logf("%+v %x", r, buf[:n])
	r2 := reflect.New(reflect.ValueOf(r).Elem().Type()).Interface()
	n2, err := decodePacket(buf[:n], r2)
	if err != nil {
		t.Errorf("decodePacket returned non-nil error %+v\n", err)
		return
	}
	if n != n2 {
		t.Errorf("sizes don't match: %d != %d", n, n2)
		return
	}
	if !reflect.DeepEqual(r, r2) {
		t.Errorf("results don't match: %+v != %+v", r, r2)
		return
	}
}

func TestEncodeShortBuffer(t *testing.T) {
	buf := make([]byte, 0)
	_, err := encodePacket(buf, &requestHeader{1, 2})
	if err != ErrShortBuffer {
		t.Errorf("encodePacket should return ErrShortBuffer on a short buffer instead of '%+v'", err)
		return
	}
}

func TestDecodeShortBuffer(t *testing.T) {
	buf := make([]byte, 0)
	_, err := decodePacket(buf, &responseHeader{})
	if err != ErrShortBuffer {
		t.Errorf("decodePacket should return ErrShortBuffer on a short buffer instead of '%+v'", err)
		return
	}
}
