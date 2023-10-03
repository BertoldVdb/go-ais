package aisnmea

import (
	"sync"

	"github.com/BertoldVdb/go-ais"
	nmea "github.com/adrianmo/go-nmea"
)

// VdmPacket is a packet that can be encoded into or decoded from a NMEA sentence
type VdmPacket struct {
	Channel     byte
	TalkerID    string
	MessageType string
	Payload     []byte
	Packet      ais.Packet
	TagBlock    nmea.TagBlock
}

type vdmAssemblyWork struct {
	expiryCounter uint64
	received      uint32
	vdms          []*nmea.VDMVDO
}

// VdmAssembler reassmbles split VDO/VDM messages
type vdmAssembler struct {
	lastChannel     byte
	msgCounter      uint64
	nextCleanup     uint64
	cleanupInterval uint64

	msgMap   map[uint32]*vdmAssemblyWork
	msgMutex sync.RWMutex
}

func (v *vdmAssembler) cleanup() {
	v.msgMutex.Lock()
	defer v.msgMutex.Unlock()

	for key, value := range v.msgMap {
		if v.msgCounter >= value.expiryCounter {
			delete(v.msgMap, key)
		}
	}
}

func (v *vdmAssembler) bufferedMessages() int {
	v.msgMutex.RLock()
	defer v.msgMutex.RUnlock()

	result := 0
	for _, k := range v.msgMap {
		result += len(k.vdms)
	}
	return result
}

func (v *vdmAssembler) process(vdm *nmea.VDMVDO) (VdmPacket, bool) {
	if vdm.NumFragments <= 0 ||
		vdm.NumFragments >= 10 ||
		vdm.FragmentNumber > vdm.NumFragments ||
		vdm.FragmentNumber <= 0 {
		return VdmPacket{}, false
	}

	v.msgCounter++

	if v.msgCounter >= v.nextCleanup {
		v.nextCleanup = v.msgCounter + v.cleanupInterval
		v.cleanup()
	}

	/* An empty channel field indicates that the channel is the same as the previous
	   message. I am not sure if this is also allowed for the number of fragments */
	if len(vdm.Channel) > 0 {
		v.lastChannel = vdm.Channel[0]
	}

	/* Is this message a single sentence? */
	if vdm.NumFragments == 1 {
		return VdmPacket{
			Channel:     v.lastChannel,
			TalkerID:    vdm.BaseSentence.TalkerID(),
			MessageType: vdm.BaseSentence.DataType(),
			Payload:     vdm.Payload,
			TagBlock:    vdm.TagBlock,
		}, true
	}

	/* Try to reassemble */
	key := uint32(vdm.NumFragments)
	key |= uint32(vdm.MessageID) << 8
	key |= uint32(v.lastChannel) << 23
	if vdm.Type == nmea.TypeVDO {
		key |= uint32(1) << 31
	}

	v.msgMutex.Lock()
	defer v.msgMutex.Unlock()

	workMsg, ok := v.msgMap[key]
	if !ok {
		workMsg = &vdmAssemblyWork{}
		workMsg.expiryCounter = v.msgCounter + v.cleanupInterval
		workMsg.vdms = make([]*nmea.VDMVDO, 0, vdm.NumFragments)
	}

	workMsg.vdms = append(workMsg.vdms, vdm)
	workMsg.received |= 1 << uint32(vdm.FragmentNumber-1)
	allMsg := uint32(1)<<uint32(vdm.NumFragments) - 1

	if !ok {
		v.msgMap[key] = workMsg
	}

	if len(workMsg.vdms) >= int(vdm.NumFragments) && (workMsg.received&allMsg == allMsg) {
		var fullPayload []byte

		/* Ok, we have all parts, reassemble */
		for i := 0; i < int(vdm.NumFragments); i++ {
			for j := 0; j < len(workMsg.vdms); j++ {
				if workMsg.vdms[j].FragmentNumber-1 == int64(i) {
					fullPayload = append(fullPayload, workMsg.vdms[j].Payload...)
					break
				}
			}
		}

		/* Merge multiple TAG Blocks into a single one */
		var composedTagBlock nmea.TagBlock
		for _, vdm := range workMsg.vdms {
			mergeTagBlocks(&composedTagBlock, &vdm.TagBlock)
		}

		delete(v.msgMap, key)

		/* Full payload is assembled */
		return VdmPacket{
			Channel:     v.lastChannel,
			TalkerID:    vdm.BaseSentence.TalkerID(),
			MessageType: vdm.BaseSentence.DataType(),
			Payload:     fullPayload,
			TagBlock:    composedTagBlock,
		}, true
	}

	return VdmPacket{}, false
}

// VdmAssemblerCreate Creates a VDM/VDO assembler
func vdmAssemblerCreate() *vdmAssembler {
	v := &vdmAssembler{}
	v.lastChannel = 'A'
	v.cleanupInterval = 32
	v.msgMap = make(map[uint32]*vdmAssemblyWork)

	return v
}
