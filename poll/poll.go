package poll

import (
	"fmt"
	"reflect"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type Event uint16

const (
	POLLIN  = Event(0x0001)
	POLLPRI = Event(0x0002)
	POLLOUT = Event(0x0004)

	POLLRDNORM = Event(0x0040)
	POLLRDBAND = Event(0x0080)
	POLLWRNORM = Event(0x0100)
	POLLWRBAND = Event(0x0200)

	POLLMSG    = Event(0x0400)
	POLLREMOVE = Event(0x1000)
	POLLRDHUP  = Event(0x2000)

	POLLERR  = Event(0x0008)
	POLLHUP  = Event(0x0010)
	POLLNVAL = Event(0x0020)
)

func (e Event) String() string {
	flags := make([]string, 0, len(eventStringMap))
	for x, str := range eventStringMap {
		if e&x > 0 {
			flags = append(flags, str)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(flags, "|"))
}

var eventStringMap = map[Event]string{
	POLLIN:     "POLLIN",
	POLLPRI:    "POLLPRI",
	POLLOUT:    "POLLOUT",
	POLLRDNORM: "POLLRDNORM",
	POLLRDBAND: "POLLRDBAND",
	POLLWRNORM: "POLLWRNORM",
	POLLWRBAND: "POLLWRBAND",
	POLLMSG:    "POLLMSG",
	POLLREMOVE: "POLLREMOVE",
	POLLRDHUP:  "POLLRDHUP",
	POLLERR:    "POLLERR",
	POLLHUP:    "POLLHUP",
	POLLNVAL:   "POLLNVAL",
}

type PollFd struct {
	Fd      int32
	Events  Event
	REvents Event
}

func Poll(fdSet []PollFd, timeout time.Duration) (int, error) {
	nFds := uintptr(len(fdSet))
	duration := int(timeout / time.Millisecond)
	if timeout < 0 {
		duration = -1
	}
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&fdSet))
	n, _, e := syscall.Syscall6(syscall.SYS_POLL, hdr.Data, nFds, uintptr(duration), 0, 0, 0)
	if n < 0 {
		return 0, e
	}
	return int(n), nil
}

func SinglePoll(fd int, events Event, timeout time.Duration) (Event, error) {
	pfd := []PollFd{
		{
			Fd:      int32(fd),
			Events:  events,
			REvents: 0,
		},
	}
	_, err := Poll(pfd, timeout)
	return pfd[0].REvents, err
}

func WaitInput(fd int, timeout time.Duration) error {
	e, err := SinglePoll(fd, POLLIN, timeout)
	if err != nil {
		return err
	}
	if e == Event(0) {
		return ErrTimeout
	}
	if e&POLLNVAL > 0 {
		return ErrInvalidFd
	}
	if e&(POLLERR|POLLHUP|POLLRDHUP|POLLREMOVE) > 0 {
		return Error(e.String())
	}
	return nil
}
