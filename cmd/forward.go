package cmd

import (
	"github.com/crosstyan/serialForwarder/log"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
	"net"
	"time"
)

const hostFlagName = "host"
const baudRateFlagName = "baudrate"

var host string
var baudRate int

var forwardCmd = cobra.Command{
	Use:   "forward [serial_port] [flags]",
	Short: "Forward serial data to a TCP socket",
	Run:   runForward,
}

func runForward(cmd *cobra.Command, args []string) {
	var serialPortPath string
	if len(args) < 1 {
		_ = cmd.Usage()
		log.Sugar().Warn("Serial port argument is required")
		return
	}
	serialPortPath = args[0]
	log.Sugar().Infow("Input arguments", "serial_port", serialPortPath, "host", host, "baudrate", baudRate)
	// TODO: custom other than 8N1
	mode := serial.Mode{
		BaudRate: baudRate,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	sp, err := serial.Open(serialPortPath, &mode)
	if err != nil {
		log.Sugar().Errorw("Failed to open serial port", "serial_port", serialPortPath, "error", err)
		return
	}
	defer func() {
		if err := sp.Close(); err != nil {
			log.Sugar().Error(err)
		}
	}()
	err = sp.SetReadTimeout(100 * time.Millisecond)
	if err != nil {
		log.Sugar().Errorw("Failed to set read timeout", "serial_port", serialPortPath, "error", err)
	}

	// connect to the server
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Sugar().Errorw("ResolveTCPAddr failed", "host", host, "error", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Sugar().Errorw("DialTCP failed", "host", host, "error", err)
		return
	}
	// TODO: find out the delimiter
	// https://docs.oracle.com/cd/E19509-01/820-5508/ghadt/index.html
	// https://docs.mulesoft.com/hl7-mllp-connector/latest/hl7-mllp-connector-examples
	// https://stackoverflow.com/questions/23988299/tcp-hl7-message-has-period-as-segment-terminator
	// https://learn.microsoft.com/en-us/biztalk/adapters-and-accelerators/accelerator-hl7/message-delimiters?redirectedfrom=MSDN
	// https://learn.microsoft.com/en-us/biztalk/adapters-and-accelerators/accelerator-hl7/processing-hl7-messages?redirectedfrom=MSDN
	spToConn := func() {
		buf := make([]byte, 1024)
		for {
			n, err := sp.Read(buf)
			if err != nil {
				log.Sugar().Error(err)
				return
			}
			b := buf[:n]
			// TODO: string to hex
			log.Sugar().Debugw("Serial port to TCP", "n", n, "data", string(b))
			_, err = conn.Write(b)
			if err != nil {
				log.Sugar().Error(err)
				return
			}
		}
	}
	connToSp := func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Sugar().Error(err)
				return
			}
			b := buf[:n]
			log.Sugar().Debugw("TCP to serial port", "n", n, "data", string(b))
			_, err = sp.Write(b)
			if err != nil {
				log.Sugar().Error(err)
				return
			}
		}
	}
	go spToConn()
	connToSp()
}

func forwardInit() {
	forwardCmd.PersistentFlags().StringVarP(&host, hostFlagName, "H", "localhost", "Host to forward to")
	forwardCmd.PersistentFlags().IntVarP(&baudRate, baudRateFlagName, "b", 115200, "Baud rate for serial port")
}