package header

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Via struct {
	Schema        string  `json:"Schema"`
	Version       float64 `json:"Version"`
	Transport     string  `json:"Transport"`
	SentByAddress string  `json:"Sent-by Address"`
	SentByPort    uint16  `json:"Sent-by port"`
	RPort         uint16  `json:"RPort"`
	Branch        string  `json:"Branch"`
	Received      string  `json:"Received"`
}

func (via *Via) SetSchema(schema string) {
	via.Schema = schema
}
func (via *Via) GetSchema() string {
	return via.Schema
}
func (via *Via) SetVersion(version float64) {
	via.Version = version
}
func (via *Via) GetVersion() float64 {
	return via.Version
}
func (via *Via) SetTransport(transport string) {
	via.Transport = transport
}
func (via *Via) GetTransport() string {
	return via.Transport
}

func (via *Via) SetSentByAddress(sentByAddress string) {
	via.SentByAddress = sentByAddress
}
func (via *Via) GetSentByAddress() string {
	return via.SentByAddress
}
func (via *Via) SetSentByPort(sentByPort uint16) {
	via.SentByPort = sentByPort
}
func (via *Via) GetSentByPort() uint16 {
	return via.SentByPort
}
func (via *Via) SetRPort(rPort uint16) {
	via.RPort = rPort
}
func (via *Via) GetRPort() uint16 {
	return via.RPort
}
func (via *Via) SetBranch(branch string) {
	via.Branch = branch
}
func (via *Via) GetBranch() string {
	return via.Branch
}
func (via *Via) SetReceived(received string) {
	via.Received = received
}
func (via *Via) GetReceived() string {
	return via.Received
}
func NewVia(schema string, version float64, transport string, sentByAddress string, sentByPort uint16, rPort uint16, branch string, received string) *Via {
	return &Via{
		Schema:        schema,
		Version:       version,
		Transport:     transport,
		SentByAddress: sentByAddress,
		SentByPort:    sentByPort,
		RPort:         rPort,
		Branch:        branch,
		Received:      received,
	}
}
func (via *Via) Raw() (string, error) {
	result := ""
	if err := via.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("Via: %s/%1.1f/%s", strings.ToUpper(via.Schema), via.Version, strings.ToUpper(via.Transport))
	if ip := net.ParseIP(via.SentByAddress); ip != nil {
		result += fmt.Sprintf(" %s:%d", via.SentByAddress, via.SentByPort)
	} else {
		result += fmt.Sprintf(" %s", via.SentByAddress)
	}
	if via.RPort == 1 {
		result += fmt.Sprintf(";%s;branch=%s", "rport", via.Branch)
	} else if via.RPort > 1 {
		result += fmt.Sprintf(";rport=%d;branch=%s;received=%s", via.RPort, via.Branch, via.Received)
	} else {
		result += fmt.Sprintf(";branch=%s", via.Branch)
	}
	result += "\r\n"
	return result, nil
}
func (via *Via) Parse(raw string) error {
	if reflect.DeepEqual(nil, via) {
		return errors.New("via caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// via field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(via).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a via header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// schema/version/transport regexp
	schemaAndVersionAndTransportRegexp := regexp.MustCompile(`(?i)(sip)/2\.0/(?i)(udp|tcp) `)
	if !schemaAndVersionAndTransportRegexp.MatchString(raw) {
		return errors.New("the values of the schema, version and transport fields cannot match")
	}
	// schema regexp
	schemaRegexp := regexp.MustCompile(`(?i)(sip)`)
	if !schemaRegexp.MatchString(schemaAndVersionAndTransportRegexp.FindString(raw)) {
		return errors.New("the values of the schema field cannot match")
	}
	via.Schema = strings.ToUpper(schemaRegexp.FindString(schemaAndVersionAndTransportRegexp.FindString(raw)))
	// version regexp
	versionRegexp := regexp.MustCompile(`2\.0`)
	if !versionRegexp.MatchString(schemaAndVersionAndTransportRegexp.FindString(raw)) {
		return errors.New("the value of the version field cannot match")
	}
	versionStr := versionRegexp.FindString(schemaAndVersionAndTransportRegexp.FindString(raw))
	version, err := strconv.ParseFloat(versionStr, 64)
	if err != nil {
		return err
	}
	via.Version = version
	// transport regexp
	transportRegexp := regexp.MustCompile(`(?i)(udp|tcp)`)
	if !transportRegexp.MatchString(schemaAndVersionAndTransportRegexp.FindString(raw)) {
		return errors.New("the value of the transport field cannot match")
	}
	via.Transport = strings.ToUpper(transportRegexp.FindString(schemaAndVersionAndTransportRegexp.FindString(raw)))
	raw = schemaAndVersionAndTransportRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the sent-by data cannot be parsed")
	}
	raw = strings.TrimLeft(raw, ";")
	raw = strings.TrimRight(raw, ";")
	raw = strings.TrimPrefix(raw, ";")
	raw = strings.TrimSuffix(raw, ";")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// port regexp
	portRegexp := regexp.MustCompile(`:\d+`)
	// response port regexp
	rPortRegexp := regexp.MustCompile(`(?i)(rport).*`)
	// branch regexp
	branchRegexp := regexp.MustCompile(`(?i)(branch).*`)
	// received regexp
	receivedRegexp := regexp.MustCompile(`(?i)(received).*`)
	rawSlice := strings.Split(raw, ";")
	for v, raws := range rawSlice {
		// sent-by address and port
		if v == 0 {
			if portRegexp.MatchString(raws) {
				ports := portRegexp.FindString(raws)
				ports = regexp.MustCompile(`:`).ReplaceAllString(ports, "")
				port, err := strconv.Atoi(ports)
				if err != nil {
					return err
				}
				via.SentByPort = uint16(port)
			}
			via.SentByAddress = portRegexp.ReplaceAllString(raws, "")
			continue
		}
		switch {
		// response port
		case rPortRegexp.MatchString(raws):
			rPorts := regexp.MustCompile(`(?i)(rport)`).ReplaceAllString(rPortRegexp.FindString(raws), "")
			rPorts = regexp.MustCompile(`=`).ReplaceAllString(rPorts, "")
			if len(strings.TrimSpace(rPorts)) == 0 {
				via.RPort = 1
			} else {
				rport, err := strconv.Atoi(rPorts)
				if err != nil {
					return err
				}
				via.RPort = uint16(rport)
			}
		// branch
		case branchRegexp.MatchString(raws):
			branchs := regexp.MustCompile(`(?i)(branch)`).ReplaceAllString(branchRegexp.FindString(raws), "")
			via.Branch = regexp.MustCompile(`=`).ReplaceAllString(branchs, "")
		// received
		case receivedRegexp.MatchString(raws):
			receivedStr := regexp.MustCompile(`(?i)(received)`).ReplaceAllString(receivedRegexp.FindString(raws), "")
			via.Received = regexp.MustCompile(`=`).ReplaceAllString(receivedStr, "")
		default:
			// other
			if strings.Contains(raws, "=") {
				fmt.Println(raws)
			}
		}
	}
	return via.Validator()
}
func (via *Via) Validator() error {
	if reflect.DeepEqual(nil, via) {
		return errors.New("via caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(via.Schema)) == 0 {
		return errors.New("the schema field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(sip)$`).MatchString(via.Schema) {
		return errors.New("the value of the schema field must be sip")
	}
	if via.Version != 2.0 {
		return errors.New("the value of the version field must be 2.0")
	}
	if len(strings.TrimSpace(via.Transport)) == 0 {
		return errors.New("the transport field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(udp|tcp)$`).MatchString(via.Transport) {
		return errors.New("the value of the transport field must be udp or tcp")
	}
	if len(strings.TrimSpace(via.SentByAddress)) == 0 {
		return errors.New("the sent-by address field is not allowed to be empty")
	}
	if ip := net.ParseIP(via.SentByAddress); ip != nil {
		if via.SentByPort == 0 {
			return errors.New("the sent-by address field gives the ip address, the sent-by port must be given")
		}
	}
	if via.RPort > 1 {
		if len(strings.TrimSpace(via.Received)) == 0 {
			return errors.New("the rport field gives the response port value, the receive must be given")
		}
	}
	if len(strings.TrimSpace(via.Branch)) == 0 {
		return errors.New("the branch field is not allowed to be empty")
	}
	return nil
}
