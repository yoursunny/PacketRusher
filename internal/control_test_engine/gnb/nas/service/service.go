package service

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/internal/control_test_engine/gnb/nas"
)

func InitServer(gnb *context.GNBContext)  {
	go gnbListen(gnb)
}

func gnbListen(gnb *context.GNBContext) {
	ln := gnb.GetInboundChannel()

	for {
		message := <- ln

		// TODO this region of the code may induces race condition.

		// new instance GNB UE context
		// store UE in UE Pool
		// store UE connection
		// select AMF and get sctp association
		// make a tun interface
		ue := gnb.NewGnBUe(message.GNBTx, message.GNBRx)
		if ue == nil {
			log.Warn("[GNB] UE has not been created")
			break
		}

		// accept and handle connection.
		go processingConn(ue, gnb)
	}
}

func processingConn(ue *context.GNBUe, gnb *context.GNBContext) {
	rx := ue.GetGnbRx()
	for {
		message, done := <- rx

		gnbUeContext, err := gnb.GetGnbUe(ue.GetRanUeId())
		if gnbUeContext == nil || err != nil {
			log.Error("[GNB][NAS] Ignoring message from UE ", ue.GetRanUeId(), " as UE Context was cleaned as requested by AMF.")
			break
		}
		if !done {
			gnbUeContext.SetStateDown()
			break
		}

		// send to dispatch.
		if message.IsNas {
			go nas.Dispatch(ue, message.Nas, gnb)
		} else {
			log.Error("[GNB] Received unknown message from UE")
		}
	}
}
