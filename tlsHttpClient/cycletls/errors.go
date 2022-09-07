package cycletls

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type errorMessage struct {
	StatusCode int
	debugger   string
	ErrorMsg   string
	Op         string
}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}

// func createErrorString(err: string) (msg, debugger string) {
func createErrorString(err error) (msg, debugger string) {
	msg = fmt.Sprintf("Request returned a Syscall Error: %s", err)
	debugger = fmt.Sprintf("%#v\n", err)
	return
}

func createErrorMessage(StatusCode int, err error, op string) errorMessage {
	msg := fmt.Sprintf("Request returned a Syscall Error: %s", err)
	debugger := fmt.Sprintf("%#v\n", err)
	return errorMessage{StatusCode: StatusCode, debugger: debugger, ErrorMsg: msg, Op: op}
}

func parseError(err error) (errormessage errorMessage) {
	var op string

	httpError := err.Error()
	status := lastString(strings.Split(httpError, "StatusCode:"))
	StatusCode, _ := strconv.Atoi(status)
	if StatusCode != 0 {
		msg, debugger := createErrorString(err)
		return errorMessage{StatusCode: StatusCode, debugger: debugger, ErrorMsg: msg}
	}
	if uErr, ok := err.(*url.Error); ok {
		if noErr, ok := uErr.Err.(*net.OpError); ok {
			op = noErr.Op
			if SyscallError, ok := noErr.Err.(*os.SyscallError); ok {
				if noErr.Timeout() {
					return createErrorMessage(408, SyscallError, op)
				}
				return createErrorMessage(401, SyscallError, op)
			} else if AddrError, ok := noErr.Err.(*net.AddrError); ok {
				return createErrorMessage(405, AddrError, op)
			} else if DNSError, ok := noErr.Err.(*net.DNSError); ok {
				return createErrorMessage(421, DNSError, op)
			} else {
				return createErrorMessage(421, noErr, op)
			}
		}
		if uErr.Timeout() {
			return createErrorMessage(408, uErr, op)
		}
	}
	return
}

type errExtensionNotExist struct {
	Context string
}

func (w *errExtensionNotExist) Error() string {
	return fmt.Sprintf("Extension {{ %s }} is not Supported by CycleTLS please raise an issue", w.Context)
}

func raiseExtensionError(info string) *errExtensionNotExist {
	return &errExtensionNotExist{
		Context: info,
	}
}
