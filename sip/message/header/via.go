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
	schema        string  // schema
	version       float64 // version
	transport     string  // transport
	sentByAddress string  // Sent-by Address
	sentByPort    uint16  // Sent-by port
	rport         uint16  // rport
	branch        string  // branch
	received      string  // received
}

func (via *Via) SetSchema(schema string) {
	via.schema = schema
}
func (via *Via) GetSchema() string {
	return via.schema
}
func (via *Via) SetVersion(version float64) {
	via.version = version
}
func (via *Via) GetVersion() float64 {
	return via.version
}
func (via *Via) SetTransport(transport string) {
	via.transport = transport
}
func (via *Via) GetTransport() string {
	return via.transport
}

func (via *Via) SetSentByAddress(sentByAddress string) {
	via.sentByAddress = sentByAddress
}
func (via *Via) GetSentByAddress() string {
	return via.sentByAddress
}
func (via *Via) SetSentByPort(sentByPort uint16) {
	via.sentByPort = sentByPort
}
func (via *Via) GetSentByPort() uint16 {
	return via.sentByPort
}
func (via *Via) SetRPort(rPort uint16) {
	via.rport = rPort
}
func (via *Via) GetRPort() uint16 {
	return via.rport
}
func (via *Via) SetBranch(branch string) {
	via.branch = branch
}
func (via *Via) GetBranch() string {
	return via.branch
}
func (via *Via) SetReceived(received string) {
	via.received = received
}
func (via *Via) GetReceived() string {
	return via.received
}
func NewVia(schema string, version float64, transport string, sentByAddress string, sentByPort uint16, rPort uint16, branch string, received string) *Via {
	return &Via{
		schema:        schema,
		version:       version,
		transport:     transport,
		sentByAddress: sentByAddress,
		sentByPort:    sentByPort,
		rport:         rPort,
		branch:        branch,
		received:      received,
	}
}
func (via *Via) Raw() (string, error) {
	result := ""
	if err := via.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("Via: %s/%1.1f/%s", strings.ToUpper(via.schema), via.version, strings.ToUpper(via.transport))
	if ip := net.ParseIP(via.sentByAddress); ip != nil {
		result += fmt.Sprintf(" %s:%d", via.sentByAddress, via.sentByPort)
	} else {
		result += fmt.Sprintf(" %s", via.sentByAddress)
	}
	if via.rport == 1 {
		result += fmt.Sprintf(";%s;branch=%s", "rport", via.branch)
	} else if via.rport > 1 {
		result += fmt.Sprintf(";rport=%d;branch=%s;received=%s", via.rport, via.branch, via.received)
	} else {
		result += fmt.Sprintf(";branch=%s", via.branch)
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
	via.schema = strings.ToUpper(schemaRegexp.FindString(schemaAndVersionAndTransportRegexp.FindString(raw)))
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
	via.version = version
	// transport regexp
	transportRegexp := regexp.MustCompile(`(?i)(udp|tcp)`)
	if !transportRegexp.MatchString(schemaAndVersionAndTransportRegexp.FindString(raw)) {
		return errors.New("the value of the transport field cannot match")
	}
	via.transport = strings.ToUpper(transportRegexp.FindString(schemaAndVersionAndTransportRegexp.FindString(raw)))
	raw = schemaAndVersionAndTransportRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the sent-by data cannot be parsed")
	}
	raw = strings.TrimPrefix(raw, ";")
	raw = strings.TrimSuffix(raw, ";")
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
				via.sentByPort = uint16(port)
			}
			via.sentByAddress = portRegexp.ReplaceAllString(raws, "")
			continue
		}
		switch {
		// response port
		case rPortRegexp.MatchString(raws):
			rPorts := regexp.MustCompile(`(?i)(rport)`).ReplaceAllString(rPortRegexp.FindString(raws), "")
			rPorts = regexp.MustCompile(`=`).ReplaceAllString(rPorts, "")
			if len(strings.TrimSpace(rPorts)) == 0 {
				via.rport = 1
			} else {
				rport, err := strconv.Atoi(rPorts)
				if err != nil {
					return err
				}
				via.rport = uint16(rport)
			}
		// branch
		case branchRegexp.MatchString(raws):
			branchs := regexp.MustCompile(`(?i)(branch)`).ReplaceAllString(branchRegexp.FindString(raws), "")
			via.branch = regexp.MustCompile(`=`).ReplaceAllString(branchs, "")
		// received
		case receivedRegexp.MatchString(raws):
			receivedStr := regexp.MustCompile(`(?i)(received)`).ReplaceAllString(receivedRegexp.FindString(raws), "")
			via.received = regexp.MustCompile(`=`).ReplaceAllString(receivedStr, "")
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
	if len(strings.TrimSpace(via.schema)) == 0 {
		return errors.New("the schema field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(sip)$`).MatchString(via.schema) {
		return errors.New("the value of the schema field must be sip")
	}
	if via.version != 2.0 {
		return errors.New("the value of the version field must be 2.0")
	}
	if len(strings.TrimSpace(via.transport)) == 0 {
		return errors.New("the transport field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(udp|tcp)$`).MatchString(via.transport) {
		return errors.New("the value of the transport field must be udp or tcp")
	}
	if len(strings.TrimSpace(via.sentByAddress)) == 0 {
		return errors.New("the sent-by address field is not allowed to be empty")
	}
	if ip := net.ParseIP(via.sentByAddress); ip != nil {
		if via.sentByPort == 0 {
			return errors.New("the sent-by address field gives the ip address, the sent-by port must be given")
		}
	}
	if via.rport > 1 {
		if len(strings.TrimSpace(via.received)) == 0 {
			return errors.New("the rport field gives the response port value, the receive must be given")
		}
	}
	if len(strings.TrimSpace(via.branch)) == 0 {
		return errors.New("the branch field is not allowed to be empty")
	}
	return nil
}

func (via *Via) String() string {
	result := ""
	if len(strings.TrimSpace(via.schema)) > 0 {
		result += fmt.Sprintf("%s/%1.1f", strings.ToUpper(via.schema), via.version)
	} else {
		result += fmt.Sprintf("%1.1f", via.version)
	}
	if len(strings.TrimSpace(via.transport)) > 0 {
		result += fmt.Sprintf("/%s", strings.ToUpper(via.transport))
	}
	if len(strings.TrimSpace(via.sentByAddress)) > 0 {
		result += fmt.Sprintf(" %s", via.sentByAddress)
	}
	if via.sentByPort > 0 {
		result += fmt.Sprintf(":%d", via.sentByPort)
	}
	if via.rport == 1 {
		result += ";rport"
	} else if via.rport > 1 {
		result += fmt.Sprintf(";rport=%d", via.rport)
	}
	if len(strings.TrimSpace(via.branch)) > 0 {
		result += fmt.Sprintf(";branch=%s", via.branch)
	}
	if len(strings.TrimSpace(via.received)) > 0 {
		result += fmt.Sprintf(";received=%s", via.received)
	}
	return result
}
