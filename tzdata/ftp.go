package tzdata

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// FTP status codes, defined in RFC 959
const (
	StatusInitiating    = 100
	StatusRestartMarker = 110
	StatusReadyMinute   = 120
	StatusAlreadyOpen   = 125
	StatusAboutToSend   = 150

	StatusCommandOK             = 200
	StatusCommandNotImplemented = 202
	StatusSystem                = 211
	StatusDirectory             = 212
	StatusFile                  = 213
	StatusHelp                  = 214
	StatusName                  = 215
	StatusReady                 = 220
	StatusClosing               = 221
	StatusDataConnectionOpen    = 225
	StatusClosingDataConnection = 226
	StatusPassiveMode           = 227
	StatusLongPassiveMode       = 228
	StatusExtendedPassiveMode   = 229
	StatusLoggedIn              = 230
	StatusLoggedOut             = 231
	StatusLogoutAck             = 232
	StatusAuthOK                = 234
	StatusRequestedFileActionOK = 250
	StatusPathCreated           = 257

	StatusUserOK             = 331
	StatusLoginNeedAccount   = 332
	StatusRequestFilePending = 350

	StatusNotAvailable             = 421
	StatusCanNotOpenDataConnection = 425
	StatusTransfertAborted         = 426
	StatusInvalidCredentials       = 430
	StatusHostUnavailable          = 434
	StatusFileActionIgnored        = 450
	StatusActionAborted            = 451
	Status452                      = 452

	StatusBadCommand              = 500
	StatusBadArguments            = 501
	StatusNotImplemented          = 502
	StatusBadSequence             = 503
	StatusNotImplementedParameter = 504
	StatusNotLoggedIn             = 530
	StatusStorNeedAccount         = 532
	StatusFileUnavailable         = 550
	StatusPageTypeUnknown         = 551
	StatusExceededStorage         = 552
	StatusBadFileName             = 553
)

// pasv will parse the PASV response into an address
func pasvToAddr(line string) (string, error) {
	// PASV response format : 227 Entering Passive Mode (h1,h2,h3,h4,p1,p2).
	start := strings.Index(line, "(")
	end := strings.LastIndex(line, ")")
	if start == -1 || end == -1 {
		return "", errors.New("invalid PASV response format")
	}

	// We have to split the response string
	pasvData := strings.Split(line[start+1:end], ",")
	if len(pasvData) < 6 {
		return "", errors.New("invalid PASV response format")
	}

	// Let's compute the port number
	portPart1, err := strconv.Atoi(pasvData[4])
	if err != nil {
		return "", err
	}

	portPart2, err := strconv.Atoi(pasvData[5])
	if err != nil {
		return "", err
	}

	// Recompose port
	port := portPart1*256 + portPart2

	// Make the IP address to connect to
	host := strings.Join(pasvData[0:4], ".")
	return net.JoinHostPort(host, strconv.Itoa(port)), nil
}

// UpdateOldNames fetches the list of old tz names and returns a mapping
func FTPDownload(target string) (bytes.Buffer, error) {
	var buf bytes.Buffer
	var err error

	// Parse the url
	u, err := url.Parse(target)
	if err != nil {
		return buf, err
	}
	port := u.Port()
	if port == "" {
		port = "21"
	}
	origin := net.JoinHostPort(u.Host, port)

	// Connect to the command server
	conn, err := net.DialTimeout("tcp", origin, 30*time.Second)
	if err != nil {
		return buf, err
	}
	defer conn.Close()
	r := bufio.NewReader(conn)
	if err != nil {
		return buf, err
	}
	resp, err := r.ReadString('\n')
	if err != nil {
		return buf, err
	}
	if !strings.HasPrefix(resp, "220") {
		return buf, fmt.Errorf("failed to connect: %v", resp)
	}

	// Send username
	_, err = conn.Write([]byte("USER anonymous\n"))
	if err != nil {
		return buf, err
	}
	resp, err = r.ReadString('\n')
	if err != nil {
		return buf, err
	}
	if !strings.HasPrefix(resp, "331") {
		return buf, fmt.Errorf("failed to login: %v", resp)
	}

	// Send password
	_, err = conn.Write([]byte("PASS anonymous\n"))
	if err != nil {
		return buf, err
	}
	resp, err = r.ReadString('\n')
	if err != nil {
		return buf, err
	}
	if !strings.HasPrefix(resp, "230") {
		return buf, fmt.Errorf("failed to login: %v", resp)
	}

	// Binary mode
	_, err = conn.Write([]byte("TYPE I\n"))
	if err != nil {
		return buf, err
	}
	resp, err = r.ReadString('\n')
	if err != nil {
		return buf, err
	}
	if !strings.HasPrefix(resp, "200") {
		return buf, fmt.Errorf("failed to switch to binary mode: %v", resp)
	}

	// Get the file transfer address
	_, err = conn.Write([]byte("PASV\n"))
	if err != nil {
		return buf, err
	}
	resp, err = r.ReadString('\n')
	if err != nil {
		return buf, err
	}
	dataAddr, err := pasvToAddr(resp)
	if err != nil {
		return buf, err
	}

	// Connect to the data transfer address
	dataConn, err := net.DialTimeout("tcp", dataAddr, 30*time.Second)
	if err != nil {
		return buf, err
	}
	defer dataConn.Close()

	// Send the command to download the file through the data transfer connection
	_, err = conn.Write([]byte(fmt.Sprintf("RETR %v\n", u.Path)))
	if err != nil {
		return buf, err
	}
	resp, err = r.ReadString('\n')
	if err != nil {
		return buf, err
	}
	if !strings.HasPrefix(resp, "150") {
		return buf, fmt.Errorf("RETR failed: %v", resp)
	}

	// Get the file response
	io.Copy(&buf, dataConn)

	// Get the transfer complete response
	resp, err = r.ReadString('\n')
	if err != nil {
		return buf, err
	}
	if !strings.HasPrefix(resp, "226") {
		return buf, fmt.Errorf("transfer failed: %v", resp)
	}

	return buf, err
}
